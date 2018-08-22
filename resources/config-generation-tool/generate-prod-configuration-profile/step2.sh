#!/bin/bash

##################################################################
# MedCo configuration generator: step 2
# generate keypair of the node or import it
##################################################################

set -e
shopt -s nullglob

if [ $# != 4 -a $# != 5 ]
then
    echo "Usage:"
    echo "Generate pair of keys:"
    echo "  bash step2.sh CONFIGURATION_PROFILE NODE_INDEX KEYSTORE_PASSWORD NODE_DNS NODE_IP"
    echo "Import pair of keys:"
    echo "  bash step2.sh CONFIGURATION_PROFILE NODE_INDEX KEYSTORE_PASSWORD KEY_FILE_PATH"
    exit
fi

SCRIPT_FOLDER="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"/..
CONF_PROFILE="$1"
CONF_FOLDER="$SCRIPT_FOLDER/../../configuration-profiles/prod/$CONF_PROFILE"
COMPOSE_FOLDER="$SCRIPT_FOLDER/../../compose-profiles/prod/$CONF_PROFILE"
NODE_IDX="$2"
KEYSTORE_PW="$3"

# check dependency
which keytool


##################################################################
# execute step 2
##################################################################

KEYSTORE="$CONF_FOLDER/srv$NODE_IDX.keystore"
KEYSTORE_PRIVATE_ALIAS="srv$NODE_IDX-private"

if [ $# == 5 ]
then
    NODE_DNS="$4"
    NODE_IP="$5"

    echo "### Generating java keystore pair of keys"
    keytool -genkeypair -keysize 2048 -alias "$KEYSTORE_PRIVATE_ALIAS" -validity 7300 \
        -dname "CN=$NODE_DNS" -ext "SAN=DNS:$NODE_DNS,IP:$NODE_IP" \
        -keyalg RSA -keypass "$KEYSTORE_PW" -storepass "$KEYSTORE_PW" -keystore "$KEYSTORE"

elif [ $# == 4 ]
then
    echo "### Importing pair of keys"
    echo "NOT IMPLEMENTED"
    exit
    # todo
fi
