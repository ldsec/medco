server {
    listen 80 default_server;
    server_name _;

    include /etc/nginx/conf.d/common/server-default.conf.inc;
    include /etc/nginx/conf.d/common/server-revproxy-base.conf.inc;
    include /etc/nginx/conf.d/common/server-revproxy-dev.conf.inc;
    include /etc/nginx/conf.d/common/security-headers.inc;
}
