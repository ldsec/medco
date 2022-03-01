FROM ghcr.io/erebe/wstunnel:latest

# WS tunnel server variables
ENV SERVER_LISTEN_URL=ws://0.0.0.0:2003 \
    SERVER_DEST_ADDRESS=medco-unlynx:2001

# WS tunnel client variables, index can range from 000 to 999 (must be padded to 3 digits)
# (variables given as examples and documentation, not actually used)
ENV CLIENT_XXX_LISTEN_ADDRESS=0.0.0.0:3000 \
    CLIENT_XXX_WS_SERVER_URL=ws://medco-unlynx-wstunnel-srv0:2003 \
    CLIENT_XXX_WS_SERVER_PATH_PREFIX="" \
    CLIENT_XXX_DEST_ADDRESS=medco-unlynx:2001

USER root
COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN apk update && apk add bash ca-certificates && rm -rf /var/cache/apk/* && \
    chmod a+x /usr/local/bin/docker-entrypoint.sh

VOLUME "/medco-configuration"
CMD docker-entrypoint.sh
