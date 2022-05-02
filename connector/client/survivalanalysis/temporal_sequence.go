package survivalclient

import (
	"fmt"

	"github.com/ldsec/medco/connector/restapi/models"
)

const (
	defaultWhichDateFirst         = models.TimingSequenceInfoWhichDateFirstSTARTDATE
	defaultWhichDateSecond        = models.TimingSequenceInfoWhichDateSecondSTARTDATE
	defaultWhichObservationFirst  = models.TimingSequenceInfoWhichObservationFirstFIRST
	defaultWhichObservationSecond = models.TimingSequenceInfoWhichObservationSecondFIRST
	defaultWhen                   = models.TimingSequenceInfoWhenLESS
)

func validateSequenceOfEvents(sequenceElements []*sequenceElement, nOfPanels int) error {
	if nOfSeqElms := len(sequenceElements); (nOfSeqElms > 0) && (nOfSeqElms != nOfPanels-1) {
		return fmt.Errorf("number of temporal sequence elements should be equal to the number of panels - 1: got %d sequence elements and %d panels", nOfSeqElms, nOfPanels)
	}
	return nil
}

// defaultedSequenceOfEvents this let the possibility to use default value in the with an empty array in YAML
// "sequence_of_events: []"
func defaultedSequenceOfEvents(sequenceElements []*sequenceElement, nOfPanels int) (defaultedSequence []*sequenceElement) {
	if (sequenceElements != nil) && (len(sequenceElements) == 0) {
		defaultedSequence = make([]*sequenceElement, nOfPanels-1)
		for i := 0; i < nOfPanels-1; i++ {
			defaultedSequence[i] = &sequenceElement{
				WhichDateFirst:         "startdate",
				WhichDateSecond:        "startdate",
				WhichObservationFirst:  "first",
				WhichObservationSecond: "first",
				When:                   "before",
			}
		}

	} else {
		defaultedSequence = sequenceElements
	}
	return

}

func convertParametersToSequenceInfo(sequenceElements []*sequenceElement) (timingSequenceInfo []*models.TimingSequenceInfo, err error) {
	if len(sequenceElements) == 0 {
		return
	}
	timingSequenceInfo = make([]*models.TimingSequenceInfo, len(sequenceElements))
	for i, elm := range sequenceElements {
		newSequenceInfo := &models.TimingSequenceInfo{
			When:                   new(string),
			WhichDateFirst:         new(string),
			WhichDateSecond:        new(string),
			WhichObservationFirst:  new(string),
			WhichObservationSecond: new(string),
		}

		switch when := elm.When; when {
		case "":
			*newSequenceInfo.When = defaultWhen
		case "before":
			*newSequenceInfo.When = models.TimingSequenceInfoWhenLESS
		case "beforeorsametime":
			*newSequenceInfo.When = models.TimingSequenceInfoWhenLESSEQUAL
		case "sametime":
			*newSequenceInfo.When = models.TimingSequenceInfoWhenEQUAL
		default:
			err = fmt.Errorf(`"%s" is not valid for "when" element of query temporal sequence`, when)
			return
		}
		switch whichDateFirst := elm.WhichDateFirst; whichDateFirst {
		case "":
			*newSequenceInfo.WhichDateFirst = defaultWhichDateFirst
		case "startdate":
			*newSequenceInfo.WhichDateFirst = models.TimingSequenceInfoWhichDateFirstSTARTDATE
		case "enddate":
			*newSequenceInfo.WhichDateFirst = models.TimingSequenceInfoWhichDateFirstENDDATE
		default:
			err = fmt.Errorf(`"%s" is not valid for "which date first" element of query temporal sequence`, whichDateFirst)
			return
		}
		switch whichDateSecond := elm.WhichDateSecond; whichDateSecond {
		case "":
			*newSequenceInfo.WhichDateSecond = defaultWhichDateSecond
		case "startdate":
			*newSequenceInfo.WhichDateSecond = models.TimingSequenceInfoWhichDateSecondSTARTDATE
		case "enddate":
			*newSequenceInfo.WhichDateSecond = models.TimingSequenceInfoWhichDateSecondENDDATE
		default:
			err = fmt.Errorf(`"%s" is not valid for "which date second" element of query temporal sequence`, whichDateSecond)
			return
		}

		switch whichObservationFirst := elm.WhichObservationFirst; whichObservationFirst {
		case "":
			*newSequenceInfo.WhichObservationFirst = defaultWhichObservationFirst
		case "first":
			*newSequenceInfo.WhichObservationFirst = models.TimingSequenceInfoWhichObservationFirstFIRST
		case "any":
			*newSequenceInfo.WhichObservationFirst = models.TimingSequenceInfoWhichObservationFirstANY
		case "last":
			*newSequenceInfo.WhichObservationFirst = models.TimingSequenceInfoWhichObservationFirstLAST
		default:
			err = fmt.Errorf(`"%s" is not valid for "which observation first" element of query temporal sequence`, whichObservationFirst)
			return
		}

		switch whichObservationSecond := elm.WhichObservationSecond; whichObservationSecond {
		case "":
			*newSequenceInfo.WhichObservationSecond = defaultWhichObservationSecond
		case "first":
			*newSequenceInfo.WhichObservationSecond = models.TimingSequenceInfoWhichObservationSecondFIRST
		case "any":
			*newSequenceInfo.WhichObservationSecond = models.TimingSequenceInfoWhichObservationSecondANY
		case "last":
			*newSequenceInfo.WhichObservationSecond = models.TimingSequenceInfoWhichObservationSecondLAST
		default:
			err = fmt.Errorf(`"%s" is not valid for "which observation second" element of query temporal sequence`, whichObservationSecond)
			return
		}

		newSequenceInfo.Spans, err = convertParametersToTimeSpan(elm.Spans)
		if err != nil {
			return
		}

		timingSequenceInfo[i] = newSequenceInfo

	}
	return
}

func convertParametersToTimeSpan(spans []*timeSpan) (timingSequenceSpan []*models.TimingSequenceSpan, err error) {
	if len(spans) == 0 {
		return
	}
	timingSequenceSpan = make([]*models.TimingSequenceSpan, len(spans))
	for i, span := range spans {
		newSpan := &models.TimingSequenceSpan{
			Operator: new(string),
			Value:    new(int64),
			Units:    new(string),
		}
		switch operator := span.Operator; operator {
		case "less":
			*newSpan.Operator = models.TimingSequenceSpanOperatorLESS
		case "lessorequal":
			*newSpan.Operator = models.TimingSequenceSpanOperatorLESSEQUAL
		case "equal":
			*newSpan.Operator = models.TimingSequenceSpanOperatorEQUAL
		case "moreorequal":
			*newSpan.Operator = models.TimingSequenceSpanOperatorGREATEREQUAL
		case "more":
			*newSpan.Operator = models.TimingSequenceSpanOperatorGREATER
		default:
			err = fmt.Errorf(`"%s" is not valid for the operator element of time span in temporal sequence query`, operator)
			return
		}

		*newSpan.Value = span.Value

		switch units := span.Units; units {
		case "hours":
			*newSpan.Units = models.TimingSequenceSpanUnitsHOUR
		case "days":
			*newSpan.Units = models.TimingSequenceSpanUnitsDAY
		case "months":
			*newSpan.Units = models.TimingSequenceSpanUnitsMONTH
		case "years":
			*newSpan.Units = models.TimingSequenceSpanUnitsYEAR
		default:
			err = fmt.Errorf(`"%s" is not valid for the units element of time span in temporal sequence query`, units)
			return

		}
		timingSequenceSpan[i] = newSpan
	}
	return
}
