#!/bin/bash
set -Eeo pipefail
shopt -s nullglob

# ===================== input parsing ==================
function example {
    echo -e "example: $0 -nn <network name> -ni <node index> -dns <node_dns_name> [-crt <certificate_path> -key <key_path>]"
}

function usage {
    echo -e "Wrong arguments, usage: bash $0 MANDATORY [OPTIONAL]"
}

function help {
    echo -e "MANDATORY:"
    echo -e "  -nn,   --network_name  VAL  Network name (e.g. test-network-deployment)"
    echo -e "  -ni,   --node_index    VAL  Node index (e.g. 0, 1, 2)"
    echo -e "  -dns,  --node_dns_name VAL  Server dns name\n"
    echo -e "OPTIONAL:"
    echo -e "  -pk,   --public_key    VAL  Unlynx node public key"
    echo -e "  -sk,   --secret_key    VAL  Unlynx node private key"
    echo -e "  -crt,  --certificate   VAL  Filepath to certificate (*.crt)"
    echo -e "  -k,    --key           VAL  Filepath to certificate key (*.key)"
    echo -e "  -h,    --help \n"
    example
}

#Declare the number of mandatory args
margs=3

# Ensures that the number of passed args are at least equals to the declared number of mandatory args.
# It also handles the special case of the -h or --help arg.
function margs_precheck {
	if [ "$2" ] && [ "$1" -lt $margs ]; then
		if [ "$2" == "--help" ] || [ "$2" == "-h" ]; then
			help
			exit
		else
	    usage
	    help
	    exit 1
		fi
	fi
}

# Ensures that all the mandatory args are not empty
function margs_check {
	if [ $# -lt $margs ]; then
	    usage
	    help
	    exit 1 # error
	fi
}

# check if no inputs where selected
if [ $# -lt 1 ]; then
  usage
  help
  exit 1
fi
margs_precheck $# "$1"

# default values
NETWORK_NAME=
NODE_IDX=
NODE_DNS_NAME=
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
	 -dns  | --node_dns_name )  shift
     						          NODE_DNS_NAME=$1
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
margs_check "$NETWORK_NAME" "$NODE_IDX" "$NODE_DNS_NAME"

set -u
if [[ ! $NETWORK_NAME =~ ^[a-zA-Z0-9-]+$ ]]; then
    echo "Network name must only contain basic characters (a-z, A-Z, 0-9, -)"
    exit 1
fi

# convenience variables
PROFILE_NAME="network-${NETWORK_NAME}-node${NODE_IDX}"
SCRIPT_FOLDER="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
MEDCO_DOCKER="ghcr.io/ldsec/medco:${MEDCO_SETUP_VER:-$(shell make --no-print-directory -C ../../ medco_version)}"
COMPOSE_FOLDER="${SCRIPT_FOLDER}/../../deployments/${PROFILE_NAME}"
CONF_FOLDER="${COMPOSE_FOLDER}/configuration"
if [[ -d ${COMPOSE_FOLDER} ]]; then
    echo "The profile folder exists. Aborting."
    exit 2
fi

read -rp "### About to generate configuration of node ${NODE_IDX} (${NODE_DNS_NAME}) for profile ${PROFILE_NAME}, <Enter> to continue, <Ctrl+C> to abort."

echo "### Dependency on Docker check, script will abort if not found"
which docker
echo "### Dependency on OpenSSL check, script will abort if not found"
which openssl

# ===================== pre-requisites ======================
mkdir "${COMPOSE_FOLDER}" "${CONF_FOLDER}"
echo -n "${NODE_DNS_NAME}" > "${CONF_FOLDER}/srv${NODE_IDX}-nodednsname.txt"


# ===================== unlynx keys =========================
echo "### Generating unlynx keys"
if [[ -z ${PUB_KEY} ]]; then
    docker run -v "$CONF_FOLDER:/medco-configuration" -u "$(id -u):$(id -g)" "${MEDCO_DOCKER}" medco-unlynx \
        server setupNonInteractive --serverBinding "${NODE_DNS_NAME}:2001" --description "${PROFILE_NAME}_medco_unlynx_server" \
        --privateTomlPath "/medco-configuration/srv${NODE_IDX}-private.toml" \
        --publicTomlPath "/medco-configuration/srv${NODE_IDX}-public.toml"
else
    docker run -v "$CONF_FOLDER:/medco-configuration" -u "$(id -u):$(id -g)" "${MEDCO_DOCKER}" medco-unlynx \
        server setupNonInteractive --serverBinding "${NODE_DNS_NAME}:2001" --description "${PROFILE_NAME}_medco_unlynx_server" \
        --privateTomlPath "/medco-configuration/srv${NODE_IDX}-private.toml" \
        --publicTomlPath "/medco-configuration/srv${NODE_IDX}-public.toml" \
        --pubKey "${PUB_KEY}" --privKey "${PRIV_KEY}"
fi
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
organizationalUnitName = EPFL LDS
organizationalUnitName_default = EPFL LDS
commonName = ${NODE_DNS_NAME}
commonName_max = 64

[ v3_req ]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = ${NODE_DNS_NAME}
EOF

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

rm -rf "${COMPOSE_FOLDER}"
echo "You must input both filepath to *.crt (-crt) and *.key (-k)."
exit 1

fi

# ===================== compose profile =====================
echo "### Generating compose profile"
cp "${SCRIPT_FOLDER}/docker-compose.yml" "${SCRIPT_FOLDER}/docker-compose.tools.yml" "${SCRIPT_FOLDER}/Makefile" "${COMPOSE_FOLDER}/"
cat > "${COMPOSE_FOLDER}/.env" <<EOF
MEDCO_NODE_DNS_NAME=${NODE_DNS_NAME}
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
