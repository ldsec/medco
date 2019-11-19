package tests

import (
	"github.com/ldsec/medco-connector/restapi/server/operations/genomic_annotations"
	utilserver "github.com/ldsec/medco-connector/util/server"
	"github.com/sirupsen/logrus"
	"testing"
)

var variant_name_query_value = "5238"
var variant_name_query_result = []string{"6:52380882:G>G", "16:75238144:C>C"}
var protein_change_query_value = "g32"
var protein_change_query_result = []string{"G325R", "G32E"}
var hugo_gene_symbol_query_value = "f22"
var hugo_gene_symbol_query_result = []string{"ZNF221", "C1ORF226", "RNF222"}

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

	//testing variant_name type annotations query
	testGenomicAnnotationsGetValues("variant_name", variant_name_query_value, variant_name_query_result, t)

	//testing protein_change type annotations query
	testGenomicAnnotationsGetValues("protein_change", protein_change_query_value, protein_change_query_result, t)

	//testing protein_change type annotations query
	testGenomicAnnotationsGetValues("hugo_gene_symbol", hugo_gene_symbol_query_value, hugo_gene_symbol_query_result, t)

	//testing query with empty result
	testGenomicAnnotationsGetValues("hugo_gene_symbol", "aaa", nil, t)

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

	query := "SELECT annotation_value FROM genomic_annotations." + params.Annotation + " WHERE annotation_value ~* $1 LIMIT $2"
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
