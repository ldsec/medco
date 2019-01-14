server {
    listen 80;
    server_name _;

    rewrite ^ https://$http_host$request_uri? permanent;
}

server {
	listen 443 ssl default_server;
    server_name _;

    include /etc/nginx/conf.d/server-common.conf.inc;

    # --- web root configuration ---
	location / {
        root   /www-data/;
        index  index.html index.htm;

        # php configuration
        location ~ \.php$ {
            #try_files $uri /index.php =404;
            fastcgi_pass php-fpm:9000;
            #fastcgi_index index.php;
            #fastcgi_buffers 16 16k;
            #fastcgi_buffer_size 32k;
            fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
            fastcgi_read_timeout 600;
            include fastcgi_params;
        }
	}

    # --- medco public key configuration ---
    location /cothority-key.pub.toml {
        alias /medco-configuration/group.toml;
    }
}

server {
    listen 443 ssl;
    server_name i2b2-medco.*;

    include /etc/nginx/conf.d/server-common.conf.inc;
    set $upstream_host i2b2-medco;

    location / {
        proxy_pass http://$upstream_host:8080;
    }

    location ^~ /wildfly-admin/ {
        proxy_pass http://$upstream_host:9990;
    }
}

server {
    listen 443 ssl;
    server_name glowing-bear.*;

    include /etc/nginx/conf.d/server-common.conf.inc;
    set $upstream_host glowing-bear;

    location / {
        proxy_pass http://$upstream_host;
    }
}

server {
    listen 443 ssl;
    server_name pic-sure.*;

    include /etc/nginx/conf.d/server-common.conf.inc;
    set $upstream_host irct;

    location / {
        proxy_pass http://$upstream_host:8080;
    }

    location ^~ /wildfly-admin/ {
        proxy_pass http://$upstream_host:9990;
    }
}

server {
    listen 443 ssl;
    server_name keycloak.*;

    include /etc/nginx/conf.d/server-common.conf.inc;
    set $upstream_host keycloak;

    location / {
        proxy_pass http://$upstream_host:8080;
    }
}

server {
    listen 443 ssl;
    server_name pg-admin.*;

    include /etc/nginx/conf.d/server-common.conf.inc;
    set $upstream_host pg-admin;

    location / {
        proxy_pass http://$upstream_host;
        proxy_redirect  http:// https://; << only for https
    }
}
