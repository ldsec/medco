package querytools

import (
	"time"

	utilserver "github.com/ldsec/medco-connector/util/server"
	_ "github.com/lib/pq"
)

func InsertQueryInstance(userId, groupId, batchMode string, startDate time.Time) (queryInstanceId int, err error) {

	rows, err := utilserver.DBConnection.Query(insertQueryInstance, userId, groupId, batchMode, startDate)
	if err != nil {
		return
	}
	if rows.Next() {
		scanErr := rows.Scan(&queryInstanceId)
		if scanErr != nil {
			err = scanErr
			return
		}
	}
	return
}
func UpdateQueryInstance(queryInstanceId, clearQueryInstanceId, cipherQueryInstanceId int, endDate time.Time, batchMode string, statusTypID StatusType, message string) (err error) {
	_, err = utilserver.DBConnection.Query(updateQueryInstance, queryInstanceId, clearQueryInstanceId, cipherQueryInstanceId, endDate, batchMode, statusTypID, message)
	return

}

func InsertResultInstance(resultType ResultType, queryInstanceId int, startDate time.Time, statusTypID StatusType) (resultInstanceID int, err error) {

	rows, err := utilserver.DBConnection.Query(insertResultInstance, resultType, queryInstanceId, startDate, statusTypID)
	if err != nil {
		return
	}
	if rows.Next() {
		scanErr := rows.Scan(&resultInstanceID)
		if scanErr != nil {
			err = scanErr
			return
		}
	}
	return
}

func UpdateResultInstance(resultInstanceID int, setSize string, endDate time.Time, statusTypID StatusType, message string, realSetSize string, obfuscationMethod string) (err error) {
	_, err = utilserver.DBConnection.Query(updateResultInstance, resultInstanceID, setSize)
	return
}

func InsertPatientSet(resultInstanceID int, patientSetIdx int64, patientNum int64) (patientSetCollectionID int64, err error) {

	rows, err := utilserver.DBConnection.Query(insertPatientSetCollection, resultInstanceID, patientSetIdx, patientNum)

	if err != nil {
		return
	}
	if rows.Next() {
		scanErr := rows.Scan(&patientSetCollectionID)
		if scanErr != nil {
			err = scanErr
			return
		}
	}
	return
}
func GetPatientList(resultInstanceID int) (patientNums []int64, err error) {
	rows, err := utilserver.DBConnection.Query(getPatientList, resultInstanceID)
	if err != nil {
		return
	}
	if rows.Next() {
		var patientNum int64
		scanErr := rows.Scan(&patientNum)
		if scanErr != nil {
			err = scanErr
			return
		}
		patientNums = append(patientNums, patientNum)
	}
	return
}
