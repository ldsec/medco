FROM jboss/wildfly:10.0.0.Final
# TODO: use the official image
# todo: configure https

# variables, JBOSS_HOME, WILDFLY_VERSION exists and set
ARG AXIS2_LOGLEVEL_ARG="INFO"

ENV AXIS2_VERSION="1.7.1" \
    AXIS2_LOGLEVEL="$AXIS2_LOGLEVEL_ARG" \
    I2B2_VERSION="v1.7.08b.0002" \
    I2B2_SRC_DIR="/opt/i2b2-src" \
    I2B2_DOMAIN_NAME="medcodeployment" \
    I2B2_DEMO_DB_NAME="i2b2demo" \
    I2B2_MEDCO_DB_NAME="i2b2medco" \
    I2B2_FR_FILES_DIR="/opt/FRC_files" \
    MEDCO_CELL_SRC_DIR="/opt/medco-src" \
    MEDCO_CELL_GIT_BRANCH="tags/MedCo_v0.1" \
    CONF_DIR="/opt/medco-configuration" \
    ADMIN_PASSWORD="prigen2017" \
    DB_PASSWORD="pFjy3EjDVwLfT2rB9xkK"

# pre-requisites
USER root
RUN $JBOSS_HOME/bin/add-user.sh admin $ADMIN_PASSWORD --silent && \
    sed -i 's/Xmx512m/Xmx2048m/g' $JBOSS_HOME/bin/standalone.conf && \
    yum update -y && yum -y install wget git ant && yum clean all && \
    mkdir "$I2B2_SRC_DIR" && chown -R jboss:jboss "$I2B2_SRC_DIR" && \
    mkdir "$MEDCO_CELL_SRC_DIR" && chown -R jboss:jboss "$MEDCO_CELL_SRC_DIR"

# install axis 2
USER jboss
WORKDIR "$JBOSS_HOME/standalone/deployments"
RUN wget "http://archive.apache.org/dist/axis/axis2/java/core/$AXIS2_VERSION/axis2-$AXIS2_VERSION-war.zip" && \
    mkdir "i2b2.war" && \
    unzip -j "axis2-$AXIS2_VERSION-war.zip" "axis2.war" && \
    unzip "axis2.war" -d "i2b2.war" && \
    rm "axis2.war" && \
    touch "i2b2.war.dodeploy" && \
    sed -i "/^log4j.rootCategory=/c\log4j.rootCategory=$AXIS2_LOGLEVEL, CONSOLE" i2b2.war/WEB-INF/classes/log4j.properties

# download, patch, compile and deploy i2b2
WORKDIR "$I2B2_SRC_DIR"
RUN git clone https://github.com/i2b2/i2b2-core-server.git . && \
    git checkout tags/$I2B2_VERSION
