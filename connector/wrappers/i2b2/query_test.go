//go:build integration_test
// +build integration_test

package i2b2

import (
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/ldsec/medco/connector/restapi/models"
	utilserver "github.com/ldsec/medco/connector/util/server"
	"github.com/stretchr/testify/assert"
)

func init() {
	utilserver.I2b2HiveURL = "http://localhost:8090/i2b2/services"
	utilserver.I2b2LoginDomain = "i2b2medcosrv0"
	utilserver.I2b2LoginProject = "MedCo"
	utilserver.I2b2LoginUser = "e2etest"
	utilserver.I2b2LoginPassword = "e2etest"
	utilserver.SetLogLevel("5")
}

// warning: all tests need the dev-local-3nodes medco deployment running locally, loaded with default data

// test ontology search query
func TestGetOntologyRootChildren(t *testing.T) {

	results, err := GetOntologyConceptChildren("/")
	if err != nil {
		t.Fail()
	}
	t.Log(*results[0])
}

func TestGetOntologyConceptChildren(t *testing.T) {

	results, err := GetOntologyConceptChildren("/E2ETEST/e2etest/")
	if err != nil {
		t.Fail()
	}
	t.Log(*results[0].MedcoEncryption)
}

func TestGetOntologyModifiers(t *testing.T) {

	results, err := GetOntologyModifiers("/E2ETEST/e2etest/1/")
	if err != nil {
		t.Fail()
	}
	t.Log(*results[0].MedcoEncryption)
}

func TestGetOntologyModifierChildren(t *testing.T) {

	results, err := GetOntologyModifierChildren("/E2ETEST/modifiers/", "/e2etest/%", "/E2ETEST/e2etest/1/")
	if err != nil {
		t.Fail()
	}
	t.Log(*results[0].MedcoEncryption)
}

func TestGetOntologyConceptInfo(t *testing.T) {

	results, err := GetOntologyConceptInfo("/E2ETEST/e2etest/1/")
	if err != nil || results[0].Metadata.ValueMetadata == nil {
		t.Fail()
	}
	t.Log(*results[0].MedcoEncryption)
}

func TestExecutePsmQuery(t *testing.T) {

	encrypted := true
	queryTerm := `/SENSITIVE_TAGGED/medco/tagged/8d3533369426ae172271e98cef8be2bbfe9919087c776083b1ea1de803fc87aa/`
	item := &models.PanelConceptItemsItems0{
		Encrypted: &encrypted,
		QueryTerm: &queryTerm,
	}

	not := false
	patientCount, patientSetID, err := ExecutePsmQuery(
		"testQuery",
		[]*models.Panel{
			{ConceptItems: []*models.PanelConceptItemsItems0{
				item,
			},
				Not: &not,
			}},
		nil,
		nil,
		models.TimingAny)

	if err != nil {
		t.Fail()
	}
	t.Log("count:"+patientCount, "set ID:"+patientSetID)
}

func TestExecutePsmQueryWithValue(t *testing.T) {

	encrypted := false
	queryTerm := `/E2ETEST/e2etest/1/`

	item := &models.PanelConceptItemsItems0{
		Encrypted: &encrypted,
		QueryTerm: &queryTerm,
		Operator:  "EQ",
		Type:      "NUMBER",
		Value:     "10",
	}

	not := false
	patientCount, patientSetID, err := ExecutePsmQuery(
		"testQuery",
		[]*models.Panel{
			{ConceptItems: []*models.PanelConceptItemsItems0{
				item,
			},
				Not: &not,
			}},
		nil,
		nil,
		models.TimingAny)

	if err != nil {
		t.Fail()
	}
	t.Log("count:"+patientCount, "set ID:"+patientSetID)
}

func TestExecutePsmQueryWithModifiers(t *testing.T) {

	encrypted := false
	queryTerm := `/E2ETEST/e2etest/1/`
	appliedPath := `/e2etest/1/`
	modifierKey := `/E2ETEST/modifiers/1/`
	modifier := models.PanelConceptItemsItems0Modifier{
		AppliedPath: &appliedPath,
		ModifierKey: &modifierKey,
	}

	item := &models.PanelConceptItemsItems0{
		Encrypted: &encrypted,
		QueryTerm: &queryTerm,
		Modifier:  &modifier,
	}

	not := false
	patientCount, patientSetID, err := ExecutePsmQuery(
		"testQuery",
		[]*models.Panel{
			{ConceptItems: []*models.PanelConceptItemsItems0{
				item,
			},
				Not: &not,
			}},
		nil,
		nil,
		models.TimingAny)

	if err != nil {
		t.Fail()
	}
	t.Log("count:"+patientCount, "set ID:"+patientSetID)

	// testing with modifier folder -------
	queryTerm = `/E2ETEST/e2etest/3/`
	appliedPath = `/e2etest/%`
	modifierKey = `/E2ETEST/modifiers/`
	modifier = models.PanelConceptItemsItems0Modifier{
		AppliedPath: &appliedPath,
		ModifierKey: &modifierKey,
	}

	item = &models.PanelConceptItemsItems0{
		Encrypted: &encrypted,
		QueryTerm: &queryTerm,
		Modifier:  &modifier,
	}

	patientCount, patientSetID, err = ExecutePsmQuery(
		"testQuery",
		[]*models.Panel{
			{ConceptItems: []*models.PanelConceptItemsItems0{
				item,
			},
				Not: &not,
			}},
		nil,
		nil,
		models.TimingAny)

	if err != nil {
		t.Fail()
	}
	t.Log("count:"+patientCount, "set ID:"+patientSetID)
}

