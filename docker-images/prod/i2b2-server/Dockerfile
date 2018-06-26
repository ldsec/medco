FROM jboss/wildfly:12.0.0.Final

ARG NODE_IDX_ARG
ARG UNLYNX_DEBUG_LEVEL_ARG

ENV CONF_DIR="/opt/medco-configuration" \
	BUILD_DIR="/opt/build-dir" \
	MEDCO_CELL_SRC_DIR="/opt/medco-src" \
    MEDCO_CELL_GIT_BRANCH="master"

# pre-requisites
USER root
RUN yum -y install git ant

## MedCo cell
WORKDIR "$MEDCO_CELL_SRC_DIR"
RUN git clone https://c4science.ch/source/medco-i2b2-cell.git . && \
    git checkout $MEDCO_CELL_GIT_BRANCH

ENV NODE_IDX="$NODE_IDX_ARG" \
    UNLYNX_DEBUG_LEVEL="$UNLYNX_DEBUG_LEVEL_ARG"

RUN sed -i "/jboss.home/c\jboss.home=$JBOSS_HOME" build.properties && \
    sed -i "/medco.unlynx.groupfilepath/c\medco.unlynx.groupfilepath=$CONF_DIR/group.toml" etc/spring/medcoapp/medco.properties && \
    sed -i "/medco.unlynx.binarypath/c\medco.unlynx.binarypath=$BUILD_DIR/medco" etc/spring/medcoapp/medco.properties && \
    sed -i "/medco.unlynx.entrypointidx/c\medco.unlynx.entrypointidx=$NODE_IDX" etc/spring/medcoapp/medco.properties && \
    sed -i "/medco.unlynx.debuglevel/c\medco.unlynx.debuglevel=$UNLYNX_DEBUG_LEVEL" etc/spring/medcoapp/medco.properties && \
    ant -f build.xml clean all deploy

# copy the build files from the build volume to the correct folder
COPY copy-i2b2-server-binary-prod.sh /root
ENTRYPOINT ["sh", "-c", "/root/copy-i2b2-server-binary-prod.sh && /opt/jboss/wildfly/bin/standalone.sh -b 0.0.0.0 -bmanagement 0.0.0.0"]