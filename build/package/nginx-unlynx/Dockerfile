FROM nginx:1.15.10

COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
COPY conf.d/servers.conf /etc/nginx/servers.conf

RUN chmod a+x /usr/local/bin/docker-entrypoint.sh

CMD docker-entrypoint.sh
