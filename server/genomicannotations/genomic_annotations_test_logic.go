package genomicannotations

import (
	"encoding/json"
	"github.com/ldsec/medco-connector/restapi/server/operations/genomic_annotations"
	utilserver "github.com/ldsec/medco-connector/util/server"
	"github.com/sirupsen/logrus"
)

type testGetValuesParameters struct {
	queryType   string
	queryValue  string
	queryResult []string
}

type testGetVariantsParameters struct {
	queryType   string
	queryValue  string
	zygosity    []string
	queryResult []string
}

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

var getValuesParams = []testGetValuesParameters{
	//testing variant_name type get values
	testGetValuesParameters{
		queryType:   "variant_name",
		queryValue:  variantNameGetValuesValue,
		queryResult: variantNameGetValuesResult,
	},
	//testing protein_change type get values
	testGetValuesParameters{
		queryType:   "protein_change",
		queryValue:  proteinChangeGetValuesValue,
		queryResult: proteinChangeGetValuesResult,
	},
	//testing protein_change type get values with value containing *
	testGetValuesParameters{
		queryType:   "protein_change",
		queryValue:  proteinChangeGetValuesValue2,
		queryResult: proteinChangeGetValuesResult2,
	},
	//testing hugo_gene_symbol type get values
	testGetValuesParameters{
		queryType:   "hugo_gene_symbol",
		queryValue:  hugoGeneSymbolGetValuesValue,
		queryResult: hugoGeneSymbolGetValuesResult,
	},
	//testing get values with empty result
	testGetValuesParameters{
		queryType:   "hugo_gene_symbol",
		queryValue:  "aaa",
		queryResult: nil,
	},
}

var getVariantsParams = []testGetVariantsParameters{
	//testing variant_name type get variants
	testGetVariantsParameters{
		queryType:   "variant_name",
		queryValue:  variantNameGetValuesResult[0],
		zygosity:    nil,
		queryResult: variantNameGetVariantsResult,
	},
	//testing protein_change type get variants
	testGetVariantsParameters{
		queryType:   "protein_change",
		queryValue:  proteinChangeGetValuesResult[0],
		zygosity:    nil,
		queryResult: proteinChangeGetVariantsResult,
	},
	//testing hugo_gene_symbol type get variants
	testGetVariantsParameters{
		queryType:   "hugo_gene_symbol",
		queryValue:  hugoGeneSymbolGetValuesResult[0],
		zygosity:    nil,
		queryResult: hugoGeneSymbolGetVariantsResult,
	},
	//testing get variants with different zygosity parameters
	testGetVariantsParameters{
		queryType:   "hugo_gene_symbol",
		queryValue:  hugoGeneSymbolGetValuesResult[0],
		zygosity:    []string{"heterozygous"},
		queryResult: hugoGeneSymbolGetVariantsResult[0:2],
	},
	testGetVariantsParameters{
		queryType:   "hugo_gene_symbol",
		queryValue:  hugoGeneSymbolGetValuesResult[0],
		zygosity:    []string{"homozygous"},
		queryResult: nil,
	},
	testGetVariantsParameters{
		queryType:   "hugo_gene_symbol",
		queryValue:  hugoGeneSymbolGetValuesResult[0],
		zygosity:    []string{"unknown"},
		queryResult: hugoGeneSymbolGetVariantsResult[2:3],
	},
	testGetVariantsParameters{
		queryType:   "hugo_gene_symbol",
		queryValue:  hugoGeneSymbolGetValuesResult[0],
		zygosity:    []string{"heterozygous", "homozygous"},
		queryResult: hugoGeneSymbolGetVariantsResult[0:2],
	},
	testGetVariantsParameters{
		queryType:   "hugo_gene_symbol",
		queryValue:  hugoGeneSymbolGetValuesResult[0],
		zygosity:    []string{"heterozygous", "unknown"},
		queryResult: hugoGeneSymbolGetVariantsResult,
	},
	testGetVariantsParameters{
		queryType:   "hugo_gene_symbol",
		queryValue:  hugoGeneSymbolGetValuesResult[0],
		zygosity:    []string{"homozygous", "unknown"},
		queryResult: hugoGeneSymbolGetVariantsResult[2:3],
	},
	testGetVariantsParameters{
		queryType:   "hugo_gene_symbol",
		queryValue:  hugoGeneSymbolGetValuesResult[0],
		zygosity:    []string{"heterozygous", "homozygous", "unknown"},
		queryResult: hugoGeneSymbolGetVariantsResult,
	},
}

