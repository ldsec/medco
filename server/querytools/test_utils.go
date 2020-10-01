package querytoolsserver

import (
	"database/sql"
	"net"
	"os"
	"strings"
)

// DBResolver attempts to connect to MC_DB_HOST address, else it connects to localhost
func DBResolver(envVarString, dbNameString string) (*sql.DB, error) {
	dbNameString = strings.Replace(dbNameString, ";", "", -1)
	dbNameString = strings.Replace(dbNameString, " ", "", -1)
	dbName := strings.Replace(dbNameString, "$", "", -1)
	if dbHost, isPresent := os.LookupEnv(envVarString); isPresent {
		addr, _ := net.LookupHost(dbHost)
		if len(addr) == 0 {
			goto localhostdb
		}
		//try the found dbHost
		sqlDB, err := sql.Open("postgres", "host="+dbHost+" port=5432 user=postgres password=postgres dbname="+dbName+" sslmode=disable")
		if err != nil {
			goto localhostdb
		} else {
			return sqlDB, nil
		}
	}

localhostdb:
	return sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname="+dbName+" sslmode=disable")

}
