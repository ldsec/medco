package survivalclient

import (
	"os"
	"testing"

	"github.com/ldsec/medco/connector/restapi/client/survival_analysis"
	"github.com/ldsec/medco/connector/restapi/models"
	utilclient "github.com/ldsec/medco/connector/util/client"
	"github.com/stretchr/testify/assert"
)

const testParamFile string = "../../../test/survival_unit_test_parameters.yaml"

func TestNewParametersFromFile(t *testing.T) {
	file, err := os.Open(testParamFile)
	if err != nil {
		if os.IsNotExist(err) {
			t.Skip("test parameter files does not exiss")
		} else {
			t.Error(err)
		}
	}
	err = file.Close()
	assert.NoError(t, err)

	_, err = NewParametersFromFile("./thisFileDoesNotExist")
	assert.Error(t, err)

	params, err := NewParametersFromFile(testParamFile)
	assert.NoError(t, err)
	assert.Equal(t, parameters, params)

}

func TestConversionOfParameters(t *testing.T) {
	convertedParams, err := convertParametersToSubGroupDefinition(parameters)
	assert.NoError(t, err)
	assert.Equal(t, parsedParameters, convertedParams)
}

var parameters = &Parameters{
	TimeResolution:   "day",
	TimeLimit:        19,
	CohortName:       "anyName",
	StartConceptPath: "/any/start/path/",
	StartModifier: &modifier{
		ModifierKey: "/any/start/modifier/key/",
		AppliedPath: "/any/start/path/%",
	},
	StartsWhen:     "earliest",
	EndConceptPath: "/any/end/path/",
	EndModifier: &modifier{
		ModifierKey: "/any/end/modifier/key/",
		AppliedPath: "/any/end/path/%",
	},
	EndsWhen: "earliest",
	SubGroups: []*subGroup{
		{
			GroupName: "AAA",
			Panels: []*panel{
				{
					PanelTiming: "any",
					Not:         false,
					Items: []*item{
						{
							Path: "/path/1/",
						},
						{
							Path: "/path/2/",
							Modifier: &modifier{
								ModifierKey: "/key1/",
								AppliedPath: "/appliedpath1/",
							},
						},
					},
				},
				{
					PanelTiming: "sameinstancenum",
					Not:         true,
					Items: []*item{
						{
							Path: "/path/3/",
						},
					},
				},
			},
			SequenceOfEvents: []*sequenceElement{
				{
					WhichDateFirst:         "startdate",
					WhichDateSecond:        "enddate",
					WhichObservationFirst:  "any",
					WhichObservationSecond: "last",
					When:                   "sametime",
					Spans: []*timeSpan{
						{Operator: "less",
							Value: 34,
							Units: "years"},
						{Operator: "equal",
							Value: 21,
							Units: "days"},
					},
				},
			},
		},
		{
			GroupName: "BBB",
			Panels: []*panel{
				{
					Not:         false,
					PanelTiming: "any",
					Items: []*item{
						{
							Path: "/path/4/",
						},
					},
				},
				{
					Not:         true,
					PanelTiming: "sameinstancenum",
					Items: []*item{
						{
							Path: "/path/3/",
						},
					},
				},
			},
			SequenceOfEvents: []*sequenceElement{
				{},
			},
		},
		{
			GroupName:   "CCC",
			GroupTiming: "sameinstancenum",
			Panels: []*panel{
				{
					Not:         false,
					PanelTiming: "any",
					Items: []*item{
						{
							Path: "/path/4/",
						},
					},
				},
			},
		},
		{
			GroupName: "DDD",
			Panels: []*panel{
				{
					Not: false,
					Items: []*item{
						{
							Path: "/path/4/",
						},
					},
				},
				{
					Not:         true,
					PanelTiming: "sameinstancenum",
					Items: []*item{
						{
							Path: "/path/3/",
						},
					},
				},
				{
					Not:         true,
					PanelTiming: "sameinstancenum",
					Items: []*item{
						{
							Path: "/path/7/",
						},
					},
				},
			},
			SequenceOfEvents: []*sequenceElement{},
		},
	},
}

func defaultSequenceInfo(length int) (ret []*models.TimingSequenceInfo) {
	for i := 0; i < length; i++ {
		ret = append(ret, &models.TimingSequenceInfo{
			When:                   utilclient.InitializeStringPointer(models.TimingSequenceInfoWhenLESS),
			WhichDateFirst:         utilclient.InitializeStringPointer(models.TimingSequenceInfoWhichDateFirstSTARTDATE),
			WhichDateSecond:        utilclient.InitializeStringPointer(models.TimingSequenceInfoWhichDateSecondSTARTDATE),
			WhichObservationFirst:  utilclient.InitializeStringPointer(models.TimingSequenceInfoWhichObservationFirstFIRST),
			WhichObservationSecond: utilclient.InitializeStringPointer(models.TimingSequenceInfoWhichObservationSecondFIRST),
			Spans:                  []*models.TimingSequenceSpan{},
		})
	}
	return
}

