FROM postgres:9.6

# variables
ARG I2B2_DOMAIN_NAME_ARG="i2b2demo"
ENV I2B2_VERSION="v1.7.08.0001" \
    I2B2_DATA_DIR="/opt/i2b2-data" \
    I2B2_DOMAIN_NAME="$I2B2_DOMAIN_NAME_ARG" \
    I2B2_FR_FILES_DIR="/opt/FRC_files"

# pre-requisites
RUN apt-get -y update && apt-get -y install git ant default-jdk && apt-get -y clean

# download i2b2 demo data and apply patches
WORKDIR "$I2B2_DATA_DIR"
RUN git clone https://github.com/i2b2/i2b2-data.git . && \
    git checkout tags/$I2B2_VERSION && \
    chmod -R 777 .
COPY patches/* ./
RUN git apply *.diff

# authorize admin connection for phppgadmin TODO: password hardcoded
ENV POSTGRES_PASSWORD prigen2017

# .sh, .sql and .sql.gz will automatically be executed in this directory
WORKDIR /docker-entrypoint-initdb.d
COPY initdb-data/* ./

