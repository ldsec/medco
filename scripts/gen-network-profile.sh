#!/bin/bash
set -Eeo pipefail
shopt -s nullglob

NETWORK_NAME="vmtest-ip"

# clean up
SCRIPT_FOLDER="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
pushd "${SCRIPT_FOLDER}/../../deployments"
rm -rf "network-${NETWORK_NAME:?}-"*
popd

# generate step 1
export MEDCO_SETUP_VER=dev
for IDX in 0 1 2; do
  bash step1.sh -nn "$NETWORK_NAME" -ni "$IDX" -ha "192.168.56.11${IDX}" -ua "192.168.57.11${IDX}:2001"
  #bash step1.sh -nn "$NETWORK_NAME" -ni "$IDX" -ha "test-medco-http-node${IDX}.misba.ch" -ua "192.168.57.11${IDX}:2001"
done

# share
pushd "${SCRIPT_FOLDER}/../../deployments"
PUB_FOLDER="$(pwd)/network-${NETWORK_NAME:?}-public"
mkdir "$PUB_FOLDER"

for profile_folder in "network-${NETWORK_NAME:?}-node"*; do
  [[ -e "$profile_folder" ]] # if no archive
  pushd "$profile_folder/configuration"
  cp srv*-public.tar.gz "$PUB_FOLDER/"
  popd
done

for profile_folder in "network-${NETWORK_NAME:?}-node"*; do
  [[ -e "$profile_folder" ]] # if no archive
  cp "$PUB_FOLDER/"* "$profile_folder/configuration/"
done
popd

# generate step 2
for IDX in 0 1 2; do
  bash step2.sh -nn "$NETWORK_NAME" -ni "$IDX" -nb 3
done
