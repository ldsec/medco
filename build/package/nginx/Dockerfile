FROM nginx:1.15.10

# run-time variables
ENV HTTP_SCHEME="http" \
    ALL_TIMEOUTS_SECONDS="600" \
    PROD_CONFIG="true"

COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN chmod a+x /usr/local/bin/docker-entrypoint.sh

CMD docker-entrypoint.sh
