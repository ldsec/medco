package i2b2

import (
	"encoding/xml"
)

// NewCrcPdoReqFromInputList returns a new request object for i2b2 pdo request
func NewCrcPdoReqFromInputList(patientSetID string) Request {

	// PDO header
	pdoHeader := PdoHeader{
		PatientSetLimit: "0",
		EstimatedTime:   "0",
		RequestType:     "getPDO_fromInputList",
	}

	// PDO request
	pdoRequest := PdoRequestFromInputList{
		Type: "crcpdons:GetPDOFromInputList_requestType",
		Xsi:  "http://www.w3.org/2001/XMLSchema-instance",
	}

	// set request for patient set ID
	pdoRequest.InputList.PatientList.Max = "1000000"
	pdoRequest.InputList.PatientList.Min = "0"
	pdoRequest.InputList.PatientList.PatientSetCollID = patientSetID
	pdoRequest.OutputOption.Name = "none"
	pdoRequest.OutputOption.PatientSet.Blob = "false"
	pdoRequest.OutputOption.PatientSet.TechData = "false"
	pdoRequest.OutputOption.PatientSet.OnlyKeys = "false"
	pdoRequest.OutputOption.PatientSet.Select = "using_input_list"

	return NewRequestWithBody(CrcPdoReqFromInputListMessageBody{
		PdoHeader:  pdoHeader,
		PdoRequest: pdoRequest,
	})
}

// --- request

// CrcPdoReqFromInputListMessageBody is an i2b2 XML message body for CRC PDO request from input list
type CrcPdoReqFromInputListMessageBody struct {
	XMLName xml.Name `xml:"message_body"`

	PdoHeader  PdoHeader               `xml:"crcpdons:pdoheader"`
	PdoRequest PdoRequestFromInputList `xml:"crcpdons:request"`
}

// PdoHeader is an i2b2 XML header for PDO requests
type PdoHeader struct {
	PatientSetLimit string `xml:"patient_set_limit"`
	EstimatedTime   string `xml:"estimated_time"`
	RequestType     string `xml:"request_type"`
}

// PdoRequestFromInputList is an i2b2 XML PDO request - from input list
type PdoRequestFromInputList struct {
	Type string `xml:"xsi:type,attr"`
	Xsi  string `xml:"xmlns:xsi,attr"`

	// todo: extend to support more queries

	InputList struct {
		PatientList struct {
			Max              string `xml:"max,attr"`
			Min              string `xml:"min,attr"`
			PatientSetCollID string `xml:"patient_set_coll_id"`
		} `xml:"patient_list,omitempty"`
	} `xml:"input_list"`

	OutputOption struct {
		Name       string `xml:"name,attr"`
		PatientSet struct {
			Select   string `xml:"select,attr"`
			OnlyKeys string `xml:"onlykeys,attr"`
			Blob     string `xml:"blob,attr"`
			TechData string `xml:"techdata,attr"`
		} `xml:"patient_set,omitempty"`
	} `xml:"output_option"`
}

// --- response

// CrcPdoRespMessageBody is an i2b2 XML message body for CRC PDO response
type CrcPdoRespMessageBody struct {
	XMLName xml.Name `xml:"message_body"`

	Response struct {
		Xsi         string `xml:"xsi,attr"`
		Type        string `xml:"type,attr"`
		PatientData struct {
			PatientSet struct {
				Patient []struct {
					PatientID string `xml:"patient_id"`
					Param     []struct {
						Text             string `xml:",chardata"`
						Type             string `xml:"type,attr"`
						ColumnDescriptor string `xml:"column_descriptor,attr"`
						Column           string `xml:"column,attr"`
					} `xml:"param"`
				} `xml:"patient"`
			} `xml:"patient_set"`
		} `xml:"patient_data"`
	} `xml:"response"`
}
