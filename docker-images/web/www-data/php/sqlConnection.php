<?php
header('Access-Control-Allow-Origin: ' getenv('CORS_ALLOW_ORIGIN')); 
header('Access-Control-Allow-Credentials: true'); 
header('Access-Control-Allow-Headers: origin, content-type, accept, authorization'); 
header('Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS, HEAD'); 

$conn = pg_connect("host=I2B2_DB_HOST port=I2B2_DB_PORT dbname=I2B2_DB_NAME user=I2B2_DB_USER password=I2B2_DB_PW");

if (!$conn) {
    echo "Error while connecting to the postgres database.";
    exit;
}
