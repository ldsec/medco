resolver 127.0.0.11 valid=30s;

error_page   500 502 503 504  /50x.html;
    location = /50x.html {
    root   /usr/share/nginx/html;
}

client_max_body_size 10m;
client_body_buffer_size 128k;
proxy_next_upstream error timeout invalid_header http_500 http_502 http_503;
send_timeout 10m;
proxy_read_timeout 240;
proxy_send_timeout 240;
proxy_connect_timeout 240;
proxy_set_header Host $host;
proxy_set_header X-Real-IP $remote_addr;
proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
proxy_set_header X-Forwarded-Proto $scheme;
proxy_set_header X-Forwarded-Host $host;
proxy_set_header X-Forwarded-Port $server_port;
proxy_set_header X-Forwarded-Ssl on;
#proxy_redirect  http://  $scheme://;
proxy_http_version 1.1;
proxy_set_header Connection "";
proxy_cache_bypass $cookie_session;
proxy_no_cache $cookie_session;
proxy_buffers 32 4k;

location /auth {
    set $upstream_host keycloak;
    proxy_pass http://$upstream_host:8080;
}

location /pgadmin {
    set $upstream_host pg-admin;
    proxy_pass http://$upstream_host;
    proxy_set_header X-Script-Name /pgadmin;
    proxy_redirect http://$upstream_host/pgadmin /pgadmin;
}

location /pic-sure-api-2 {

    add_header 'Access-Control-Allow-Origin' "$http_origin";
    add_header 'Access-Control-Allow-Methods' 'GET, POST, PUT, HEAD, DELETE, OPTIONS';
    add_header 'Access-Control-Allow-Headers' 'Authorization,DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range';
    add_header 'Access-Control-Expose-Headers' 'Content-Length,Content-Range';

    if ($request_method = OPTIONS) {
        return 200;
    }

    set $upstream_host picsure;
    proxy_pass http://$upstream_host:8080;
}

location /i2b2 {
    set $upstream_host i2b2;
    proxy_pass http://$upstream_host:8080;
}

location /glowing-bear {
    set $upstream_host glowing-bear-medco;
    proxy_pass http://$upstream_host;
}

location /medco-connector {
    set $upstream_host medco-connector;
    proxy_pass http://$upstream_host;
}
