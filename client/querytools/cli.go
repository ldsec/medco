package querytoolsclient

import (
	"fmt"

	utilclient "github.com/ldsec/medco-connector/util/client"

	"github.com/sirupsen/logrus"
)

// ExecuteGetCohorts executes a get cohorts query and display its results
func ExecuteGetCohorts(token, username, password string, disableTLSCheck bool, resultFile string) error {
	accessToken, err := utilclient.RetrieveOrGetNewAccessToken(token, username, password, disableTLSCheck)
	if err != nil {
		err = fmt.Errorf("while retrieving access token: %s", err.Error())
		logrus.Error(err)
		return err
	}
	logrus.Debug("access token received")
	logrus.Tracef("token %s", accessToken)

	logrus.Debug("creating get cohorts request")
	getCohorts, err := NewGetCohorts(accessToken, disableTLSCheck)
	if err != nil {
		err = fmt.Errorf("while crafting new get cohorts request: %s", err.Error())
		logrus.Error(err)
		return err
	}

	logrus.Debug("executing get cohorts request")
	_, err = getCohorts.Execute()
	if err != nil {
		err = fmt.Errorf("cohorts request execution: %s", err.Error())
		logrus.Error(err)
		return err
	}

	return fmt.Errorf("NOT IMPLEMENTED")

}
