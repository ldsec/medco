FROM debian:jessie-backports

ARG NODE_IDX_ARG

ENV SHRINE_SRC_DIR="/opt/shrine-src" \
    LIGHTTPD_WEB_ROOT="/var/www/html" \
    I2B2_DOMAIN_NAME="medcodeployment" \
    I2B2_DEMO_DB_NAME="i2b2demo" \
    I2B2_MEDCO_DB_NAME="i2b2medco" \
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

# copy settings
USER root
COPY i2b2-web-writeconfig.sh shrine-webclient-update.sh /opt/
# copy the build files from the build volume to the correct folder
COPY copy-i2b2-web-binary-prod.sh /root

ENTRYPOINT ["sh", "-c", "/root/copy-i2b2-web-binary-prod.sh"]