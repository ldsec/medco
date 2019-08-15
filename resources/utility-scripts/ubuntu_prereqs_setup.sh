#!/usr/bin/env bash
set -Eeuo pipefail

apt-get -y update
apt-get -y install software-properties-common
add-apt-repository ppa:git-core/ppa -y
curl -s https://packagecloud.io/install/repositories/github/git-lfs/script.deb.sh | bash
apt-get -y install git-lfs

apt-get -y install git apt-transport-https curl software-properties-common screen nano

apt-get remove docker docker-engine docker.io
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -

apt-key fingerprint 0EBFCD88
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -

add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"
add-apt-repository "deb http://apt.postgresql.org/pub/repos/apt/ $(lsb_release -cs)-pgdg main"

apt-get -y update
apt-get -y install docker-ce postgresql-client-10

docker run hello-world

curl -L "https://github.com/docker/compose/releases/download/1.23.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose
mkdir /opt/medco
git clone https://github.com/lca1/medco-deployment.git /opt/medco/medco-deployment

mkdir /disk
mount /dev/sdb /disk
