<?php
header('Access-Control-Allow-Origin: ' getenv('CORS_ALLOW_ORIGIN')); 
header('Access-Control-Allow-Credentials: true'); 
header('Access-Control-Allow-Headers: origin, content-type, accept, authorization'); 
header('Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS, HEAD'); 

// case insensitive with regex
//select annotation_value
//from Protein_position
//where annotation_value ~* '.*a1.*'
//LIMIT 20

include 'sqlConnection.php';

// get the row which contains all the values of the passed annotation
$query =
    "SELECT annotation_value 
FROM protein_position
WHERE annotation_value ~* '.*" . $_GET["proteinPosition"] .".*'
LIMIT " . $_GET["limit"];

$result = pg_query($conn, $query);
if (!$result) {
    echo "An error occurred while querying the database.\n";
    exit;
}

// In json format return the list of annotation names
$proteinPositions = "";
while ($row = pg_fetch_row($result)) {
    $proteinPositions .= "\"$row[0]\",";
}
// drop the last comma and concatenate in json format
echo "[" . substr($proteinPositions, 0, -1) . "]";