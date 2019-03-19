package i2b2

import (
	"os"
	"testing"
)

func TestGetOntologyChildren(t *testing.T) {
	os.Setenv("I2B2_ONT_URL", "http://localhost:8090/i2b2/services/OntologyService")
	os.Setenv("I2B2_LOGIN_DOMAIN", "i2b2medcosrv0")
	os.Setenv("I2B2_LOGIN_PROJECT", "MedCo")
	os.Setenv("I2B2_LOGIN_USER", "e2etest")
	os.Setenv("I2B2_LOGIN_PASSWORD", "e2etest")

	results, err := GetOntologyChildren("/")
	if err != nil {
		t.Fail()
	}
	t.Log(*results[0])
}
