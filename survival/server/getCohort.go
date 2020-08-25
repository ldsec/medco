package survivalserver

import "database/sql"

func GetPatientList(db *sql.DB, cohortID int64, userID string) (list []int64, err error) {
	row := db.QueryRow(getPatientList, string(cohortID), userID)

	err = row.Scan(list)
	return
}

const getPatientList string = `
SELECT enc_result_set FROM explore_query_results
WHERE query_id IN
(SELECT query_id FROM cohorts
WHERE user_id=$1 AND cohort_id =$2)
`
