package querytools

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	cohortscommon "github.com/ldsec/medco-connector/cohorts/common"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var ConnectorDB *sql.DB

func init() {
	var err error
	ConnectorDB, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("MC_DB_HOST"), os.Getenv("MC_DB_PORT"), os.Getenv("MC_DB_USER"), os.Getenv("MC_DB_PW"), os.Getenv("MC_DB_NAME")))
	if err != nil {
		logrus.Error(err)
	}

	err = ConnectorDB.Ping()
	if err != nil {
		logrus.Error(err)
	}

}

func GetPatientList(db *sql.DB, user_id string, resultInstanceID int64) (patientNums []int64, err error) {
	row := db.QueryRow(getPatientList, user_id, resultInstanceID)
	patientNumsString := new(string)
	err = row.Scan(patientNumsString)
	var pNum int64
	for _, pID := range strings.Split(strings.Trim(*patientNumsString, "{}"), ",") {

		pNum, err = strconv.ParseInt(pID, 10, 64)
		if err != nil {
			return
		}
		patientNums = append(patientNums, pNum)
	}

	return
}

func GetSavedCohorts(db *sql.DB, userID string) ([]cohortscommon.Cohort, error) {
	rows, err := db.Query(getCohorts, userID)
	if err != nil {
		return nil, err
	}
	var id int
	var qid int
	var name string
	var createDateString string
	var createDate time.Time
	var updateDateString string
	var updateDate time.Time
	var cohorts = make([]cohortscommon.Cohort, 0)
	for rows.Next() {
		err = rows.Scan(&id, &qid, &name, &createDateString, &updateDateString)
		if err != nil {
			return nil, err
		}
		createDate, err = time.Parse(time.RFC3339, createDateString)
		if err != nil {
			return nil, err
		}
		updateDate, err = time.Parse(time.RFC3339, updateDateString)
		if err != nil {
			return nil, err
		}
		cohorts = append(cohorts, cohortscommon.Cohort{
			CohortId:     id,
			QueryId:      qid,
			CohortName:   name,
			CreationDate: createDate,
			UpdateDate:   updateDate,
		})
	}
	err = rows.Close()
	if err != nil {
		return nil, err
	}

	logrus.Infof("Got %d cohorts", len(cohorts))
	return cohorts, nil
}

func GetDate(db *sql.DB, userID string, cohortID int) (time.Time, error) {
	row := db.QueryRow(getDate, userID, cohortID)
	timeString := new(string)
	err := row.Scan(timeString)
	if err != nil {
		return time.Now(), err
	}

	timeParsed, err := time.Parse(time.RFC3339, *timeString)

	return timeParsed, err

}

func InsertCohort(db *sql.DB, userID string, queryID int, cohortName string, createDate, updateDate time.Time) (int, error) {
	row := db.QueryRow(insertCohort, userID, queryID, cohortName, createDate, updateDate)
	res := new(string)
	err := row.Scan(res)
	if err != nil {
		return -1, err
	}
	cohortID, err := strconv.Atoi(*res)
	return cohortID, err
}
