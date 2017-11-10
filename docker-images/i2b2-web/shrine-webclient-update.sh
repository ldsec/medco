#!/bin/bash
set -e

LIGHTTPD_WEB_ROOT="/var/www/html"
SHRINE_SRC_DIR="/opt/shrine-src"

cd "$SHRINE_SRC_DIR"
git pull
rm -R "$LIGHTTPD_WEB_ROOT/shrine-client"
cp -R "$SHRINE_SRC_DIR/shrine-webclient/src/main/html" "$LIGHTTPD_WEB_ROOT/shrine-client"
bash /opt/i2b2-web-writeconfig.sh
chown -R www-data:www-data "$LIGHTTPD_WEB_ROOT/shrine-client"
chmod -R +r "$LIGHTTPD_WEB_ROOT/shrine-client"
chmod -R +x "$LIGHTTPD_WEB_ROOT/shrine-client"