// TestGenomicAnnotations runs the genomic-annotations tests.
func TestGenomicAnnotations() (testPassed bool) {

	err := utilserver.GaDBConnection.Ping()
	if err != nil {
		logrus.Error("Impossible to connect to genomic annotations DB: " + err.Error())
		return false
	}

	for _, testParams := range getValuesParams {
		if !testGetValues(testParams) {
			log := "test failed: "
			text, err := json.Marshal(testParams)
			if err == nil {
				log += string(text)
			} else {
				log += err.Error()
			}
			logrus.Warn(log)
			return false
		}
	}

	for _, testParams := range getVariantsParams {
		if !testGetVariants(testParams) {
			log := "test failed: "
			text, err := json.Marshal(testParams)
			if err == nil {
				log += string(text)
			} else {
				log += err.Error()
			}
			logrus.Warn(log)
			return false
		}
	}

	return true
}

func testDBConnection() (testPassed bool) {
	var err error
	utilserver.GaDBConnection, err = utilserver.InitializeConnectionToDB(utilserver.GaDBHost, utilserver.GaDBPort, utilserver.GaDBName, utilserver.GaDBLoginUser, utilserver.GaDBLoginPassword)
	if err != nil {
		return false
	}

	err = utilserver.GaDBConnection.Ping()
	if err != nil {
		logrus.Error("Impossible to connect to genomic annotations DB: " + err.Error())
		return false
	}

	return true
}

func testGetValues(testParams testGetValuesParameters) (testPassed bool) {

	var annotations []string
	var annotation string
	params := genomic_annotations.NewGetValuesParams()
	var err error

	params.Annotation = testParams.queryType
	params.Value = testParams.queryValue

	query, _ := BuildGetValuesQuery(params)
	rows, err := utilserver.GaDBConnection.Query(query, params.Annotation, params.Value, *params.Limit, true)
	if err != nil {
		logrus.Error("Query execution error " + err.Error())
		return false
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&annotation)
		if err != nil {
			logrus.Error("Query result reading error " + err.Error())
			return false
		}
		annotations = append(annotations, annotation)
	}

	if !areEqual(annotations, testParams.queryResult) {
		logrus.Error("Wrong " + testParams.queryType + " query result")
		return false
	}

	return true

}

func testGetVariants(testParams testGetVariantsParameters) (testPassed bool) {

	var variants []string
	var variant string
	var err error

	params := genomic_annotations.NewGetVariantsParams()

	params.Annotation = testParams.queryType
	params.Value = testParams.queryValue
	params.Zygosity = testParams.zygosity

	zygosityStr := ""
	if len(params.Zygosity) > 0 {
		zygosityStr = params.Zygosity[0]

		for i := 1; i < len(params.Zygosity); i++ {
			zygosityStr += "|" + params.Zygosity[i]
		}
	}

	query, _ := BuildGetVariantsQuery(params)
	rows, err := utilserver.GaDBConnection.Query(query, params.Annotation, params.Value, zygosityStr, false, true)
	if err != nil {
		logrus.Error("Query execution error " + err.Error())
		return false
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&variant)
		if err != nil {
			logrus.Error("Query result reading error " + err.Error())
			return false
		}
		variants = append(variants, variant)
	}

	if !areEqual(variants, testParams.queryResult) {
		logrus.Error("Wrong " + testParams.queryType + " query result")
		return false
	}

	return true

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
