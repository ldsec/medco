#!/bin/bash

##################################################################
# MedCo configuration generator: step 3
# generate certificate of the node or import it
##################################################################

set -e
shopt -s nullglob

if [ $# != 4 -a $# != 5 ]
then
    echo "Usage:"
    echo "Generate certificate with the generated CA:"
    echo "  bash step3.sh CONFIGURATION_PROFILE NODE_INDEX KEYSTORE_PASSWORD NODE_DNS NODE_IP"
    echo "Import certificate of previously imported keypair:"
    echo "  bash step3.sh CONFIGURATION_PROFILE NODE_INDEX KEYSTORE_PASSWORD CERT_FILE_PATH"
    exit
fi

SCRIPT_FOLDER="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"/..
CONF_PROFILE="$1"
CONF_FOLDER="$SCRIPT_FOLDER/../../configuration-profiles/$CONF_PROFILE"
COMPOSE_FOLDER="$SCRIPT_FOLDER/../../compose-profiles/$CONF_PROFILE"
NODE_IDX="$2"
KEYSTORE_PW="$3"

# check dependency
which keytool openssl


##################################################################
# execute step 3
##################################################################

KEYSTORE="$CONF_FOLDER/srv$NODE_IDX.keystore"
KEYSTORE_PRIVATE_ALIAS="srv$NODE_IDX-private"

if [ $# == 5 ]
then
    NODE_DNS="$4"
    NODE_IP="$5"

    echo "### Generating certificate signature request"
    keytool -certreq -alias "$KEYSTORE_PRIVATE_ALIAS" -keyalg RSA -file "$SCRIPT_FOLDER/newreq.pem" -keypass "$KEYSTORE_PW" \
        -storepass "$KEYSTORE_PW" -keystore "$KEYSTORE" -ext "SAN=DNS:$NODE_DNS,IP:$NODE_IP"

    # openssl additional configuration
    cat > "$SCRIPT_FOLDER/openssl.ext.tmp.cnf" <<EOL
        basicConstraints=CA:FALSE
        subjectAltName=@alt_names
        subjectKeyIdentifier = hash

        [ alt_names ]
        IP.1 = $NODE_IP
        DNS.1 = $NODE_DNS
EOL

    echo "###$NODE_IDX### Signing it with the CA"
    CATOP="$CONF_FOLDER/srv$NODE_IDX-CA" SSLEAY_CONFIG="-extfile $SCRIPT_FOLDER/openssl.ext.tmp.cnf" "$SCRIPT_FOLDER"/CA.sh -sign

    echo "###$NODE_IDX### Importing in keystore own certificate signed by CA (chained to the private key)"
    keytool -noprompt -import -v -alias "$KEYSTORE_PRIVATE_ALIAS" -file "$SCRIPT_FOLDER"/newcert.pem -keystore "$KEYSTORE" \
        -storepass "$KEYSTORE_PW" -keypass "$KEYSTORE_PW" -trustcacerts

    echo "###$NODE_IDX### Generating pem certificates (lighttpd)"
    keytool -noprompt -importkeystore -srckeystore "$KEYSTORE" -srcalias "$KEYSTORE_PRIVATE_ALIAS" -destkeystore "$KEYSTORE".p12 \
        -deststoretype PKCS12 -srcstorepass "$KEYSTORE_PW" -deststorepass "$KEYSTORE_PW"
    openssl pkcs12 -in "$KEYSTORE".p12 -out "$CONF_FOLDER/srv$NODE_IDX.pem" -password pass:"$KEYSTORE_PW" -nodes

    # cleanup
    rm "$SCRIPT_FOLDER/newreq.pem" "$SCRIPT_FOLDER/openssl.ext.tmp.cnf" "$KEYSTORE.p12" "$SCRIPT_FOLDER/newcert.pem"

elif [ $# == 4 ]
then
    echo "### Importing certificate"
    echo "NOT IMPLEMENTED"
    exit
    # todo
fi
