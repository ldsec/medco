#!/bin/bash
set -Eeuo pipefail

sed -i "/medco.unlynx.entrypointidx/c\medco.unlynx.entrypointidx=$NODE_IDX" "$JBOSS_HOME/standalone/configuration/medcoapp/medco.properties"
sed -i "/medco.unlynx.debuglevel/c\medco.unlynx.debuglevel=$UNLYNX_DEBUG_LEVEL" "$JBOSS_HOME/standalone/configuration/medcoapp/medco.properties"
