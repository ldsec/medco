resolver 127.0.0.11 valid=30s;

error_page   500 502 503 504  /50x.html;
    location = /50x.html {
    root   /usr/share/nginx/html;
}

# --- HTTPS configuration ---
ssl_certificate /etc/nginx/conf.d/cert.crt;
ssl_certificate_key /etc/nginx/conf.d/cert.key;
ssl_protocols TLSv1.2;
ssl_ecdh_curve secp384r1;
ssl_ciphers EECDH+AESGCM:EECDH+CHACHA20:EECDH+AES;
ssl_prefer_server_ciphers on;

# --- reverse proxies configuration ---
client_max_body_size 10m;
client_body_buffer_size 128k;
proxy_next_upstream error timeout invalid_header http_500 http_502 http_503;
send_timeout 10m;
proxy_read_timeout 240;
proxy_send_timeout 240;
proxy_connect_timeout 240;
proxy_set_header Host $http_host;
proxy_set_header X-Real-IP $remote_addr;
proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
proxy_set_header X-Forwarded-Proto https;
proxy_set_header X-Forwarded-Host $http_host;
proxy_set_header X-Forwarded-Port $server_port;
proxy_set_header X-Forwarded-Ssl on;
#proxy_redirect  http://  $scheme://;
proxy_http_version 1.1;
proxy_set_header Connection "";
proxy_cache_bypass $cookie_session;
proxy_no_cache $cookie_session;
proxy_buffers 32 4k;
