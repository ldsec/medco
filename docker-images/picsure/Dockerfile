FROM maven:3.6 AS build

ENV GIT_REPO="https://github.com/hms-dbmi/pic-sure.git" \
    VERSION=7c3222318e97db72418d96e44b6eda2d491934eb

# get dependencies and sources
WORKDIR /src
COPY patches /patches
RUN apt-get update && \
    apt-get install -y git && \
    git clone $GIT_REPO . && \
    git checkout ${VERSION} && \
    git apply /patches/*.diff

# build
RUN mvn -pl '!pic-sure-api-wildfly' package

# -------------------------------------------
FROM jboss/wildfly:16.0.0.Final as release

COPY --from=build --chown=jboss:jboss /src/pic-sure-api-war/src/main/resources/wildfly-configuration/standalone.xml wildfly/standalone/configuration/
COPY --from=build --chown=jboss:jboss /src/pic-sure-api-war/target/pic-sure-api-war.war wildfly/standalone/deployments/pic-sure-api-2.war

# pre-requisites
USER root
COPY docker-entrypoint.sh /usr/local/bin/
COPY sql/* /sql/
RUN yum update -y && yum -y install postgresql wget && yum clean all && \
    mkdir /modules && chmod -R a+w /modules && \
    chmod a+x /usr/local/bin/docker-entrypoint.sh
USER jboss

# set up postgresql database and CORS configuration
COPY conf/postgresql.cli conf/jwks.cli /
RUN wget https://jdbc.postgresql.org/download/postgresql-42.2.5.jar -O /modules/postgresql-jdbc.jar && \
    /opt/jboss/wildfly/bin/jboss-cli.sh --file=/postgresql.cli && \
    /opt/jboss/wildfly/bin/jboss-cli.sh --file=/jwks.cli && \
    rm -rf /opt/jboss/wildfly/standalone/configuration/standalone_xml_history/current

# run-time variables used directly by PIC-SURE
ENV PICSURE2_PGDB_ADDRESS="postgresql" \
    PICSURE2_PGDB_PORT="5432" \
    PICSURE2_PGDB_DB="picsure" \
    PICSURE2_PGDB_USER="picsure" \
    PICSURE2_PGDB_PW="picsure" \
    VERIFY_USER_METHOD="local_jwks" \
    TOKEN_INTROSPECTION_URL="http://keycloak:8080/auth/realms/master/protocol/openid-connect/token/introspect" \
    TOKEN_INTROSPECTION_TOKEN="none" \
    PICSURE2_USER_ID_CLAIM="email" \
    PICSURE2_CLIENT_SECRET="none" \
    PICSURE_2_TOKEN="none" \
    PICSURE2_JWKS_URL="http://keycloak:8080/auth/realms/master/protocol/openid-connect/certs"

# additional run-time variables
ENV PICSURE2_CLIENT_ID="picsure" \
    MEDCO_NETWORK_NAME="testnetwork" \
    MEDCO_NODES_NAME="a,b,c" \
    MEDCO_NODES_CONNECTOR_URL="http://medco-connector-srv0/picsure2,http://medco-connector-srv1/picsure2,http://medco-connector-srv2/picsure2"

ENTRYPOINT ["docker-entrypoint.sh"]
