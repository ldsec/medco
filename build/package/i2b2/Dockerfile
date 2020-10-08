FROM jboss/wildfly:10.0.0.Final as build

# build-time variables
ENV JBOSS_WAR_DEPLOYMENTS="$JBOSS_HOME/standalone/deployments" \
    AXIS2_VERSION="1.7.1" \
    I2B2_SRC_DIR="/i2b2-src" \
    I2B2_PATCHES_DIR="/i2b2-patches" \
    I2B2_VERSION="tags/v1.7.10.0001"

# pre-requisites: packages & wildfly config
USER root
RUN yum update -y && \
    yum -y install wget git ant && \
    yum clean all
COPY install-i2b2.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/install-i2b2.sh
COPY patches "$I2B2_PATCHES_DIR"
RUN mkdir -p "$I2B2_SRC_DIR" && \
    chown -R jboss:jboss "$I2B2_SRC_DIR"
USER jboss

# install i2b2
RUN install-i2b2.sh

# install i2b2 password hasher
RUN mkdir -p "$JBOSS_HOME/I2b2PasswordHash" && \
    javac -d "$JBOSS_HOME/I2b2PasswordHash/" "$I2B2_PATCHES_DIR/I2b2PasswordHash.java"

# -------------------------------------------
FROM jboss/wildfly:10.0.0.Final as release

# run-time variables
ENV I2B2_DB_HOST="postgresql" \
    I2B2_DB_PORT="5432" \
    I2B2_DB_USER="i2b2" \
    I2B2_DB_PW="i2b2" \
    I2B2_DB_NAME="i2b2" \
    WILDFLY_ADMIN_PASSWORD="admin" \
    I2B2_DOMAIN_NAME="i2b2demo" \
    I2B2_SERVICE_PASSWORD="pFjy3EjDVwLfT2rB9xkK" \
    DEFAULT_USER_PASSWORD="demouser" \
    AXIS2_LOGLEVEL="INFO"

# build-time variables
ENV I2B2_SQL_DIR="/i2b2-sql" \
    I2B2_COMPRESSED_DATA_DIR="/i2b2-compressed-data" \
    I2B2_DATA_DIR="/i2b2-data" \
    PRE_INIT_SCRIPT_DIR="/pre-init-scripts" \
    I2B2_DATA_VERSION="v1.7.10.0002" \
    I2B2_FR_FILES_DIR="$JBOSS_HOME/standalone/data/i2b2_FR_files" \
    PG_JDBC_JAR="postgresql-42.1.4.jar" \
    JBOSS_WAR_DEPLOYMENTS="$JBOSS_HOME/standalone/deployments"

COPY --from=build --chown=jboss:jboss "$JBOSS_HOME/standalone/configuration" "$JBOSS_HOME/standalone/configuration"
COPY --from=build --chown=jboss:jboss "$JBOSS_HOME/standalone/deployments" "$JBOSS_HOME/standalone/deployments"
COPY --from=build --chown=jboss:jboss "$JBOSS_HOME/I2b2PasswordHash" "$JBOSS_HOME/I2b2PasswordHash"
COPY docker-entrypoint.sh /usr/local/bin/
COPY sql "$I2B2_SQL_DIR"
COPY pre-init-scripts "$PRE_INIT_SCRIPT_DIR"

USER root
RUN sed -i 's/Xmx512m/Xmx2048m/g' $JBOSS_HOME/bin/standalone.conf && \
    sed -i 's/MaxMetaspaceSize=256m/MaxMetaspaceSize=1024m/g' $JBOSS_HOME/bin/standalone.conf && \
    chmod +x /usr/local/bin/docker-entrypoint.sh && \
    mkdir -p "$I2B2_COMPRESSED_DATA_DIR" "$I2B2_DATA_DIR" "$I2B2_FR_FILES_DIR" && \
    chown -R jboss:jboss "$I2B2_DATA_DIR" "$JBOSS_HOME" && \
    yum update -y && \
    yum -y install wget git ant postgresql && \
    yum clean all
RUN git clone --depth 1 --branch "$I2B2_DATA_VERSION" https://github.com/i2b2/i2b2-data.git "$I2B2_DATA_DIR" && \
    GZIP=-9 tar -cvzf "$I2B2_COMPRESSED_DATA_DIR"/i2b2-data.tar.gz "$I2B2_DATA_DIR" && \
    rm -rf "$I2B2_DATA_DIR"/{*,.*} || true
USER jboss

# run
EXPOSE 8080 9990
ENTRYPOINT docker-entrypoint.sh
