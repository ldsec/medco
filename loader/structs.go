package loader

import (
	"time"
	"github.com/lca1/unlynx/lib"
	"gopkg.in/dedis/crypto.v0/abstract"
	"gopkg.in/dedis/crypto.v0/base64"
	"gopkg.in/dedis/onet.v1/log"
	"strconv"
)

// HELPER STRUCTS

// ListSensitiveConcepts list all the sensitive concepts (paths)
var ListSensitiveConcepts []string

// IDModifier used to assign IDs to the modifiers concepts
var IDModifiers int64

// IDConcepts used to assign IDs to the different concepts
var IDConcepts  int64

// DATA TYPES

// TableShrineOntologyClear is the shrine_ontology table (it maps the concept path to a concept) with only the NON_SENSITIVE concepts (it INCLUDES MODIFIER NON-SENSITIVE concepts)
var TableShrineOntologyClear map[string]*ShrineOntology

// TableShrineOntologyClear is the shrine_ontology table (it maps the concept path to a concept) with only the SENSITIVE concepts (NO MODIFIER SENSITIVE concepts)
var TableShrineOntologyEnc map[string]*ShrineOntology

// TableShrineOntologyClear is the shrine_ontology table (it maps the concept path to a concept) with only the SENSITIVE concepts (it INCLUDES MODIFIER SENSITIVE concepts)
var TableShrineOntologyModifierEnc map[string][]*ShrineOntology

// HeaderShrineOntology contains all the headers for the shrine table
var HeaderShrineOntology []string

// ShrineOntology is the table that contains all concept codes from the shrine ontology
type ShrineOntology struct {
	NodeEncryptID       int64
	ChildrenEncryptIDs  []int64

	HLevel				string
	Fullname			string
	Name 				string
	SynonymCD			string
	VisualAttributes 	string
	TotalNum			string
	BaseCode			string
	MetadataXML			string
	FactTableColumn		string
	Tablename			string
	ColumnName			string
	ColumnDataType      string
	Operator            string
	DimCode 			string
	Comment 			string
	Tooltip 			string
	AdminColumns    	AdministrativeColumns
	ValueTypeCD   		string
	AppliedPath 		string
	ExclusionCD 		string
}

