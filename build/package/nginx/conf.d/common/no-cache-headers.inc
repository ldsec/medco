# disable browser caching
add_header 'Pragma' 'no-cache' always;
add_header 'Cache-Control' 'private, max-age=0, no-cache' always;
add_header 'Expires' '01/01/99 20:00:00 GMT' always;