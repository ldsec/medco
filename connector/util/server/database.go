package utilserver

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	// Registering postgres driver
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

// InitializeConnectionToDB initializes the database handle
func InitializeConnectionToDB(dbmsHost string, dbmsPort int, dbName string, dbLoginUser string, dbLoginPassword string) (*sql.DB, error) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbmsHost, dbmsPort, dbLoginUser, dbLoginPassword, dbName)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return db, nil

}

// ConvertIntListToString is used to convert a list of int into a list of integer argument for a sql query.
// For instance the output of this function will be accepted by ANY($n::integer[]) where `n` is the index of the parameter in a SQL argument.
func ConvertIntListToString(intList []int64) string {
	strList := make([]string, len(intList))
	for i, num := range intList {
		strList[i] = strconv.FormatInt(num, 10)
	}
	return "{" + strings.Join(strList, ",") + "}"
}
