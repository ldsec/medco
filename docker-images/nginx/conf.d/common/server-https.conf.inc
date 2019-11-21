ssl_certificate /medco-configuration/certificate.crt;
ssl_certificate_key /medco-configuration/certificate.key;
ssl_protocols TLSv1.2;
ssl_ecdh_curve secp384r1;
ssl_ciphers EECDH+AESGCM:EECDH+CHACHA20:EECDH+AES;
ssl_prefer_server_ciphers on;
