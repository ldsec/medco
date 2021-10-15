package querytoolsclient

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	medcoclient "github.com/ldsec/medco/connector/client"

	"github.com/ldsec/medco/connector/restapi/models"
	utilclient "github.com/ldsec/medco/connector/util/client"

	"github.com/sirupsen/logrus"
)

// ExecuteGetCohorts executes a get cohorts query and displays its results
func ExecuteGetCohorts(token, username, password string, disableTLSCheck bool, resultFile string, limit int) error {
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
	resultCSV.Write([]string{"node_index", "cohort_name", "cohort_id", "query_id", "creation_date", "update_date", "query_timing", "query_timing_sequence", "panels"})

	for nodeIndex, nodeResult := range cohorts {
		for _, cohortInfo := range nodeResult {
			logrus.Debugf("Writing result %d", nodeIndex)
			panelJSONs, err := marshalPanels(cohortInfo.QueryDefinition.Panels)
			if err != nil {
				err = fmt.Errorf("cohorts request writing results: %s", err.Error())
				logrus.Error(err)
				return err
			}
			queryTimingSequenceJSONs, err := marshalSequenceInfo(cohortInfo.QueryDefinition.QueryTimingSequence)
			if err != nil {
				err = fmt.Errorf("cohorts request writing results: %s", err.Error())
				logrus.Error(err)
				return err
			}
			// removing the quotes from the marshalling process eases parsing of the produced file
			panelJSONs = strings.Replace(panelJSONs, `"`, "", -1)
			queryTimingSequenceJSONs = strings.Replace(queryTimingSequenceJSONs, `"`, "", -1)
			err = resultCSV.Write([]string{
				strconv.Itoa(nodeIndex),
				cohortInfo.CohortName,
				strconv.Itoa(cohortInfo.CohortID),
				strconv.Itoa(cohortInfo.QueryID),
				cohortInfo.CreationDate.Format(time.RFC3339),
				cohortInfo.UpdateDate.Format(time.RFC3339),
				string(cohortInfo.QueryDefinition.QueryTiming),
				queryTimingSequenceJSONs,
				panelJSONs,
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
func ExecutePostCohorts(token, username, password, cohortName, queryIDsString string, disableTLSCheck bool) error {

	accessToken, err := utilclient.RetrieveOrGetNewAccessToken(token, username, password, disableTLSCheck)
	if err != nil {
		err = fmt.Errorf("while retrieving access token: %s", err.Error())
		logrus.Error(err)
		return err
	}
	logrus.Debug("access token received")
	logrus.Tracef("token %s", accessToken)

	logrus.Debug("creating post cohorts request")
	logrus.Tracef("query IDs %v , cohort name %s", queryIDsString, cohortName)

	queryIDs := make([]int, 0)
	for _, queryID := range strings.Split(queryIDsString, ",") {
		id, err := strconv.Atoi(strings.TrimSpace(queryID))
		if err != nil {
			err = fmt.Errorf("while parsing int from string %s in parameters: %s", queryID, err.Error())
			logrus.Error(err)
			return err
		}
		queryIDs = append(queryIDs, id)
	}
	postCohorts, err := NewPostCohorts(accessToken, queryIDs, cohortName, disableTLSCheck)
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
func ExecutePutCohorts(token, username, password, cohortName, queryIDsString string, disableTLSCheck bool) error {

	accessToken, err := utilclient.RetrieveOrGetNewAccessToken(token, username, password, disableTLSCheck)
	if err != nil {
		err = fmt.Errorf("while retrieving access token: %s", err.Error())
		logrus.Error(err)
		return err
	}
	logrus.Debug("access token received")
	logrus.Tracef("token %s", accessToken)

	logrus.Debug("creating post cohorts request")
	logrus.Tracef("query IDs %v , cohort name %s", queryIDsString, cohortName)

	queryIDs := make([]int, 0)
	for _, queryID := range strings.Split(queryIDsString, ",") {
		id, err := strconv.Atoi(strings.TrimSpace(queryID))
		if err != nil {
			err = fmt.Errorf("while parsing int from string %s in parameters: %s", queryID, err.Error())
			logrus.Error(err)
			return err
		}
		queryIDs = append(queryIDs, id)
	}
	putCohorts, err := NewPutCohorts(accessToken, queryIDs, cohortName, disableTLSCheck)
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

// ExecuteCohortsPatientList executes a cohorts patient list query
func ExecuteCohortsPatientList(token, username, password, cohortName, resultFile, timerFile string, disableTLSCheck bool) error {
	accessToken, err := utilclient.RetrieveOrGetNewAccessToken(token, username, password, disableTLSCheck)
	if err != nil {
		err = fmt.Errorf("while retrieving access token: %s", err.Error())
		logrus.Error(err)
		return err
	}
	logrus.Debug("access token received")
	logrus.Tracef("token %s", accessToken)

	// calling API
	cohortsPatientList, err := NewCohortsPatientList(accessToken, cohortName, disableTLSCheck)
	if err != nil {
		err = fmt.Errorf("while creating new cohorts patient list request: %s", err.Error())
		logrus.Error(err.Error())
		return err
	}
	patientLists, nodeTimers, localTimers, err := cohortsPatientList.Execute()
	if err != nil {
		err = fmt.Errorf("while executing patient list request, request ID %s: %s", cohortsPatientList.id, err.Error())
		logrus.Error(err.Error())
		return err
	}

	// displaying results
	resultCSV, err := utilclient.NewCSV(resultFile)
	if err != nil {
		err = fmt.Errorf("cohorts patient list request writing results: %s", err.Error())
		logrus.Error(err)
		return err
	}
	for i, list := range patientLists {
		resultCSV.Write([]string{fmt.Sprintf("Node idx %d", i)})

		// the list is not required to be sorted, but it is guaranteed here for testing purpose
		sort.Slice(list, func(a int, b int) bool { return list[a] < list[b] })

		listString := make([]string, len(list))
		for j, pNum := range list {
			listString[j] = strconv.FormatInt(pNum, 10)
		}
		resultCSV.Write(listString)
	}
	err = resultCSV.Flush()
	if err != nil {
		err = fmt.Errorf("cohorts patient list request flushing result file: %s", err.Error())
		logrus.Error(err)
		return err
	}
	err = resultCSV.Close()
	if err != nil {
		err = fmt.Errorf("cohorts patient list request closing result file: %s", err.Error())
		logrus.Error(err)
		return err
	}

	// dumping timers
	err = medcoclient.DumpTimers(timerFile, nodeTimers, localTimers)
	if err != nil {
		err = fmt.Errorf("while dumping timers: %s", err.Error())
		logrus.Error(err)
		return err
	}

	return nil
}
func marshalPanels(panels []*models.Panel) (res string, err error) {
	panelStrings := make([]string, 0)
	var marshalledPanel []byte
	for _, panel := range panels {
		if panel != nil {
			marshalledPanel, err = panel.MarshalBinary()
			if err != nil {
				err = fmt.Errorf("while marshalling I2B2 panel: %s", err.Error())
				return "", err
			}
			panelStrings = append(panelStrings, string(marshalledPanel))
		}
	}
	return fmt.Sprintf(`{"panels":[%s]}`, strings.Join(panelStrings, ",")), nil
}

func marshalSequenceInfo(sequences []*models.TimingSequenceInfo) (res string, err error) {
	sequenceStrings := make([]string, 0)
	var marshalledSequence []byte
	for _, sequence := range sequences {
		if sequence != nil {
			marshalledSequence, err = sequence.MarshalBinary()
			if err != nil {
				err = fmt.Errorf("while marshalling temporal sequence info: %s", err.Error())
				return "", err
			}
			sequenceStrings = append(sequenceStrings, string(marshalledSequence))
		}
	}
	return fmt.Sprintf(`{"temporalSequence":[%s]}`, strings.Join(sequenceStrings, ",")), nil
}
