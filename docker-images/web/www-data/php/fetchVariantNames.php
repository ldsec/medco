<?php
header('Access-Control-Allow-Origin: '.getenv('CORS_ALLOW_ORIGIN')); 
header('Access-Control-Allow-Credentials: true');
header('Access-Control-Allow-Headers: origin, content-type, accept, authorization');
header('Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS, HEAD');

// case insensitive with regex

//select variant_name
//from variant_names
//where variant_name ~* '.*a1.*'
//LIMIT 20

include 'sqlConnection.php';

// escape the ? for regex query
$variant_name = str_replace("?", "\?", $_GET["variant_name"]);
// get the row which contains all the values of the passed annotation
$stmt = $pdo->prepare("SELECT variant_name FROM genomic_annotations WHERE variant_name ~* '.*?.*' LIMIT ?");
$stmt->execute([$variant_name, $_GET["limit"]]);

// In json format return the list of genes
$variantNames = "";
while ($row = $stmt->fetch()) {
    $variantNames .= "\"$row[0]\",";
}
// drop the last comma and concatenate in json format
echo "[" . substr($variantNames, 0, -1) . "]";
?>