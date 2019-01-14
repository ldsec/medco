server {
    listen 80 default_server;
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
        add_header 'Access-Control-Allow-Origin' '*';
        add_header 'Access-Control-Allow-Methods' 'GET, OPTIONS';
        add_header 'Access-Control-Allow-Headers' 'DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range';
        add_header 'Access-Control-Expose-Headers' 'Content-Length,Content-Range';
    }
}
