FROM jboss/wildfly:10.0.0.Final as dev

# build-time variables
ENV MEDCO_CONF_DIR="/medco-configuration" \
    MEDCO_CELL_SRC_DIR="/medco-cell-src" \
    MEDCO_CELL_VERSION="master"

USER root
COPY install-medco-cell.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/install-medco-cell.sh && \
    mkdir -p "$MEDCO_CELL_SRC_DIR" && \
    chown -R jboss:jboss "$MEDCO_CELL_SRC_DIR"
RUN yum update -y && \
    yum -y install git ant && \
    yum clean all
USER jboss

RUN install-medco-cell.sh

# -------------------------------------------
FROM medco/i2b2:latest as prod

# run-time variables
ENV NODE_IDX="0" \
    UNLYNX_DEBUG_LEVEL="1"

# build-time variables
ENV MEDCO_CONF_DIR="/medco-configuration"

COPY --from=dev --chown=jboss:jboss "$JBOSS_HOME/standalone/configuration" "$JBOSS_HOME/standalone/configuration"
COPY --from=dev --chown=jboss:jboss "$JBOSS_HOME/standalone/deployments" "$JBOSS_HOME/standalone/deployments"
COPY sql/* "$I2B2_SQL_DIR/"
COPY pre-init-scripts/* "$PRE_INIT_SCRIPT_DIR/"

VOLUME $MEDCO_CONF_DIR
