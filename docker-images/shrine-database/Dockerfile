FROM mysql:5.5

# env variablebles
ENV ADMIN_PASSWORD="prigen2017" \
    DB_PASSWORD="pFjy3EjDVwLfT2rB9xkK" \
    SHRINE_VERSION="fork/1.22.8-medco" \
    SHRINE_SRC_DIR="/opt/shrine-src"
ENV MYSQL_ROOT_PASSWORD="$ADMIN_PASSWORD" \
    MYSQL_DATABASE="" \
    MYSQL_USER="" \
    MYSQL_PASSWORD=""

# install some requirements
RUN apt-get -y update && \
    apt-get -y install git wget && \
    apt-get -y clean

# get shrine sources (for sql files)
WORKDIR "$SHRINE_SRC_DIR"
RUN git clone https://c4science.ch/source/shrine-medco.git . && \
    git checkout $SHRINE_VERSION

# .sh, .sql and .sql.gz will automatically be executed in this directory
WORKDIR /docker-entrypoint-initdb.d
COPY initdb-data/* ./
