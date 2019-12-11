package utilserver

import (
	"database/sql"
	"fmt"
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

// GetGenomicAnnotationTypes get the genomic annotation types available in the database
func GetGenomicAnnotationTypes() []string {
	//TODO: Make this dynamic
	return []string{"variant_name", "protein_change", "hugo_gene_symbol"}
}
