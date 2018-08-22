FROM tomcat:9.0-jre8

ARG SHRINE_DEBUG_LEVEL_ARG
ARG NODE_IDX_ARG

# pre-existing variables: CATALINA_HOME
ENV CONF_DIR="/opt/medco-configuration" \
	BUILD_DIR="/opt/build-dir" \
	SHRINE_DEBUG_LEVEL="$SHRINE_DEBUG_LEVEL_ARG" \
    NODE_IDX="$NODE_IDX_ARG"

# system and tomcat prerequisites
RUN apt-get -y update && \
    apt-get -y install git maven zip wget unzip openjdk-8-jdk-headless && \
    apt-get -y clean && \
    echo "<?xml version='1.0' encoding='utf-8'?><tomcat-users><role rolename=\"manager-gui\" /><role rolename=\"admin-gui\" />" \
        "<user username=\"admin\" password=\"$ADMIN_PASSWORD\" roles=\"manager-gui,admin-gui\" /></tomcat-users>" > \
        "$CATALINA_HOME/conf/tomcat-users.xml" && \
    echo 'export CATALINA_OPTS=" -Dakka.daemonic=on "' > "$CATALINA_HOME/bin/setenv.sh" && \
    echo '{ "allow_root": true }' > /root/.bowerrc

WORKDIR $CONF_DIR
# copy the build files from the build volume to the correct folder
COPY copy-shrine-server-binary-prod.sh /root
ENTRYPOINT ["sh", "-c", "/root/copy-shrine-server-binary-prod.sh && $CATALINA_HOME/bin/catalina.sh run"]