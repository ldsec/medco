#!/usr/bin/env bash
set -Eeuo pipefail

if [ ! -f /etc/lighttpd/conf-enabled/15-fastcgi-php-copy.conf ]; then
	cp /etc/lighttpd/conf-enabled/15-fastcgi-php.conf /etc/lighttpd/conf-enabled/15-fastcgi-php-copy.conf

	sed -i "s/\"bin-environment\" => (/\"bin-environment\" => (\n\t\t\t\"CORS_ALLOW_ORIGIN\" => \"$CORS_ALLOW_ORIGIN\",/g" /etc/lighttpd/conf-enabled/15-fastcgi-php-copy.conf
	sed -i "s/\"bin-environment\" => (/\"bin-environment\" => (\n\t\t\t\"I2B2_DB_HOST\" => \"$I2B2_DB_HOST\",/g" /etc/lighttpd/conf-enabled/15-fastcgi-php-copy.conf
	sed -i "s/\"bin-environment\" => (/\"bin-environment\" => (\n\t\t\t\"I2B2_DB_PORT\" => \"$I2B2_DB_PORT\",/g" /etc/lighttpd/conf-enabled/15-fastcgi-php-copy.conf
	sed -i "s/\"bin-environment\" => (/\"bin-environment\" => (\n\t\t\t\"I2B2_DB_USER\" => \"$I2B2_DB_USER\",/g" /etc/lighttpd/conf-enabled/15-fastcgi-php-copy.conf
	sed -i "s/\"bin-environment\" => (/\"bin-environment\" => (\n\t\t\t\"I2B2_DB_PW\" => \"$I2B2_DB_PW\",/g" /etc/lighttpd/conf-enabled/15-fastcgi-php-copy.conf
	sed -i "s/\"bin-environment\" => (/\"bin-environment\" => (\n\t\t\t\"I2B2_DB_NAME\" => \"$I2B2_DB_NAME\",/g" /etc/lighttpd/conf-enabled/15-fastcgi-php-copy.conf
	sed -i "s/\"bin-environment\" => (/\"bin-environment\" => (\n\t\t\t\"I2B2_DOMAIN_NAME\" => \"$I2B2_DOMAIN_NAME\",/g" /etc/lighttpd/conf-enabled/15-fastcgi-php-copy.conf

	cp /etc/lighttpd/conf-enabled/15-fastcgi-php-copy.conf /etc/lighttpd/conf-enabled/15-fastcgi-php.conf
fi

su "www-data" -p -s /bin/bash -c /i2b2-web-writeconfig.sh
exec lighttpd -D -f /etc/lighttpd/lighttpd.conf