// To CSV text writes the ShrineOntology object in a way that can be added to a .csv file - "","","", etc.
func (so ShrineOntology) ToCSVText() string{
	if so.NodeEncryptID != int64(-1) { // sensitive
		metadata := ""

		if so.VisualAttributes[:1]=="C" { 			// if concept_parent_node
			metadata += "<?xml version=\"\"1.0\"\"?><ValueMetadata><Version>MedCo-0.1</Version><EncryptedType>CONCEPT_PARENT_NODE</EncryptedType>"
		} else if so.VisualAttributes[:1]=="F" { 	// else if concept_internal_node
			metadata += "<?xml version=\"\"1.0\"\"?><ValueMetadata><Version>MedCo-0.1</Version><EncryptedType>CONCEPT_INTERNAL_NODE</EncryptedType><NodeEncryptID>" + strconv.FormatInt(so.NodeEncryptID,10) + "</NodeEncryptId>"
		} else if so.VisualAttributes[:1]=="L" { 	// else if concept_leaf
			metadata += "<?xml version=\"\"1.0\"\"?><ValueMetadata><Version>MedCo-0.1</Version><EncryptedType>CONCEPT_LEAF</EncryptedType><NodeEncryptID>" + strconv.FormatInt(so.NodeEncryptID,10) + "</NodeEncryptId>"
		} else if so.VisualAttributes[:1]=="O" { 	// else if modifier_parent_node
			metadata += "<?xml version=\"\"1.0\"\"?><ValueMetadata><Version>MedCo-0.1</Version><EncryptedType>MODIFIER_PARENT_NODE</EncryptedType>"
		} else if so.VisualAttributes[:1]=="D" { 	// else if modifier_internal_node
			metadata += "<?xml version=\"\"1.0\"\"?><ValueMetadata><Version>MedCo-0.1</Version><EncryptedType>MODIFIER_INTERNAL_NODE</EncryptedType><NodeEncryptID>" + strconv.FormatInt(so.NodeEncryptID,10) + "</NodeEncryptId>"
		} else if so.VisualAttributes[:1]=="R" { 	// else if modifier_leaf
			metadata += "<?xml version=\"\"1.0\"\"?><ValueMetadata><Version>MedCo-0.1</Version><EncryptedType>MODIFIER_LEAF</EncryptedType><NodeEncryptID>" + strconv.FormatInt(so.NodeEncryptID,10) + "</NodeEncryptId>"
		} else if so.VisualAttributes[:1]=="M" {
			log.Fatal("Not supported go fuck yourself!")
		} else {
			log.Fatal("Wrong VisualAttribute")
		}

		if len(so.ChildrenEncryptIDs) > 0 {
			metadata += "<ChildrenEncryptIDs>"
			for _,childID := range so.ChildrenEncryptIDs {
				metadata += "<ChildEncryptID>" + strconv.FormatInt(childID,10) + "</ChildEncryptID>"
			}

			metadata += "</ChildrenEncryptIDs>"
		}
		so.MetadataXML = metadata + "</ValueMetadata>"
	}

	// i do not call the AdminColumns ToCSVText because some fields are empty (damn you shrine)
	return "\"" + so.HLevel + "\"," + "\"" + so.Fullname + "\"," + "\"" + so.Name + "\"," + "\"" + so.SynonymCD + "\"," + "\"" + so.VisualAttributes + "\"," + "\"" + so.TotalNum + "\"," +
		"\"" + so.BaseCode + "\"," + "\"" + so.MetadataXML + "\"," + "\"" + so.FactTableColumn + "\"," + "\"" + so.Tablename + "\"," + "\"" + so.ColumnDataType + "\"," + "\"" + so.Operator + "\"," +
		"\"" + so.DimCode + "\"," + "\"" + so.Comment + "\"," + "\"" + so.Tooltip + "\"," + "\"" + so.AdminColumns.UpdateDate + "\"," + "\"" + so.AdminColumns.DownloadDate + "\"," + "\"" + so.AdminColumns.ImportDate + "\"," +
		"\"" + so.AdminColumns.SourceSystemCD + "\"," + "\"" + so.ValueTypeCD + "\"," + "\"" + so.AppliedPath + "\"," + "\"" + so.ExclusionCD + "\""
}

// TableObservationFact is observation_fact table
var TableObservationFact map[ObservationFactPK]ObservationFact

// ObservationFact is the fact table of the CRC-I2B2 start schema
type ObservationFact struct{
	ValTypeCD		string
	TValChar		string
	NValNum			float64
	ValueFlagCD 	string
	QuantityNum 	float64
	UnitsCD 		string
	EndDate     	time.Time
	LocationCD  	string
	ObservationBlob string
	ConfidenceNum   float64
	AdminColumns    AdministrativeColumns

}

// ObservationFactPK is the primary key of ObservationFact
type ObservationFactPK struct {
	Encounter 		*VisitDimension
	Patient 		*PatientDimension
	Concept   		*ConceptDimension
	Provider		*ProviderDimension
	StartDate		time.Time
	Modifier      	*ModifierDimension
	InstanceNum     int
}

// AdministrativeColumns are a set of columns that exist in every i2b2 table
type AdministrativeColumns struct {
	UpdateDate      string
	DownloadDate	string
	ImportDate      string
	SourceSystemCD	string
	UploadID		string
	TextSearchIndex string
}

// To CSV text writes the AdministrativeColumns object in a way that can be added to a .csv file - "","","", etc.
func (ac AdministrativeColumns) ToCSVText() string{
	return "\"" + ac.UpdateDate + "\"," + "\"" + ac.DownloadDate + "\"," + "\"" + ac.ImportDate + "\"," + "\"" + ac.SourceSystemCD + "\"," + "\"" + ac.UploadID + "\"," + "\"" + ac.TextSearchIndex + "\""
}

// TablePatientDimension is patient_dimension table
var TablePatientDimension map[*PatientDimensionPK]PatientDimension

// HeaderPatientDimension contains all the headers for the Patient_Dimension table
var HeaderPatientDimension []string

