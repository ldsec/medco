#!/usr/bin/env bash
set -Eeuo pipefail

# apply configuration from environment variables
pushd /etc/nginx/conf.d/

# append stream directive to default configuration of nginx
envsubst '${MEDCO_NODE_IDX} ${UNLYNX_PORT_0} ${UNLYNX_PORT_1}' < nginx.conf.template > nginx.conf.template_temp

cat ../nginx.conf nginx.conf.template_temp > nginx_new.conf
mv nginx_new.conf ../nginx.conf
rm nginx.conf.template_temp
cat ../nginx.conf

envsubst '$HTTP_SCHEME' < servers.conf.template > servers.conf
envsubst '$ALL_TIMEOUTS_SECONDS' < common/server-revproxy-base.conf.inc.template > common/server-revproxy-base.conf.inc

if [[ ${PROD_CONFIG} == "false" ]]; then
  cp common/server-revproxy-dev.conf.inc.template common/server-revproxy-dev.conf.inc
else
  touch common/server-revproxy-dev.conf.inc
fi

popd

exec nginx -g 'daemon off;'
