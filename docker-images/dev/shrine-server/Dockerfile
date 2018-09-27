FROM tomcat:9.0-jre8

# pre-existing variables: CATALINA_HOME
ENV SHRINE_VERSION="fork/1.22.8-medco" \
    SHRINE_SRC_DIR="/opt/shrine-src" \
    ADMIN_PASSWORD="prigen2017" \
    BUILD_DIR="/opt/build-dir"

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

COPY copy-shrine-server-binary-dev.sh /root
ENTRYPOINT ["sh", "-c", "/root/copy-shrine-server-binary-dev.sh"]
