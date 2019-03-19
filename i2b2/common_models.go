package i2b2

import (
	"encoding/xml"
	"github.com/lca1/medco-connector/util"
	"strconv"
	"time"
)

// create a new ready-to-use i2b2 request, with a nil message body
func NewRequest() Request {
	now := time.Now()
	return Request{
		XMLNSMSG:    "http://www.i2b2.org/xsd/hive/msg/1.1/",
		XMLNSONT:    "http://www.i2b2.org/xsd/cell/ont/1.1/",
		XMLNSPDO:    "http://www.i2b2.org/xsd/hive/pdo/1.1/",
		XMLNSCRCPDO: "http://www.i2b2.org/xsd/cell/crc/pdo/1.1/",
		XMLNSCRCPSM: "http://www.i2b2.org/xsd/cell/crc/psm/1.1/",

		MessageHeader: MessageHeader{
			I2b2VersionCompatible: "0.3",
			Hl7VersionCompatible: "2.4",
			SendingApplicationApplicationName: "MedCo Connector",
			SendingApplicationApplicationVersion: "0.2",
			SendingFacilityFacilityName: "MedCo",
			ReceivingApplicationApplicationName: "i2b2 cell",
			ReceivingApplicationApplicationVersion: "1.7",
			ReceivingFacilityFacilityName: "i2b2 hive",
			DatetimeOfMessage: now.Format(time.RFC3339),
			SecurityDomain: util.I2b2LoginDomain(),
			SecurityUsername: util.I2b2LoginUser(),
			SecurityPassword: util.I2b2LoginPassword(),
			MessageTypeMessageCode: "EQQ",
			MessageTypeEventType: "Q04",
			MessageTypeMessageStructure: "EQQ_Q04",
			MessageControlIDSessionID: now.Format(time.RFC3339),
			MessageControlIDMessageNum: strconv.FormatInt(now.Unix(), 10),
			MessageControlIDInstanceNum: "0",
			ProcessingIDProcessingID: "P",
			ProcessingIDProcessingMode: "I",
			AcceptAcknowledgementType: "messageId",
			ApplicationAcknowledgementType: "",
			CountryCode: "CH",
			ProjectID: util.I2b2LoginProject(),
		},
		RequestHeader: RequestHeader{
			ResultWaittimeMs: strconv.Itoa(util.I2b2TimeoutSeconds() * 1000),
		},
	}
}

// create a new ready-to-use i2b2 request, with a message body
func NewRequestWithBody(body MessageBody) (req Request) {
	req = NewRequest()
	req.MessageBody = body
	return
}

type Request struct {
	XMLName        xml.Name `xml:"msgns:request"`
	XMLNSMSG       string `xml:"xmlns:msgns,attr"`
	XMLNSPDO       string `xml:"xmlns:pdons,attr"`
	XMLNSONT       string `xml:"xmlns:ontns,attr"`
	XMLNSCRCPDO    string `xml:"xmlns:crcpdons,attr"`
	XMLNSCRCPSM    string `xml:"xmlns:crcpsmns,attr"`

	MessageHeader  MessageHeader `xml:"message_header"`
	RequestHeader  RequestHeader `xml:"request_header"`
	MessageBody    MessageBody   `xml:"message_body"`
}

type Response struct {
	XMLName        xml.Name `xml:"response"`
	MessageHeader  MessageHeader `xml:"message_header"`
	RequestHeader  RequestHeader `xml:"request_header"`
	ResponseHeader ResponseHeader `xml:"response_header"`
	MessageBody    MessageBody   `xml:"message_body"`
}

type MessageHeader struct {
	XMLName               xml.Name `xml:"message_header"`

	I2b2VersionCompatible string   `xml:"i2b2_version_compatible"`
	Hl7VersionCompatible  string   `xml:"hl7_version_compatible"`

	SendingApplicationApplicationName string `xml:"sending_application>application_name"`
	SendingApplicationApplicationVersion string `xml:"sending_application>application_version"`

	SendingFacilityFacilityName string `xml:"sending_facility>facility_name"`

	ReceivingApplicationApplicationName string `xml:"receiving_application>application_name"`
	ReceivingApplicationApplicationVersion string `xml:"receiving_application>application_version"`

	ReceivingFacilityFacilityName string `xml:"receiving_facility>facility_name"`

	DatetimeOfMessage string `xml:"datetime_of_message"`

	SecurityDomain string `xml:"security>domain"`
	SecurityUsername string `xml:"security>username"`
	SecurityPassword string `xml:"security>password"`

	MessageTypeMessageCode string `xml:"message_type>message_code"`
	MessageTypeEventType string `xml:"message_type>event_type"`
	MessageTypeMessageStructure string `xml:"message_type>message_structure"`

	MessageControlIDSessionID   string `xml:"message_control_id>session_id"`
	MessageControlIDMessageNum  string `xml:"message_control_id>message_num"`
	MessageControlIDInstanceNum string `xml:"message_control_id>instance_num"`

	ProcessingIDProcessingID   string `xml:"processing_id>processing_id"`
	ProcessingIDProcessingMode string `xml:"processing_id>processing_mode"`

	AcceptAcknowledgementType      string `xml:"accept_acknowledgement_type"`
	ApplicationAcknowledgementType string `xml:"application_acknowledgement_type"`
	CountryCode                    string `xml:"country_code"`
	ProjectID                      string `xml:"project_id"`
}

type RequestHeader struct {
	XMLName          xml.Name `xml:"request_header"`
	ResultWaittimeMs string   `xml:"result_waittime_ms"`
}

type ResponseHeader struct {
	XMLName xml.Name `xml:"response_header"`
	Info    struct {
		Text string `xml:",chardata"`
		URL  string `xml:"url,attr"`
	} `xml:"info"`
	ResultStatus struct {
		Status struct {
			Text string `xml:",chardata"`
			Type string `xml:"type,attr"`
		} `xml:"status"`
		PollingURL struct {
			Text       string `xml:",chardata"`
			IntervalMs string `xml:"interval_ms,attr"`
		} `xml:"polling_url"`
		Conditions struct {
			Condition []struct {
				Text         string `xml:",chardata"`
				Type         string `xml:"type,attr"`
				CodingSystem string `xml:"coding_system,attr"`
			} `xml:"condition"`
		} `xml:"conditions"`
	} `xml:"result_status"`
}

type MessageBody interface{}

