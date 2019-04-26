#!/bin/bash
set -Eeuo pipefail
shopt -s nullglob

# command-line arguments
if [[ $# -ne 3 ]]; then
    echo "Wrong number of arguments, usage: bash $0 <network_name> <node index> <node DNS name>"
    exit 1
fi
NETWORK_NAME="$1"
NODE_IDX="$2"
NODE_DNS_NAME="$3"

# convenience variables
PROFILE_NAME="test-network-$NETWORK_NAME-node$NODE_IDX"
SCRIPT_FOLDER="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
MEDCO_UNLYNX_VER="v0.2.0-rc1"
CONF_FOLDER="$SCRIPT_FOLDER/../../../configuration-profiles/$PROFILE_NAME"
COMPOSE_FOLDER="$SCRIPT_FOLDER/../../../compose-profiles/$PROFILE_NAME"
if [[ -d ${CONF_FOLDER} ]] || [[ -d ${COMPOSE_FOLDER} ]]; then
    echo "The compose and/or configuration profile folder exists. Aborting."
    exit 2
fi

read -p "### About to generate configuration of node $NODE_IDX ($NODE_DNS_NAME) for profile $PROFILE_NAME, <Enter> to continue, <Ctrl+C> to abort."

echo "### Dependencies check, script will abort if dependency if not found"
which docker openssl

# ===================== pre-requisites ======================
mkdir "$CONF_FOLDER" "$COMPOSE_FOLDER"
echo -n "$NODE_DNS_NAME" > "$CONF_FOLDER/srv$NODE_IDX-nodednsname.txt"


# ===================== unlynx keys =========================
echo "### Generating unlynx keys"
docker run -v "$CONF_FOLDER":/medco-configuration -u $(id -u):$(id -g) medco/medco-unlynx:${MEDCO_UNLYNX_VER} \
    server setupNonInteractive --serverBinding "$NODE_DNS_NAME:2000" --description "${PROFILE_NAME}_medco_unlynx_server" \
    --privateTomlPath "/medco-configuration/srv$NODE_IDX-private.toml" \
    --publicTomlPath "/medco-configuration/srv$NODE_IDX-public.toml"
echo "### Unlynx keys generated!"


# ===================== HTTPS cert ==========================
echo "### Generating self-signed HTTPS certificate"
cat > "$SCRIPT_FOLDER/openssl.cnf" <<EOF
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
organizationalUnitName = EPFL LCA1
organizationalUnitName_default = EPFL LCA1
commonName = ${NODE_DNS_NAME}
commonName_max = 64

[ v3_req ]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = ${NODE_DNS_NAME}
EOF

openssl genrsa -out "$CONF_FOLDER/certificate.key" 2048
echo -e "\n\n\n\n\n" | openssl req -new -out "$CONF_FOLDER/certificate.csr" \
    -key "$CONF_FOLDER/certificate.key" -config "$SCRIPT_FOLDER/openssl.cnf"
openssl x509 -req -days 3650 -in "$CONF_FOLDER/certificate.csr" -signkey "$CONF_FOLDER/certificate.key" \
    -out "$CONF_FOLDER/certificate.crt" -extensions v3_req -extfile "$SCRIPT_FOLDER/openssl.cnf"
cp "$CONF_FOLDER/certificate.crt" "$CONF_FOLDER/srv$NODE_IDX-certificate.crt"
rm "$SCRIPT_FOLDER/openssl.cnf"
echo "### Certificate generated!"


# ===================== compose profile =====================
echo "### Generating compose profile"
cp "$SCRIPT_FOLDER/docker-compose.common.yml" "$SCRIPT_FOLDER/docker-compose.node.yml" "$COMPOSE_FOLDER/"
cat > "$COMPOSE_FOLDER/.env" <<EOF
MEDCO_NODE_URL=https://${NODE_DNS_NAME}
MEDCO_NODE_IDX=${NODE_IDX}
MEDCO_PROFILE_NAME=${PROFILE_NAME}
MEDCO_NETWORK_NAME=${NETWORK_NAME}

I2B2_WILDFLY_PASSWORD=admin
I2B2_SERVICE_PASSWORD=pFjy3EjDVwLfT2rB9xkK
I2B2_USER_PASSWORD=demouser
POSTGRES_PASSWORD=postgres1
PGADMIN_DEFAULT_PASSWORD=admin

KEYCLOAK_CLIENT_ID=medco
KEYCLOAK_USER_CLAIM=preferred_username

EOF

if [[ ${NODE_IDX} -eq 0 ]]; then
    echo "### Generating additional configuration for central node"
    cp "$SCRIPT_FOLDER/docker-compose.central.yml" "$COMPOSE_FOLDER/"
    cat >> "$COMPOSE_FOLDER/.env" <<EOF
KEYCLOAK_PASSWORD=keycloak

EOF
fi
echo "### Compose profile generated!"


# ===================== public archive ======================
echo "### Generating archive to be shared"
tar czvf "$CONF_FOLDER/srv$NODE_IDX-public.tar.gz" -C "$CONF_FOLDER" \
    "srv$NODE_IDX-certificate.crt" "srv$NODE_IDX-public.toml" "srv$NODE_IDX-nodednsname.txt"
echo "### DONE! srv$NODE_IDX-public.tar.gz generated, ready for step 2"

