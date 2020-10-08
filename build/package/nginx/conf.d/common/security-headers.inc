# do not send out web server info
server_tokens off;

# various HTTP headers for enhanced security
add_header 'X-Frame-Options' 'sameorigin' always;
add_header 'X-XSS-Protection' '1; mode=block' always;
add_header 'X-Content-Type-Options' 'nosniff' always;
add_header 'Content-Security-Policy' 'default-src \'self\' \'unsafe-inline\' \'unsafe-eval\'; connect-src *' always;

# prevent disallowed HTTP methods
if ($request_method !~* (GET|POST|OPTIONS|HEAD|DELETE|PUT)) {
return 405;
}

# TODO
## Content-Security-Policy: allow connect-src to other MedCo node;
## Feature-Policy: blacklist unused features (most of them)
