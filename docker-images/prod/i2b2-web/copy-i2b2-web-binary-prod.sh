#!/bin/bash
set -e

# copy pm cell deployment data to shared volume
#cp -a -R $BUILD_DIR/i2b2-web/html/* $LIGHTTPD_WEB_ROOT

#chgrp -R www-data /opt /etc/lighttpd
#chmod -R g+rwx /opt /etc/lighttpd
#chown -R www-data:www-data "$LIGHTTPD_WEB_ROOT"
#chmod -R +rx "$LIGHTTPD_WEB_ROOT"
#install -d -o www-data -g www-data -m 0750 "/var/run/lighttpd"
#su "www-data" -p -s /bin/bash -c "/opt/i2b2-web-writeconfig.sh"

lighttpd -D -f /etc/lighttpd/lighttpd.conf