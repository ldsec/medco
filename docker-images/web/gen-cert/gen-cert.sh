#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_FOLDER="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $SCRIPT_FOLDER

openssl genrsa -out ../nginx-conf.d/cert.key 2048
echo -e "\n\n\n\n\n" | openssl req -new -out ../nginx-conf.d/cert.csr -key ../nginx-conf.d/cert.key -config ./openssl.cnf
openssl x509 -req -days 3650 -in ../nginx-conf.d/cert.csr -signkey ../nginx-conf.d/cert.key -out ../nginx-conf.d/cert.crt -extensions v3_req -extfile ./openssl.cnf
