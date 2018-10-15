#!/bin/bash
set -Eeuo pipefail
shopt -s nullglob
# todo: somehow configure as well pic-sure?

# dependencies: openssl, docker
# usage: bash generate-dev-configuration-profile.sh CONFIGURATION_PROFILE KEYSTORE_PASSWORD NODE_IP_1 NODE_IP_2 NODE_IP_3 ...
if [ $# -lt 4 ]
then
    echo "Wrong number of arguments, usage: bash generate-dev-configuration-profile.sh CONFIGURATION_PROFILE NODE_DNS_1 NODE_IP_1 NODE_DNS_2 NODE_IP_2 NODE_DNS_3 NODE_IP_3 ..."
    exit
fi

echo "### Dependencies check, script will abort if dependency if not found"
which openssl docker

# variables & arguments
SCRIPT_FOLDER="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
CONF_PROFILE="$1"
COMPOSE_FOLDER="$SCRIPT_FOLDER/../../compose-profiles/$CONF_PROFILE"
CONF_FOLDER="$SCRIPT_FOLDER/../../configuration-profiles/$CONF_PROFILE"
shift
shift

# clean up previous entries
mkdir -p "$CONF_FOLDER" "$COMPOSE_FOLDER"
rm -f "$CONF_FOLDER"/*.pem "$CONF_FOLDER"/*.toml "$CONF_FOLDER"/unlynxMedCo

echo "### Producing Unlynx binary with Docker"
docker build -t lca1/unlynx:medco-deployment "$SCRIPT_FOLDER"/../../docker-images/unlynx/
docker run -v "$CONF_FOLDER":/medco-configuration --entrypoint sh lca1/unlynx:medco-deployment copy-unlynx-binary.sh

# generate configuration for each node
NODE_IDX="-1"
while [ $# -gt 0 ]
do
    NODE_DNS="$1"
    NODE_IP="$2"
    shift
    shift

    NODE_IDX=$((NODE_IDX+1))

    echo "###$NODE_IDX### Generating unlynx keys"
    "$CONF_FOLDER"/unlynxMedCo server setupNonInteractive --serverBinding "$NODE_IP:2000" --description "Unlynx Server $NODE_IDX" \
        --privateTomlPath "$CONF_FOLDER/srv$NODE_IDX-private.toml" --publicTomlPath "$CONF_FOLDER/srv$NODE_IDX-public.toml"

    echo "###$NODE_IDX### Generating docker-compose file"
    TARGET_COMPOSE_FILE="$COMPOSE_FOLDER/docker-compose-srv$NODE_IDX.yml"
    cp "$SCRIPT_FOLDER/docker-compose-template.yml" "$TARGET_COMPOSE_FILE"
    sed -i "s#_NODE_INDEX_#$NODE_IDX#g" "$TARGET_COMPOSE_FILE"
    sed -i "s#_CONF_PROFILE_#$CONF_PROFILE#g" "$TARGET_COMPOSE_FILE"
done

echo "### Generating group.toml file"
cat "$CONF_FOLDER"/srv*-public.toml > "$CONF_FOLDER/group.toml"

echo "### Configuration generated!"
