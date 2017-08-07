FROM tomcat:8.0-jre8

# variables, ? exists and set
ARG I2B2_DOMAIN_NAME_ARG="i2b2demo"
ENV SHRINE_VERSION="1.22.8" \
    SHRINE_SRC_DIR="/opt/shrine-src" \
    I2B2_DOMAIN_NAME="$I2B2_DOMAIN_NAME_ARG"


# pre-requisites TODO: hardcoded password
RUN echo "deb http://httpredir.debian.org/debian jessie-backports main contrib non-free" >> /etc/apt/sources.list && \
#$JBOSS_HOME/bin/add-user.sh admin prigen2017 --silent && \ # /usr/local/tomcat/conf/tomcat-users.xml
    apt-get -y update && \
    apt-get -y install git maven zip wget unzip && \
    apt-get -t jessie-backports install postgresql-client-9.6 && \
    apt-get -y clean && \
#    mkdir "$I2B2_SRC_DIR" && chown -R jboss:jboss "$I2B2_SRC_DIR"


# download and compile shrine
WORKDIR "$SHRINE_SRC_DIR"
RUN git clone https://open.med.harvard.edu/stash/git/SHRINE/shrine.git . && \
    git checkout tags/$SHRINE_VERSION

RUN mvn clean install

## i2b2 psql configuration
COPY shrine-i2b2-db.sh .
RUN chmod +x ./shrine-i2b2-db.sh && /bin/bash -c ./shrine-i2b2-db.sh