// PatientDimension table represents a patient in the database
type PatientDimension struct {
	PK 				*PatientDimensionPK
	VitalStatusCD   string
	BirthDate       string
	DeathDate		string
	OptionalFields	map[string]string
	AdminColumns    AdministrativeColumns
	EncryptedFlag   lib.CipherText
}

// PatientDimensionPK is the primary key of the Patient_Dimension table
type PatientDimensionPK struct {
	PatientNum		string
}

// To CSV text writes the PatientDimensionPK object in a way that can be added to a .csv file - "","","", etc.
func (pdk *PatientDimensionPK) ToCSVText() string{
	return "\"" + pdk.PatientNum + "\""
}

// To CSV text writes the PatientDimension object in a way that can be added to a .csv file - "","","", etc.
func (pd PatientDimension) ToCSVText() string{
	b := pd.EncryptedFlag.ToBytes()
	encodedEncryptedFlag := "\"" + base64.StdEncoding.EncodeToString(b) + "\""

	return pd.PK.ToCSVText() + ",\"" + pd.VitalStatusCD + "\"," + "\"" + pd.BirthDate + "\"," + "\"" + pd.DeathDate + "\"," + OptionalFieldsMapToCSVText(pd.OptionalFields) + "," + pd.AdminColumns.ToCSVText() + "," + encodedEncryptedFlag
}

// OptionalFields table contains the optional fields
type OptionalFields struct {
	ValType	string
	Value 	string
}

func OptionalFieldsMapToCSVText(of map[string]string) string{
	ofString := ""
	for i:=0; i<len(of); i++{
		// +4 because there is on pk field and 3 mandatory fields
		ofString += "\"" + of[HeaderPatientDimension[i+4]] + "\","
	}
	return ofString[:len(ofString)-1]
}

// TableConceptDimension is concept_dimension table
var TableConceptDimension map[ConceptDimensionPK]ConceptDimension

// ConceptDimension table contains one row for each concept
type ConceptDimension struct {
	ConceptCD   string
	NameChar    string
	ConceptBlob string
	AdminColumns    AdministrativeColumns
}

// ConceptDimensionPK is the primary key of the Concept_Dimension table
type ConceptDimensionPK struct {
	ConceptPath string
}

// TableVisitDimension is visit_dimension table
var TableVisitDimension map[VisitDimensionPK]VisitDimension

// VisitDimension table represents the sessions where observations were made
type VisitDimension struct {
	ActiveStatusCD		string
	StartDate			time.Time
	EndDate				time.Time
	OptionalFields	map[string]OptionalFields
	AdminColumns    AdministrativeColumns
}

// VisitDimensionPK is the primary key of the Visit_Dimension table
type VisitDimensionPK struct {
	EncounterNum		int
	PatientDimension 	*PatientDimension
}

// TableProviderDimension is provider_dimension table
var TableProviderDimension map[ProviderDimensionPK]ProviderDimension

// ProviderDimension table represents a physician or provider at an institution
type ProviderDimension struct {
	NameChar		string
	ProviderBlob 	string
	AdminColumns    AdministrativeColumns
}

// ProviderDimensionPK is the primary key of the Provider_Dimension table
type ProviderDimensionPK struct {
	ProviderID		string
	ProviderPath 	string
}

// TableModifierDimension is modifier_dimension table
var TableModifierDimension map[ModifierDimensionPK]ModifierDimension

// ModifierDimension table contains one row for each modifier
type ModifierDimension struct {
	ModifierCD  	string
	NameChar		string
	ModifierBlob 	string
	AdminColumns    AdministrativeColumns
}

// ModifierDimensionPK is the primary key of the Modifier_Dimension table
type ModifierDimensionPK struct {
	ModifierPath	string
}

// TablePatientMapping is patient_mapping table
var TablePatientMapping map[PatientMappingPK]PatientMapping

// PatientMapping table maps the i2b2 patient_num to an encrypted number
type PatientMapping struct {
	Patient				*PatientDimension
	PatientIDEStatus 	string
	ProjectID			string
	AdminColumns    	AdministrativeColumns
}

// PatientMappingPK is the primary key of the Patient_Mapping table
type PatientMappingPK struct {
	PatientIDE		 string
	PatientIDESource string
}

