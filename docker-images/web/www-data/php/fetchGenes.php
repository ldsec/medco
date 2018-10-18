<?php
header('Access-Control-Allow-Origin: ' getenv('CORS_ALLOW_ORIGIN')); 
header('Access-Control-Allow-Credentials: true'); 
header('Access-Control-Allow-Headers: origin, content-type, accept, authorization'); 
header('Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS, HEAD'); 

// case sensitive % stands for .*
//select *
//from gene_values
//where gene_value LIKE '%A1%' //
//LIMIT 20

// case insensitive with regex
//select *
//from gene_values
//where gene_value ~* '.*a1.*'
//LIMIT 20

include 'sqlConnection.php';

// get the row which contains all the values of the passed annotation
$query =
    "SELECT gene_value 
FROM gene_values
WHERE gene_value ~* '.*" . $_GET["gene"] .".*'
LIMIT " . $_GET["limit"];

$result = pg_query($conn, $query);
if (!$result) {
    echo "An error occurred while querying the database.\n";
    exit;
}

// In json format return the list of genes
$geneList = "";
while ($row = pg_fetch_row($result)) {
    $geneList .= "\"$row[0]\",";
}
// drop the last comma and concatenate in json format
echo "[" . substr($geneList, 0, -1) . "]";
