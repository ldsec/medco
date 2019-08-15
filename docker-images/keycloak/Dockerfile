FROM jboss/keycloak:6.0.1

USER root

COPY proxy-update.cli /opt/jboss/tools/cli/

# invoke configuration commands to take into account the change in build-time environment variable
ENV PROXY_ADDRESS_FORWARDING="true"
RUN cd /opt/jboss/keycloak && \
    bin/jboss-cli.sh --file=/opt/jboss/tools/cli/proxy-update.cli && \
    rm -rf /opt/jboss/keycloak/standalone/configuration/standalone_xml_history && \
    chown -R jboss:0 /opt/jboss/keycloak && \
    chmod -R g+rw /opt/jboss/keycloak

USER 1000
