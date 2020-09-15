package querytools

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq" //postgres driver
	"github.com/sirupsen/logrus"
)

// ConnectorDB refers to medco connector postgres database
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

// GetPatientList runs a SQL query on db and returns the list of patient IDs for given queryID and userID
func GetPatientList(db *sql.DB, userID string, resultInstanceID int64) (patientNums []int64, err error) {
	row := db.QueryRow(getPatientList, userID, resultInstanceID)
	patientNumsString := new(string)
	err = row.Scan(patientNumsString)
	var pNum int64
	logrus.Tracef("Got response %s", *patientNumsString)
	for _, pID := range strings.Split(strings.Trim(*patientNumsString, "{}"), ",") {

		pNum, err = strconv.ParseInt(pID, 10, 64)
		if err != nil {
			return
		}
		patientNums = append(patientNums, pNum)
	}

	return
}

// GetSavedCohorts runs a SQL query on db and returns the list of saved cohorts for given queryID and userID
func GetSavedCohorts(db *sql.DB, userID string) ([]Cohort, error) {
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
	var cohorts = make([]Cohort, 0)
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
		cohorts = append(cohorts, Cohort{
			CohortID:     id,
			QueryID:      qid,
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

// GetDate runs a SQL query on db and returns the update date of cohort corresponding to  cohortID
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

// InsertCohort runs a SQL query to either insert a new cohort or update an existing one
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
