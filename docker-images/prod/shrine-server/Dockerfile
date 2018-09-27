FROM tomcat:9.0-jre8

ARG SHRINE_DEBUG_LEVEL_ARG="INFO"
ARG NODE_IDX_ARG="0"

# pre-existing variables: CATALINA_HOME
ENV SHRINE_ADAPTER_MAPPINGS_URL="https://documents.epfl.ch/users/m/mi/misbach/www/AdapterMappings.xml" \
	SHRINE_MYSQL_JAR_URL="http://central.maven.org/maven2/mysql/mysql-connector-java/5.1.40/mysql-connector-java-5.1.40.jar" \
	I2B2_DOMAIN_NAME="medcodeployment" \
	CONF_DIR="/opt/medco-configuration" \
	ADMIN_PASSWORD="prigen2017" \
	DB_PASSWORD="pFjy3EjDVwLfT2rB9xkK" \
	NODE_IDX="$NODE_IDX_ARG" \
	SHRINE_DEBUG_LEVEL="$SHRINE_DEBUG_LEVEL_ARG" \
	CONF_DIR="/opt/medco-configuration" \
    BUILD_DIR="/opt/build-dir"

RUN apt-get -y update && \
apt-get -y install git maven zip wget unzip openjdk-8-jdk-headless && \
apt-get -y clean && \
echo '{ "allow_root": true }' > /root/.bowerrc

COPY conf/shrine.conf "/root/"
COPY conf/server.xml conf/context.xml "/root/"

WORKDIR $CONF_DIR
# copy the build files from the build volume to the correct folder
COPY copy-shrine-server-binary-prod.sh /root
ENTRYPOINT ["sh", "-c", "/root/copy-shrine-server-binary-prod.sh && $CATALINA_HOME/bin/catalina.sh run"]