<?php
// case insensitive with regex
//select variant_name
//from variant_names
//where variant_name ~* '.*a1.*'
//LIMIT 20

include 'sqlConnection.php';

// escape the ? for regex query
$variant_name = str_replace("?", "\?", $_GET["variant_name"]);
// get the row which contains all the values of the passed annotation
$query =
    "SELECT variant_name 
FROM genomic_annotations
WHERE variant_name ~* '.*" . $variant_name .".*'
LIMIT " . $_GET["limit"];

$result = pg_query($conn, $query);
if (!$result) {
    echo "An error occurred while querying the database.\n";
    exit;
}

// In json format return the list of genes
$variantNames = "";
while ($row = pg_fetch_row($result)) {
    $variantNames .= "\"$row[0]\",";
}
// drop the last comma and concatenate in json format
echo "[" . substr($variantNames, 0, -1) . "]";
