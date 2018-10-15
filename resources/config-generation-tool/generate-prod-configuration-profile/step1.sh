#!/bin/bash

##################################################################
# MedCo configuration generator: step 1
# init configuration + generate own CA or import CA certificate
##################################################################

set -e
shopt -s nullglob

if [ $# != 4 -a $# != 5 ]
then
    echo "Usage:"
    echo "Generate a certificate authority:"
    echo "  bash step1.sh CONFIGURATION_PROFILE NODE_INDEX KEYSTORE_PASSWORD NODE_DNS"
    echo "Import a certificate authority certificate (PEM file):"
    echo "  bash step1.sh CONFIGURATION_PROFILE NODE_INDEX KEYSTORE_PASSWORD NODE_DNS CA_PUBLIC_KEY_PATH"
    exit
fi

SCRIPT_FOLDER="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"/..
CONF_PROFILE="$1"
CONF_FOLDER="$SCRIPT_FOLDER/../../configuration-profiles/$CONF_PROFILE"
COMPOSE_FOLDER="$SCRIPT_FOLDER/../../compose-profiles/$CONF_PROFILE"
NODE_IDX="$2"
KEYSTORE="$CONF_FOLDER/srv$NODE_IDX.keystore"
KEYSTORE_PW="$3"
NODE_DNS="$4"

# check dependency
which keytool


##################################################################
# execute step 1
##################################################################


echo "### Init configuration"
mkdir "$CONF_FOLDER" "$COMPOSE_FOLDER"
echo "\"Hospital $NODE_IDX\" = \"https://$NODE_DNS:6443/shrine/rest/adapter/requests\"" >> "$CONF_FOLDER/srv$NODE_IDX-shrine_downstream_nodes.conf"

TARGET_COMPOSE_FILE="$COMPOSE_FOLDER/docker-compose-srv$NODE_IDX.yml"
cp "$SCRIPT_FOLDER/docker-compose-template.yml" "$TARGET_COMPOSE_FILE"
sed -i "s#_NODE_INDEX_#$NODE_IDX#g" "$TARGET_COMPOSE_FILE"
sed -i "s#_CONF_PROFILE_#$CONF_PROFILE#g" "$TARGET_COMPOSE_FILE"

if [ $# == 4 ]
then
    echo "### Generating certificate authority"

    # execute CA.sh with -newca, user has the option to import existing CA certificate (with the priv. key only though)
    CATOP="$CONF_FOLDER/srv$NODE_IDX-CA" "$SCRIPT_FOLDER"/CA.sh -newca

    # import CA into the keystore
    keytool -noprompt -import -v -alias "shrine-ca-srv$NODE_IDX" -file "$CONF_FOLDER/srv$NODE_IDX-CA/cacert.pem" \
        -keystore "$KEYSTORE" -storepass "$KEYSTORE_PW"

elif [ $# == 5 ]
then
    echo "### Importing certificate authority certificate"
    cp "$5" "$CONF_FOLDER/srv$NODE_IDX-CA/cacert.pem"

fi