func TestExecutePsmQueryWithModifierAndValue(t *testing.T) {

	encrypted := false
	queryTerm := `/E2ETEST/e2etest/1/`

	appliedPath := `/e2etest/1/`
	modifierKey := `/E2ETEST/modifiers/1/`
	modifier := &models.PanelConceptItemsItems0Modifier{
		AppliedPath: &appliedPath,
		ModifierKey: &modifierKey,
	}

	item := &models.PanelConceptItemsItems0{
		Encrypted: &encrypted,
		QueryTerm: &queryTerm,
		Operator:  "EQ",
		Value:     "15",
		Modifier:  modifier,
	}

	not := false
	patientCount, patientSetID, err := ExecutePsmQuery(
		"testQuery",
		[]*models.Panel{
			{ConceptItems: []*models.PanelConceptItemsItems0{
				item,
			},
				Not: &not,
			}},
		nil,
		nil,
		models.TimingAny)

	if err != nil {
		t.Fail()
	}
	t.Log("count:"+patientCount, "set ID:"+patientSetID)
}

func TestExecutePsmQueryWithSequence(t *testing.T) {

	encrypted1 := false
	queryTerm1 := `/E2ETEST/e2etest/1/`

	item1 := &models.PanelConceptItemsItems0{
		Encrypted: &encrypted1,
		QueryTerm: &queryTerm1,
	}

	encrypted2 := false
	queryTerm2 := `/E2ETEST/e2etest/2/`

	item2 := &models.PanelConceptItemsItems0{
		Encrypted: &encrypted2,
		QueryTerm: &queryTerm2,
	}

	when := models.TimingSequenceInfoWhenLESS
	whichDateFirst := models.TimingSequenceInfoWhichDateFirstSTARTDATE
	whichDateSecond := models.TimingSequenceInfoWhichDateSecondSTARTDATE
	whichObservationFirst := models.TimingSequenceInfoWhichObservationFirstFIRST
	whichObservationSecond := models.TimingSequenceInfoWhichObservationSecondFIRST

	not := false
	selectingPanels := []*models.Panel{
		{ConceptItems: []*models.PanelConceptItemsItems0{
			item1,
		},
			Not: &not,
		}}
	sequencePanels := []*models.Panel{
		{ConceptItems: []*models.PanelConceptItemsItems0{
			item1,
		},
			Not: &not,
		},
		{ConceptItems: []*models.PanelConceptItemsItems0{
			item2,
		},
			Not: &not,
		}}
	timingSequenceOperators := []*models.TimingSequenceInfo{
		{When: &when,
			WhichDateFirst:         &whichDateFirst,
			WhichDateSecond:        &whichDateSecond,
			WhichObservationFirst:  &whichObservationFirst,
			WhichObservationSecond: &whichObservationSecond},
	}
	patientCount, patientSetID, err := ExecutePsmQuery(
		"testQuery",
		selectingPanels,
		timingSequenceOperators,
		sequencePanels,

		models.TimingAny)

	assert.NoError(t, err)
	t.Log("count:"+patientCount, "set ID:"+patientSetID)

	// not a correct number of sequence panels for the number of sequence operators
	patientCount, patientSetID, err = ExecutePsmQuery(
		"testQuery",
		sequencePanels,
		timingSequenceOperators,
		selectingPanels,

		models.TimingAny)

	assert.Error(t, err)
	t.Log("count:"+patientCount, "set ID:"+patientSetID)
}

func TestGetPatientSet(t *testing.T) {

	patientIDs, patientDummyFlags, err := GetPatientSet("9", false)
	if err != nil {
		t.Fail()
	}
	t.Log(patientIDs)
	t.Log(patientDummyFlags)
}

func TestGetOntologyElements(t *testing.T) {

	utilserver.SetForTesting()
	utilserver.TestI2B2DBConnection(t)

	result, err := GetOntologyElements("gender", 10)
	assert.NoError(t, err)

	n := 0
	for _, element := range result {
		logrus.Info(element.Path)
		n++
	}
	assert.Equal(t, n, 3)

}

