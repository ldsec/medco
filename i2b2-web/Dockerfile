FROM debian:jessie

# variables
ARG I2B2_DOMAIN_NAME_ARG="i2b2demo"
ENV I2B2_VERSION="v1.7.08a" \
    I2B2_WEB_DIR="/opt/i2b2-webclient" \
    SHRINE_SRC_DIR="/opt/shrine-src" \
    SHRINE_VERSION="fork/1.22.8-medco" \
    I2B2_DOMAIN_NAME="$I2B2_DOMAIN_NAME_ARG" \
    LIGHTTPD_WEB_ROOT="/var/www/html" \
    DB_PASSWORD="pFjy3EjDVwLfT2rB9xkK"


# requirements
RUN echo "deb http://httpredir.debian.org/debian jessie-backports main contrib non-free" >> /etc/apt/sources.list && \
    apt-get -y update && \
    apt-get -y install lighttpd php5-cgi curl libcurl3 php5-common php5-cli php5-curl git php5-pgsql wget tar && \
    apt-get -y -t jessie-backports install postgresql-client && \
    apt-get -y clean && \
    lighttpd-enable-mod fastcgi-php
    #ssl

# download i2b2 webclient
WORKDIR "$I2B2_WEB_DIR"
RUN git clone https://github.com/i2b2/i2b2-webclient.git . && \
    git checkout tags/$I2B2_VERSION && \
    chmod -R 777 .

# download shrine webclient
WORKDIR "$SHRINE_SRC_DIR"
RUN git clone https://c4science.ch/source/shrine-medco.git . && \
    git checkout $SHRINE_VERSION && \
    cp "$I2B2_WEB_DIR"/index.php "$SHRINE_SRC_DIR"/shrine-webclient/src/main/html/ && \
    chmod -R 777 .


# copy clients
RUN cp -R "$I2B2_WEB_DIR" "$LIGHTTPD_WEB_ROOT/i2b2-admin" && \
    cp -R "$I2B2_WEB_DIR" "$LIGHTTPD_WEB_ROOT/i2b2-client" && \
    cp -R "$SHRINE_SRC_DIR/shrine-webclient/src/main/html" "$LIGHTTPD_WEB_ROOT/shrine-client"

# copy and apply settings
COPY i2b2-web-writeconfig.sh .
RUN chmod -R +x ./i2b2-web-writeconfig.sh && sleep 1 && /bin/bash -c ./i2b2-web-writeconfig.sh

# phpPgAdmin
WORKDIR "$LIGHTTPD_WEB_ROOT"
RUN wget http://downloads.sourceforge.net/phppgadmin/phpPgAdmin-5.1.tar.gz?download && \
    tar xvzf phpPgAdmin-5.1.tar.gz?download && \
    rm phpPgAdmin-5.1.tar.gz\?download && \
    mv phpPgAdmin-5.1 phppgadmin && \
    sed -i "s/'host'] = ''/'host'] = 'i2b2-database'/" phppgadmin/conf/config.inc.php && \
    sed -i "s/ = 'PostgreSQL'/ = '$I2B2_DOMAIN_NAME PostgreSQL Server'/" phppgadmin/conf/config.inc.php && \
    sed -i "s/'extra_login_security'] = true/'extra_login_security'] = false/" phppgadmin/conf/config.inc.php && \
    sed -i "/\$cmd = \$exe/c\$cmd = \$exe;" phppgadmin/dbexport.php

# todo: phpmyadmin https://files.phpmyadmin.net/phpMyAdmin/4.7.4/phpMyAdmin-4.7.4-all-languages.zip
# RUN wget http://downloads.sourceforge.net/phppgadmin/phpPgAdmin-5.1.tar.gz?download && \
# /phpMyAdmin-4.7.4-all-languages/

# fix permissions + bug fix
RUN chown -R www-data:www-data "$LIGHTTPD_WEB_ROOT" && \
    chmod -R +r "$LIGHTTPD_WEB_ROOT" && \
    chmod -R +x "$LIGHTTPD_WEB_ROOT" && \
    install -d -o www-data -g www-data -m 0750 "/var/run/lighttpd"

# run
EXPOSE 80
CMD ["lighttpd", "-D", "-f", "/etc/lighttpd/lighttpd.conf"]
