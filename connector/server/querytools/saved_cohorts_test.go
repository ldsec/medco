//go:build integration_test
// +build integration_test

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
	const expectedPatientNumber = 228
	utilserver.TestDBConnection(t)

	pList, err := GetPatientList("test", "testCohort")
	assert.NoError(t, err)
	assert.Equal(t, expectedPatientNumber, len(pList))

	_, err = GetPatientList("thisUserDoesNotExist", "testCohort")
	assert.Error(t, err)
	_, err = GetPatientList("test", "thisCohortDoesNotExist")
	assert.Error(t, err)

}

func TestGetSavedCohorts(t *testing.T) {
	utilserver.TestDBConnection(t)

	cohorts, err := GetSavedCohorts("test", 0)
	assert.NoError(t, err)
	assert.Equal(t, true, len(cohorts) > 0)
	//change user_id

	cohorts, err = GetSavedCohorts("testestest", 0)
	assert.Equal(t, len(cohorts), 0)

}

func TestGetDate(t *testing.T) {
	utilserver.TestDBConnection(t)

	updateDate, err := GetDate("test", -1)
	if err != nil {
		t.Fatal(err)
	}
	expectedDate, err := time.Parse(time.RFC3339, "2020-08-25T13:57:00Z")
	assert.NoError(t, err)
	assert.Equal(t, expectedDate, updateDate)

	//change cohort_id
	updateDate, err = GetDate("test", -2)
	assert.Error(t, err)
	//change user_id
	updateDate, err = GetDate("testestest", -1)
	assert.Error(t, err)
}

func TestInsertCohortAndUpdateCohortAndRemoveCohort(t *testing.T) {
	utilserver.TestDBConnection(t)

	now := time.Now()
	_, err := UpdateCohort("testCohort2", "test", -1, now)
	assert.Error(t, err)

	cohortID, err := InsertCohort("test", -1, "testCohort2", now, now)
	assert.NoError(t, err)

	_, err = InsertCohort("test", -1, "testCohort2", now, now)
	assert.Error(t, err)

	cohorts, err := GetSavedCohorts("test", 1)
	assert.NoError(t, err)

	cohorts, err = GetSavedCohorts("test", 0)
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

	assert.Equal(t, true, found)

	time.Sleep(10 * time.Second)
	now2 := time.Now()
	cohortID, err = UpdateCohort("testCohort2", "test", -1, now2)
	assert.NoError(t, err)
	cohorts, err = GetSavedCohorts("test", 0)
	assert.NoError(t, err)

	cohorts2, err := GetSavedCohorts("test", 0)
	assert.NoError(t, err)

	found = false
	for _, cohort := range cohorts2 {
		if cohort.CohortID == cohortID {
			found = true
			assert.Equal(t, cohortID, cohort.CohortID)
			assert.Equal(t, "testCohort2", cohort.CohortName)
			assert.Equal(t, now2.Format(time.Stamp), cohort.UpdateDate.Format(time.Stamp))
			break
		}
	}
	assert.Equal(t, true, found)

	cohortID, err = UpdateCohort("testCohort2", "test", -1, now2)
	assert.NoError(t, err)

	err = RemoveCohort("test", "testCohort2")
	assert.NoError(t, err)

	err = RemoveCohort("test", "testCohort2")
	assert.Error(t, err)

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
