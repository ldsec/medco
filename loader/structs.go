package loader

import (
	"time"
	"github.com/lca1/unlynx/lib"
	"gopkg.in/dedis/crypto.v0/abstract"
)

// DATA TYPES

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
var TablePatientDimension map[PatientDimensionPK]PatientDimension

// HeaderPatientDimension contains all the headers for the Patient_Dimension table
var HeaderPatientDimension []string

// PatientDimension table represents a patient in the database
type PatientDimension struct {
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

// OptionalFields table contains the optional fields
type OptionalFields struct {
	ValType	string
	Value 	string
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
func ObservationFactFromString(line string) ObservationFact{
	of := ObservationFact{}
	return of
}

func PatientDimensionFromString(line []string, pk abstract.Point) PatientDimension{
	pdk := PatientDimensionPK{
		PatientNum: line[0],
	}

	pd := PatientDimension{
		VitalStatusCD: line[1],
		BirthDate:     line[2],
		DeathDate:     line[3],
	}

	size := len(line)

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

	TablePatientDimension[pdk] = pd

	return pd
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