// TableEncounterMapping is encounter_mapping table
var TableEncounterMapping map[EncounterMappingPK]EncounterMapping

// EncounterMapping table maps i2b2 encounter_num to an encrypted number
type EncounterMapping struct {
	Encounter 			*VisitDimension
	PatientIDE			PatientMapping
	EncounterIDEStatus 	string
	AdminColumns    	AdministrativeColumns
}

// EncounterMappingPK is the primary key of the Encounter_Mapping table
type EncounterMappingPK struct {
	EncounterIDE		string
	EncounterIDESource 	string
	ProjectID			string
}

// AdapterMappings is the xml pre-generated struct to parse the AdapterMappings.xml
type AdapterMappings struct {
	Hostname string `xml:"hostname"`
	TimeStamp string`xml:"timestamp"`
	ListEntries []Entry `xml:"mappings>entry"`
}

// Entry is part of the AdapterMappings.xml
type Entry struct {
	Key string `xml:"key"`
	ListLocalKeys []string `xml:"value>local_key"`
}

// SUPPORT FUNCTIONS
func ShrineOntologyFromString(line []string) *ShrineOntology {
	size := len(line)

	ac := AdministrativeColumns{
		UpdateDate:     	line[size-7],
		DownloadDate:		line[size-6],
		ImportDate:      	line[size-5],
		SourceSystemCD:		line[size-4],
	}

	so := &ShrineOntology {
		NodeEncryptID:      int64(-1), //signals that this shrine ontology element is not sensitive so no need for an encrypt ID
		ChildrenEncryptIDs: nil,	   //same thing as before
		HLevel:				line[0],
		Fullname:			line[1],
		Name:				line[2],
		SynonymCD: 			line[3],
		VisualAttributes:	line[4],
		TotalNum:			line[5],
		BaseCode:			line[6],
		MetadataXML: 		line[7],
		FactTableColumn:    line[8],
		Tablename:          line[9],
		ColumnName:         line[10],
		ColumnDataType: 	line[11],
		Operator: 			line[12],
		DimCode: 			line[13],
		Comment: 			line[14],
		Tooltip: 			line[15],
		AdminColumns: 		ac,
		ValueTypeCD:        line[20],
		AppliedPath:        line[21],
		ExclusionCD:        line[22],
	}

	return so
}


func ObservationFactFromString(line string) ObservationFact{
	of := ObservationFact{}
	return of
}

func PatientDimensionFromString(line []string, pk abstract.Point) (*PatientDimensionPK, PatientDimension){
	pdk := &PatientDimensionPK{
		PatientNum: line[0],
	}

	pd := PatientDimension{
		PK: 			pdk,
		VitalStatusCD: 	line[1],
		BirthDate:     	line[2],
		DeathDate:     	line[3],
	}

	size := len(line)

	// optional fields
	of := make(map[string]string)

	for i:=4; i<size-6; i++{
		of[HeaderPatientDimension[i]] = line[i]
	}

	ac := AdministrativeColumns{
		UpdateDate:     	line[size-6],
		DownloadDate:		line[size-5],
		ImportDate:      	line[size-4],
		SourceSystemCD:		line[size-3],
		UploadID:			line[size-2],
		TextSearchIndex: 	line[size-1],
	}

	// TODO: right now we do not have fake patients
	ef := lib.EncryptInt(pk,1)

	pd.OptionalFields = of
	pd.AdminColumns = ac
	pd.EncryptedFlag = *ef

	return pdk, pd
}

func ConceptDimensionFromString(line string) ConceptDimension{
	cd := ConceptDimension{}
	return cd
}

func VisitDimensionFromString(line string) VisitDimension{
	vd := VisitDimension{}
	return vd
}

func ProviderDimensionFromString(line string) ProviderDimension{
	pd := ProviderDimension{}
	return pd
}

func ModifierDimensionFromString(line string) ModifierDimension{
	md := ModifierDimension{}
	return md
}

func PatientMappingFromString(line string) PatientMapping{
	pm := PatientMapping{}
	return pm
}

func EncounterMappingFromString(line string) EncounterMapping{
	em := EncounterMapping{}
	return em
}
