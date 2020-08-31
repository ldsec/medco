package survivalserver

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

func GetPatientList(db *sql.DB, queryID int64, userID string) (list []int64, err error) {
	logrus.Debugf("patient list argument: cohortID %d, FomatInt(cohortID) %s, userID %s", queryID, strconv.FormatInt(queryID, 10), userID)
	row := db.QueryRow(getPatientList, userID, strconv.FormatInt(queryID, 10))
	listString := new(string)
	err = row.Scan(listString)
	var val int64
	for _, elm := range strings.Split(strings.Trim(*listString, "{}"), ",") {
		val, err = strconv.ParseInt(elm, 10, 64)
		if err != nil {
			return
		}
		list = append(list, val)
	}

	return
}

const getPatientList string = `
SELECT clear_result_set FROM query_tools.explore_query_results
WHERE user_id=$1 AND query_id =$2;
`
