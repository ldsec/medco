FROM tomcat:9.0-jre8

ARG SHRINE_DEBUG_LEVEL_ARG
ARG NODE_IDX_ARG

# pre-existing variables: CATALINA_HOME
#ENV CONF_DIR="/opt/medco-configuration" \
#	BUILD_DIR="/opt/build-dir" \
#	SHRINE_DEBUG_LEVEL="$SHRINE_DEBUG_LEVEL_ARG" \
#    NODE_IDX="$NODE_IDX_ARG"

# pre-existing variables: CATALINA_HOME
ENV SHRINE_VERSION="fork/1.22.8-medco" \
SHRINE_SRC_DIR="/opt/shrine-src" \
SHRINE_ADAPTER_MAPPINGS_URL="https://documents.epfl.ch/users/m/mi/misbach/www/AdapterMappings.xml" \
SHRINE_MYSQL_JAR_URL="http://central.maven.org/maven2/mysql/mysql-connector-java/5.1.40/mysql-connector-java-5.1.40.jar" \
I2B2_DOMAIN_NAME="medcodeployment" \
CONF_DIR="/opt/medco-configuration" \
ADMIN_PASSWORD="prigen2017" \
DB_PASSWORD="pFjy3EjDVwLfT2rB9xkK"

# system and tomcat prerequisites
RUN apt-get -y update && \
    apt-get -y install git maven zip wget unzip openjdk-8-jdk-headless && \
    apt-get -y clean && \
    echo "<?xml version='1.0' encoding='utf-8'?><tomcat-users><role rolename=\"manager-gui\" /><role rolename=\"admin-gui\" />" \
        "<user username=\"admin\" password=\"$ADMIN_PASSWORD\" roles=\"manager-gui,admin-gui\" /></tomcat-users>" > \
        "$CATALINA_HOME/conf/tomcat-users.xml" && \
    echo 'export CATALINA_OPTS=" -Dakka.daemonic=on "' > "$CATALINA_HOME/bin/setenv.sh" && \
    echo '{ "allow_root": true }' > /root/.bowerrc

# download sources
WORKDIR "$SHRINE_SRC_DIR"
RUN git clone https://c4science.ch/source/shrine-medco.git . && \
    git checkout $SHRINE_VERSION

# compilation and installation
RUN mvn -e -pl commons/test-commons install -DskipTests && \
    mvn -e -N install -DskipTests && \
    mvn -e -pl commons/util install -DskipTests && \
    mvn -e -pl commons/config install -DskipTests && \
    mvn -e -pl commons/data-commons install -DskipTests && \
    mvn -e -pl commons/protocol-query install -DskipTests && \
    mvn -e -pl commons/protocol install -DskipTests && \
    mvn -e -pl tools install -DskipTests && \
    mvn -e -pl tools/utility-commons install -DskipTests && \
    mvn -e -pl apps/meta-app install -DskipTests && \
    mvn -e -pl apps/meta-war install -DskipTests && \
    mvn -e -pl commons/crypto install -DskipTests && \
    mvn -e -pl commons/client install -DskipTests && \
    mvn -e -pl tools/batch-querier install -DskipTests && \
    mvn -e -pl commons/auth install -DskipTests && \
    mvn -e -pl adapter/adapter-api install -DskipTests && \
    mvn -e -pl apps/dashboard-app install -DskipTests && \
    mvn -e -pl apps/dashboard-war install -DskipTests && \
    mvn -e -pl commons/email install -DskipTests && \
    mvn -e -pl apps/steward-app install -DskipTests && \
    mvn -e -pl apps/steward-war install -DskipTests && \
    mvn -e -pl apps/proxy install -DskipTests && \
    mvn -e -pl hms-support install -DskipTests && \
    mvn -e -pl hub/broadcaster-aggregator install -DskipTests && \
    mvn -e -pl tools/monitor install -DskipTests && \
    mvn -e -pl hub/broadcaster-service install -DskipTests && \
    mvn -e -pl adapter/adapter-service install -DskipTests && \
    mvn -e -pl commons/ont-support install -DskipTests && \
    mvn -e -pl tools/mapping-automation install -DskipTests && \
    mvn -e -pl install install -DskipTests && \
    mvn -e -pl shrine-webclient install -DskipTests && \
    mvn -e install -Dmaven.test.skip=true

RUN cp  "$SHRINE_SRC_DIR/apps/steward-war/target/steward.war" \
        "$SHRINE_SRC_DIR/apps/dashboard-war/target/shrine-dashboard.war" \
        "$CATALINA_HOME/webapps/" && \
    cp  "$SHRINE_SRC_DIR/apps/war/target/shrine-cell.war" "$CATALINA_HOME/webapps/shrine.war" && \
    cp  "$SHRINE_SRC_DIR/apps/meta-war/target/shrine-metadata.war" "$CATALINA_HOME/webapps/shrine-meta.war"

# configuration
COPY conf/shrine.conf "$CATALINA_HOME/lib/"
COPY conf/server.xml conf/context.xml "$CATALINA_HOME/conf/"
RUN wget "$SHRINE_MYSQL_JAR_URL" -P "$CATALINA_HOME/lib/" && \
    wget "$SHRINE_ADAPTER_MAPPINGS_URL" -O "$CATALINA_HOME/lib/AdapterMappings.xml" && \
    sed -i "s#SHRINE_DOWNSTREAM_NODES_FILE_PATH#$CONF_DIR/shrine_downstream_nodes.conf#g" "$CATALINA_HOME/lib/shrine.conf" && \
    sed -i "s#SHRINE_CA_CERT_ALIASES_FILE_PATH#$CONF_DIR/shrine_ca_cert_aliases.conf#g" "$CATALINA_HOME/lib/shrine.conf" && \
    sed -i "s/SHRINE_KEYSTORE_PASSWORD/$ADMIN_PASSWORD/g" "$CATALINA_HOME/conf/server.xml" && \
    sed -i "s/SHRINE_DB_PASSWORD/$DB_PASSWORD/g" "$CATALINA_HOME/conf/context.xml"

# configuration bis (dependent on the arguments)
ARG SHRINE_DEBUG_LEVEL_ARG="INFO"
ARG NODE_IDX_ARG="0"
ENV SHRINE_DEBUG_LEVEL="$SHRINE_DEBUG_LEVEL_ARG" \
    NODE_IDX="$NODE_IDX_ARG"
RUN sed -i "s/SHRINE_KEYSTORE_PRIVATE_KEY_ALIAS/srv$NODE_IDX-private/g" "$CATALINA_HOME/conf/server.xml" && \
    sed -i "s#SHRINE_KEYSTORE_FILE_PATH#$CONF_DIR/srv$NODE_IDX.keystore#g" "$CATALINA_HOME/conf/server.xml" && \
    sed -i "s#FINE#$SHRINE_DEBUG_LEVEL#g" "$CATALINA_HOME/conf/logging.properties" && \
    sed -i "s#INFO#$SHRINE_DEBUG_LEVEL#g" "$CATALINA_HOME/conf/logging.properties"

EXPOSE 6060 6443
VOLUME $CONF_DIR
WORKDIR $CONF_DIR

#WORKDIR $CONF_DIR
# copy the build files from the build volume to the correct folder
#COPY copy-shrine-server-binary-prod.sh /root
#ENTRYPOINT ["sh", "-c", "/root/copy-shrine-server-binary-prod.sh && $CATALINA_HOME/bin/catalina.sh run"]