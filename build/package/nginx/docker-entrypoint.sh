#!/usr/bin/env bash
set -Eeuo pipefail

# apply configuration from environment variables
pushd /etc/nginx/conf.d/

# append stream directive to default configuration of nginx
if [[ ${PROD_CONFIG} == "false" ]]; then
  envsubst '${MEDCO_NODE_IDX} ${UNLYNX_PORT_0} ${UNLYNX_PORT_1}' < nginx.conf.template > ../nginx.conf
else
  envsubst '${UNLYNX_PORT_0} ${UNLYNX_PORT_1}' < nginx.conf.prod.template > ../nginx.conf
fi

envsubst '$HTTP_SCHEME' < servers.conf.template > servers.conf
envsubst '$ALL_TIMEOUTS_SECONDS' < common/server-revproxy-base.conf.inc.template > common/server-revproxy-base.conf.inc

if [[ ${PROD_CONFIG} == "false" ]]; then
  cp common/server-revproxy-dev.conf.inc.template common/server-revproxy-dev.conf.inc
else
  touch common/server-revproxy-dev.conf.inc
fi

popd

exec nginx -g 'daemon off;'
