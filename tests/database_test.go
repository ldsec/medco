package tests

import (
	"github.com/ldsec/medco-connector/restapi/server/operations/genomic_annotations"
	utilserver "github.com/ldsec/medco-connector/util/server"
	"github.com/sirupsen/logrus"
	"testing"
)

var variant_name_get_values_value = "5238"
var variant_name_get_values_result = []string{"16:75238144:C>C", "6:52380882:G>G"}
var protein_change_get_values_value = "g32"
var protein_change_get_values_result = []string{"G325R", "G32E"}
var hugo_gene_symbol_get_values_value = "f22"
var hugo_gene_symbol_get_values_result = []string{"C1ORF226", "RNF222", "ZNF221"}

var variant_name_get_variants_result = []string{"qiaK9Vu2h49eGbtTPhRBOj_KdwZRpfdHNn8U8G8rry7rMeHSa9Ipooiog8fAdFGLrE_rCSTv2I9c4jwtebULQg=="}
var protein_change_get_variants_result = []string{"nrHv8spIGMZEGLn3GY5niuPD7U8z2E0FPcNcDJOCjCMoLLS86M5HE46PkpjttmY4rop5ugTojIjBsWBLuXmqMQ=="}
var hugo_gene_symbol_get_variants_result = []string{"uvrEA3sLdzG9x5vysj6SuGlWV2BCHVREAtxBT7GzE_3LKuczCzL3tfA2mRYnq3JrgTv8FZXKjk7-RcCq4PGppg==",
	"6mmVnBWNiAF3BaSDLEHihdpB4Atc68XTXakFJXDBoc9TEl_GXyQ9Bx-joN4g3izS11GSElZmNnzt0lAtni7p3w==",
	"MyPkqzW1cv1BqMhQe578veTRnQG3gPrifoQPR0sKKG-VO0wVC7dQ1M5qA9l1LCaS5IIAQFE7jJKi2vfOfJyRTw=="}

func init() {
	utilserver.DBMSHost = "localhost"
	utilserver.DBMSPort = 5432
	utilserver.DBName = "i2b2medcosrv0"
	utilserver.DBLoginUser = "postgres"
	utilserver.DBLoginPassword = "postgres"
	utilserver.SetLogLevel("5")
}

func TestDBConnection(t *testing.T) {
	var err error
	utilserver.DBConnection, err = utilserver.InitializeConnectionToDB(utilserver.DBMSHost, utilserver.DBMSPort, utilserver.DBName, utilserver.DBLoginUser, utilserver.DBLoginPassword)
	if err != nil {
		t.Fail()
	}

	err = utilserver.DBConnection.Ping()
	if err != nil {
		logrus.Error("Impossible to connect to DB " + err.Error())
		t.Fail()
	}
}

func TestGenomicAnnotationsGetValues(t *testing.T) {

	//testing variant_name type get values
	testGenomicAnnotationsGetValues("variant_name", variant_name_get_values_value, variant_name_get_values_result, t)
	//testing protein_change type get values
	testGenomicAnnotationsGetValues("protein_change", protein_change_get_values_value, protein_change_get_values_result, t)
	//testing hugo_gene_symbol type get values
	testGenomicAnnotationsGetValues("hugo_gene_symbol", hugo_gene_symbol_get_values_value, hugo_gene_symbol_get_values_result, t)
	//testing get values with empty result
	testGenomicAnnotationsGetValues("hugo_gene_symbol", "aaa", nil, t)

}

