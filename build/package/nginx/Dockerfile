FROM nginx:1.15.10

# run-time variables
ENV HTTP_SCHEME="http" \
    ALL_TIMEOUTS_SECONDS="600" \
    PROD_CONFIG="true" \
    MEDCO_NODE_IDX=0 \
    UNLYNX_PORT_0=2001 \
    UNLYNX_PORT_1=2002

RUN sed  -i '1i include stream.conf;' /etc/nginx/nginx.conf

COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN chmod a+x /usr/local/bin/docker-entrypoint.sh

CMD docker-entrypoint.sh
