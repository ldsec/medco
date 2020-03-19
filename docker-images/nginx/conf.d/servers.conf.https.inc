server {
    listen 80;
    server_name _;

    rewrite ^ https://$http_host$request_uri? permanent;
}

server {
	listen 443 ssl default_server;
    server_name _;

    include /etc/nginx/conf.d/common/server-default.conf.inc;
    include /etc/nginx/conf.d/common/server-revproxy.conf.inc;
    include /etc/nginx/conf.d/common/server-https.conf.inc;
    include /etc/nginx/conf.d/common/security-headers.inc;
}
