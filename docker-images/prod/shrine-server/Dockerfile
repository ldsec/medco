FROM tomcat:9.0-jre8

ARG SHRINE_DEBUG_LEVEL_ARG
ARG NODE_IDX_ARG

# pre-existing variables: CATALINA_HOME
ENV CONF_DIR="/opt/medco-configuration" \
	BUILD_DIR="/opt/build-dir" \
	SHRINE_DEBUG_LEVEL="$SHRINE_DEBUG_LEVEL_ARG" \
    NODE_IDX="$NODE_IDX_ARG"

WORKDIR $CONF_DIR
# copy the build files from the build volume to the correct folder
COPY copy-shrine-server-binary-prod.sh /root
ENTRYPOINT ["sh", "-c", "/root/copy-shrine-server-binary-prod.sh && $CATALINA_HOME/bin/catalina.sh run"]