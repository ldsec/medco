<?php
header('Access-Control-Allow-Origin: '.getenv('CORS_ALLOW_ORIGIN')); 
header('Access-Control-Allow-Credentials: true'); 
header('Access-Control-Allow-Headers: origin, content-type, accept, authorization'); 
header('Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS, HEAD'); 

// query_type define is we are searching by:
// - old is the old version that uses the first version of the schema
// - Gene + Zygosity
// - ...

include 'sqlConnection.php';

$query = "";
switch($_GET["query_type"]){
    case "gene_and_zygosity":
        // case insensitive regex query
        // select variant_id
        // from genomic_annotations_new
        // where annotations ~* '^CTBP2P1;(Homozygous|Unknown);'

        $zigosity = $_GET["zygosity"]; // array of zygosity options, put them in the query separated by | (or)

        $stmt = $pdo->prepare("SELECT variant_id FROM genomic_annotations WHERE hugo_gene_symbol=? AND annotations ~* '^(?)");
        $stmt->execute([$_GET["gene_value"], join('|', $zigosity)]);
        break;

    case "protein_position_and_zygosity":
        $zigosity = $_GET["zygosity"]; // array of zygosity options, put them in the query separated by | (or)

        $stmt = $pdo->prepare("SELECT variant_id FROM genomic_annotations WHERE protein_change=? AND annotations ~* '^(?)");
        $stmt->execute([$_GET["protein_change_value"], join('|', $zigosity)]);
        break;

    case "variantName_and_zygosity":
        // select variant_id
        // from genomic_annotations_new
        /* where variant_name='Y:59022489:?>A'*/

        $zigosity = $_GET["zygosity"];

        $stmt = $pdo->prepare("SELECT variant_id FROM genomic_annotations WHERE variant_name=? AND annotations ~* '^(?)");
        $stmt->execute([$_GET["variant_name"], join('|', $zigosity)]);
        break;

    case "annotation_and_zygosity":
        //select variant_id
        //from genomic_annotations_new
        //where annotations ~* '(Homozygous|Unknown);.*VARIANT_CLASS=DELETION'
        $zigosity = $_GET["zygosity"];

        $stmt = $pdo->prepare("SELECT variant_id FROM genomic_annotations WHERE annotations ~* '(?); ?=?");
        $stmt->execute([join('|', $zigosity)], $_GET["annotation_name"], $_GET["annotation_value"]);
        break;

    default:
        echo "Error: query type not recognized";
        return;

}

// In json format return both the panel number and the list of variants
echo "{ \"variants\" : ";
$variantList = "";
while ($row = $stmt->fetch()) {
    $variantList .= "\"$row[0]\",";
}
// drop the last comma
echo "[" . substr($variantList, 0, -1) . "]}";
?>