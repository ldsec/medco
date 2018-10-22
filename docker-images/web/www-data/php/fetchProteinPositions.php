<?php
header('Access-Control-Allow-Origin: '.getenv('CORS_ALLOW_ORIGIN'));
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
$stmt = $pdo->prepare("SELECT annotation_value FROM protein_position WHERE annotation_value ~* '.*?.*' LIMIT ?");
$stmt->execute([$_GET["proteinPosition"], $_GET["limit"]]);

// In json format return the list of annotation names
$proteinPositions = "";
while ($row = $stmt->fetch()) {
    $proteinPositions .= "\"$row[0]\",";
}
// drop the last comma and concatenate in json format
echo "[" . substr($proteinPositions, 0, -1) . "]";
?>