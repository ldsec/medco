<?php
header('Access-Control-Allow-Origin: ' getenv('CORS_ALLOW_ORIGIN')); 
header('Access-Control-Allow-Credentials: true'); 
header('Access-Control-Allow-Headers: origin, content-type, accept, authorization'); 
header('Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS, HEAD'); 

// query_type define is we are searching by:
// - old is the old version that uses the first version of the schema
// - Gene + Zygosity
// - ...

$query = "";
switch($_GET["query_type"]){
    case "gene_and_zygosity":
        // case insensitive regex query
        // select variant_id
        // from genomic_annotations_new
        // where annotations ~* '^CTBP2P1;(Homozygous|Unknown);'

        $zigosity = $_GET["zygosity"]; // array of zygosity options, put them in the query separated by | (or)

        $query =
            "SELECT variant_id " .
            "FROM genomic_annotations " .
            "WHERE hugo_gene_symbol='" . $_GET["gene_value"] . "' AND annotations ~* '^(" . join('|', $zigosity) .");'";

        break;

    case "protein_position_and_zygosity":
        $zigosity = $_GET["zygosity"]; // array of zygosity options, put them in the query separated by | (or)

        $query =
            "SELECT variant_id " .
            "FROM genomic_annotations " .
            "WHERE 	protein_change='" . $_GET["protein_change_value"] . "' AND annotations ~* '^(" . join('|', $zigosity) .");'";
        break;

    case "variantName_and_zygosity":
        // select variant_id
        // from genomic_annotations_new
        /* where variant_name='Y:59022489:?>A'*/

        $zigosity = $_GET["zygosity"];

        $query = "SELECT variant_id " .
        "FROM genomic_annotations " .
        "WHERE variant_name='" . $_GET["variant_name"] ."' AND " .
        "annotations ~* '" . "(" . join('|', $zigosity) .")'";
        break;

    case "annotation_and_zygosity":
        //select variant_id
        //from genomic_annotations_new
        //where annotations ~* '(Homozygous|Unknown);.*VARIANT_CLASS=DELETION'
        $zigosity = $_GET["zygosity"];
        $query = "SELECT variant_id " .
            "FROM genomic_annotations " .
            "WHERE annotations ~* '" . "(" . join('|', $zigosity) .");.*" . $_GET["annotation_name"] . "=" . $_GET["annotation_value"] . "'";

        break;

    default:
        echo "Error: query type not recognized";
        return;

}

//echo $query;
//return;

include 'sqlConnection.php';

$result = pg_query($conn, $query);
if (!$result) {
    echo "An error occurred while querying the database.\n";
    exit;
}

// In json format return both the panel number and the list of variants
echo "{ \"variants\" : ";
$variantList = "";
while ($row = pg_fetch_row($result)) {
    $variantList .= "\"$row[0]\",";
}
// drop the last comma
echo "[" . substr($variantList, 0, -1) . "]}";