func TestGetOntologyTermInfo(t *testing.T) {
	results, err := GetOntologyConceptInfo("E2ETEST/wrongFormat/")
	assert.Error(t, err)

	results, err = GetOntologyConceptInfo("/E2ETEST/e2etestthisIsNotAnExistingPathForSure/")
	assert.NoError(t, err)
	assert.Empty(t, results)

	results, err = GetOntologyConceptInfo("/E2ETEST/e2etest/")
	assert.NoError(t, err)
	assert.NotEmpty(t, results)
	res := results[0]
	assert.Equal(t, termInfo1.Code, res.Code)
	assert.Equal(t, termInfo1.DisplayName, res.DisplayName)
	assert.Equal(t, *termInfo1.Leaf, *res.Leaf)
	assert.Equal(t, termInfo1.Name, res.Name)
	assert.Equal(t, termInfo1.Path, res.Path)
	assert.Equal(t, termInfo1.Type, res.Type)
	assert.NotNil(t, res.MedcoEncryption)

}

func TestGetOntologyModifierInfo(t *testing.T) {
	results, err := GetOntologyModifierInfo("/E2ETEST/modifiers/1/", "/e2etest/1/")

	assert.NoError(t, err)
	assert.NotEmpty(t, results)
	assert.NotEqual(t, nil, results[0].Metadata.ValueMetadata)

	t.Log(*results[0].MedcoEncryption)

	results, err = GetOntologyModifierInfo("/E2ETEST/modifiersNotAnExistingModifier/", "/e2etest/%")

	assert.Empty(t, results)

	results, err = GetOntologyModifierInfo("/E2ETEST/modifiers/", "/e2etest/%")
	assert.NoError(t, err)
	assert.NotEmpty(t, results)

	res1 := results[0]
	assert.Equal(t, modifierInfo1.Code, res1.Code)
	assert.Equal(t, modifierInfo1.DisplayName, res1.DisplayName)
	assert.Equal(t, *modifierInfo1.Leaf, *res1.Leaf)
	assert.Equal(t, modifierInfo1.Name, res1.Name)
	assert.Equal(t, modifierInfo1.Path, res1.Path)
	assert.Equal(t, modifierInfo1.AppliedPath, res1.AppliedPath)
	assert.Equal(t, modifierInfo1.Type, res1.Type)
	assert.NotNil(t, res1.MedcoEncryption)

	// If the path is found and the applied path has no match, modifiers are returned regardless of the applied path
	results, err = GetOntologyModifierInfo("/E2ETEST/modifiers/", "/e2etest/1/")
	assert.NoError(t, err)
	assert.NotEmpty(t, results)
	res2 := results[0]
	assert.Equal(t, res1, res2)

	results, err = GetOntologyModifierInfo("/E2ETEST/modifiers/1/", "/e2etest/2/")
	assert.NoError(t, err)
	assert.NotEmpty(t, results)
	res3 := results[0]
	assert.Equal(t, modifierInfo2.Code, res3.Code)
	assert.Equal(t, modifierInfo2.DisplayName, res3.DisplayName)
	assert.Equal(t, *modifierInfo2.Leaf, *res3.Leaf)
	assert.Equal(t, modifierInfo2.Name, res3.Name)
	assert.Equal(t, modifierInfo2.Path, res3.Path)
	assert.Equal(t, modifierInfo2.AppliedPath, res3.AppliedPath)
	assert.Equal(t, modifierInfo2.Type, res3.Type)
	assert.NotNil(t, res3.MedcoEncryption)

}

var falseBool bool = false
var trueBool bool = true

var termInfo1 = &models.ExploreSearchResultElement{
	Code:        "",
	DisplayName: "End-To-End Test",
	Leaf:        &falseBool,
	Name:        "End-To-End Test",
	Path:        "/E2ETEST/e2etest/",
	Type:        "concept_container",
}

var modifierInfo1 = &models.ExploreSearchResultElement{
	Code:        "ENC_ID:4",
	DisplayName: "E2E Modifiers test",
	Leaf:        &falseBool,
	Name:        "E2E Modifiers test",
	Path:        "/E2ETEST/modifiers/",
	AppliedPath: "/e2etest/%",
	Type:        "modifier_folder",
}

var modifierInfo2 = &models.ExploreSearchResultElement{
	Code:        "ENC_ID:5",
	DisplayName: "E2E Modifier 1",
	Leaf:        &trueBool,
	Name:        "E2E Modifier 1",
	Path:        "/E2ETEST/modifiers/1/",
	AppliedPath: "/e2etest/1/",
	Type:        "modifier",
}
