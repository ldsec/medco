package i2b2

import (
	"errors"

	"github.com/ldsec/medco/connector/restapi/models"
	utilserver "github.com/ldsec/medco/connector/util/server"
	"github.com/ldsec/medco/connector/wrappers/unlynx"
	"github.com/sirupsen/logrus"
)

// ExecutePsmQuery executes an i2b2 PSM query and returns the corresponding patient set id
func ExecutePsmQuery(queryName string, panels []*models.Panel, sequences []*models.TimingSequenceInfo, queryTiming models.Timing) (patientCount string, patientSetID string, err error) {

	// craft and execute request
	xmlResponse := &Response{
		MessageBody: &CrcPsmRespMessageBody{},
	}
	xmlQuery, err := NewCrcPsmReqFromQueryDef(
		queryName,
		panels,
		sequences,
		[]ResultOutputName{Patientset, PatientCountXML},
		queryTiming,
	)
	if err != nil {
		return
	}
	err = i2b2XMLRequest(
		utilserver.I2b2HiveURL+"/QueryToolService/request",
		xmlQuery,
		xmlResponse,
	)

	if err != nil {
		return
	}

	// check error in CRC response
	err = xmlResponse.MessageBody.(*CrcPsmRespMessageBody).checkStatus()
	if err != nil {
		logrus.Error("i2b2 CRC response error:", err)
		return
	}

	// extract results from result instances
	for _, resultInstance := range xmlResponse.MessageBody.(*CrcPsmRespMessageBody).Response.QueryResultInstances {

		// check error in result instance
		err = resultInstance.checkStatus()
		if err != nil {
			logrus.Error("i2b2 instance error:", err)
			return
		}

		// extract results
		if resultInstance.QueryResultType.Name == string(Patientset) {
			patientSetID = resultInstance.ResultInstanceID
		} else if resultInstance.QueryResultType.Name == string(PatientCountXML) {
			patientCount = resultInstance.SetSize
		}
	}

	if patientCount == "" || patientSetID == "" {
		err = errors.New("i2b2 results not found")
		logrus.Error(err)
		return
	}

	return
}

// GetPatientSet retrieves an i2b2 patient set
func GetPatientSet(patientSetID string, generateDummyFlags bool) (patientIDs []string, patientDummyFlags []string, err error) {

	// craft and execute request
	xmlResponse := &Response{
		MessageBody: &CrcPdoRespMessageBody{},
	}

	err = i2b2XMLRequest(
		utilserver.I2b2HiveURL+"/QueryToolService/pdorequest",
		NewCrcPdoReqFromInputList(patientSetID),
		xmlResponse,
	)

	if err != nil {
		return
	}

	// extract patient data
	for _, patient := range xmlResponse.MessageBody.(*CrcPdoRespMessageBody).Response.PatientData.PatientSet.Patient {

		dummyFlagFound := false
		for _, patientColumn := range patient.Param {
			if patientColumn.Column == "enc_dummy_flag_cd" && len(patientColumn.Text) > 0 {
				patientIDs = append(patientIDs, patient.PatientID)
				patientDummyFlags = append(patientDummyFlags, patientColumn.Text)
				dummyFlagFound = true
				break
			}
		}
		if !dummyFlagFound {
			patientIDs = append(patientIDs, patient.PatientID)
			if generateDummyFlags {
				var realPatientFlag string
				realPatientFlag, err = unlynx.EncryptWithCothorityKey(int64(1))
				if err != nil {
					return
				}
				patientDummyFlags = append(patientDummyFlags, realPatientFlag)
				logrus.Warn("GetPatientSet: patient ", patient.PatientID, " misses dummy flag. Setting it as a real patient")
			} else {
				logrus.Warn("GetPatientSet: patient ", patient.PatientID, " misses dummy flag.")
			}
		}
	}
	return
}
