// +build integration_test
package referenceintervalserver

import (
	"strconv"
	"testing"
	"time"

	"github.com/ldsec/medco/connector/restapi/models"
	"github.com/ldsec/medco/connector/restapi/server/operations/explore_statistics"
	utilserver "github.com/ldsec/medco/connector/util/server"
	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq"
)

func init() {
	utilserver.SetForTesting()
}

func getParamBody(bucketSize float64, minObservation float64) explore_statistics.ExploreStatisticsBody {
	return explore_statistics.ExploreStatisticsBody{
		ID: "1234",
		CohortDefinition: &explore_statistics.ExploreStatisticsParamsBodyCohortDefinition{
			IsPanelEmpty: true,
			Panels:       make([]*models.Panel, 0),
		},
		Concepts:       make([]string, 0),
		Modifiers:      make([]*explore_statistics.ExploreStatisticsParamsBodyModifiersItems0, 0),
		BucketSize:     bucketSize,
		UserPublicKey:  "whatever",
		MinObservation: minObservation,
	}
}

func newQueryResults(observations []float64) (queryResults []QueryResult) {
	queryResults = make([]QueryResult, 0)
	for i := 0; i < len(observations); i++ {
		queryResults = append(queryResults,
			QueryResult{
				NumericValue:  observations[i],
				Unit:          "a",
				PatientNumber: int64(i),
			},
		)
	}
	return
}

func newQuery(t *testing.T, bucketWidth float64, minObs float64) (query *Query) {
	body := getParamBody(bucketWidth, minObs)

	query, err := NewQuery(
		"1234",
		body)

	if err != nil {
		t.Fail()
	}
	return
}

func TestProcessObservations(t *testing.T) {
	bucketWidth := 1.
	query := newQuery(t, bucketWidth, 0.)

	queryResults := newQueryResults([]float64{
		1, 1.5, 1.7,

		2.354, 2.65,

		3, 3.3, 3.6, 3.8,
	})

	encCounts, _, err := query.locallyProcessObservations(bucketWidth, queryResults, time.Now(), "fake_concept_code", fakeEncrypt)

	if err != nil {
		t.Fail()
	}

	assert.Equal(t, 4, len(encCounts))

	assertCount := func(expectedCount string, intervalIndex int) {
		assert.Equal(t, expectedCount, encCounts[intervalIndex])
	}

	assertCount("0", 0) //[0, 1[
	assertCount("3", 1) //[1, 2[
	assertCount("2", 2) //[2, 3[
	assertCount("4", 3) //[3, 4[
}

func TestProcessObservations2(t *testing.T) {
	bucketWidth := 1.5
	query := newQuery(t, bucketWidth, 0.)

	queryResults := newQueryResults([]float64{
		0, .3, .6, .9, 1, 1.2,

		1.5, 1.7, 2.5, 2.354, 2.65,

		3, 3.3, 3.6, 3.8,
	})

	encCounts, _, err := query.locallyProcessObservations(bucketWidth, queryResults, time.Now(), "fake_concept_code", fakeEncrypt)

	if err != nil {
		t.Fail()
	}

	assert.Equal(t, 3, len(encCounts))

	assertCount := func(expectedCount string, intervalIndex int) {
		assert.Equal(t, expectedCount, encCounts[intervalIndex])
	}

	assertCount("6", 0) //[0, 1.5[
	assertCount("5", 1) //[1.5, 3[
	assertCount("4", 2) //[3, 4.5[
}

func TestProcessObservations3(t *testing.T) {
	bucketWidth := 2.
	query := newQuery(t, bucketWidth, 1)

	queryResults := newQueryResults([]float64{
		0, .3, .6, .9, //ignored in principle

		1.5, 2.9,

		4.3, 3.3, 3.1,
	})

	encCounts, _, err := query.locallyProcessObservations(bucketWidth, queryResults, time.Now(), "fake_concept_code", fakeEncrypt)

	if err != nil {
		t.Fail()
	}

	assert.Equal(t, 2, len(encCounts))

	assertCount := func(expectedCount string, intervalIndex int) {
		assert.Equal(t, expectedCount, encCounts[intervalIndex])
	}

	assertCount("2", 0) //[1, 3[
	assertCount("3", 1) //[3, 5[
}

func TestProcessObservations4(t *testing.T) {
	bucketWidth := 1.
	query := newQuery(t, bucketWidth, -3)

	queryResults := newQueryResults([]float64{
		-3, -2.5, -2.1, //[-3, -2[
		//[-2, -1[
		-1, -.5, //[-1, 0[

		0.1, 0.6, 0.72134243, 0.81, //[0, 1[
	})

	encCounts, _, err := query.locallyProcessObservations(bucketWidth, queryResults, time.Now(), "fake_concept_code", fakeEncrypt)

	if err != nil {
		t.Fail()
	}

	assert.Equal(t, 4, len(encCounts))

	assertCount := func(expectedCount string, intervalIndex int) {
		assert.Equal(t, expectedCount, encCounts[intervalIndex])
	}

	assertCount("3", 0) //[-3, -2[
	assertCount("0", 1) //[-2, -1[
	assertCount("2", 2) //[-1, 0[
	assertCount("4", 3) //[0, 1[
}

func fakeEncrypt(input int64) (string, error) {
	return strconv.Itoa(int(input)), nil
}
