# --- web root configuration ---
location / {
    root   /www-data/;

    # redirect root to glowing-bear
    rewrite ^(/)$ /glowing-bear;
}
