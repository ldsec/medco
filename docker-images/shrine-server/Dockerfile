FROM tomcat:8.0-jre8

# pre-existing variables: CATALINA_HOME
ARG I2B2_DOMAIN_NAME_ARG="i2b2demo"
ENV SHRINE_VERSION="fork/1.22.8-medco" \
    SHRINE_SRC_DIR="/opt/shrine-src" \
    SHRINE_ADAPTER_MAPPINGS_URL="https://open.med.harvard.edu/svn/shrine-ontology/SHRINE_Demo_Downloads/trunk/AdapterMappings_i2b2_DemoData.xml" \
    SHRINE_MYSQL_JAR_URL="http://central.maven.org/maven2/mysql/mysql-connector-java/5.1.40/mysql-connector-java-5.1.40.jar" \
    I2B2_DOMAIN_NAME="$I2B2_DOMAIN_NAME_ARG" \
    CONF_DIR="/opt/medco-configuration" \
    ADMIN_PASSWORD="prigen2017" \
    DB_PASSWORD="pFjy3EjDVwLfT2rB9xkK"

# system and tomcat prerequisites
RUN apt-get -y update && \
    apt-get -y install git maven zip wget unzip openjdk-8-jdk-headless && \
    apt-get -y clean && \
    echo "<?xml version='1.0' encoding='utf-8'?><tomcat-users><role rolename=\"manager-gui\" /><user username=\"admin\"" \
        "password=\"$ADMIN_PASSWORD\" roles=\"manager-gui\" /></tomcat-users>" > "$CATALINA_HOME/conf/tomcat-users.xml" && \
    echo 'export CATALINA_OPTS=" -Dakka.daemonic=on "' > "$CATALINA_HOME/bin/setenv.sh" && \
    echo '{ "allow_root": true }' > /root/.bowerrc

# download sources
WORKDIR "$SHRINE_SRC_DIR"
RUN git clone https://c4science.ch/source/shrine-medco.git . && \
    git checkout $SHRINE_VERSION

# compilation and installation
RUN mvn -e install -DskipTests -Dmaven.test.skip=true
RUN cp  "$SHRINE_SRC_DIR/apps/steward-war/target/steward.war" \
        "$SHRINE_SRC_DIR/apps/dashboard-war/target/shrine-dashboard.war" \
        "$CATALINA_HOME/webapps/" && \
    cp  "$SHRINE_SRC_DIR/apps/war/target/shrine-cell.war" "$CATALINA_HOME/webapps/shrine.war" && \
    cp  "$SHRINE_SRC_DIR/apps/meta-war/target/shrine-metadata.war" "$CATALINA_HOME/webapps/shrine-meta.war"

# webclient [disabled]
#RUN cp "$SHRINE_SRC_DIR/apps/proxy/target/shrine-proxy.war" "$CATALINA_HOME/webapps/" && \
#    cp -r "$SHRINE_SRC_DIR/shrine-webclient/src/main/html" "$CATALINA_HOME/webapps/shrine-client"
#COPY conf/i2b2_config_data.js "$CATALINA_HOME/webapps/shrine-client/"
#COPY conf/cell_config_data.js "$CATALINA_HOME/webapps/shrine-client/js-i2b2/cells/SHRINE/"
#RUN sed -i "s/SHRINE_WEBCLIENT_DOMAIN/$I2B2_DOMAIN_NAME/g" "$CATALINA_HOME/webapps/shrine-client/i2b2_config_data.js" && \
#    sed -i "s/SHRINE_WEBCLIENT_NAME/Domain $I2B2_DOMAIN_NAME/g" "$CATALINA_HOME/webapps/shrine-client/i2b2_config_data.js"

# configuration
COPY conf/shrine.conf "$CATALINA_HOME/lib/"
COPY conf/server.xml conf/context.xml "$CATALINA_HOME/conf/"
RUN wget "$SHRINE_MYSQL_JAR_URL" -P "$CATALINA_HOME/lib/" && \
    wget "$SHRINE_ADAPTER_MAPPINGS_URL" -O "$CATALINA_HOME/lib/AdapterMappings.xml" && \
    sed -i "s#SHRINE_DOWNSTREAM_NODES_FILE_PATH#$CONF_DIR/shrine_downstream_nodes.conf#g" "$CATALINA_HOME/lib/shrine.conf" && \
    sed -i "s/SHRINE_KEYSTORE_PASSWORD/$ADMIN_PASSWORD/g" "$CATALINA_HOME/conf/server.xml" && \
    sed -i "s/SHRINE_DB_PASSWORD/$DB_PASSWORD/g" "$CATALINA_HOME/conf/context.xml"

# configuration bis (dependent on the arguments) todo: tomcat/shrine logs level controlable
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

