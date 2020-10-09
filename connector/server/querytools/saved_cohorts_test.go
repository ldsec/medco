package querytoolsserver

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetPatientList(t *testing.T) {
	testDB, err := DBResolver("MC_DB_HOST", "medcoconnectorsrv0")
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
	testDB, err := DBResolver("MC_DB_HOST", "medcoconnectorsrv0")
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
	assert.Equal(t, true, len(cohorts) > 0)
	//change user_id

	cohorts, err = GetSavedCohorts(testDB, "testestest")
	assert.Equal(t, len(cohorts), 0)
	err = testDB.Close()
	if err != nil {
		t.Fatal(err)
	}

}

func TestGetDate(t *testing.T) {
	testDB, err := DBResolver("MC_DB_HOST", "medcoconnectorsrv0")
	if err != nil {
		t.Fatal(err)
	}

	err = testDB.Ping()
	if err != nil {
		t.Fatal(err)
	}

	updateDate, err := GetDate(testDB, "test", -1)
	if err != nil {
		t.Fatal(err)
	}
	expectedDate, err := time.Parse(time.RFC3339, "2020-08-25T13:57:00Z")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expectedDate, updateDate)

	//change cohort_id
	updateDate, err = GetDate(testDB, "test", -2)
	assert.Error(t, err)
	//change user_id
	updateDate, err = GetDate(testDB, "testestest", -1)
	assert.Error(t, err)

	err = testDB.Close()
	if err != nil {
		t.Fatal(err)
	}

}

func TestInsertCohortAndRemoveCohort(t *testing.T) {
	testDB, err := DBResolver("MC_DB_HOST", "medcoconnectorsrv0")
	if err != nil {
		t.Fatal(err)
	}
	err = testDB.Ping()
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	cohortID, err := InsertCohort(testDB, "test", -1, "testCohort2", now, now)
	if err != nil {
		t.Fatal(err)
	}

	cohorts, err := GetSavedCohorts(testDB, "test")
	assert.NoError(t, err)

	found := false
	for _, cohort := range cohorts {
		if cohort.CohortID == cohortID {
			found = true
			assert.Equal(t, cohortID, cohort.CohortID)
			assert.Equal(t, "testCohort2", cohort.CohortName)
			assert.Equal(t, now.Format(time.Stamp), cohort.UpdateDate.Format(time.Stamp))
			break
		}
	}

	assert.Equal(t, found, true)

	err = RemoveCohort(testDB, "test", "testCohort2")
	assert.NoError(t, err)

	err = RemoveCohort(testDB, "test", "testCohort2")
	assert.NoError(t, err)

}

func TestDoesCohortExist(t *testing.T) {
	testDB, err := DBResolver("MC_DB_HOST", "medcoconnectorsrv0")
	if err != nil {
		t.Fatal(err)
	}
	err = testDB.Ping()
	if err != nil {
		t.Fatal(err)
	}

	exists, err := DoesCohortExist(testDB, "test", "testCohort")
	assert.NoError(t, err)
	assert.Equal(t, true, exists)

	exists, err = DoesCohortExist(testDB, "test", "IForSureDoNotExist")
	assert.NoError(t, err)
	assert.Equal(t, false, exists)

	exists, err = DoesCohortExist(testDB, "IForSureDoNotExist", "testCohort")
	assert.NoError(t, err)
	assert.Equal(t, false, exists)
}