var parsedParameters = []*survival_analysis.SurvivalAnalysisParamsBodySubGroupDefinitionsItems0{
	{
		GroupName:      "AAA",
		SubGroupTiming: defaultTiming,
		Panels: []*models.Panel{
			{
				Not:         utilclient.InitializeBoolPointer(false),
				PanelTiming: models.TimingAny,
				ConceptItems: []*models.PanelConceptItemsItems0{
					{
						Encrypted: utilclient.InitializeBoolPointer(false),
						QueryTerm: utilclient.InitializeStringPointer("/path/1/")},
					{
						Encrypted: utilclient.InitializeBoolPointer(false),
						QueryTerm: utilclient.InitializeStringPointer("/path/2/"),
						Modifier: &models.PanelConceptItemsItems0Modifier{
							AppliedPath: utilclient.InitializeStringPointer("/appliedpath1/"),
							ModifierKey: utilclient.InitializeStringPointer("/key1/"),
						},
					},
				},
			},
			{
				Not:         utilclient.InitializeBoolPointer(true),
				PanelTiming: models.TimingSameinstancenum,
				ConceptItems: []*models.PanelConceptItemsItems0{
					{
						Encrypted: utilclient.InitializeBoolPointer(false),
						QueryTerm: utilclient.InitializeStringPointer("/path/3/"),
					},
				},
			},
		},
		QueryTimingSequence: []*models.TimingSequenceInfo{
			{
				When:                   utilclient.InitializeStringPointer(models.TimingSequenceInfoWhenEQUAL),
				WhichDateFirst:         utilclient.InitializeStringPointer(models.TimingSequenceInfoWhichDateFirstSTARTDATE),
				WhichDateSecond:        utilclient.InitializeStringPointer(models.TimingSequenceInfoWhichDateSecondENDDATE),
				WhichObservationFirst:  utilclient.InitializeStringPointer(models.TimingSequenceInfoWhichObservationFirstANY),
				WhichObservationSecond: utilclient.InitializeStringPointer(models.TimingSequenceInfoWhichObservationSecondLAST),
				Spans: []*models.TimingSequenceSpan{
					{
						Operator: utilclient.InitializeStringPointer(models.TimingSequenceInfoWhenLESS),
						Value:    utilclient.InitializeInt64Pointer(34),
						Units:    utilclient.InitializeStringPointer(models.TimingSequenceSpanUnitsYEAR),
					},
					{
						Operator: utilclient.InitializeStringPointer(models.TimingSequenceInfoWhenEQUAL),
						Value:    utilclient.InitializeInt64Pointer(21),
						Units:    utilclient.InitializeStringPointer(models.TimingSequenceSpanUnitsDAY),
					},
				},
			},
		},
	},
	{
		GroupName:      "BBB",
		SubGroupTiming: defaultTiming,
		Panels: []*models.Panel{
			{
				Not:         utilclient.InitializeBoolPointer(false),
				PanelTiming: models.TimingAny,
				ConceptItems: []*models.PanelConceptItemsItems0{
					{
						QueryTerm: utilclient.InitializeStringPointer("/path/4/"),
						Encrypted: utilclient.InitializeBoolPointer(false),
					},
				},
			},
			{
				Not:         utilclient.InitializeBoolPointer(true),
				PanelTiming: models.TimingSameinstancenum,
				ConceptItems: []*models.PanelConceptItemsItems0{
					{
						QueryTerm: utilclient.InitializeStringPointer("/path/3/"),
						Encrypted: utilclient.InitializeBoolPointer(false),
					},
				},
			},
		},
		QueryTimingSequence: defaultSequenceInfo(1),
	},
	{
		GroupName:      "CCC",
		SubGroupTiming: models.TimingSameinstancenum,
		Panels: []*models.Panel{
			{
				Not:         utilclient.InitializeBoolPointer(false),
				PanelTiming: models.TimingAny,
				ConceptItems: []*models.PanelConceptItemsItems0{
					{
						Encrypted: utilclient.InitializeBoolPointer(false),
						QueryTerm: utilclient.InitializeStringPointer("/path/4/"),
					},
				},
			},
		},
		QueryTimingSequence: []*models.TimingSequenceInfo{},
	},
	{
		GroupName:      "DDD",
		SubGroupTiming: defaultTiming,
		Panels: []*models.Panel{
			{
				PanelTiming: defaultTiming,
				Not:         utilclient.InitializeBoolPointer(false),
				ConceptItems: []*models.PanelConceptItemsItems0{
					{
						Encrypted: utilclient.InitializeBoolPointer(false),
						QueryTerm: utilclient.InitializeStringPointer("/path/4/"),
					},
				},
			},
			{
				PanelTiming: models.TimingSameinstancenum,
				Not:         utilclient.InitializeBoolPointer(true),
				ConceptItems: []*models.PanelConceptItemsItems0{
					{
						Encrypted: utilclient.InitializeBoolPointer(false),
						QueryTerm: utilclient.InitializeStringPointer("/path/3/"),
					},
				},
			},
			{
				PanelTiming: models.TimingSameinstancenum,
				Not:         utilclient.InitializeBoolPointer(true),
				ConceptItems: []*models.PanelConceptItemsItems0{
					{
						Encrypted: utilclient.InitializeBoolPointer(false),
						QueryTerm: utilclient.InitializeStringPointer("/path/7/"),
					},
				},
			},
		},
		QueryTimingSequence: defaultSequenceInfo(2),
	},
}
