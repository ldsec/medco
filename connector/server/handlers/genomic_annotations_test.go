//go:build integration_test

package handlers

import (
	"sort"
	"testing"

	"github.com/ldsec/medco/connector/restapi/server/operations/genomic_annotations"
	utilserver "github.com/ldsec/medco/connector/util/server"
	"github.com/sirupsen/logrus"
)

var variantNameGetValuesValue = "78"
var variantNameGetValuesResult = []string{"6:35786830:GGGACC>TAATAC", "7:78098298:TCTTTA>AACGGA", "8:57873164:GCTGTG>GGCT"}
var proteinChangeGetValuesValue = "y"
var proteinChangeGetValuesResult = []string{"N232Y", "H277Y", "Y1062H"}
var proteinChangeGetValuesValue2 = "l61sfs*"
var proteinChangeGetValuesResult2 = []string{"L61Sfs*54"}
var hugoGeneSymbolGetValuesValue = "nav"
var hugoGeneSymbolGetValuesResult = []string{"NAV3"}

var variantNameGetVariantsResult = []string{"-7455563962931223533"}
var proteinChangeGetVariantsResult = []string{"-2823470849823937376"}
var hugoGeneSymbolGetVariantsResult = []string{"-7121901993980174104", "-4898572880864589696", "-6271408487767448598"}

func init() {
	utilserver.SetForTesting()
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

func testGenomicAnnotationsGetValues(queryType string, queryValue string, queryResult []string, t *testing.T) {

	utilserver.TestDBConnection(t)

	var annotations []string
	var annotation string
	params := genomic_annotations.NewGetValuesParams()
	var err error

	params.Annotation = queryType
	params.Value = queryValue

	query, _ := buildGetValuesQuery(params)
	rows, err := utilserver.DBConnection.Query(query, params.Annotation, params.Value, *params.Limit)
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

	if !areEqual(annotations, queryResult) {
		logrus.Error("Wrong " + queryType + " query result")
		t.Fail()
	}

}

func testGenomicAnnotationsGetVariants(queryType string, queryValue string, zygosity []string, queryResult []string, t *testing.T) {

	utilserver.TestDBConnection(t)

	var variants []string
	var variant string
	var err error

	params := genomic_annotations.NewGetVariantsParams()

	params.Annotation = queryType
	params.Value = queryValue
	params.Zygosity = zygosity

	zygosityStr := ""
	if len(params.Zygosity) > 0 {
		zygosityStr = params.Zygosity[0]

		for i := 1; i < len(params.Zygosity); i++ {
			zygosityStr += "|" + params.Zygosity[i]
		}
	}

	query, _ := buildGetVariantsQuery(params)
	rows, err := utilserver.DBConnection.Query(query, params.Annotation, params.Value, zygosityStr, false)
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

	if !areEqual(variants, queryResult) {
		logrus.Error("Wrong " + queryType + " query result")
		t.Fail()
	}

}

func areEqual(slice1, slice2 []string) bool {

	if len(slice1) != len(slice2) {
		return false
	}

	slice1C := make([]string, len(slice1))
	slice2C := make([]string, len(slice2))
	copy(slice1C, slice1)
	copy(slice2C, slice2)
	sort.Strings(slice1C)
	sort.Strings(slice2C)

	for i, element := range slice1C {
		if element != slice2C[i] {
			return false
		}
	}

	return true
}
