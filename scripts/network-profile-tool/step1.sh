#!/bin/bash
set -Eeo pipefail
shopt -s nullglob

source common.sh

# ===================== input parsing ==================
function example {
    echo -e "example: $0 -nn <network name> -ni <node index> -ha <HTTPS address> [-ua <unlynx address> -crt <certificate_path> -k <key_path>]"
}

function help {
    echo -e "MANDATORY:"
    echo -e "  -nn,   --network_name    VAL  Network name (e.g. test-network-deployment)"
    echo -e "  -ni,   --node_index      VAL  Node index (e.g. 0, 1, 2)"
    echo -e "  -ha,   --https_address   VAL  Node HTTPS address, either DNS name or IP address (e.g. test.medco.com or 192.168.43.22)\n"
    echo -e "OPTIONAL:"
    echo -e "  -ua,   --unlynx_address  VAL  Unlynx address, either DNS name or IP address, if different from node HTTPS address (e.g. test-unlynx.medco.com or 192.168.65.55)"
    echo -e "  -pk,   --public_key      VAL  Unlynx node public key, if it is not to be generated"
    echo -e "  -sk,   --secret_key      VAL  Unlynx node private key, if it is not to be generated"
    echo -e "  -crt,  --certificate     VAL  Filepath to certificate (*.crt), if it is not to be generated"
    echo -e "  -k,    --key             VAL  Filepath to certificate key (*.key), if it is not to be generated"
    echo -e "  -h,    --help \n"
    example
}

margs=3 # number of mandatory args
margs_precheck $# "$1" $margs

# default values
NETWORK_NAME=
NODE_IDX=
HTTPS_ADDRESS=
UNLYNX_ADDRESS=
PUB_KEY=
PRIV_KEY=
CRT=
KEY=

# Args while-loop
while [ "$1" != "" ];
do
   case $1 in
   -nn  | --network_name )  shift
                          NETWORK_NAME=$1
                		      ;;
   -ni  | --node_index )  shift
   						            NODE_IDX=$(printf "%03d" "$1")
			                    ;;
	 -ha  | --https_address )  shift
     						          HTTPS_ADDRESS=$1
  			                  ;;
   -ua  | --unlynx_address  )  shift
                          UNLYNX_ADDRESS=$1
                          ;;
   -pk  | --public_key  )  shift
                          PUB_KEY=$1
                          ;;
   -sk  | --secret_key  )  shift
                          PRIV_KEY=$1
                          ;;
   -crt  | --certificate  )  shift
                          CRT=$1
                          ;;
   -k  | --key  )  shift
                          KEY=$1
                          ;;
   -h   | --help )        help
                          exit
                          ;;
   *)
                          echo "$0: illegal option $1"
                          usage
                          help
						              exit 1
                          ;;
    esac
    shift
done

# Check if all mandatory args have assigned values
margs_check $margs "$NETWORK_NAME" "$NODE_IDX" "$HTTPS_ADDRESS"
set -u
check_network_name "$NETWORK_NAME"

