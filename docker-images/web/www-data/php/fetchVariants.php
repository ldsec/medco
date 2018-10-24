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
$zigosity = "^(".join('|', $_GET["zygosity"]).");"; // array of zygosity options, put them in the query separated by | (or)
switch($_GET["query_type"]){
    case "gene_and_zygosity":
        // case insensitive regex query
        // select variant_id
        // from genomic_annotations
        // where annotations ~* '^CTBP2P1;(Homozygous|Unknown);'

        //fetchVariants.php?query_type=gene_and_zygosity&gene_value=ACACB&zygosity[]=Heterozygous

        $stmt = $pdo->prepare("SELECT variant_id FROM genomic_annotations.genomic_annotations WHERE hugo_gene_symbol=? AND annotations ~* ?");
        $stmt->bindValue(1, $_GET["gene_value"], PDO::PARAM_STR);
        $stmt->bindValue(2, $zigosity, PDO::PARAM_STR);
        $stmt->execute();
        break;

    case "protein_position_and_zygosity":

        //fetchVariants.php?query_type=protein_position_and_zygosity&protein_change_value=A443V&zygosity[]=Heterozygous
        
        $stmt = $pdo->prepare("SELECT variant_id FROM genomic_annotations.genomic_annotations WHERE protein_change=? AND annotations ~* ?");
        $stmt->bindValue(1, $_GET["protein_change_value"], PDO::PARAM_STR);
        $stmt->bindValue(2, $zigosity, PDO::PARAM_STR);
        $stmt->execute();
        break;

    case "variantName_and_zygosity":
        // select variant_id
        // from genomic_annotations
        /* where variant_name='Y:59022489:?>A'*/

        //fetchVariants.php?query_type=variantName_and_zygosity&variant_name=12:109613959:C>C&zygosity[]=Heterozygous

        $stmt = $pdo->prepare("SELECT variant_id FROM genomic_annotations.genomic_annotations WHERE variant_name=? AND annotations ~* ?");
        $stmt->bindValue(1, $_GET["variant_name"], PDO::PARAM_STR);
        $stmt->bindValue(2, $zigosity, PDO::PARAM_STR);
        $stmt->execute();
        break;

    case "annotation_and_zygosity":
        //select variant_id
        //from genomic_annotations
        //where annotations ~* '(Homozygous|Unknown);.*VARIANT_CLASS=DELETION'
        $zigosity = $zigosity.".*" . $_GET["annotation_name"] . "=" . $_GET["annotation_value"];
        echo $zigosity."\n"

        //fetchVariants.php?query_type=annotation_and_zygosity&annotation_name=Variant%20Type&annotation_value=SNP&zygosity[]=Heterozygous

        $stmt = $pdo->prepare("SELECT variant_id FROM genomic_annotations.genomic_annotations WHERE annotations ~* ?");
        $stmt->bindValue(1, $zigosity, PDO::PARAM_STR);
        $stmt->execute();
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