#!/bin/bash

##################################################################
# MedCo configuration generator: step 5
# aggregation of the files
##################################################################

set -e
shopt -s nullglob

if [ $# -lt 4 ]
then
    echo "Usage:"
    echo "Aggregation of the configuration:"
    echo "  bash step5.sh CONFIGURATION_PROFILE NODE_INDEX KEYSTORE_PASSWORD PUBLIC_DATA_ARCHIVE..."
    exit
fi

SCRIPT_FOLDER="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"/..
CONF_PROFILE="$1"
CONF_FOLDER="$SCRIPT_FOLDER/../../configuration-profiles/prod/$CONF_PROFILE"
COMPOSE_FOLDER="$SCRIPT_FOLDER/../../compose-profiles/prod/$CONF_PROFILE"
NODE_IDX="$2"
KEYSTORE="$CONF_FOLDER/srv$NODE_IDX.keystore"
KEYSTORE_PW="$3"

# check dependency
which keytool


##################################################################
# execute step 5
##################################################################

echo "### Extracting public data of other nodes"
shift
shift
shift
while [ $# -gt 0 ]
do
    tar -xvzf "$1" "$CONF_FOLDER"/
    shift
done

echo "### Aggregating files"
cat "$CONF_FOLDER"/srv*-shrine_downstream_nodes.conf > "$CONF_FOLDER/shrine_downstream_nodes.conf"
cat "$CONF_FOLDER"/srv*-public.toml > "$CONF_FOLDER/group.toml"

echo -n "caCertAliases = [" > "$CONF_FOLDER/shrine_ca_cert_aliases.conf"
I="-1"
for CA_FOLDER in "$CONF_FOLDER"/srv*-CA
do
    I=$((I+1))
    echo -n "\"shrine-ca-srv$I\", " >> "$CONF_FOLDER/shrine_ca_cert_aliases.conf"
    keytool -noprompt -import -v -alias "shrine-ca-srv$I" -file "$CA_FOLDER/cacert.pem" -keystore "$KEYSTORE" -storepass "$KEYSTORE_PW"

done
echo "]" >> "$CONF_FOLDER/shrine_ca_cert_aliases.conf"

echo "### Configuration generated! MedCo is ready to run."
