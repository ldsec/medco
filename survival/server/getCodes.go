package survivalserver

import (
	"database/sql"
	"strings"

	"github.com/sirupsen/logrus"
)

//TODO use I2B2 ONT messaging

func GetCode(db *sql.DB, path string) (string, error) {
	tokens := strings.Split(path, `\`)
	table := tokens[2]
	logrus.Debug("table", strings.ToLower(table))
	fullName := `\` + strings.Join(tokens[3:], `\`)
	var res string
	logrus.Debug("fullName", fullName)
	row := db.QueryRow(sqlConcept, fullName)
	err := row.Scan(&res)

	return res, err

}

const sqlConcept = `
SELECT c_basecode FROM medco_ont.e2etest
WHERE c_fullname = $1;
`
