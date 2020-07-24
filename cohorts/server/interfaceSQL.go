package cohortsserver

import (
	utilserver "github.com/ldsec/medco-connector/util/server"
	_ "github.com/lib/pq"

	cohortscommon "github.com/ldsec/medco-connector/cohorts/common"
)

func GetCohorts(userId string) (cohorts []cohortscommon.Cohort, err error) {

	rows, err := utilserver.DBConnection.Query(getCohorts, userId)

	for rows.Next() {
		cohort := &cohortscommon.Cohort{}
		scanErr := rows.Scan(cohort)
		if scanErr != nil {
			err = scanErr
			return
		}
		cohorts = append(cohorts, *cohort)
	}

	return
}

func GetDates(userId, cohortName string) (createDate int64, updateDate int64, err error) {
	row := utilserver.DBConnection.QueryRow(getDates, userId, cohortName)
	err = row.Scan(&createDate, &updateDate)
	return
}

func InsertCohorts(userId, cohortName string, resultInstanceID int, createDate, updateDate int64) (err error) {
	_, err = utilserver.DBConnection.Exec(insertCohorts, userId, cohortName, resultInstanceID, createDate, updateDate)
	return
}
