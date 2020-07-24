package survivalserver

import (
	utilserver "github.com/ldsec/medco-connector/util/server"
	"github.com/sirupsen/logrus"
)

type SqlTimePoint struct {
	timePoint      int
	localAggregate int
}

func BuildTimePoints(patientList []int64, startConceptCode string, startConceptColumn string, endConceptCode string, endConceptColumn string) (timePoints []SqlTimePoint, err error) {
	logrus.Debug("SQL query : " + sql6)
	rows, err := utilserver.DBConnection.Query(sql6)
	if err != nil {
		return
	}

	for rows.Next() {
		sqlTimePoint := SqlTimePoint{}
		scanErr := rows.Scan(&(sqlTimePoint.timePoint), &(sqlTimePoint.localAggregate))
		if scanErr != nil {
			err = scanErr
			return
		}
		timePoints = append(timePoints, sqlTimePoint)
	}
	return

}
