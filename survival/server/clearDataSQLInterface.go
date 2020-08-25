package survivalserver

import (
	"database/sql"

	"github.com/sirupsen/logrus"
)

type SqlTimePoint struct {
	timePoint      int
	localAggregate int
}

func BuildTimePoints(db *sql.DB, patientList []int64, startConceptCode string, startConceptColumn string, startConceptModifier string, endConceptCode string, endConceptColumn string, endConceptModifier string) (timePoints []SqlTimePoint, err error) {
	logrus.Debug("SQL query : " + sql6)
	rows, err := db.Query(sql6)
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
