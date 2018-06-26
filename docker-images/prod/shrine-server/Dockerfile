FROM tomcat:9.0-jre8

ARG SHRINE_DEBUG_LEVEL_ARG
ARG NODE_IDX_ARG

# pre-existing variables: CATALINA_HOME
ENV CONF_DIR="/opt/medco-configuration" \
	BUILD_DIR="/opt/build-dir" \
	SHRINE_DEBUG_LEVEL="$SHRINE_DEBUG_LEVEL_ARG" \
    NODE_IDX="$NODE_IDX_ARG"

COPY shrine-server $CATALINA_HOME

RUN sed -i "s/SHRINE_KEYSTORE_PRIVATE_KEY_ALIAS/srv$NODE_IDX-private/g" "$CATALINA_HOME/conf/server.xml" && \
	sed -i "s#SHRINE_KEYSTORE_FILE_PATH#$CONF_DIR/srv$NODE_IDX.keystore#g" "$CATALINA_HOME/conf/server.xml" && \
	sed -i "s#FINE#$SHRINE_DEBUG_LEVEL#g" "$CATALINA_HOME/conf/logging.properties" && \
	sed -i "s#INFO#$SHRINE_DEBUG_LEVEL#g" "$CATALINA_HOME/conf/logging.properties"

WORKDIR $CONF_DIR
# copy the build files from the build volume to the correct folder
#COPY copy-shrine-server-binary-prod.sh /root
#ENTRYPOINT ["sh", "-c", "/root/copy-shrine-server-binary-prod.sh"]