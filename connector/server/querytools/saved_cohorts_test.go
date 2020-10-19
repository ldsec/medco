package querytoolsserver

import (
	"testing"
	"time"

	utilserver "github.com/ldsec/medco/connector/util/server"
	"github.com/stretchr/testify/assert"
)

func init() {
	utilserver.SetForTesting()
}

func TestGetPatientList(t *testing.T) {
	utilserver.TestDBConnection(t)

	pList, err := GetPatientList("test", "testCohort")
	if err != nil {
		t.Fatal(err)
	}
	if len(pList) != 228 {
		t.Fatalf("Expected 228 patients, got: %d", len(pList))
	}

}

func TestGetSavedCohorts(t *testing.T) {
	utilserver.TestDBConnection(t)

	cohorts, err := GetSavedCohorts("test")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, len(cohorts) > 0)
	//change user_id

	cohorts, err = GetSavedCohorts("testestest")
	assert.Equal(t, len(cohorts), 0)

}

func TestGetDate(t *testing.T) {
	utilserver.TestDBConnection(t)

	updateDate, err := GetDate("test", -1)
	if err != nil {
		t.Fatal(err)
	}
	expectedDate, err := time.Parse(time.RFC3339, "2020-08-25T13:57:00Z")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expectedDate, updateDate)

	//change cohort_id
	updateDate, err = GetDate("test", -2)
	assert.Error(t, err)
	//change user_id
	updateDate, err = GetDate("testestest", -1)
	assert.Error(t, err)
}

func TestInsertCohortAndRemoveCohort(t *testing.T) {
	utilserver.TestDBConnection(t)

	now := time.Now()
	cohortID, err := InsertCohort("test", -1, "testCohort2", now, now)
	if err != nil {
		t.Fatal(err)
	}

	cohorts, err := GetSavedCohorts("test")
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

	err = RemoveCohort("test", "testCohort2")
	assert.NoError(t, err)

	err = RemoveCohort("test", "testCohort2")
	assert.NoError(t, err)

}

func TestDoesCohortExist(t *testing.T) {
	utilserver.TestDBConnection(t)

	exists, err := DoesCohortExist("test", "testCohort")
	assert.NoError(t, err)
	assert.Equal(t, true, exists)

	exists, err = DoesCohortExist("test", "IForSureDoNotExist")
	assert.NoError(t, err)
	assert.Equal(t, false, exists)

	exists, err = DoesCohortExist("IForSureDoNotExist", "testCohort")
	assert.NoError(t, err)
	assert.Equal(t, false, exists)
}
