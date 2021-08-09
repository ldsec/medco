package i2b2

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ldsec/medco/connector/restapi/models"
	utilserver "github.com/ldsec/medco/connector/util/server"
)

// NewCrcPsmReqFromQueryDef returns a new request object for i2b2 psm request
func NewCrcPsmReqFromQueryDef(queryName string, queryPanels []*models.Panel, querySequences []*models.TimingSequenceInfo, resultOutputs []ResultOutputName, queryTiming models.Timing) (Request, error) {

	// PSM header
	psmHeader := PsmHeader{
		PatientSetLimit: "0",
		EstimatedTime:   "0",
		QueryMode:       "optimize_without_temp_table",
		RequestType:     "CRC_QRY_runQueryInstance_fromQueryDefinition",
	}
	psmHeader.User.Text = utilserver.I2b2LoginUser
	psmHeader.User.Group = utilserver.I2b2LoginDomain
	psmHeader.User.Login = utilserver.I2b2LoginUser

	// PSM request

	psmRequest := PsmRequestFromQueryDef{
		Type: "crcpsmns:query_definition_requestType",
		Xsi:  "http://www.w3.org/2001/XMLSchema-instance",

		QueryName:        queryName,
		QueryID:          queryName,
		QueryDescription: "Query from MedCo connector (" + queryName + ")",
		QueryTiming:      strings.ToUpper(string(queryTiming)),
		SpecificityScale: "0",
	}

	var panelsOnlyConcepts []Panel

	// embed query in request
	for p, queryPanel := range queryPanels {

		invert := "0"
		if *queryPanel.Not {
			invert = "1"
		}

		i2b2Panel := Panel{
			PanelNumber:          strconv.Itoa(p + 1),
			PanelAccuracyScale:   "100",
			Invert:               invert,
			PanelTiming:          strings.ToUpper(string(queryPanel.PanelTiming)),
			TotalItemOccurrences: "1",
		}
		i2b2PanelNoCohort := i2b2Panel

		for _, queryItem := range queryPanel.ConceptItems {

			i2b2Item := Item{
				ItemKey: convertPathToI2b2Format(*queryItem.QueryTerm),
			}
			if queryItem.Operator != "" && queryItem.Modifier == nil {
				i2b2Item.ConstrainByValue = &ConstrainByValue{
					ValueType:       queryItem.Type,
					ValueOperator:   queryItem.Operator,
					ValueConstraint: queryItem.Value,
				}
			}
			if queryItem.Modifier != nil {
				i2b2Item.ConstrainByModifier = &ConstrainByModifier{
					AppliedPath: strings.ReplaceAll(*queryItem.Modifier.AppliedPath, "/", `\`),
					ModifierKey: convertPathToI2b2Format(*queryItem.Modifier.ModifierKey),
				}
				if queryItem.Operator != "" {
					i2b2Item.ConstrainByModifier.ConstrainByValue = &ConstrainByValue{
						ValueType:       queryItem.Type,
						ValueOperator:   queryItem.Operator,
						ValueConstraint: queryItem.Value,
					}
				}
			}
			i2b2Panel.Items = append(i2b2Panel.Items, i2b2Item)
			i2b2PanelNoCohort.Items = append(i2b2Panel.Items, i2b2Item)
		}
		panelsOnlyConcepts = append(panelsOnlyConcepts, i2b2PanelNoCohort)

		for _, cohort := range queryPanel.CohortItems {

			i2b2Item := Item{
				ItemKey: cohort,
			}
			i2b2Panel.Items = append(i2b2Panel.Items, i2b2Item)
		}

		psmRequest.Panels = append(psmRequest.Panels, i2b2Panel)

	}

	// embed subqueries and subquery constraint if sequences are in use
	if querySequences != nil && len(querySequences) > 0 {
		if len(querySequences)+1 != len(panelsOnlyConcepts) {
			err := fmt.Errorf("the number of items in query sequence info + 1 is not equal to this of panels")
			return NewRequest(), err
		}

		for i, queryPanel := range panelsOnlyConcepts {
			subQueryStringID := queryName + "_SUBQUERY_" + strconv.Itoa(i)
			subquery := Subquery{
				QueryType:   "EVENT",
				QueryName:   subQueryStringID,
				QueryID:     subQueryStringID,
				QueryTiming: "SAMEINSTANCENUM",
				Panels:      []Panel{queryPanel},
			}
			psmRequest.Subqueries = append(psmRequest.Subqueries, subquery)
		}

		for _, querySequence := range querySequences {
			subqueryConstraint := SubqueryConstraint{
				Operator: *querySequence.When,
				FirstQuery: SubqueryConstraintOperand{
					AggregateOperator: *querySequence.WhichObservationFirst,
					JoinColumn:        *querySequence.WhichDateFirst,
				},
				SecondQuery: SubqueryConstraintOperand{
					AggregateOperator: *querySequence.WhichDateSecond,
					JoinColumn:        *querySequence.WhichObservationSecond,
				},
			}
			psmRequest.SubqueryConstraints = append(psmRequest.SubqueryConstraints, subqueryConstraint)
		}
	}

	// embed result outputs
	for i, resultOutput := range resultOutputs {
		psmRequest.ResultOutputs = append(psmRequest.ResultOutputs, ResultOutput{
			PriorityIndex: strconv.Itoa(i + 1),
			Name:          string(resultOutput),
		})

	}

	return NewRequestWithBody(CrcPsmReqFromQueryDefMessageBody{
		PsmHeader:  psmHeader,
		PsmRequest: psmRequest,
	}), nil
}

// --- request

// CrcPsmReqFromQueryDefMessageBody is an i2b2 XML message body for CRC PSM request from query definition
type CrcPsmReqFromQueryDefMessageBody struct {
	XMLName xml.Name `xml:"message_body"`

	PsmHeader  PsmHeader              `xml:"crcpsmns:psmheader"`
	PsmRequest PsmRequestFromQueryDef `xml:"crcpsmns:request"`
}

// PsmHeader is an i2b2 XML header for PSM request
type PsmHeader struct {
	User struct {
		Text  string `xml:",chardata"`
		Group string `xml:"group,attr"`
		Login string `xml:"login,attr"`
	} `xml:"user"`

	PatientSetLimit string `xml:"patient_set_limit"`
	EstimatedTime   string `xml:"estimated_time"`
	QueryMode       string `xml:"query_mode"`
	RequestType     string `xml:"request_type"`
}

// PsmRequestFromQueryDef is an i2b2 XML PSM request from query definition
type PsmRequestFromQueryDef struct {
	Type string `xml:"xsi:type,attr"`
	Xsi  string `xml:"xmlns:xsi,attr"`

	QueryName           string               `xml:"query_definition>query_name"`
	QueryDescription    string               `xml:"query_definition>query_description"`
	QueryID             string               `xml:"query_definition>query_id"`
	QueryTiming         string               `xml:"query_definition>query_timing"`
	SpecificityScale    string               `xml:"query_definition>specificity_scale"`
	Panels              []Panel              `xml:"query_definition>panel"`
	Subqueries          []Subquery           `xml:"query_definition>subquery"`
	SubqueryConstraints []SubqueryConstraint `xml:"query_definition>subquery_constraint"`

	ResultOutputs []ResultOutput `xml:"result_output_list>result_output"`
}

// Panel is an i2b2 XML panel
type Panel struct {
	PanelNumber          string `xml:"panel_number"`
	PanelAccuracyScale   string `xml:"panel_accuracy_scale"`
	Invert               string `xml:"invert"`
	PanelTiming          string `xml:"panel_timing"`
	TotalItemOccurrences string `xml:"total_item_occurrences"`

	Items []Item `xml:"item"`
}

// Item is an i2b2 XML item
type Item struct {
	Hlevel              string               `xml:"hlevel"`
	ItemName            string               `xml:"item_name"`
	ItemKey             string               `xml:"item_key"`
	Tooltip             string               `xml:"tooltip"`
	Class               string               `xml:"class"`
	ConstrainByValue    *ConstrainByValue    `xml:"constrain_by_value,omitempty"`
	ConstrainByModifier *ConstrainByModifier `xml:"constrain_by_modifier,omitempty"`
	ItemIcon            string               `xml:"item_icon"`
	ItemIsSynonym       string               `xml:"item_is_synonym"`
}

// ConstrainByModifier is an i2b2 XML constrain_by_modifier element
type ConstrainByModifier struct {
	AppliedPath      string            `xml:"applied_path"`
	ModifierKey      string            `xml:"modifier_key"`
	ConstrainByValue *ConstrainByValue `xml:"constrain_by_value"`
}

// ConstrainByValue is an i2b2 XML constrain_by_value element
type ConstrainByValue struct {
	ValueType       string `xml:"value_type"`
	ValueOperator   string `xml:"value_operator"`
	ValueConstraint string `xml:"value_constraint"`
}

// Subquery is an i2b2 XML subquery
type Subquery struct {
	QueryType        string  `xml:"query_type"`
	QueryName        string  `xml:"query_name"`
	QueryDescription string  `xml:"query_description"`
	QueryID          string  `xml:"query_id"`
	QueryTiming      string  `xml:"query_timing"`
	SpecificityScale string  `xml:"specificity_scale"`
	Panels           []Panel `xml:"panel"`
}

// SubqueryConstraint is an i2b2 XML suquery_constraint
type SubqueryConstraint struct {
	FirstQuery  SubqueryConstraintOperand `xml:"first_query"`
	Operator    string                    `xml:"operator"`
	SecondQuery SubqueryConstraintOperand `xml:"second_query"`
}

// SubqueryConstraintOperand is a helper structure for SubqueryConstraint
type SubqueryConstraintOperand struct {
	QueryID           string `xml:"query_id"`
	JoinColumn        string `xml:"join_column"`
	AggregateOperator string `xml:"aggregate_operator"`
}

// ResultOutput is an i2b2 XML requested result type
type ResultOutput struct {
	PriorityIndex string `xml:"priority_index,attr"`
	Name          string `xml:"name,attr"`
}

// ResultOutputName is an i2b2 XML requested result type value
type ResultOutputName string

// enumerated values of ResultOutputName
const (
	Patientset                 ResultOutputName = "PATIENTSET"
	PatientEncounterSet        ResultOutputName = "PATIENT_ENCOUNTER_SET"
	PatientCountXML            ResultOutputName = "PATIENT_COUNT_XML"
	PatientGenderCountXML      ResultOutputName = "PATIENT_GENDER_COUNT_XML"
	PatientAgeCountXML         ResultOutputName = "PATIENT_AGE_COUNT_XML"
	PatientVitalstatusCountXML ResultOutputName = "PATIENT_VITALSTATUS_COUNT_XML"
	PatientRaceCountXML        ResultOutputName = "PATIENT_RACE_COUNT_XML"
)

// --- response

// CrcPsmRespMessageBody is an i2b2 XML message body for CRC PSM response
type CrcPsmRespMessageBody struct {
	XMLName xml.Name `xml:"message_body"`

	Response struct {
		Type string `xml:"type,attr"`

		Status []struct {
			Text string `xml:",chardata"`
			Type string `xml:"type,attr"`
		} `xml:"status>condition"`

		QueryMasters []struct {
			QueryMasterID string `xml:"query_master_id"`
			Name          string `xml:"name"`
			UserID        string `xml:"user_id"`
			GroupID       string `xml:"group_id"`
			CreateDate    string `xml:"create_date"`
			DeleteDate    string `xml:"delete_date"`
			RequestXML    string `xml:"request_xml"`
			GeneratedSQL  string `xml:"generated_sql"`
		} `xml:"query_master"`

		QueryInstances []struct {
			QueryInstanceID string `xml:"query_instance_id"`
			QueryMasterID   string `xml:"query_master_id"`
			UserID          string `xml:"user_id"`
			GroupID         string `xml:"group_id"`
			BatchMode       string `xml:"batch_mode"`
			StartDate       string `xml:"start_date"`
			EndDate         string `xml:"end_date"`
			QueryStatusType struct {
				StatusTypeID string `xml:"status_type_id"`
				Name         string `xml:"name"`
				Description  string `xml:"description"`
			} `xml:"query_status_type"`
		} `xml:"query_instance"`

		QueryResultInstances []QueryResultInstance `xml:"query_result_instance"`
	} `xml:"response"`
}

func (responseBody *CrcPsmRespMessageBody) checkStatus() error {
	var errorMessages []string
	for _, status := range responseBody.Response.Status {
		if status.Type == "ERROR" || status.Type == "FATAL_ERROR" {
			errorMessages = append(errorMessages, status.Text)
		}
	}

	if len(errorMessages) != 0 {
		return errors.New(strings.Join(errorMessages, "; "))
	}
	return nil
}

// QueryResultInstance is an i2b2 XML query result instance
type QueryResultInstance struct {
	ResultInstanceID string `xml:"result_instance_id"`
	QueryInstanceID  string `xml:"query_instance_id"`
	QueryResultType  struct {
		ResultTypeID string `xml:"result_type_id"`
		Name         string `xml:"name"`
		Description  string `xml:"description"`
	} `xml:"query_result_type"`
	SetSize         string `xml:"set_size"`
	StartDate       string `xml:"start_date"`
	EndDate         string `xml:"end_date"`
	QueryStatusType struct {
		StatusTypeID string `xml:"status_type_id"`
		Name         string `xml:"name"`
		Description  string `xml:"description"`
	} `xml:"query_status_type"`
}

func (instance *QueryResultInstance) checkStatus() error {
	if instance.QueryStatusType.StatusTypeID != "3" {
		return errors.New("i2b2 result instance does not have finished status")
	}
	return nil
}
