package utilserver

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func InitializeConnectionToDB(dbmsHost string, dbmsPort int, dbName string, dbLoginUser string, dbLoginPassword string) (*sql.DB, error) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		DBMSHost, DBMSPort, DBLoginUser, DBLoginPassword, DBName)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return db, nil

}
