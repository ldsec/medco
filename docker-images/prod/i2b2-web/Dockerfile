FROM debian:jessie-backports

ARG NODE_IDX_ARG

ENV SHRINE_SRC_DIR="/opt/shrine-src" \
    I2B2_DOMAIN_NAME="medcodeployment" \
    I2B2_MEDCO_DB_NAME="i2b2medco" \
    LIGHTTPD_WEB_ROOT="/var/www/html" \
    DB_PASSWORD="pFjy3EjDVwLfT2rB9xkK" \
    CONF_DIR="/opt/medco-configuration" \
    BUILD_DIR="/opt/build-dir" \
    NODE_IDX="$NODE_IDX_ARG"

# requirements
RUN echo "deb http://httpredir.debian.org/debian jessie-backports main contrib non-free" >> /etc/apt/sources.list && \
    apt-get -y update && \
    apt-get -y install lighttpd php5-cgi curl libcurl3 php5-common php5-cli php5-curl git php5-pgsql php5-mysql wget tar && \
    apt-get -y -t jessie-backports install postgresql-client && \
    apt-get -y clean && \
    lighttpd-enable-mod fastcgi-php ssl

# requirements
RUN echo "deb http://httpredir.debian.org/debian jessie-backports main contrib non-free" >> /etc/apt/sources.list && \
    apt-get -y update && \
    apt-get -y install lighttpd php5-cgi curl libcurl3 php5-common php5-cli php5-curl git php5-pgsql php5-mysql wget tar && \
    apt-get -y -t jessie-backports install postgresql-client && \
    apt-get -y clean && \
    lighttpd-enable-mod fastcgi-php ssl

RUN mkdir -p "$I2B2_WEB_DIR" "$SHRINE_SRC_DIR" && \
    chgrp -R www-data /opt /etc/lighttpd && \
    chmod -R g+rwx /opt /etc/lighttpd && \
    chown -R www-data:www-data "$LIGHTTPD_WEB_ROOT" && \
    chmod -R +rx "$LIGHTTPD_WEB_ROOT" && \
    install -d -o www-data -g www-data -m 0750 "/var/run/lighttpd"
USER www-data

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

# phpMyAdmin
RUN wget https://files.phpmyadmin.net/phpMyAdmin/4.8.2/phpMyAdmin-4.8.2-english.tar.gz && \
    tar xvzf phpMyAdmin-4.8.2-english.tar.gz && \
    rm phpMyAdmin-4.8.2-english.tar.gz && \
    mv phpMyAdmin-4.8.2-english phpmyadmin && \
    cp phpmyadmin/config.sample.inc.php phpmyadmin/config.inc.php && \
    sed -i "s#localhost#shrine-database#g" phpmyadmin/config.inc.php && \
    sed -i "s#\$cfg['blowfish_secret'] = '';#\$cfg['blowfish_secret'] = 'ac28vb01sfv0vcs1v96fv06fv0df62fv';#g" phpmyadmin/config.inc.php

# download i2b2 webclient
WORKDIR "$I2B2_WEB_DIR"
RUN git clone https://github.com/i2b2/i2b2-webclient.git . && \
    git checkout tags/$I2B2_VERSION

# download shrine webclient
WORKDIR "$SHRINE_SRC_DIR"
RUN git clone https://c4science.ch/source/shrine-medco.git . && \
    git checkout $SHRINE_VERSION && \
    cp "$I2B2_WEB_DIR"/index.php "$SHRINE_SRC_DIR"/shrine-webclient/src/main/html/

# copy clients
RUN cp -R "$I2B2_WEB_DIR" "$LIGHTTPD_WEB_ROOT/i2b2-admin" && \
    cp -R "$I2B2_WEB_DIR" "$LIGHTTPD_WEB_ROOT/i2b2-client" && \
    cp -R "$SHRINE_SRC_DIR/shrine-webclient/src/main/html" "$LIGHTTPD_WEB_ROOT/shrine-client"

# copy settings
USER root
COPY i2b2-web-writeconfig.sh shrine-webclient-update.sh /opt/
# copy the build files from the build volume to the correct folder
#COPY copy-i2b2-web-binary-prod.sh /root

ENTRYPOINT ["sh", "-c", "/root/copy-i2b2-web-binary-prod.sh"]