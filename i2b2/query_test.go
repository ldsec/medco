package i2b2

import (
	"os"
	"testing"
)

func setupEnv() {
	os.Setenv("I2B2_ONT_URL", "http://localhost:8090/i2b2/services/OntologyService")
	os.Setenv("I2B2_LOGIN_DOMAIN", "i2b2medcosrv0")
	os.Setenv("I2B2_LOGIN_PROJECT", "MedCo")
	os.Setenv("I2B2_LOGIN_USER", "e2etest")
	os.Setenv("I2B2_LOGIN_PASSWORD", "e2etest")
}

// test ontology search query
// warning: needs the dev-local-3nodes medco deployment running locally
func TestGetOntologyChildrenRoot(t *testing.T) {
	setupEnv()

	results, err := GetOntologyChildren("/")
	if err != nil {
		t.Fail()
	}
	t.Log(*results[0])
}

func TestGetOntologyChildrenNode(t *testing.T) {
	setupEnv()

	results, err := GetOntologyChildren("/E2ETEST/e2etest/")
	if err != nil {
		t.Fail()
	}
	t.Log(*results[0].MedcoEncryption)
}