COPY patches/* ./
RUN git apply *.diff

## i2b2 wildfly datasources configuration
USER root
COPY i2b2-server-writedatasources.sh .
RUN chmod +x ./i2b2-server-writedatasources.sh && sleep 1 && /bin/bash -c ./i2b2-server-writedatasources.sh
USER jboss

## core cell
WORKDIR "$I2B2_SRC_DIR/edu.harvard.i2b2.server-common"
RUN sed -i "/jboss.home/c\jboss.home=$JBOSS_HOME" build.properties && \
    ant clean dist deploy jboss_pre_deployment_setup

## PM cell
WORKDIR "$I2B2_SRC_DIR/edu.harvard.i2b2.pm"
RUN sed -i "/jboss.home/c\jboss.home=$JBOSS_HOME" build.properties && \
    ant -f master_build.xml clean build-all deploy

## ONT cell
WORKDIR "$I2B2_SRC_DIR/edu.harvard.i2b2.ontology"
RUN sed -i "/jboss.home/c\jboss.home=$JBOSS_HOME" build.properties && \
    sed -i "/edu.harvard.i2b2.ontology.applicationdir/c\edu.harvard.i2b2.ontology.applicationdir=$JBOSS_HOME/standalone/configuration/ontologyapp" etc/spring/ontology_application_directory.properties && \
    sed -i "/ontology.ws.pm.url/c\ontology.ws.pm.url=http://localhost:8080/i2b2/services/PMService/getServices" etc/spring/ontology.properties && \
    sed -i "/edu.harvard.i2b2.ontology.ws.fr.url/c\edu.harvard.i2b2.ontology.ws.fr.url=http://localhost:8080/i2b2/services/FRService/" etc/spring/ontology.properties && \
    sed -i "/edu.harvard.i2b2.ontology.ws.crc.url/c\edu.harvard.i2b2.ontology.ws.crc.url=http://localhost:8080/i2b2/services/QueryToolService" etc/spring/ontology.properties && \
    sed -i "/edu.harvard.i2b2.ontology.pm.serviceaccount.password/c\edu.harvard.i2b2.ontology.pm.serviceaccount.password=$DB_PASSWORD" etc/spring/ontology.properties && \
    ant -f master_build.xml clean build-all deploy

## CRC cell
WORKDIR "$I2B2_SRC_DIR/edu.harvard.i2b2.crc"
RUN sed -i "/jboss.home/c\jboss.home=$JBOSS_HOME" build.properties && \
    sed -i "/edu.harvard.i2b2.crc.applicationdir/c\edu.harvard.i2b2.crc.applicationdir=$JBOSS_HOME/standalone/configuration/crcapp" etc/spring/crc_application_directory.properties && \
    sed -i "/edu.harvard.i2b2.crc.loader.ws.fr.url/c\edu.harvard.i2b2.crc.loader.ws.fr.url=http://localhost:8080/i2b2/services/FRService/" etc/spring/edu.harvard.i2b2.crc.loader.properties && \
    sed -i "/edu.harvard.i2b2.crc.loader.ws.pm.url/c\edu.harvard.i2b2.crc.loader.ws.pm.url=http://localhost:8080/i2b2/services/PMService/getServices" etc/spring/edu.harvard.i2b2.crc.loader.properties && \
    sed -i "/edu.harvard.i2b2.crc.loader.ds.lookup.servertype/c\edu.harvard.i2b2.crc.loader.ds.lookup.servertype=PostgreSQL" etc/spring/edu.harvard.i2b2.crc.loader.properties && \

    sed -i "/queryprocessor.ws.pm.url/c\queryprocessor.ws.pm.url=http://localhost:8080/i2b2/services/PMService/getServices" etc/spring/crc.properties && \
    sed -i "/queryprocessor.ds.lookup.servertype/c\queryprocessor.ds.lookup.servertype=PostgreSQL" etc/spring/crc.properties && \
    sed -i "/queryprocessor.ws.ontology.url/c\queryprocessor.ws.ontology.url=http://localhost:8080/i2b2/services/OntologyService/getTermInfo" etc/spring/crc.properties && \
    sed -i "/edu.harvard.i2b2.crc.delegate.ontology.url/c\edu.harvard.i2b2.crc.delegate.ontology.url=http://localhost:8080/i2b2/services/OntologyService" etc/spring/crc.properties && \
    sed -i "/edu.harvard.i2b2.crc.lockout.setfinderquery.count/c\edu.harvard.i2b2.crc.lockout.setfinderquery.count=-1" etc/spring/crc.properties && \
    sed -i "/edu.harvard.i2b2.crc.pm.serviceaccount.password/c\edu.harvard.i2b2.crc.pm.serviceaccount.password=$DB_PASSWORD" etc/spring/crc.properties && \
    ant -f master_build.xml clean build-all deploy

## WORK cell
WORKDIR "$I2B2_SRC_DIR/edu.harvard.i2b2.workplace"
RUN sed -i "/jboss.home/c\jboss.home=$JBOSS_HOME" build.properties && \
    sed -i "/edu.harvard.i2b2.workplace.applicationdir/c\edu.harvard.i2b2.workplace.applicationdir=$JBOSS_HOME/standalone/configuration/workplaceapp" etc/spring/workplace_application_directory.properties && \
    sed -i "/workplace.ws.pm.url/c\workplace.ws.pm.url=http://localhost:8080/i2b2/services/PMService/getServices" etc/spring/workplace.properties && \
    ant -f master_build.xml clean build-all deploy

## FR cell
WORKDIR "$I2B2_SRC_DIR/edu.harvard.i2b2.fr"
RUN sed -i "/jboss.home/c\jboss.home=$JBOSS_HOME" build.properties && \
    sed -i "/edu.harvard.i2b2.fr.applicationdir/c\edu.harvard.i2b2.fr.applicationdir=$JBOSS_HOME/standalone/configuration/frapp" etc/spring/fr_application_directory.properties && \
    sed -i "/edu.harvard.i2b2.fr.ws.pm.url/c\edu.harvard.i2b2.fr.ws.pm.url=http://localhost:8080/i2b2/services/PMService/getServices" etc/spring/edu.harvard.i2b2.fr.properties && \
    ant -f master_build.xml clean build-all deploy
USER root
RUN mkdir "$I2B2_FR_FILES_DIR" && chown -R jboss:jboss "$I2B2_FR_FILES_DIR"
USER jboss

## IM cell
WORKDIR "$I2B2_SRC_DIR/edu.harvard.i2b2.im"
RUN sed -i "/jboss.home/c\jboss.home=$JBOSS_HOME" build.properties && \
    sed -i "/edu.harvard.i2b2.im.applicationdir/c\edu.harvard.i2b2.im.applicationdir=$JBOSS_HOME/standalone/configuration/imapp" etc/spring/im_application_directory.properties && \
    sed -i "/im.ws.pm.url/c\im.ws.pm.url=http://localhost:8080/i2b2/services/PMService/getServices" etc/spring/im.properties && \
    ant -f master_build.xml clean build-all deploy

# MedCo cell
## download
WORKDIR "$MEDCO_CELL_SRC_DIR"
RUN git clone https://c4science.ch/source/medco-i2b2-cell.git . && \
    git checkout $MEDCO_CELL_GIT_BRANCH

## install
ARG NODE_IDX_ARG="0"
ARG UNLYNX_DEBUG_LEVEL_ARG="1"
ENV NODE_IDX="$NODE_IDX_ARG" \
    UNLYNX_DEBUG_LEVEL="$UNLYNX_DEBUG_LEVEL_ARG"

RUN sed -i "/jboss.home/c\jboss.home=$JBOSS_HOME" build.properties && \
    sed -i "/medco.unlynx.groupfilepath/c\medco.unlynx.groupfilepath=$CONF_DIR/group.toml" etc/spring/medcoapp/medco.properties && \
    sed -i "/medco.unlynx.binarypath/c\medco.unlynx.binarypath=$CONF_DIR/unlynxMedCo" etc/spring/medcoapp/medco.properties && \
    sed -i "/medco.unlynx.entrypointidx/c\medco.unlynx.entrypointidx=$NODE_IDX" etc/spring/medcoapp/medco.properties && \
    sed -i "/medco.unlynx.debuglevel/c\medco.unlynx.debuglevel=$UNLYNX_DEBUG_LEVEL" etc/spring/medcoapp/medco.properties && \
    ant -f build.xml clean all deploy

# run
VOLUME $I2B2_FR_FILES_DIR $CONF_DIR
EXPOSE 8080 9990
ENTRYPOINT exec /opt/jboss/wildfly/bin/standalone.sh -b 0.0.0.0 -bmanagement 0.0.0.0
