#!/bin/bash

##################################################################
# MedCo configuration generator: step 4
# generate unlynx keys & package files to share
##################################################################

set -e
shopt -s nullglob

if [ $# != 3 ]
then
    echo "Usage:"
    echo "Generate certificate with the generated CA:"
    echo "  bash step4.sh CONFIGURATION_PROFILE NODE_INDEX NODE_IP"
    exit
fi

SCRIPT_FOLDER="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"/..
CONF_PROFILE="$1"
CONF_FOLDER="$SCRIPT_FOLDER/../../configuration-profiles/prod/$CONF_PROFILE"
COMPOSE_FOLDER="$SCRIPT_FOLDER/../../compose-profiles/prod/$CONF_PROFILE"
NODE_IDX="$2"
NODE_IP="$3"

# check dependency
which docker


##################################################################
# execute step 4
##################################################################

echo "### Producing Unlynx binary with Docker"
docker build -t lca1/unlynx:medco-deployment "$SCRIPT_FOLDER"/../../docker-images/unlynx/
docker run -v "$CONF_FOLDER":/opt/medco-configuration --entrypoint sh lca1/unlynx:medco-deployment /copy-unlynx-binary.sh

echo "### Generating unlynx keys"
"$CONF_FOLDER"/medco server setupNonInteractive --serverBinding "$NODE_IP:2000" --description "Unlynx Server $NODE_IDX" \
    --privateTomlPath "$CONF_FOLDER/srv$NODE_IDX-private.toml" --publicTomlPath "$CONF_FOLDER/srv$NODE_IDX-public.toml"

echo "### Packaging files to share"
tar -cvzf "$CONF_FOLDER/srv$NODE_IDX-publicdata.tar.gz" \
    -C "$CONF_FOLDER" \
    "srv$NODE_IDX-public.toml" \
    "srv$NODE_IDX-shrine_downstream_nodes.conf" \
    "srv$NODE_IDX-CA/cacert.pem"

echo "### Done! Share the archive srv$NODE_IDX-publicdata.tar.gz with the responsible of the other nodes"
