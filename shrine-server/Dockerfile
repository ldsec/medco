FROM tomcat:8.0-jre8

# variables CATALINA_HOME exists and set
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

# download and compile shrine wars
ENV NOCACHE="2"
WORKDIR "$SHRINE_SRC_DIR"
RUN git clone https://c4science.ch/source/shrine-medco.git . && \
    git checkout $SHRINE_VERSION

# install
# todo: reorganize properly this part (cache and build shrine parts)
# todo: cleanup file (+ remove from conf the .js files)

# cache
RUN mvn -e -pl commons/test-commons install -DskipTests
RUN mvn -e -N install -DskipTests
RUN mvn -e -pl commons/util install -DskipTests
RUN mvn -e -pl commons/config install -DskipTests
RUN mvn -e -pl commons/data-commons install -DskipTests
RUN mvn -e -pl commons/protocol-query install -DskipTests
#RUN mvn dependency:resolve

# breakdown modules installation
RUN mvn -e -pl commons/protocol install -DskipTests
RUN mvn -e -pl tools install -DskipTests
RUN mvn -e -pl tools/utility-commons install -DskipTests
RUN mvn -e -pl apps/meta-app install -DskipTests
RUN mvn -e -pl apps/meta-war install -DskipTests
RUN mvn -e -pl commons/crypto install -DskipTests
RUN mvn -e -pl commons/client install -DskipTests
RUN mvn -e -pl tools/batch-querier install -DskipTests
RUN mvn -e -pl commons/auth install -DskipTests
RUN mvn -e -pl adapter/adapter-api install -DskipTests
RUN mvn -e -pl apps/dashboard-app install -DskipTests
RUN mvn -e -pl apps/dashboard-war install -DskipTests
RUN mvn -e -pl commons/email install -DskipTests
RUN mvn -e -pl apps/steward-app install -DskipTests
RUN mvn -e -pl apps/steward-war install -DskipTests
RUN mvn -e -pl apps/proxy install -DskipTests
RUN mvn -e -pl hms-support install -DskipTests
RUN mvn -e -pl hub/broadcaster-aggregator install -DskipTests
RUN mvn -e -pl tools/monitor install -DskipTests
#RUN cd hms-support && mvn -e -pl hms-core/shrine-hms-core install && cd ..
#RUN mvn -e -pl qep/service install
RUN mvn -e -pl hub/broadcaster-service install -DskipTests
RUN mvn -e -pl adapter/adapter-service install -DskipTests
#RUN mvn -e -pl tools/adapter-queries-to-qep install
RUN mvn -e -pl commons/ont-support install -DskipTests
#RUN mvn -e -pl tools/scanner install
RUN mvn -e -pl tools/mapping-automation install -DskipTests
#RUN mvn -e -pl apps/shrine-app install
#RUN mvn -e -pl apps/war install
RUN mvn -e -pl install install -DskipTests
#RUN mvn -e -pl integration install
RUN mvn -e -pl shrine-webclient install -DskipTests

# compilation and installation
RUN mvn -e install -Dmaven.test.skip=true
RUN cp  "$SHRINE_SRC_DIR/apps/steward-war/target/steward.war" \
        "$SHRINE_SRC_DIR/apps/dashboard-war/target/shrine-dashboard.war" \
        #"$SHRINE_SRC_DIR/apps/proxy/target/shrine-proxy.war" \
        "$CATALINA_HOME/webapps/" && \
    cp  "$SHRINE_SRC_DIR/apps/war/target/shrine-cell.war" "$CATALINA_HOME/webapps/shrine.war" && \
    cp  "$SHRINE_SRC_DIR/apps/meta-war/target/shrine-metadata.war" "$CATALINA_HOME/webapps/shrine-meta.war"
    #&& \
    #cp -r "$SHRINE_SRC_DIR/shrine-webclient/src/main/html" "$CATALINA_HOME/webapps/shrine-client"

# shrine webclient
#WORKDIR "$SHRINE_SRC_DIR-webclient"
#RUN git clone https://c4science.ch/source/shrine-medco.git . && \
#    git checkout $SHRINE_VERSION && \
#    mvn -e -pl shrine-webclient install -DskipTests && \
#    rm -rf "$CATALINA_HOME/webapps/shrine-client" && \
#    cp -r "$SHRINE_SRC_DIR-webclient/shrine-webclient/src/main/html" "$CATALINA_HOME/webapps/shrine-client"

# configuration
RUN wget "$SHRINE_ADAPTER_MAPPINGS_URL" -O "$CATALINA_HOME/lib/AdapterMappings.xml" && \
    wget "$SHRINE_MYSQL_JAR_URL" -P "$CATALINA_HOME/lib/"
COPY conf/shrine.conf "$CATALINA_HOME/lib/"
COPY conf/server.xml conf/context.xml "$CATALINA_HOME/conf/"
#COPY conf/i2b2_config_data.js "$CATALINA_HOME/webapps/shrine-client/"
#COPY conf/cell_config_data.js "$CATALINA_HOME/webapps/shrine-client/js-i2b2/cells/SHRINE/"
RUN sed -i "s#SHRINE_DOWNSTREAM_NODES_FILE_PATH#$CONF_DIR/shrine_downstream_nodes.conf#g" "$CATALINA_HOME/lib/shrine.conf" && \
    #sed -i "s#SHRINE_ALIAS_MAP_FILE_PATH#$CONF_DIR/shrine_alias_map.conf#g" "$CATALINA_HOME/lib/shrine.conf" && \

    sed -i "s/SHRINE_KEYSTORE_PASSWORD/$ADMIN_PASSWORD/g" "$CATALINA_HOME/conf/server.xml" && \

    sed -i "s/SHRINE_DB_PASSWORD/$DB_PASSWORD/g" "$CATALINA_HOME/conf/context.xml"
    #&& \

    #sed -i "s/SHRINE_WEBCLIENT_DOMAIN/$I2B2_DOMAIN_NAME/g" "$CATALINA_HOME/webapps/shrine-client/i2b2_config_data.js" && \
    #sed -i "s/SHRINE_WEBCLIENT_NAME/Domain $I2B2_DOMAIN_NAME/g" "$CATALINA_HOME/webapps/shrine-client/i2b2_config_data.js"

# configuration dependent on the arguments todo: log
ARG SHRINE_DEBUG_LEVEL_ARG="?"
ARG SHRINE_NODE_NAME_ARG="MedCo node 1"
ARG SHRINE_KEYSTORE_NODE_ALIAS_ARG="TODO"
ENV SHRINE_DEBUG_LEVEL="$SHRINE_DEBUG_LEVEL_ARG" \
    SHRINE_NODE_NAME="$SHRINE_NODE_NAME_ARG" \
    SHRINE_KEYSTORE_NODE_ALIAS="$SHRINE_KEYSTORE_NODE_ALIAS_ARG"
RUN sed -i "s/SHRINE_KEYSTORE_PRIVATE_KEY_ALIAS/$SHRINE_KEYSTORE_NODE_ALIAS-private/g" "$CATALINA_HOME/conf/server.xml" && \
    sed -i "s#SHRINE_KEYSTORE_FILE_PATH#$CONF_DIR/$SHRINE_KEYSTORE_NODE_ALIAS.keystore#g" "$CATALINA_HOME/conf/server.xml"


EXPOSE 6060 6443
VOLUME $CONF_DIR
WORKDIR $CONF_DIR

