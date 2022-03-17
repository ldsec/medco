//go:build integration_test
// +build integration_test

package survivalserver

import (
	"testing"

	utilserver "github.com/ldsec/medco/connector/util/server"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func init() {
	utilserver.SetForTesting()
}

func TestGetTableName(t *testing.T) {

	expected := "e2etest"
	utilserver.TestI2B2DBConnection(t)
	res, err := getTableName("SPHN")
	assert.NoError(t, err)
	assert.Equal(t, expected, res)

	_, err = getTableName("this table does not exist")
	assert.Error(t, err)
}

func TestGetCodes(t *testing.T) {
	expectedList := []string{"ENC_ID:1", "ENC_ID:2", "ENC_ID:3"}

	utilserver.TestI2B2DBConnection(t)

	res, err := getConceptCodes("/E2ETEST/e2etest/%")
	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedList, res)

	res, err = getConceptCodes("/E2ETEST/e2etest/")
	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedList, res)

	res, err = getConceptCodes("/E2ETEST/e2etest")
	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedList, res)

}

func TestGetModifierCodes(t *testing.T) {
	expectedList1 := []string{"ENC_ID:4"}
	expectedList2 := []string{"ENC_ID:5"}
	utilserver.TestI2B2DBConnection(t)

	res, err := getModifierCodes(`/E2ETEST/modifiers/%`, `/e2etest/%`)
	assert.NoError(t, err)
	assert.ElementsMatch(t, res, expectedList1)

	res, err = getModifierCodes(`/E2ETEST/modifiers/%`, `/e2etest/1/`)
	assert.NoError(t, err)
	assert.ElementsMatch(t, res, expectedList2)
	res, err = getModifierCodes(`/E2ETEST/modifiers/1/`, `/e2etest/1/`)
	assert.NoError(t, err)
	assert.ElementsMatch(t, res, expectedList2)

}
