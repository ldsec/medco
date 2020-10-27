package i2b2

import (
	"encoding/xml"
)

// NewOntReqGetCategoriesMessageBody returns a new request object for i2b2 categories (ontology root nodes)
func NewOntReqGetCategoriesMessageBody() Request {
	body := OntReqGetCategoriesMessageBody{}

	body.GetCategories.Hiddens = "false"
	body.GetCategories.Blob = "true"
	body.GetCategories.Synonyms = "false"
	body.GetCategories.Type = "core"

	return NewRequestWithBody(body)
}

// NewOntReqGetChildrenMessageBody returns a new request object for i2b2 children of a node
func NewOntReqGetChildrenMessageBody(parent string) Request {
	body := OntReqGetChildrenMessageBody{}

	body.GetChildren.Hiddens = "false"
	body.GetChildren.Blob = "true"
	body.GetChildren.Synonyms = "false"
	body.GetChildren.Max = "200"
	body.GetChildren.Type = "core"

	body.GetChildren.Parent = parent

	return NewRequestWithBody(body)
}

// NewOntReqGetModifiersMessageBody returns a new request object to get the i2b2 modifiers that apply to the concept path
func NewOntReqGetModifiersMessageBody(self string) Request {
	body := OntReqGetModifiersMessageBody{}

	body.GetModifiers.Hiddens = "false"
	body.GetModifiers.Synonyms = "false"

	body.GetModifiers.Self = self

	return NewRequestWithBody(body)
}

// NewOntReqGetModifierChildrenMessageBody returns a new request object to get the i2b2 modifiers that apply to the concept path
func NewOntReqGetModifierChildrenMessageBody(parent, appliedPath, appliedConcept string) Request {
	body := OntReqGetModifierChildrenMessageBody{}

	body.GetModifierChildren.Blob = "false"
	body.GetModifierChildren.Type = "limited"
	body.GetModifierChildren.Max = "200"
	body.GetModifierChildren.Synonyms = "false"
	body.GetModifierChildren.Hiddens = "false"

	body.GetModifierChildren.Parent = parent
	body.GetModifierChildren.AppliedPath = appliedPath
	body.GetModifierChildren.AppliedConcept = appliedConcept

	return NewRequestWithBody(body)
}

// --- request

type baseMessageBody struct {
	Hiddens  string `xml:"hiddens,attr,omitempty"`
	Synonyms string `xml:"synonyms,attr,omitempty"`
	Type     string `xml:"type,attr,omitempty"`
	Blob     string `xml:"blob,attr,omitempty"`
	Max      string `xml:"max,attr,omitempty"`
}

// OntReqGetCategoriesMessageBody is an i2b2 XML message body for ontology categories request
type OntReqGetCategoriesMessageBody struct {
	XMLName       xml.Name `xml:"message_body"`
	GetCategories struct {
		baseMessageBody
	} `xml:"ontns:get_categories"`
}

// OntReqGetChildrenMessageBody is an i2b2 XML message for ontology children request
type OntReqGetChildrenMessageBody struct {
	XMLName     xml.Name `xml:"message_body"`
	GetChildren struct {
		baseMessageBody
		Parent string `xml:"parent"`
	} `xml:"ontns:get_children"`
}

// OntReqGetModifiersMessageBody is an i2b2 XML message for ontology modifiers request
type OntReqGetModifiersMessageBody struct {
	XMLName      xml.Name `xml:"message_body"`
	GetModifiers struct {
		baseMessageBody
		Self string `xml:"self"`
	} `xml:"ontns:get_modifiers"`
}

// OntReqGetModifierChildrenMessageBody is an i2b2 XML message for ontology modifier children request
type OntReqGetModifierChildrenMessageBody struct {
	XMLName             xml.Name `xml:"message_body"`
	GetModifierChildren struct {
		baseMessageBody
		Parent         string `xml:"parent"`
		AppliedPath    string `xml:"applied_path"`
		AppliedConcept string `xml:"applied_concept"`
	} `xml:"ontns:get_modifier_children"`
}

// --- response

// OntRespConceptsMessageBody is the message_body of the i2b2 get_children response message
type OntRespConceptsMessageBody struct {
	XMLName  xml.Name  `xml:"message_body"`
	Concepts []Concept `xml:"concepts>concept"`
}

// Concept is an i2b2 XML concept
type Concept struct {
	Level            string      `xml:"level"`
	Key              string      `xml:"key"`
	Name             string      `xml:"name"`
	SynonymCd        string      `xml:"synonym_cd"`
	Visualattributes string      `xml:"visualattributes"`
	Totalnum         string      `xml:"totalnum"`
	Basecode         string      `xml:"basecode"`
	Metadataxml      Metadataxml `xml:"metadataxml"`
	Facttablecolumn  string      `xml:"facttablecolumn"`
	Tablename        string      `xml:"tablename"`
	Columnname       string      `xml:"columnname"`
	Columndatatype   string      `xml:"columndatatype"`
	Operator         string      `xml:"operator"`
	Dimcode          string      `xml:"dimcode"`
	Comment          string      `xml:"comment"`
	Tooltip          string      `xml:"tooltip"`
	UpdateDate       string      `xml:"update_date"`
	DownloadDate     string      `xml:"download_date"`
	ImportDate       string      `xml:"import_date"`
	SourcesystemCd   string      `xml:"sourcesystem_cd"`
	ValuetypeCd      string      `xml:"valuetype_cd"`
}

// Modifier is an i2b2 XML modifier
type Modifier struct {
	Concept
	AppliedPath string `xml:"applied_path"`
}

// OntRespModifiersMessageBody is the message_body of the i2b2 get_modifiers response message
type OntRespModifiersMessageBody struct {
	XMLName   xml.Name   `xml:"message_body"`
	Modifiers []Modifier `xml:"modifiers>modifier"`
}

// Metadataxml is an i2b2 XML metadata entity
type Metadataxml struct {
	XMLName       xml.Name `xml:"metadataxml"`
	ValueMetadata struct {
		EncryptedType      string `xml:"EncryptedType"`
		NodeEncryptID      string `xml:"NodeEncryptID"`
		ChildrenEncryptIDs string `xml:"ChildrenEncryptIDs"`

		// todo: other elements not unmarshaled, add it to add support for other types
	} `xml:"ValueMetadata"`
}
