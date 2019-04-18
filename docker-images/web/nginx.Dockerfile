FROM nginx:1.15.10

# run-time variables
ENV HTTP_SCHEME="http"

# run
CMD /bin/bash -c "envsubst < /etc/nginx/conf.d/servers.conf.template > /etc/nginx/conf.d/servers.conf && exec nginx -g 'daemon off;'"
