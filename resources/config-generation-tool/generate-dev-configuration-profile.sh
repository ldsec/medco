#!/bin/bash
set -e
shopt -s nullglob

# dependencies: openssl, keytool (java), docker
# usage: bash generate-dev-configuration-profile.sh CONFIGURATION_PROFILE KEYSTORE_PASSWORD NODE_DNS_1 NODE_IP_1 NODE_DNS_2 NODE_IP_2 NODE_DNS_3 NODE_IP_3 ...
if [ $# -lt 5 ]
then
    echo "Wrong number of arguments, usage: bash generate-dev-configuration-profile.sh CONFIGURATION_PROFILE KEYSTORE_PASSWORD NODE_DNS_1 NODE_IP_1 NODE_DNS_2 NODE_IP_2 NODE_DNS_3 NODE_IP_3 ..."
    exit
fi

echo "### Dependencies check, script will abort if dependency if not found"
which openssl keytool docker

# variables & arguments
SCRIPT_FOLDER="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
CONF_PROFILE="$1"
COMPOSE_FOLDER="$SCRIPT_FOLDER/../../compose-profiles/prod/$CONF_PROFILE"
CONF_FOLDER="$SCRIPT_FOLDER/../../configuration-profiles/prod/$CONF_PROFILE"
BUILD_FOLDER="$SCRIPT_FOLDER/../../configuration-profiles/dev"
KEYSTORE_PW="$2"
shift
shift

# clean up previous entries
mkdir -p "$CONF_FOLDER" "$COMPOSE_FOLDER"
rm -f "$CONF_FOLDER"/*.keystore "$CONF_FOLDER"/shrine_ca_cert_aliases.conf "$CONF_FOLDER"/shrine_downstream_nodes.conf \
    "$CONF_FOLDER"/*.pem "$CONF_FOLDER"/*.toml "$CONF_FOLDER"/unlynx
rm -rf "$CONF_FOLDER"/srv*-CA

echo "### Producing Unlynx binary with Docker"
docker build -t lca1/unlynx:medco-deployment "$SCRIPT_FOLDER"/../../docker-images/dev/unlynx/
docker run -v "$BUILD_FOLDER":/opt/medco-configuration --entrypoint sh lca1/unlynx:medco-deployment /copy-unlynx-binary-dev.sh

echo "caCertAliases = [\"shrine-ca\"]" > "$CONF_FOLDER/shrine_ca_cert_aliases.conf"
echo "### Producing CA"
CATOP="$CONF_FOLDER/CA" "$SCRIPT_FOLDER"/CA.sh -newca
echo "unique_subject = no" > "$CONF_FOLDER/CA/index.txt.attr"

# generate configuration for each node
NODE_IDX="-1"
while [ $# -gt 0 ]
do
    NODE_DNS="$1"
    NODE_IP="$2"
    shift
    shift

    NODE_IDX=$((NODE_IDX+1))
    KEYSTORE="$CONF_FOLDER/srv$NODE_IDX.keystore"
    KEYSTORE_PRIVATE_ALIAS="srv$NODE_IDX-private"

    echo "### Setting up certificate authority and import it in keystore"
    keytool -noprompt -import -v -alias "shrine-ca" -file "$CONF_FOLDER/CA/cacert.pem" -keystore "$KEYSTORE" -storepass "$KEYSTORE_PW"

    echo "###$NODE_IDX### Generating java keystore pair of keys"
    keytool -genkeypair -keysize 2048 -alias "$KEYSTORE_PRIVATE_ALIAS" -validity 7300 \
        -dname "CN=$NODE_DNS, OU=LCA1, O=EPFL, L=Lausanne, S=VD, C=CH" \
        -ext "SAN=DNS:$NODE_DNS,IP:$NODE_IP" \
        -keyalg RSA -keypass "$KEYSTORE_PW" -storepass "$KEYSTORE_PW" -keystore "$KEYSTORE"

    echo "###$NODE_IDX### Generating certificate signature request"
    keytool -certreq -alias "$KEYSTORE_PRIVATE_ALIAS" -keyalg RSA -file "$SCRIPT_FOLDER/newreq.pem" -keypass "$KEYSTORE_PW" \
        -storepass "$KEYSTORE_PW" -keystore "$KEYSTORE" -ext "SAN=DNS:$NODE_DNS,IP:$NODE_IP"
    cat > "$SCRIPT_FOLDER/openssl.ext.tmp.cnf" <<EOL
        basicConstraints=CA:FALSE
        subjectAltName=@alt_names
        subjectKeyIdentifier = hash

        [ alt_names ]
        IP.1 = $NODE_IP
        DNS.1 = $NODE_DNS
EOL

    echo "###$NODE_IDX### Signing it with the CA"
    CATOP="$CONF_FOLDER/CA" SSLEAY_CONFIG="-extfile $SCRIPT_FOLDER/openssl.ext.tmp.cnf" "$SCRIPT_FOLDER"/CA.sh -sign

    echo "###$NODE_IDX### Importing in keystore own certificate signed by CA (chained to the private key)"
    keytool -noprompt -import -v -alias "$KEYSTORE_PRIVATE_ALIAS" -file "$SCRIPT_FOLDER"/newcert.pem -keystore "$KEYSTORE" -storepass "$KEYSTORE_PW" \
        -keypass "$KEYSTORE_PW" -trustcacerts

    echo "###$NODE_IDX### Generating pem certificates (lighttpd)"
    keytool -noprompt -importkeystore -srckeystore "$KEYSTORE" -srcalias "$KEYSTORE_PRIVATE_ALIAS" -destkeystore "$KEYSTORE".p12 \
        -deststoretype PKCS12 -srcstorepass "$KEYSTORE_PW" -deststorepass "$KEYSTORE_PW"
    openssl pkcs12 -in "$KEYSTORE".p12 -out "$CONF_FOLDER/srv$NODE_IDX.pem" -password pass:"$KEYSTORE_PW" -nodes

    echo "###$NODE_IDX### Adding entry in the downstream nodes config file"
    echo "\"Hospital $NODE_IDX\" = \"https://$NODE_DNS:6443/shrine/rest/adapter/requests\"" >> "$CONF_FOLDER/shrine_downstream_nodes.conf"

    echo "###$NODE_IDX### Generating unlynx keys"
    "$CONF_FOLDER"/unlynx server setupNonInteractive --serverBinding "$NODE_IP:2000" --description "Unlynx Server $NODE_IDX" \
        --privateTomlPath "$CONF_FOLDER/srv$NODE_IDX-private.toml" --publicTomlPath "$CONF_FOLDER/srv$NODE_IDX-public.toml"

    echo "###$NODE_IDX### Generating docker-compose file"
    TARGET_COMPOSE_FILE="$COMPOSE_FOLDER/docker-compose-srv$NODE_IDX.yml"
    cp "$SCRIPT_FOLDER/docker-compose-template.yml" "$TARGET_COMPOSE_FILE"
    sed -i "s#_NODE_INDEX_#$NODE_IDX#g" "$TARGET_COMPOSE_FILE"
    sed -i "s#_CONF_PROFILE_#$CONF_PROFILE#g" "$TARGET_COMPOSE_FILE"

    echo "###$NODE_IDX### Cleaning up"
    rm "$SCRIPT_FOLDER/newreq.pem" "$SCRIPT_FOLDER/openssl.ext.tmp.cnf" "$KEYSTORE".p12 "$SCRIPT_FOLDER/newcert.pem"

    # keytool -list -v -keystore "$KEYSTORE" -storepass "$KEYSTORE_PW" # list content of keystore (disabled)
done

echo "### Generating group.toml file and finalizing shrine config file"
cat "$CONF_FOLDER"/srv*-public.toml > "$CONF_FOLDER/group.toml"
"$CONF_FOLDER"/unlynx server getAggregateKey --file "$CONF_FOLDER/group.toml"

echo "### Configuration generated!"
