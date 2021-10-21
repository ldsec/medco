package survivalclient

import (
	"fmt"

	"github.com/ldsec/medco/connector/restapi/models"
)

func convertParametersToSequenceInfo(sequenceElements []*sequenceElement) (timingSequenceInfo []*models.TimingSequenceInfo, err error) {
	timingSequenceInfo = make([]*models.TimingSequenceInfo, len(sequenceElements))
	for i, elm := range sequenceElements {
		newSequenceInfo := &models.TimingSequenceInfo{
			When:                   new(string),
			WhichDateFirst:         new(string),
			WhichDateSecond:        new(string),
			WhichObservationFirst:  new(string),
			WhichObservationSecond: new(string),
		}

		switch when := elm.When; {
		case (when == "") || (when == "before"):
			*newSequenceInfo.When = models.TimingSequenceInfoWhenLESS
		case when == "beforeorsametime":
			*newSequenceInfo.When = models.TimingSequenceInfoWhenLESSEQUAL
		case when == "sametime":
			*newSequenceInfo.When = models.TimingSequenceInfoWhenEQUAL
		default:
			err = fmt.Errorf(`"%s" is not valid for "when" element of query temporal sequence`, when)
			return
		}
		switch whichDateFirst := elm.WhichDateFirst; {
		case (whichDateFirst == "") || (whichDateFirst == "startdate"):
			*newSequenceInfo.WhichDateFirst = models.TimingSequenceInfoWhichDateFirstSTARTDATE
		case whichDateFirst == "enddate":
			*newSequenceInfo.WhichDateFirst = models.TimingSequenceInfoWhichDateFirstENDDATE
		default:
			err = fmt.Errorf(`"%s" is not valid for "which date first" element of query temporal sequence`, whichDateFirst)
			return
		}
		switch whichDateSecond := elm.WhichDateSecond; {
		case (whichDateSecond == "") || (whichDateSecond == "startdate"):
			*newSequenceInfo.WhichDateSecond = models.TimingSequenceInfoWhichDateSecondSTARTDATE
		case whichDateSecond == "enddate":
			*newSequenceInfo.WhichDateSecond = models.TimingSequenceInfoWhichDateSecondENDDATE
		default:
			err = fmt.Errorf(`"%s" is not valid for "which date second" element of query temporal sequence`, whichDateSecond)
			return
		}

		switch whichObservationFirst := elm.WhichObservationFirst; {
		case (whichObservationFirst == "") || (whichObservationFirst == "first"):
			*newSequenceInfo.WhichObservationFirst = models.TimingSequenceInfoWhichObservationFirstFIRST
		case whichObservationFirst == "any":
			*newSequenceInfo.WhichObservationFirst = models.TimingSequenceInfoWhichObservationFirstANY
		case whichObservationFirst == "last":
			*newSequenceInfo.WhichObservationFirst = models.TimingSequenceInfoWhichObservationFirstLAST
		default:
			err = fmt.Errorf(`"%s" is not valid for "which observation first" element of query temporal sequence`, whichObservationFirst)
			return
		}

		switch whichObservationSecond := elm.WhichObservationSecond; {
		case (whichObservationSecond == "") || (whichObservationSecond == "first"):
			*newSequenceInfo.WhichObservationSecond = models.TimingSequenceInfoWhichObservationSecondFIRST
		case whichObservationSecond == "any":
			*newSequenceInfo.WhichObservationSecond = models.TimingSequenceInfoWhichObservationSecondANY
		case whichObservationSecond == "last":
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
