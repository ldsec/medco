package querytools

import (
	"database/sql"
	"log"
	"testing"
)

func TestGetPatientList(t *testing.T) {
	testDB, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=medcoconnectorsrv0 sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	err = testDB.Ping()
	if err != nil {
		t.Fatal(err)
	}

	pList, err := GetPatientList(testDB, "test", int64(-1))
	if err != nil {
		t.Fatal(err)
	}
	if len(pList) != 228 {
		t.Fatalf("Expected 228 patients, got: %d", len(pList))
	}
	err = testDB.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetSavedCohorts(t *testing.T) {
	testDB, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=medcoconnectorsrv0 sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	err = testDB.Ping()
	if err != nil {
		t.Fatal(err)
	}

	cohorts, err := GetSavedCohorts(testDB, "test")
	if err != nil {
		t.Fatal(err)
	}
	log.Print(cohorts)
	err = testDB.Close()
	if err != nil {
		t.Fatal(err)
	}

	if len(cohorts) != 1 {
		t.Fatalf("Expected 1 row, got: %d", len(cohorts))
	}
}
