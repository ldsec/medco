FROM nginx:1.20.1

# run-time variables
ENV ALL_TIMEOUTS_SECONDS="600" \
    PROD_CONFIG="true"

COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN chmod a+x /usr/local/bin/docker-entrypoint.sh

CMD docker-entrypoint.sh
