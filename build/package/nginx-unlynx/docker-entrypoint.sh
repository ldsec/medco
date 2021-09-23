#!/usr/bin/env bash
set -Eeuo pipefail

# apply configuration from environment variables
pushd /etc/nginx

exec nginx -g 'daemon off;'
