FROM debian:stretch as dev

# build-time variables
ENV I2B2_VERSION="tags/v1.7.10.0001"

# requirements
RUN apt-get -y update && \
    apt-get -y install git wget tar && \
    apt-get -y clean

# download i2b2 webclient
WORKDIR "/i2b2-webclient"
RUN git init . && \
    git remote add origin https://github.com/i2b2/i2b2-webclient.git && \
    git pull --depth=1 origin $I2B2_VERSION && \
    rm -rf ./.git && \
    mkdir /www && \
    cp -R . /www/i2b2-admin && \
    cp -R . /www/i2b2-client

# -------------------------------------------
FROM debian:stretch as prod

# run-time variables
ENV I2B2_DB_HOST="postgresql" \
    I2B2_DB_PORT="5432" \
    I2B2_DB_USER="i2b2" \
    I2B2_DB_PW="i2b2" \
    I2B2_DB_NAME="i2b2medco" \
    I2B2_DOMAIN_NAME="i2b2medco"

# build-time variables
ENV LIGHTTPD_WEB_ROOT="/var/www/html"

ENV CORS_ALLOW_ORIGIN="http://localhost:4200"

COPY --from=dev --chown=www-data:www-data "/www" "$LIGHTTPD_WEB_ROOT"
COPY --chown=www-data:www-data www-data/* "$LIGHTTPD_WEB_ROOT/"
COPY --chown=www-data:www-data lighttpd-conf/* /etc/lighttpd/conf-enabled/
COPY --chown=www-data:www-data i2b2-web-writeconfig.sh /
COPY docker-entrypoint.sh /usr/local/bin/

# install lighttpd & write config
RUN apt-get -y update && \
    apt-get -y install wget apt-transport-https lsb-release ca-certificates && \
    wget -O /etc/apt/trusted.gpg.d/php.gpg https://packages.sury.org/php/apt.gpg && \
    sh -c 'echo "deb https://packages.sury.org/php/ $(lsb_release -sc) main" > /etc/apt/sources.list.d/php.list' && \
    apt-get -y update && \
    apt-get -y install lighttpd postgresql-client libterm-readline-gnu-perl curl php5.6 php5.6-cgi php5.6-curl && \
    apt-get -y clean && \
    sed -i 's/display_errors = Off/display_errors = On/g' /etc/php/5.6/cgi/php.ini && \
    chgrp -R www-data /etc/lighttpd && \
    chmod -R g+rwx /etc/lighttpd && \
    chown -R www-data:www-data "$LIGHTTPD_WEB_ROOT" && \
    chmod -R +rx "$LIGHTTPD_WEB_ROOT" && \
    install -d -o www-data -g www-data -m 0750 "/var/run/lighttpd" && \
    chmod +x /i2b2-web-writeconfig.sh /usr/local/bin/docker-entrypoint.sh && \
    lighttpd-enable-mod fastcgi-php

# set global variables (to be available for php scripts)
RUN sed -i 's/"bin-environment" => (/"bin-environment" => (\n\t\t\t"CORS_ALLOW_ORIGIN" => "$CORS_ALLOW_ORIGIN",/g' /etc/lighttpd/conf-enabled/15-fastcgi-php.conf && \
    sed -i 's/"bin-environment" => (/"bin-environment" => (\n\t\t\t"I2B2_DB_HOST" => "$I2B2_DB_HOST",/g' /etc/lighttpd/conf-enabled/15-fastcgi-php.conf && \
    sed -i 's/"bin-environment" => (/"bin-environment" => (\n\t\t\t"I2B2_DB_PORT" => "$I2B2_DB_PORT",/g' /etc/lighttpd/conf-enabled/15-fastcgi-php.conf && \
    sed -i 's/"bin-environment" => (/"bin-environment" => (\n\t\t\t"I2B2_DB_USER" => "$I2B2_DB_USER",/g' /etc/lighttpd/conf-enabled/15-fastcgi-php.conf && \
    sed -i 's/"bin-environment" => (/"bin-environment" => (\n\t\t\t"I2B2_DB_PW" => "$I2B2_DB_PW",/g' /etc/lighttpd/conf-enabled/15-fastcgi-php.conf && \
    sed -i 's/"bin-environment" => (/"bin-environment" => (\n\t\t\t"I2B2_DB_NAME" => "$I2B2_DB_NAME",/g' /etc/lighttpd/conf-enabled/15-fastcgi-php.conf && \
    sed -i 's/"bin-environment" => (/"bin-environment" => (\n\t\t\t"I2B2_DOMAIN_NAME" => "$I2B2_DOMAIN_NAME",/g' /etc/lighttpd/conf-enabled/15-fastcgi-php.conf

# run
EXPOSE 80
ENTRYPOINT docker-entrypoint.sh
# todo: factor appart installed web services (i2b2 web client and others)