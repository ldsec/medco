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

func contains(arr []QueryResult, x float64) bool {
	for _, r := range arr {
		if r.NumericValue == x {
			return true
		}
	}
	return false
}

//to generate sample with specific sample std and sample mean: https://stackoverflow.com/questions/51515423/generate-sample-data-with-an-exact-mean-and-standard-deviation
func TestOutlierRemovalNormal0Mean1Std(t *testing.T) {
	observations := []float64{
		0.7084059931947092, -1.3938866529750165, 0.8381589889349612, -0.15854619536495415, -1.648164911610732, -0.6054119313504469, 0.4572723633358847, 0.4186652857269435, -0.24690588110536288, -2.2574927403514495, 0.13992139803328693, 1.2442811370376787, 0.8162841137009632, 0.44723107362341913, 1.1804799163129334, 0.09032358192615982, 1.742408910847399, -0.9388871506088993, -0.19148360032025327, -0.6426536989872252,
	}
	outliers := []float64{3.4, 7, 9, -4, -3.1, -1000}

	mean_val := 0.
	std_val := 1.
	OutlierRemovalTester(t, mean_val, std_val, outliers, observations)
}

func TestOutlierRemovalNormal3Mean5Std(t *testing.T) {
	observations := []float64{
		-4.504332807814968, -3.588438949993705, -1.4253477722082, 4.481350113178383, -0.05858067548906476, 8.695616991288777, 0.8911308392424653, 9.391114371637403, 8.094555884981943, -4.290838531517852, 10.79099641135041, 1.3150492298392589, 10.902284491587732, -0.8873436991601062, 2.228352198498256, 3.60927053212485, -1.1011040410767086, 10.284889522654836, 2.635879581548595, 2.535496309327691,
	}
	outliers := []float64{-12.1, -14, -30, -1000, 18.3, 19, 34, 24, 10000}

	mean_val := 3.
	std_val := 5.
	OutlierRemovalTester(t, mean_val, std_val, outliers, observations)
}

func OutlierRemovalTester(t *testing.T, expected_mean float64, expected_std float64, outliers []float64, observations []float64) {
	// observations taken from the normal distribution with mean 0 and std 1. Ajusted so the sample mean and std exactly equal 0 and 1 respectively.

	assert.InDelta(t, expected_mean, mean(newQueryResults(observations)), 0.00001)
	assert.InDelta(t, expected_std, std(newQueryResults(observations), expected_mean), 0.00001)

	everyObs := append(observations, outliers...)
	everyResults := newQueryResults(everyObs)
	trimmed, err := outlierRemovalHelper(everyResults, expected_mean, expected_std)

	assert.Nil(t, err)

	for _, out := range outliers {
		assert.False(t, contains(trimmed, out))
	}

	for _, obs := range observations {
		assert.True(t, contains(trimmed, obs))
	}

}
