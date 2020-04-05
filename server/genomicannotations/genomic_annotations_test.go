package genomicannotations

import (
	utilserver "github.com/ldsec/medco-connector/util/server"
	"testing"
)

func init() {
	utilserver.GaDBHost = "localhost"
	utilserver.GaDBPort = 5432
	utilserver.GaDBName = "gamedcosrv0"
	utilserver.GaDBLoginUser = "genomicannotations"
	utilserver.GaDBLoginPassword = "genomicannotations"
	utilserver.SetLogLevel("5")
}

func TestDBConnection(t *testing.T) {
	if !testDBConnection() {
		t.Failed()
	}
}

// warning: this test needs the dev-local-3nodes medco deployment running locally, loaded with default data
func TestGetValues(t *testing.T) {

	TestDBConnection(t)

	for _, testParams := range getValuesParams {
		if !testGetValues(testParams) {
			t.Fail()
		}
	}

}

// warning: this test needs the dev-local-3nodes medco deployment running locally, loaded with default data
func TestGetVariants(t *testing.T) {

	TestDBConnection(t)

	for _, testParams := range getVariantsParams {
		if !testGetVariants(testParams) {
			t.Fail()
		}
	}

}
