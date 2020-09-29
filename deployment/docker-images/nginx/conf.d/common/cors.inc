add_header 'Access-Control-Allow-Origin' "$http_origin" always;
add_header 'Access-Control-Allow-Methods' 'GET, POST, PUT, HEAD, DELETE, OPTIONS' always;
add_header 'Access-Control-Allow-Headers' 'Authorization,DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range' always;
add_header 'Access-Control-Expose-Headers' 'Content-Length,Content-Range' always;

if ($request_method = OPTIONS) {
return 200;
}
