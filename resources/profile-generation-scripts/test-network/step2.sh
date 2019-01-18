#!/bin/bash
set -Eeuo pipefail
shopt -s nullglob

if [[ $# -ne 2 ]]; then
    echo "Wrong number of arguments, usage: bash $0 <profile_name> <node index>"
    exit 1
fi

NODE_IDX="$2"
PROFILE_NAME="test-network-$1-node$NODE_IDX"

SCRIPT_FOLDER="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
CONF_FOLDER="$SCRIPT_FOLDER/../../../configuration-profiles/$PROFILE_NAME"
COMPOSE_FOLDER="$SCRIPT_FOLDER/../../../compose-profiles/$PROFILE_NAME"

if [[ ! -d ${CONF_FOLDER} ]] || [[ ! -d ${COMPOSE_FOLDER} ]] || [[ -f ${CONF_FOLDER/group.toml} ]]; then
    echo "The compose or configuration profile folder does not exist, or the step 2 has already been executed. Aborting."
    exit 2
fi

read -p "### About to finalize the configuration of node $NODE_IDX for profile $PROFILE_NAME, <Enter> to continue, <Ctrl+C> to abort."

echo "### Extracting archives from other nodes"
pushd "${CONF_FOLDER}"
for archive in $(ls -1 *-public.tar.gz); do tar -zxvf ${archive}; done
popd
echo "### Archives extracted"

echo "### Generating group.toml file"
cat "$CONF_FOLDER"/srv*-public.toml > "$CONF_FOLDER/group.toml"
docker run -v "$CONF_FOLDER":/medco-configuration -u $(id -u):$(id -g) medco/medco-unlynx \
    server getAggregateKey --file "/medco-configuration/group.toml"
docker run -v "$CONF_FOLDER":/medco-configuration -u $(id -u):$(id -g) medco/medco-unlynx \
    server generateTaggingSecrets --file "/medco-configuration/group.toml"
echo "### DONE! MedCo profile ready"
