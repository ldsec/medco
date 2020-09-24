package survivalserver

import (
	"database/sql"
	"strconv"
	"strings"

	survivalcommon "github.com/ldsec/medco-connector/survival/common"

	"github.com/sirupsen/logrus"
)

// BuildTimePoints execute a SQL query that returns event counts per time point, for given input patient set, start and end  concept codes and modifiers
func BuildTimePoints(db *sql.DB, patientList []int64, startConceptCode string, startConceptModifier string, endConceptCode string, endConceptModifier string, timeLimit int) (timePoints survivalcommon.TimePoints, err error) {
	logrus.Debug("SQL query : " + sql6)
	pList := make([]string, len(patientList))
	for i, pNum := range patientList {
		pList[i] = strconv.FormatInt(pNum, 10)
	}
	patients := "{" + strings.Join(pList, ",") + "}"
	rows, err := db.Query(sql6, startConceptCode, startConceptModifier, patients, endConceptCode, endConceptModifier, timeLimit)
	if err != nil {
		return
	}
	timePointString := new(string)
	eventsString := new(string)
	censoringString := new(string)
	for rows.Next() {
		sqlTimePoint := survivalcommon.TimePoint{}
		scanErr := rows.Scan(timePointString, eventsString, censoringString)
		if scanErr != nil {
			err = scanErr
			return
		}
		sqlTimePoint.Events.EventsOfInterest, err = strconv.ParseInt(*eventsString, 10, 64)
		if err != nil {

			return
		}
		sqlTimePoint.Events.CensoringEvents, err = strconv.ParseInt(*censoringString, 10, 64)
		if err != nil {

			return
		}
		sqlTimePoint.Time, err = strconv.Atoi(*timePointString)
		if err != nil {

			return
		}
		timePoints = append(timePoints, sqlTimePoint)
	}
	return

}
