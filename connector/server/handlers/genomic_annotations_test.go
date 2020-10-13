package handlers

import (
	"github.com/ldsec/medco/connector/restapi/server/operations/genomic_annotations"
	utilserver "github.com/ldsec/medco/connector/util/server"
	"github.com/sirupsen/logrus"
	"testing"
)

var variantNameGetValuesValue = "5238"
var variantNameGetValuesResult = []string{"16:75238144:C>C", "6:52380882:G>G"}
var proteinChangeGetValuesValue = "g32"
var proteinChangeGetValuesResult = []string{"G325R", "G32E"}
var proteinChangeGetValuesValue2 = "7cfs*"
var proteinChangeGetValuesResult2 = []string{"S137Cfs*28"}
var hugoGeneSymbolGetValuesValue = "tr5"
var hugoGeneSymbolGetValuesResult = []string{"HTR5A"}

var variantNameGetVariantsResult = []string{"-4530899676219565056"}
var proteinChangeGetVariantsResult = []string{"-2429151887266669568"}
var hugoGeneSymbolGetVariantsResult = []string{"-7039476204566471680", "-7039476580443220992", "-7039476780159200256"}

func init() {
	utilserver.MedcoDBHost = "localhost"
	utilserver.MedcoDBPort = 5432
	utilserver.MedcoDBName = "medcoconnectorsrv0"
	utilserver.MedcoDBLoginUser = "medcoconnector"
	utilserver.MedcoDBLoginPassword = "medcoconnector"
	utilserver.SetLogLevel("5")
}

func TestDBConnection(t *testing.T) {

	var err error
	utilserver.MedcoDBConnection, err = utilserver.InitializeConnectionToDB(utilserver.MedcoDBHost, utilserver.MedcoDBPort, utilserver.MedcoDBName, utilserver.MedcoDBLoginUser, utilserver.MedcoDBLoginPassword)
	if err != nil {
		t.Fail()
	}

	err = utilserver.MedcoDBConnection.Ping()
	if err != nil {
		logrus.Error("Impossible to connect to DB " + err.Error())
		t.Fail()
	}
}

// warning: this test needs the dev-local-3nodes medco deployment running locally, loaded with default data
func TestGenomicAnnotationsGetValues(t *testing.T) {

	//testing variant_name type get values
	testGenomicAnnotationsGetValues("variant_name", variantNameGetValuesValue, variantNameGetValuesResult, t)
	//testing protein_change type get values
	testGenomicAnnotationsGetValues("protein_change", proteinChangeGetValuesValue, proteinChangeGetValuesResult, t)
	//testing protein_change type get values with value containing *
	testGenomicAnnotationsGetValues("protein_change", proteinChangeGetValuesValue2, proteinChangeGetValuesResult2, t)
	//testing hugo_gene_symbol type get values
	testGenomicAnnotationsGetValues("hugo_gene_symbol", hugoGeneSymbolGetValuesValue, hugoGeneSymbolGetValuesResult, t)
	//testing get values with empty result
	testGenomicAnnotationsGetValues("hugo_gene_symbol", "aaa", nil, t)

}

// warning: this test needs the dev-local-3nodes medco deployment running locally, loaded with default data
func TestGenomicAnnotationsGetVariants(t *testing.T) {

	//testing variant_name type get variants
	testGenomicAnnotationsGetVariants("variant_name", variantNameGetValuesResult[0], nil, variantNameGetVariantsResult, t)
	//testing protein_change type get variants
	testGenomicAnnotationsGetVariants("protein_change", proteinChangeGetValuesResult[0], nil, proteinChangeGetVariantsResult, t)
	//testing hugo_gene_symbol type get variants
	testGenomicAnnotationsGetVariants("hugo_gene_symbol", hugoGeneSymbolGetValuesResult[0], nil, hugoGeneSymbolGetVariantsResult, t)

	//testing get variants with different zygosity parameters
	testGenomicAnnotationsGetVariants("hugo_gene_symbol", hugoGeneSymbolGetValuesResult[0], []string{"heterozygous"}, hugoGeneSymbolGetVariantsResult[0:2], t)
	testGenomicAnnotationsGetVariants("hugo_gene_symbol", hugoGeneSymbolGetValuesResult[0], []string{"homozygous"}, nil, t)
	testGenomicAnnotationsGetVariants("hugo_gene_symbol", hugoGeneSymbolGetValuesResult[0], []string{"unknown"}, hugoGeneSymbolGetVariantsResult[2:3], t)
	testGenomicAnnotationsGetVariants("hugo_gene_symbol", hugoGeneSymbolGetValuesResult[0], []string{"heterozygous", "homozygous"}, hugoGeneSymbolGetVariantsResult[0:2], t)
	testGenomicAnnotationsGetVariants("hugo_gene_symbol", hugoGeneSymbolGetValuesResult[0], []string{"heterozygous", "unknown"}, hugoGeneSymbolGetVariantsResult, t)
	testGenomicAnnotationsGetVariants("hugo_gene_symbol", hugoGeneSymbolGetValuesResult[0], []string{"homozygous", "unknown"}, hugoGeneSymbolGetVariantsResult[2:3], t)
	testGenomicAnnotationsGetVariants("hugo_gene_symbol", hugoGeneSymbolGetValuesResult[0], []string{"heterozygous", "homozygous", "unknown"}, hugoGeneSymbolGetVariantsResult, t)

}

func testGenomicAnnotationsGetValues(query_type string, query_value string, query_result []string, t *testing.T) {

	TestDBConnection(t)

	var annotations []string
	var annotation string
	params := genomic_annotations.NewGetValuesParams()
	var err error

	params.Annotation = query_type
	params.Value = query_value

	query, _ := buildGetValuesQuery(params)
	rows, err := utilserver.MedcoDBConnection.Query(query, params.Annotation, params.Value, *params.Limit)
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

	query, _ := buildGetVariantsQuery(params)
	rows, err := utilserver.MedcoDBConnection.Query(query, params.Annotation, params.Value, zygosityStr, false)
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
