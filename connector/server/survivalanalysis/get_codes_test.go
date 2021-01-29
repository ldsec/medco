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

func TestGetCodes(t *testing.T) {
	expectedList := []string{"ENC_ID:1", "ENC_ID:2", "ENC_ID:3"}

	utilserver.TestI2B2DBConnection(t)

	res, err := getCodes(`\\e2etest\\%`)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedList, res)

	res, err = getCodes(`\\e2etest\\`)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedList, res)

	res, err = getCodes(`\\e2etest`)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedList, res)

}

func TestGetModifierCodes(t *testing.T) {
	expectedList1 := []string{"ENC_ID:4"}
	expectedList2 := []string{"ENC_ID:5"}
	utilserver.TestI2B2DBConnection(t)

	res, err := getModifierCodes(`\\%`, `\\e2etest\\%`)
	assert.NoError(t, err)
	assert.ElementsMatch(t, res, expectedList1)
	res, err = getModifierCodes(`\\modifiers\\%`, `\\e2etest\\%`)
	assert.NoError(t, err)
	assert.ElementsMatch(t, res, expectedList1)

	res, err = getModifierCodes(`\\%`, `\\e2etest\\1\\`)
	assert.NoError(t, err)
	assert.ElementsMatch(t, res, expectedList2)

	res, err = getModifierCodes(`\\%`, `\\%`)
	assert.NoError(t, err)
	assert.Empty(t, res)

}
