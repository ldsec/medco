#!/usr/bin/env bash
set -Eeuo pipefail

# get medco cell sources
pushd "$MEDCO_CELL_SRC_DIR"
git clone --depth 1 --branch $MEDCO_CELL_VERSION https://c4science.ch/source/medco-i2b2-cell.git .

# compile and deploy
sed -i "/jboss.home/c\jboss.home=$JBOSS_HOME" build.properties
sed -i "/medco.unlynx.groupfilepath/c\medco.unlynx.groupfilepath=$MEDCO_CONF_DIR/group.toml" etc/spring/medcoapp/medco.properties
sed -i "/medco.unlynx.binarypath/c\medco.unlynx.binarypath=$MEDCO_CONF_DIR/unlynxMedCo" etc/spring/medcoapp/medco.properties
ant -f build.xml clean all deploy
popd
