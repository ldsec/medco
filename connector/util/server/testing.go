package utilserver

import (
	"testing"

	"github.com/sirupsen/logrus"
)

// TestDBConnection tests the connection to medco connector postgresql data base
func TestDBConnection(t *testing.T) {

	var err error
	DBConnection, err = InitializeConnectionToDB(DBHost, DBPort, DBName, DBLoginUser, DBLoginPassword)
	if err != nil {
		t.Fail()
	}

	err = DBConnection.Ping()
	if err != nil {
		logrus.Error("Impossible to connect to DB " + err.Error())
		t.Fail()
	}
}

// TestI2B2DBConnection tests the connection to I2B2 data base
func TestI2B2DBConnection(t *testing.T) {

	var err error
	I2B2DBConnection, err = InitializeConnectionToDB(I2B2DBHost, I2B2DBPort, I2B2DBName, I2B2DBLoginUser, I2B2DBLoginPassword)
	if err != nil {
		t.Fail()
	}

	err = I2B2DBConnection.Ping()
	if err != nil {
		logrus.Error("Impossible to connect to DB " + err.Error())
		t.Fail()
	}
}

// SetForTesting changes utility variables for unit tests. It is to be called in init() of test packages.
func SetForTesting() {
	DBHost = "localhost"
	DBPort = 5432
	DBName = "medcoconnectorsrv0"
	DBLoginUser = "medcoconnector"
	DBLoginPassword = "medcoconnector"

	I2B2DBHost = "localhost"
	I2B2DBPort = 5432
	I2B2DBName = "i2b2medcosrv0"
	I2B2DBLoginUser = "i2b2"
	I2B2DBLoginPassword = "i2b2"

	SetLogLevel("5")
}
