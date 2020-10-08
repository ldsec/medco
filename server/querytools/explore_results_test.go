package querytoolsserver

import (
	"testing"

	_ "github.com/lib/pq"

	"github.com/stretchr/testify/assert"
)

func TestExploreResults(t *testing.T) {
	testDB, err := DBResolver("MC_DB_HOST", "medcoconnectorsrv0")
	if err != nil {
		t.Fatal(err)
	}
	err = testDB.Ping()
	if err != nil {
		t.Fatal(err)
	}
	queryID, err := InsertExploreResultInstance(testDB, "test", "name1", "")
	assert.NoError(t, err)

	defer testDB.Exec(exploreResultDeletion, queryID)

	// both sets undefined
	err = UpdateExploreResultInstance(testDB, queryID, 0, []int{1, 2, 3}, nil, nil)
	assert.Error(t, err)
	set := new(int)
	*set = -1
	err = UpdateExploreResultInstance(testDB, queryID, 0, []int{1, 2, 3}, set, nil)
	assert.NoError(t, err)

	// cannot call more than once for the same query id
	err = UpdateExploreResultInstance(testDB, queryID, 0, []int{1, 2, 3}, set, nil)
	assert.Error(t, err)

	queryIDError, err := InsertExploreResultInstance(testDB, "test", "name2", "")
	assert.NoError(t, err)
	defer testDB.Exec(exploreResultDeletion, queryIDError)

	err = UpdateErrorExploreResultInstance(testDB, queryIDError)
	assert.NoError(t, err)
	// cannot call more than once for the same query id
	err = UpdateErrorExploreResultInstance(testDB, queryIDError)
	assert.Error(t, err)

}

func TestCheckQueryID(t *testing.T) {
	testDB, err := DBResolver("MC_DB_HOST", "medcoconnectorsrv0")
	if err != nil {
		t.Fatal(err)
	}
	err = testDB.Ping()
	if err != nil {
		t.Fatal(err)
	}

	hasID, err := CheckQueryID(testDB, "test", -1)
	assert.Equal(t, true, hasID)
	assert.NoError(t, err)

	hasID, err = CheckQueryID(testDB, "test", -10)
	assert.Equal(t, false, hasID)
	assert.NoError(t, err)

}

const exploreResultDeletion = `
DELETE FROM query_tools.explore_query_results
WHERE query_id = $1
`
