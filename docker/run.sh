#!/usr/bin/env bash

SCRIPT_FOLDER="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

docker build -t lca1/medco-loader:medco-loader .
docker run -v "$SCRIPT_FOLDER":/opt/build-dir --entrypoint sh lca1/medco-loader:medco-loader /copy-medco-loader-binary.sh
mv medco-loader ../app/medco-loader