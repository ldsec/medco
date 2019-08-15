#!/usr/bin/env bash
set -Eeuo pipefail

# trust the certificates of other nodes
if [[ `ls -1 /medco-configuration/srv*-certificate.crt 2>/dev/null | wc -l` != 0 ]]; then
    /bin/cp -f /medco-configuration/srv*-certificate.crt /usr/local/share/ca-certificates/
    update-ca-certificates
fi

exec $@
