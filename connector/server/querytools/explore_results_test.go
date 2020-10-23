// +build integration_test

package querytoolsserver

import (
	"testing"

	utilserver "github.com/ldsec/medco/connector/util/server"

	_ "github.com/lib/pq"

	"github.com/stretchr/testify/assert"
)

func init() {
	utilserver.SetForTesting()
}

func TestExploreResults(t *testing.T) {
	utilserver.TestDBConnection(t)

	queryID, err := InsertExploreResultInstance("test", "name1", "")
	assert.NoError(t, err)

	defer utilserver.DBConnection.Exec(exploreResultDeletion, queryID)

	// both sets undefined
	err = UpdateExploreResultInstance(queryID, 0, []int{1, 2, 3}, nil, nil)
	assert.Error(t, err)
	set := new(int)
	*set = -1
	err = UpdateExploreResultInstance(queryID, 0, []int{1, 2, 3}, set, nil)
	assert.NoError(t, err)

	// cannot call more than once for the same query id
	err = UpdateExploreResultInstance(queryID, 0, []int{1, 2, 3}, set, nil)
	assert.Error(t, err)

	queryIDError, err := InsertExploreResultInstance("test", "name2", "")
	assert.NoError(t, err)
	defer utilserver.DBConnection.Exec(exploreResultDeletion, queryIDError)

	err = UpdateErrorExploreResultInstance(queryIDError)
	assert.NoError(t, err)
	// cannot call more than once for the same query id
	err = UpdateErrorExploreResultInstance(queryIDError)
	assert.Error(t, err)

}

func TestCheckQueryID(t *testing.T) {
	utilserver.TestDBConnection(t)

	hasID, err := CheckQueryID("test", -1)
	assert.Equal(t, true, hasID)
	assert.NoError(t, err)

	hasID, err = CheckQueryID("test", -10)
	assert.Equal(t, false, hasID)
	assert.NoError(t, err)

}

const exploreResultDeletion = `
DELETE FROM query_tools.explore_query_results
WHERE query_id = $1
`
