<?php
header('Access-Control-Allow-Origin: '.getenv('CORS_ALLOW_ORIGIN')); 
header('Access-Control-Allow-Credentials: true'); 
header('Access-Control-Allow-Headers: origin, content-type, accept, authorization'); 
header('Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS, HEAD'); 

include 'sqlConnection.php';

// get the row which contains all the values of the passed annotation
$stmt = $pdo->prepare("SELECT concept_path FROM i2b2demodata_i2b2.concept_dimension WHERE concept_cd = ?");
$stmt->execute([$_GET['concept_cd']]);

$concept_path = "";
while ($row = $stmt->fetch()) {
    $concept_path .= "\"$row[0]\",";
}
// drop the last comma and concatenate in json format
echo "[" . substr($concept_path, 0, -1) . "]";
?> 