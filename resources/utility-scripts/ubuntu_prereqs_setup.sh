#!/usr/bin/env bash

apt-get -y update
apt-get -y install git apt-transport-https ca-certificates curl software-properties-common screen default-jre nano unzip

apt-get remove docker docker-engine docker.io
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -

apt-key fingerprint 0EBFCD88
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -

add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"
add-apt-repository "deb http://apt.postgresql.org/pub/repos/apt/ trusty-pgdg main"

apt-get -y update
apt-get -y install docker-ce postgresql-client-10

docker run hello-world

curl -L "https://github.com/docker/compose/releases/download/1.23.0-rc1/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose
git clone https://github.com/lca1/medco-deployment.git
