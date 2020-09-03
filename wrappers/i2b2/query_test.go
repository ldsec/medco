package i2b2

import (
	"testing"

	utilserver "github.com/ldsec/medco-connector/util/server"
)

func init() {
	utilserver.I2b2HiveURL = "http://localhost:8090/i2b2/services"
	utilserver.I2b2LoginDomain = "i2b2medcosrv0"
	utilserver.I2b2LoginProject = "MedCo"
	utilserver.I2b2LoginUser = "e2etest"
	utilserver.I2b2LoginPassword = "e2etest"
	utilserver.SetLogLevel("5")
}

// warning: all tests need the dev-local-3nodes medco deployment running locally, loaded with default data
func TestGetTermInfo(t *testing.T) {
	results, err := GetOntologyTermInfo("/SPHN/SPHNv2020.1/FophDiagnosis/")
	if err != nil {
		t.Error(err)
	}
	if results[0].Code != "A168" {
		t.Errorf("Expected basecode A168, got %s", results[0].Code)
	}
}

// test ontology search query
func TestGetOntologyChildrenRoot(t *testing.T) {

	results, err := GetOntologyChildren("/")
	if err != nil {
		t.Fail()
	}
	t.Log(*results[0])
}

func TestGetOntologyChildrenNode(t *testing.T) {

	results, err := GetOntologyChildren("/E2ETEST/e2etest/")
	if err != nil {
		t.Fail()
	}
	t.Log(*results[0].MedcoEncryption)
}

func TestExecutePsmQuery(t *testing.T) {

	patientCount, patientSetID, err := ExecutePsmQuery(
		"testQuery",
		[][]string{{`\\SENSITIVE_TAGGED\medco\tagged\09bc15e0d90046c102199f1b4d20eef9ee91b2ea3fd4608303775d000dd1248c\`}},
		[]bool{false},
	)
	if err != nil {
		t.Fail()
	}
	t.Log("count:"+patientCount, "set ID:"+patientSetID)
}

func TestGetPatientSet(t *testing.T) {

	patientIDs, patientDummyFlags, err := GetPatientSet("9")
	if err != nil {
		t.Fail()
	}
	t.Log(patientIDs)
	t.Log(patientDummyFlags)
}