func TestGenomicAnnotationsGetVariants(t *testing.T) {

	//testing variant_name type get variants
	testGenomicAnnotationsGetVariants("variant_name", variant_name_get_values_result[0], nil, variant_name_get_variants_result, t)
	//testing protein_change type get variants
	testGenomicAnnotationsGetVariants("protein_change", protein_change_get_values_result[0], nil, protein_change_get_variants_result, t)
	//testing hugo_gene_symbol type get variants
	testGenomicAnnotationsGetVariants("hugo_gene_symbol", hugo_gene_symbol_get_values_result[0], nil, hugo_gene_symbol_get_variants_result, t)

	//testing get variants with different zygosity parameters
	testGenomicAnnotationsGetVariants("hugo_gene_symbol", hugo_gene_symbol_get_values_result[0], []string{"heterozygous"}, hugo_gene_symbol_get_variants_result[1:3], t)
	testGenomicAnnotationsGetVariants("hugo_gene_symbol", hugo_gene_symbol_get_values_result[0], []string{"homozygous"}, nil, t)
	testGenomicAnnotationsGetVariants("hugo_gene_symbol", hugo_gene_symbol_get_values_result[0], []string{"unknown"}, hugo_gene_symbol_get_variants_result[0:1], t)
	testGenomicAnnotationsGetVariants("hugo_gene_symbol", hugo_gene_symbol_get_values_result[0], []string{"heterozygous", "homozygous"}, hugo_gene_symbol_get_variants_result[1:3], t)
	testGenomicAnnotationsGetVariants("hugo_gene_symbol", hugo_gene_symbol_get_values_result[0], []string{"heterozygous", "unknown"}, hugo_gene_symbol_get_variants_result, t)
	testGenomicAnnotationsGetVariants("hugo_gene_symbol", hugo_gene_symbol_get_values_result[0], []string{"homozygous", "unknown"}, hugo_gene_symbol_get_variants_result[0:1], t)
	testGenomicAnnotationsGetVariants("hugo_gene_symbol", hugo_gene_symbol_get_values_result[0], []string{"heterozygous", "homozygous", "unknown"}, hugo_gene_symbol_get_variants_result, t)

}

func testGenomicAnnotationsGetValues(query_type string, query_value string, query_result []string, t *testing.T) {

	var annotations []string
	var annotation string
	params := genomic_annotations.NewGetValuesParams()

	var err error
	utilserver.DBConnection, err = utilserver.InitializeConnectionToDB(utilserver.DBMSHost, utilserver.DBMSPort, utilserver.DBName, utilserver.DBLoginUser, utilserver.DBLoginPassword)
	if err != nil {
		t.Fail()
	}

	err = utilserver.DBConnection.Ping()
	if err != nil {
		logrus.Error("Impossible to connect to DB " + err.Error())
		t.Fail()
	}

	params.Annotation = query_type
	params.Value = query_value

	query := "SELECT annotation_value FROM genomic_annotations." + params.Annotation + " WHERE annotation_value ~* $1 ORDER BY annotation_value LIMIT $2"
	rows, err := utilserver.DBConnection.Query(query, params.Value, *params.Limit)
	if err != nil {
		logrus.Error("Query execution error " + err.Error())
		t.Fail()
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&annotation)
		if err != nil {
			logrus.Error("Query result reading error " + err.Error())
			t.Fail()
		}
		annotations = append(annotations, annotation)
	}

	if !areEqual(annotations, query_result) {
		logrus.Error("Wrong " + query_type + " query result")
		t.Fail()
	}

}

func testGenomicAnnotationsGetVariants(query_type string, query_value string, zygosity []string, query_result []string, t *testing.T) {

	TestDBConnection(t)

	var variants []string
	var variant string
	var err error

	params := genomic_annotations.NewGetVariantsParams()

	params.Annotation = query_type
	params.Value = query_value
	params.Zygosity = zygosity

	zygosityStr := ""
	if len(params.Zygosity) > 0 {
		zygosityStr = params.Zygosity[0]

		for i := 1; i < len(params.Zygosity); i++ {
			zygosityStr += "|" + params.Zygosity[i]
		}
	}

	query := "SELECT variant_id FROM genomic_annotations.genomic_annotations WHERE " + params.Annotation + " = $1 AND annotations ~* $2 ORDER BY variant_id"
	rows, err := utilserver.DBConnection.Query(query, params.Value, zygosityStr)
	if err != nil {
		logrus.Error("Query execution error " + err.Error())
		t.Fail()
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&variant)
		if err != nil {
			logrus.Error("Query result reading error " + err.Error())
			t.Fail()
		}
		variants = append(variants, variant)
	}

	if !areEqual(variants, query_result) {
		logrus.Error("Wrong " + query_type + " query result")
		t.Fail()
	}

}

func areEqual(slice1, slice2 []string) bool {

	if len(slice1) != len(slice2) {
		return false
	}

	for i, element := range slice1 {
		if element != slice2[i] {
			return false
		}
	}

	return true
}
