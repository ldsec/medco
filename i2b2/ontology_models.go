package i2b2

import (
	"encoding/xml"
)

// returns a new request object for i2b2 categories (ontology root nodes)
func NewOntReqGetCategoriesRequest() Request {
	body := OntReqGetCategoriesMessageBody{}
	body.GetCategories.Hiddens = "false"
	body.GetCategories.Blob = "true"
	body.GetCategories.Synonyms = "false"
	body.GetCategories.Type = "core"

	return NewRequestWithBody(body)
}

// returns a new request object for i2b2 children of a node
func NewOntReqGetChildrenMessageBody(path string) Request {
	body := OntReqGetChildrenMessageBody{}
	body.GetChildren.Parent = path
	body.GetChildren.Hiddens= "false"
	body.GetChildren.Blob= "true"
	body.GetChildren.Synonyms= "false"
	body.GetChildren.Max= ""
	body.GetChildren.Type= "core"

	return NewRequestWithBody(body)
}

type OntReqGetCategoriesMessageBody struct {
	XMLName       xml.Name `xml:"message_body"`
	GetCategories struct {
		Hiddens  string `xml:"hiddens,attr"`
		Synonyms string `xml:"synonyms,attr"`
		Type     string `xml:"type,attr"`
		Blob     string `xml:"blob,attr"`
	}`xml:"ontns:get_categories"`

}

type OntReqGetChildrenMessageBody struct {
	XMLName     xml.Name `xml:"message_body"`
	GetChildren struct {
		Max      string `xml:"max,attr"`
		Hiddens  string `xml:"hiddens,attr"`
		Synonyms string `xml:"synonyms,attr"`
		Type     string `xml:"type,attr"`
		Blob     string `xml:"blob,attr"`
		Parent   string `xml:"parent"`
	}`xml:"ontns:get_children"`
}

type OntRespConceptsMessageBody struct {
	XMLName  xml.Name `xml:"message_body"`
	Concepts []Concept `xml:"concepts>concept"`
}

type Concept struct {
	Level            string `xml:"level"`
	Key              string `xml:"key"`
	Name             string `xml:"name"`
	SynonymCd        string `xml:"synonym_cd"`
	Visualattributes string `xml:"visualattributes"`
	Totalnum         string `xml:"totalnum"`
	Basecode         string `xml:"basecode"`
	Metadataxml      Metadataxml `xml:"metadataxml"`
	Facttablecolumn  string `xml:"facttablecolumn"`
	Tablename        string `xml:"tablename"`
	Columnname       string `xml:"columnname"`
	Columndatatype   string `xml:"columndatatype"`
	Operator         string `xml:"operator"`
	Dimcode          string `xml:"dimcode"`
	Comment          string `xml:"comment"`
	Tooltip          string `xml:"tooltip"`
	UpdateDate       string `xml:"update_date"`
	DownloadDate     string `xml:"download_date"`
	ImportDate       string `xml:"import_date"`
	SourcesystemCd   string `xml:"sourcesystem_cd"`
	ValuetypeCd      string `xml:"valuetype_cd"`
}

type Metadataxml struct {
	XMLName       xml.Name `xml:"metadataxml"`
	ValueMetadata struct {
		EncryptedType string `xml:"EncryptedType"`
		NodeEncryptID string `xml:"NodeEncryptID"`
		ChildrenEncryptIDs string `xml:"ChildrenEncryptIDs"`

		// todo: other elements not unmarshaled, add it to add support for other types
	} `xml:"ValueMetadata"`
}
