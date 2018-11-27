FROM php:7.0-fpm-stretch as prod

# build time env
ENV WWW_DATA="/www-data"

# run-time variables
ENV I2B2_DB_HOST="postgresql" \
    I2B2_DB_PORT="5432" \
    I2B2_DB_USER="i2b2" \
    I2B2_DB_PW="i2b2" \
    I2B2_DB_NAME="i2b2medco" \
    I2B2_DOMAIN_NAME="i2b2medco" \
    CORS_ALLOW_ORIGIN="http://localhost:4200"

# install additional packages
RUN apt-get -y update && \
    apt-get -y install wget libpq-dev && \
    docker-php-ext-install pdo pdo_pgsql && \
    apt-get -y clean

# run
VOLUME $WWW_DATA
