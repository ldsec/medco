#!/usr/bin/env bash
set -Eeuo pipefail

# apply configuration from environment variables
pushd /etc/nginx/conf.d/
envsubst '$HTTP_SCHEME' < servers.conf.template > servers.conf
envsubst '$ALL_TIMEOUTS_SECONDS' < common/server-revproxy.conf.inc.template > common/server-revproxy.conf.inc
popd

exec nginx -g 'daemon off;'
