package querytoolsclient

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	utilclient "github.com/ldsec/medco/connector/util/client"

	"github.com/sirupsen/logrus"
)

// ExecuteGetCohorts executes a get cohorts query and displays its results
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
	cohorts, err := getCohorts.Execute()
	if err != nil {
		err = fmt.Errorf("cohorts request execution: %s", err.Error())
		logrus.Error(err)
		return err
	}

	resultCSV, err := utilclient.NewCSV(resultFile)
	if err != nil {
		err = fmt.Errorf("cohorts request writing results: %s", err.Error())
		logrus.Error(err)
		return err
	}
	logrus.Debug("Writing headers")
	resultCSV.Write([]string{"node_index", "cohort_name", "cohort_id", "query_id", "creation_date", "update_date"})

	for nodeIndex, nodeResult := range cohorts {
		for _, cohortInfo := range nodeResult {
			logrus.Debugf("Writing result %d", nodeIndex)
			err = resultCSV.Write([]string{
				strconv.Itoa(nodeIndex),
				cohortInfo.CohortName,
				strconv.Itoa(cohortInfo.CohortID),
				strconv.Itoa(cohortInfo.QueryID),
				cohortInfo.CreationDate.Format(time.RFC3339),
				cohortInfo.UpdateDate.Format(time.RFC3339),
			})
			if err != nil {
				err = fmt.Errorf("cohorts request writing results: %s", err.Error())
				logrus.Error(err)
				return err
			}
		}
	}

	logrus.Debug("Flushing result file")
	err = resultCSV.Flush()
	if err != nil {
		err = fmt.Errorf("cohorts request flushing result file: %s", err.Error())
		logrus.Error(err)
		return err
	}

	logrus.Debug("Closing result file")
	err = resultCSV.Close()
	if err != nil {
		err = fmt.Errorf("cohorts request closing file: %s", err.Error())
		logrus.Error(err)
		return err
	}

	return nil
}

// ExecutePostCohorts executes a post cohorts query
func ExecutePostCohorts(token, username, password, cohortName, patientSetIDsString string, disableTLSCheck bool) error {

	accessToken, err := utilclient.RetrieveOrGetNewAccessToken(token, username, password, disableTLSCheck)
	if err != nil {
		err = fmt.Errorf("while retrieving access token: %s", err.Error())
		logrus.Error(err)
		return err
	}
	logrus.Debug("access token received")
	logrus.Tracef("token %s", accessToken)

	logrus.Debug("creating post cohorts request")
	logrus.Tracef("patient set IDs %v , cohort name %s", patientSetIDsString, cohortName)

	patientSetIDs := make([]int, 0)
	for _, setID := range strings.Split(patientSetIDsString, ",") {
		id, err := strconv.Atoi(strings.TrimSpace(setID))
		if err != nil {
			err = fmt.Errorf("while parsing int from string %s in parameters: %s", setID, err.Error())
			logrus.Error(err)
			return err
		}
		patientSetIDs = append(patientSetIDs, id)
	}
	postCohorts, err := NewPostCohorts(accessToken, patientSetIDs, cohortName, disableTLSCheck)
	if err != nil {
		err = fmt.Errorf("while crafting new post cohorts request: %s", err.Error())
		logrus.Error(err)
		return err
	}

	err = postCohorts.Execute()
	if err != nil {
		err = fmt.Errorf("cohorts request execution: %s", err.Error())
		logrus.Error(err)
		return err
	}
	return nil
}

// ExecutePutCohorts executes a put cohorts query
func ExecutePutCohorts(token, username, password, cohortName, patientSetIDsString string, disableTLSCheck bool) error {

	accessToken, err := utilclient.RetrieveOrGetNewAccessToken(token, username, password, disableTLSCheck)
	if err != nil {
		err = fmt.Errorf("while retrieving access token: %s", err.Error())
		logrus.Error(err)
		return err
	}
	logrus.Debug("access token received")
	logrus.Tracef("token %s", accessToken)

	logrus.Debug("creating post cohorts request")
	logrus.Tracef("patient set IDs %v , cohort name %s", patientSetIDsString, cohortName)

	patientSetIDs := make([]int, 0)
	for _, setID := range strings.Split(patientSetIDsString, ",") {
		id, err := strconv.Atoi(strings.TrimSpace(setID))
		if err != nil {
			err = fmt.Errorf("while parsing int from string %s in parameters: %s", setID, err.Error())
			logrus.Error(err)
			return err
		}
		patientSetIDs = append(patientSetIDs, id)
	}
	putCohorts, err := NewPutCohorts(accessToken, patientSetIDs, cohortName, disableTLSCheck)
	if err != nil {
		err = fmt.Errorf("while crafting new put cohorts request: %s", err.Error())
		logrus.Error(err)
		return err
	}

	err = putCohorts.Execute()
	if err != nil {
		err = fmt.Errorf("cohorts request execution: %s", err.Error())
		logrus.Error(err)
		return err
	}
	return nil
}

// ExecuteRemoveCohorts executes a post cohorts query
func ExecuteRemoveCohorts(token, username, password, cohortName string, disableTLSCheck bool) error {
	accessToken, err := utilclient.RetrieveOrGetNewAccessToken(token, username, password, disableTLSCheck)
	if err != nil {
		err = fmt.Errorf("while retrieving access token: %s", err.Error())
		logrus.Error(err)
		return err
	}
	logrus.Debug("access token received")
	logrus.Tracef("token %s", accessToken)

	logrus.Debug("creating remove cohorts request")
	logrus.Tracef(" cohort name %s", cohortName)
	removeCohorts, err := NewRemoveCohorts(accessToken, cohortName, disableTLSCheck)
	if err != nil {
		err = fmt.Errorf("while crafting new remove cohorts request: %s", err.Error())
		logrus.Error(err)
		return err
	}

	err = removeCohorts.Execute()
	if err != nil {
		err = fmt.Errorf("cohorts request execution: %s", err.Error())
		logrus.Error(err)
		return err
	}
	return nil

}
