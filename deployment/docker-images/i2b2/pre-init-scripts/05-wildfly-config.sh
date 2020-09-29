#!/bin/bash
set -Eeuo pipefail

$JBOSS_HOME/bin/add-user.sh admin $WILDFLY_ADMIN_PASSWORD --silent || true
