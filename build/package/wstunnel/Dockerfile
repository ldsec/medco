FROM ghcr.io/erebe/wstunnel:latest

# WS tunnel server variables
ENV SERVER_LISTEN_URL=ws://0.0.0.0:2003 \
    SERVER_DEST_ADDRESS=medco-unlynx:2001

# WS tunnel client variables, values are separated by a comma, all must be the same number of elements
ENV CLIENTS_LISTEN_ADDRESSES=0.0.0.0:3000,0.0.0.0:3001,0.0.0.0:3002 \
    CLIENTS_WS_SERVER_URLS=wss://nginx:443,wss://nginx:443,wss://nginx:443 \
    CLIENTS_WS_SERVER_PATH_PREFIXES=local-3nodes/wstunnel-unlynx-0,local-3nodes/wstunnel-unlynx-1,local-3nodes/wstunnel-unlynx-2 \
    CLIENTS_DEST_ADDRESSES=medco-unlynx-srv0:2001,medco-unlynx-srv1:2003,medco-unlynx-srv2:2005

USER root
COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN apk update && apk add bash ca-certificates && rm -rf /var/cache/apk/* && \
    chmod a+x /usr/local/bin/docker-entrypoint.sh

VOLUME "/medco-configuration"
CMD docker-entrypoint.sh