# parse addresses
HTTPS_IP_ADDRESS=
HTTPS_DNS_NAME=
if [[ $HTTPS_ADDRESS =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "### HTTPS address is an IP address"
  HTTPS_IP_ADDRESS=$HTTPS_ADDRESS
else
  echo "### HTTPS address is a DNS name"
  HTTPS_DNS_NAME=$HTTPS_ADDRESS
fi

if [[ -z "$UNLYNX_ADDRESS" ]]; then
  echo "### Unlynx address defaults to HTTPS address"
  UNLYNX_ADDRESS="${HTTPS_ADDRESS}:2001"
else
  echo "### Unlynx address was provided and is different from HTTPS address"
  UNLYNX_ADDRESS="${UNLYNX_ADDRESS}:2001"
fi

# generate convenience variables
export_variables "$NETWORK_NAME" "$NODE_IDX"
if [[ -d ${COMPOSE_FOLDER} ]]; then
    echo "The profile folder exists. Aborting."
    exit 2
fi

echo "### About to generate configuration of node ${NODE_IDX} (${HTTPS_ADDRESS}) for profile ${PROFILE_NAME}"
read -rp "### <Enter> to continue, <Ctrl+C> to abort."
dependency_check

# ===================== pre-requisites ======================
mkdir "${COMPOSE_FOLDER}" "${CONF_FOLDER}"
echo -n "${HTTPS_ADDRESS}" > "${CONF_FOLDER}/srv${NODE_IDX}-nodednsname.txt"

# ===================== unlynx keys =========================
unlynx_setup_args=(
  medco-unlynx server setupNonInteractive
  --serverBinding "$UNLYNX_ADDRESS"
  --wsUrl "http://medco-unlynx:2002"
  --description "${PROFILE_NAME}_medco_unlynx_server"
  --privateTomlPath "/medco-configuration/srv${NODE_IDX}-private.toml"
  --publicTomlPath "/medco-configuration/srv${NODE_IDX}-public.toml"
)

echo "### Generating unlynx keys with address ${UNLYNX_ADDRESS}"
if [[ -n "$PUB_KEY" ]]; then
  echo "### Using pre-generated unlynx key ${PUB_KEY}"
  unlynx_setup_args=("${unlynx_setup_args[@]}" --pubKey "${PUB_KEY}" --privKey "${PRIV_KEY}")
fi

"${MEDCO_BIN[@]}" "${unlynx_setup_args[@]}"
echo "### Unlynx keys generated!"

if [ -z "$CRT" ] && [ -z "$KEY" ]; then

# ===================== self-signed HTTPS cert ==========================
echo "### Generating self-signed HTTPS certificate"
cat > "${SCRIPT_FOLDER}/openssl.cnf" <<EOF
[req]
distinguished_name = req_distinguished_name
req_extensions = v3_req

[req_distinguished_name]
countryName = CH
countryName_default = CH
stateOrProvinceName = Vaud
stateOrProvinceName_default = Vaud
localityName = Lausanne
localityName_default = Lausanne
organizationalUnitName = MedCo
organizationalUnitName_default = MedCo
commonName = ${HTTPS_ADDRESS}
commonName_max = 64

[ v3_req ]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
subjectAltName = @alt_names

[alt_names]
EOF

if [[ -n "$HTTPS_IP_ADDRESS" ]]; then
  echo "IP.1 = ${HTTPS_IP_ADDRESS}" >> "${SCRIPT_FOLDER}/openssl.cnf"
else
  echo "DNS.1 = ${HTTPS_DNS_NAME}" >> "${SCRIPT_FOLDER}/openssl.cnf"
fi

openssl genrsa -out "${CONF_FOLDER}/certificate.key" 2048
echo -e "\n\n\n\n\n" | openssl req -new -out "${CONF_FOLDER}/certificate.csr" \
    -key "${CONF_FOLDER}/certificate.key" -config "${SCRIPT_FOLDER}/openssl.cnf"
openssl x509 -req -days 3650 -in "${CONF_FOLDER}/certificate.csr" -signkey "${CONF_FOLDER}/certificate.key" \
    -out "${CONF_FOLDER}/certificate.crt" -extensions v3_req -extfile "${SCRIPT_FOLDER}/openssl.cnf"
cp "${CONF_FOLDER}/certificate.crt" "${CONF_FOLDER}/srv${NODE_IDX}-certificate.crt"
rm "${SCRIPT_FOLDER}/openssl.cnf"
echo "### Self-signed certificate generated!"

elif [ -n "$CRT" ] && [ -n "$KEY" ]; then

# ===================== HTTPS cert ==========================
cp "$CRT" "${CONF_FOLDER}/certificate.crt"
cp "${CONF_FOLDER}/certificate.crt" "${CONF_FOLDER}/srv${NODE_IDX}-certificate.crt"
echo "### Certificate selected!"
cp "$KEY" "${CONF_FOLDER}/certificate.key"
echo "### Key selected!"

else

echo "You must input both filepath to *.crt (-crt) and *.key (-k)."
exit 1

fi

# ===================== compose profile =====================
echo "### Generating compose profile"
cp "${SCRIPT_FOLDER}/docker-compose.yml" "${SCRIPT_FOLDER}/docker-compose.tools.yml" "${SCRIPT_FOLDER}/Makefile" "${COMPOSE_FOLDER}/"
cat > "${COMPOSE_FOLDER}/.env" <<EOF
MEDCO_NODE_DNS_NAME=${HTTPS_ADDRESS}
MEDCO_NODE_IDX=${NODE_IDX}
MEDCO_PROFILE_NAME=${PROFILE_NAME}
MEDCO_NETWORK_NAME=${NETWORK_NAME}

I2B2_WILDFLY_PASSWORD=replaceme
I2B2_SERVICE_PASSWORD=replaceme
I2B2_USER_PASSWORD=replaceme
POSTGRES_PASSWORD=replaceme
KEYCLOAK_PASSWORD=replaceme

KEYCLOAK_REALM=master
KEYCLOAK_CLIENT_ID=medco
KEYCLOAK_USER_CLAIM=preferred_username

EOF
echo "### Compose profile generated!"


# ===================== public archive ======================
echo "### Generating archive to be shared"
tar czvf "${CONF_FOLDER}/srv${NODE_IDX}-public.tar.gz" -C "${CONF_FOLDER}" \
    "srv${NODE_IDX}-certificate.crt" "srv${NODE_IDX}-public.toml" "srv${NODE_IDX}-nodednsname.txt"
echo "### DONE! srv${NODE_IDX}-public.tar.gz generated, ready for step 2"
