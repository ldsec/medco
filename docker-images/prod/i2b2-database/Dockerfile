FROM postgres:9.6

# variables
ENV I2B2_VERSION="v1.7.08.0001" \
    I2B2_DATA_DIR="/opt/i2b2-data" \
    I2B2_DOMAIN_NAME="medcodeployment" \
    I2B2_DEMO_DB_NAME="i2b2demo" \
    I2B2_MEDCO_DB_NAME="i2b2medco" \
    I2B2_FR_FILES_DIR="/opt/FRC_files" \
    SHRINE_DEMO_ONTOLOGY_SQL="/opt/ShrineDemo.sql" \
    ADMIN_PASSWORD="prigen2017" \
    DB_PASSWORD="pFjy3EjDVwLfT2rB9xkK"

# pre-requisites
RUN apt-get -y update && apt-get -y install git ant default-jdk wget && apt-get -y clean

# download i2b2 demo data, apply patches, and download shrine demo ontology
WORKDIR "$I2B2_DATA_DIR"
RUN git clone https://github.com/i2b2/i2b2-data.git . && \
    git checkout tags/$I2B2_VERSION && \
    chmod -R 777 . && \
    wget https://documents.epfl.ch/users/m/mi/misbach/www/ShrineDemo.sql -O "$SHRINE_DEMO_ONTOLOGY_SQL" && \
    sed -i '1s/^/SET search_path TO shrine_ont;\n/' "$SHRINE_DEMO_ONTOLOGY_SQL"
# hack postgresql docker image to allow transmart loading (allow TCP/IP connections during init phase)
RUN sed -i "s/listen_addresses=''/listen_addresses='*'/g" /usr/local/bin/docker-entrypoint.sh

COPY patches/* ./
RUN git apply *.diff

# authorize admin connection for phppgadmin
ENV POSTGRES_PASSWORD="$ADMIN_PASSWORD"

# .sh, .sql and .sql.gz will automatically be executed in this directory
WORKDIR /docker-entrypoint-initdb.d
COPY initdb-data/* ./
