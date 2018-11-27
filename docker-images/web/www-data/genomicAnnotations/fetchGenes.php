<?php
header('Access-Control-Allow-Origin: '.getenv('CORS_ALLOW_ORIGIN')); 
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

//fetchGenes.php?gene=AA&limit=10

// get the row which contains all the values of the passed annotation
$gene = ".*".$_GET["gene"].".*";
$stmt = $pdo->prepare("SELECT gene_value FROM genomic_annotations.gene_values WHERE gene_value ~* ? LIMIT ?");
$stmt->bindValue(1, $gene, PDO::PARAM_STR);
$stmt->bindValue(2, $_GET["limit"], PDO::PARAM_STR);
$stmt->execute();

// In json format return the list of genes
$geneList = "";
while ($row = $stmt->fetch()) {
    $geneList .= "\"$row[0]\",";
}
// drop the last comma and concatenate in json format
echo "[" . substr($geneList, 0, -1) . "]";
?>