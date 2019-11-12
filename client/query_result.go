package medcoclient

import (
	"github.com/ldsec/medco-connector/restapi/models"
	"github.com/ldsec/medco-connector/wrappers/unlynx"
	"github.com/sirupsen/logrus"
	"time"
)

// ExploreQueryResult contains the decrypted results of a node
type ExploreQueryResult struct {
	Count int64
	PatientList []int64
	Times map[string]time.Duration
}

// newQueryResult parses a query result from a node and decrypts its fields
func newQueryResult(nodeResult *models.ExploreQueryResultElement, privateKey string) (parsedResult *ExploreQueryResult, err error) {
	parsedResult = &ExploreQueryResult{
		Times: make( map[string]time.Duration),
	}

	// decrypt count
	parsedResult.Count, err = unlynx.Decrypt(nodeResult.EncryptedCount, privateKey)
	if err != nil {
		logrus.Error("error decrypting count: ", err)
		return
	}

	// decrypt patient list
	for _, patientID := range nodeResult.EncryptedPatientList {
		decryptedPatientID, err := unlynx.Decrypt(patientID, privateKey)
		if err != nil {
			logrus.Error("error decrypting patient ID: ", err)
			return nil, err
		}

		if decryptedPatientID != 0 {
			parsedResult.PatientList = append(parsedResult.PatientList, decryptedPatientID)
		}
	}

	// parse times
	for _, timer := range nodeResult.Timers {
		parsedResult.Times[timer.Name] = time.Duration(timer.Milliseconds) * time.Millisecond
	}

	logrus.Info("Node result: count=", parsedResult.Count, ", # patient IDs decrypted=",
		len(nodeResult.EncryptedPatientList), ", # non dummy patients=", len(parsedResult.PatientList))
	return
}
