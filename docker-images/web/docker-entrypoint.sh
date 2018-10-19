#!/usr/bin/env bash
set -Eeuo pipefail

su "www-data" -p -s /bin/bash -c /i2b2-web-writeconfig.sh
exec lighttpd -D -f /etc/lighttpd/lighttpd.conf