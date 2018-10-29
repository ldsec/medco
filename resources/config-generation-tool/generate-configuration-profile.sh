#!/bin/bash
set -e
shopt -s nullglob

# dependencies: openssl, keytool (java), docker
# usage: bash generate-dev-configuration-profile.sh CONFIGURATION_PROFILE NODE_IP_1 NODE_IP_2 NODE_IP_3 ...
if [ $# -lt 3 ]
then
    echo "Wrong number of arguments, usage: bash generate-dev-configuration-profile.sh CONFIGURATION_PROFILE NODE_IP_1 NODE_IP_2 NODE_IP_3 ..."
    exit
fi

echo "### Dependencies check, script will abort if dependency if not found"
which docker

# variables & arguments
SCRIPT_FOLDER="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
CONF_PROFILE="$1"
CONF_FOLDER="$SCRIPT_FOLDER/../../configuration-profiles/$CONF_PROFILE"
shift

# clean up previous entries
mkdir -p "$CONF_FOLDER"
rm -f  "$CONF_FOLDER"/*.toml "$CONF_FOLDER"/unlynxMedCo

echo "### Producing medco-unlynx binary with Docker"
docker pull medco/i2b2-unlynx:latest
docker run -v "$CONF_FOLDER":/medco-configuration --entrypoint="" medco/i2b2-unlynx:latest cp -a /go/bin/unlynxMedCo /medco-configuration/unlynxMedCo

# generate configuration for each node
NODE_IDX="-1"
while [ $# -gt 0 ]
do
    NODE_IP="$1"
    shift

    NODE_IDX=$((NODE_IDX+1))
    KEYSTORE_PRIVATE_ALIAS="srv$NODE_IDX-private"

    echo "###$NODE_IDX### Generating unlynx keys"
    "$BUILD_FOLDER"/unlynxMedCo server setupNonInteractive --serverBinding "$NODE_IP:2000" --description "MedCo-UnLynx Server $NODE_IDX" \
        --privateTomlPath "$CONF_FOLDER/srv$NODE_IDX-private.toml" --publicTomlPath "$CONF_FOLDER/srv$NODE_IDX-public.toml"
done

echo "### Generating group.toml file"
cat "$CONF_FOLDER"/srv*-public.toml > "$CONF_FOLDER/group.toml"
"$CONF_FOLDER"/unlynx server getAggregateKey --file "$CONF_FOLDER/group.toml"

echo "### Configuration generated!"