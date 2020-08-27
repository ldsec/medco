package survivalserver

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

type SqlTimePoint struct {
	timePoint             int
	localEventAggregate   int
	localCensoringAggrete int
}

func BuildTimePoints(db *sql.DB, patientList []int64, startConceptCode string, startConceptColumn string, startConceptModifier string, endConceptCode string, endConceptColumn string, endConceptModifier string) (timePoints []SqlTimePoint, err error) {
	logrus.Print("SQL query : " + sql6)
	pList := make([]string, len(patientList))
	for i, pNum := range patientList {
		pList[i] = strconv.FormatInt(pNum, 10)
	}
	patients := "{" + strings.Join(pList, ",") + "}"
	rows, err := db.Query(sql6, startConceptCode, startConceptModifier, patients, endConceptCode, endConceptModifier)
	if err != nil {
		return
	}
	timePointString := new(string)
	eventsString := new(string)
	censoringString := new(string)
	for rows.Next() {
		sqlTimePoint := SqlTimePoint{}
		scanErr := rows.Scan(timePointString, eventsString, censoringString)
		if scanErr != nil {
			err = scanErr
			return
		}
		sqlTimePoint.localEventAggregate, err = strconv.Atoi(*eventsString)
		if err != nil {

			return
		}
		sqlTimePoint.localCensoringAggrete, err = strconv.Atoi(*censoringString)
		if err != nil {

			return
		}
		sqlTimePoint.timePoint, err = strconv.Atoi(*timePointString)
		if err != nil {

			return
		}
		timePoints = append(timePoints, sqlTimePoint)
	}
	return

}
