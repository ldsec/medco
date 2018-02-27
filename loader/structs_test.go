package loader_test

import (
	"encoding/base64"
	"encoding/csv"
	"github.com/lca1/medco-loader/loader"
	"github.com/lca1/unlynx/lib"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

// ----------------------------------------------------------------------------------------------------------- //
// ---------------------------------------- TO STRING -------------------------------------------------------- //
// ----------------------------------------------------------------------------------------------------------- //

func TestShrineOntology_ToCSVText(t *testing.T) {

	ac := loader.AdministrativeColumns{
		UpdateDate:     "\\N",
		DownloadDate:   "\\N",
		ImportDate:     "\\N",
		SourceSystemCD: "SHRINE",
	}

	so := loader.ShrineOntology{
		NodeEncryptID:      -1,
		ChildrenEncryptIDs: nil,

		HLevel:           "0",
		Fullname:         "\\SHRINE\\",
		Name:             "SHRINE",
		SynonymCD:        "N",
		VisualAttributes: "CA ",
		TotalNum:         "\\N",
		BaseCode:         "\\N",
		MetadataXML:      "",
		FactTableColumn:  "concept_cd",
		Tablename:        "concept_dimension",
		ColumnName:       "concept_path",
		ColumnDataType:   "T",
		Operator:         "LIKE",
		DimCode:          "\\SHRINE\\",
		Comment:          "",
		Tooltip:          "\\N",
		AdminColumns:     ac,
		ValueTypeCD:      "\\N",
		AppliedPath:      "@",
		ExclusionCD:      "\\N",
	}
	assert.Equal(t, so.ToCSVText(), `"0","\SHRINE\","SHRINE","N","CA ","\N","\N","","concept_cd","concept_dimension","concept_path","T","LIKE","\SHRINE\","","\N","\N","\N","\N","SHRINE","\N","@","\N"`)

	so.NodeEncryptID = 1
	assert.Equal(t, so.ToCSVText(), `"0","\SHRINE\","SHRINE","N","CA ","\N","\N","<?xml version=""1.0""?><ValueMetadata><Version>MedCo-0.1</Version><EncryptedType>CONCEPT_PARENT_NODE</EncryptedType></ValueMetadata>","concept_cd","concept_dimension","concept_path","T","LIKE","\SHRINE\","","\N","\N","\N","\N","SHRINE","\N","@","\N"`)

	so.VisualAttributes = "LA "
	assert.Equal(t, so.ToCSVText(), `"0","\SHRINE\","SHRINE","N","LA ","\N","\N","<?xml version=""1.0""?><ValueMetadata><Version>MedCo-0.1</Version><EncryptedType>CONCEPT_LEAF</EncryptedType><NodeEncryptID>1</NodeEncryptID></ValueMetadata>","concept_cd","concept_dimension","concept_path","T","LIKE","\SHRINE\","","\N","\N","\N","\N","SHRINE","\N","@","\N"`)

	so.ChildrenEncryptIDs = append(so.ChildrenEncryptIDs, 2)
	so.ChildrenEncryptIDs = append(so.ChildrenEncryptIDs, 3)
	assert.Equal(t, so.ToCSVText(), `"0","\SHRINE\","SHRINE","N","LA ","\N","\N","<?xml version=""1.0""?><ValueMetadata><Version>MedCo-0.1</Version><EncryptedType>CONCEPT_LEAF</EncryptedType><NodeEncryptID>1</NodeEncryptID></ValueMetadata>","concept_cd","concept_dimension","concept_path","T","LIKE","\SHRINE\","","\N","\N","\N","\N","SHRINE","\N","@","\N"`)

	so.VisualAttributes = "FA "
	assert.Equal(t, so.ToCSVText(), `"0","\SHRINE\","SHRINE","N","FA ","\N","\N","<?xml version=""1.0""?><ValueMetadata><Version>MedCo-0.1</Version><EncryptedType>CONCEPT_INTERNAL_NODE</EncryptedType><NodeEncryptID>1</NodeEncryptID><ChildrenEncryptIDs><ChildEncryptID>2</ChildEncryptID><ChildEncryptID>3</ChildEncryptID></ChildrenEncryptIDs></ValueMetadata>","concept_cd","concept_dimension","concept_path","T","LIKE","\SHRINE\","","\N","\N","\N","\N","SHRINE","\N","@","\N"`)

	so.VisualAttributes = "DA "
	assert.Equal(t, so.ToCSVText(), `"0","\SHRINE\","SHRINE","N","DA ","\N","\N","<?xml version=""1.0""?><ValueMetadata><Version>MedCo-0.1</Version><EncryptedType>MODIFIER_INTERNAL_NODE</EncryptedType><NodeEncryptID>1</NodeEncryptID><ChildrenEncryptIDs><ChildEncryptID>2</ChildEncryptID><ChildEncryptID>3</ChildEncryptID></ChildrenEncryptIDs></ValueMetadata>","concept_cd","concept_dimension","concept_path","T","LIKE","\SHRINE\","","\N","\N","\N","\N","SHRINE","\N","@","\N"`)

}

func TestLocalOntology_ToCSVText(t *testing.T) {

	ac := loader.AdministrativeColumns{
		UpdateDate:     "2007-04-10 00:00:00",
		DownloadDate:   "2007-04-10 00:00:00",
		ImportDate:     "2007-04-10 00:00:00",
		SourceSystemCD: "DEMO",
	}

	lo := loader.LocalOntology{
		HLevel:           "4",
		Fullname:         "\\i2b2\\Demographics\\Zip codes\\Arkansas\\Parkdale\\",
		Name:             "Parkdale",
		SynonymCD:        "N",
		VisualAttributes: "FA ",
		TotalNum:         "\\N",
		BaseCode:         "\\N",
		MetadataXML:      "\\N",
		FactTableColumn:  "concept_cd",
		Tablename:        "concept_dimension",
		ColumnName:       "concept_path",
		ColumnDataType:   "T",
		Operator:         "LIKE",
		DimCode:          "\\i2b2\\Demographics\\Zip codes\\Arkansas\\Parkdale\\",
		Comment:          "\\N",
		Tooltip:          "Demographics \\ Zip codes \\ Arkansas \\ Parkdale",
		AppliedPath:      "@",
		AdminColumns:     ac,
		ValueTypeCD:      "\\N",
		ExclusionCD:      "\\N",
		Path:             "\\N",
		Symbol:           "\\N",

		PCoriBasecode: "\\N",
	}

	assert.Equal(t, lo.ToCSVText(), `"4","\i2b2\Demographics\Zip codes\Arkansas\Parkdale\","Parkdale","N","FA ","\N","\N","\N","concept_cd","concept_dimension","concept_path","T","LIKE","\i2b2\Demographics\Zip codes\Arkansas\Parkdale\","\N","Demographics \ Zip codes \ Arkansas \ Parkdale","@","2007-04-10 00:00:00","2007-04-10 00:00:00","2007-04-10 00:00:00","DEMO","\N","\N","\N","\N"`)

	tag := libunlynx.GroupingKey("1")
	assert.Equal(t, loader.LocalOntologySensitiveConceptToCSVText(&tag, 20), `"3", "\medco\tagged\concept\1\", "", "N", "LA ", "\N", "TAG_ID:20", "\N", "concept_cd", "concept_dimension", "concept_path", "T", "LIKE", "\medco\tagged\concept\1\", "\N", "\N", "NOW()", "\N", "\N", "\N", "TAG_ID", "@", "\N", "\N", "\N", "\N"`)

}

func TestPatientDimension_ToCSVText(t *testing.T) {

	ac := loader.AdministrativeColumns{
		UpdateDate:     "2010-11-04 10:43:00",
		DownloadDate:   "2010-08-18 09:50:00",
		ImportDate:     "2010-11-04 10:43:00",
		SourceSystemCD: "DEMO",
		UploadID:       "\\N",
	}

	pdk := &loader.PatientDimensionPK{
		PatientNum: "1000000001",
	}

	op := make([]loader.OptionalFields, 0)
	op = append(op, loader.OptionalFields{ValType: "sex_cd", Value: "F"})
	op = append(op, loader.OptionalFields{ValType: "age_in_years_num", Value: "24"})
	op = append(op, loader.OptionalFields{ValType: "language_cd", Value: "english"})
	op = append(op, loader.OptionalFields{ValType: "race_cd", Value: "black"})
	op = append(op, loader.OptionalFields{ValType: "marital_status_cd", Value: "married"})
	op = append(op, loader.OptionalFields{ValType: "religion_cd", Value: "roman catholic"})
	op = append(op, loader.OptionalFields{ValType: "zip_cd", Value: "02140"})
	op = append(op, loader.OptionalFields{ValType: "statecityzip_path", Value: "Zip codes\\Massachusetts\\Cambridge\\02140\\"})
	op = append(op, loader.OptionalFields{ValType: "income_cd", Value: "Low"})
	op = append(op, loader.OptionalFields{ValType: "patient_blob", Value: ""})

	_, pubKey := libunlynx.GenKey()
	enc := libunlynx.EncryptInt(pubKey, int64(2))

	pd := loader.PatientDimension{
		PK:             pdk,
		VitalStatusCD:  "D",
		BirthDate:      "1985-11-17 00:00:00",
		DeathDate:      "\\N",
		OptionalFields: op,
		AdminColumns:   ac,
		EncryptedFlag:  *enc,
	}

	b := pd.EncryptedFlag.ToBytes()
	encodedEncryptedFlag := "\"" + base64.StdEncoding.EncodeToString(b) + "\""

	assert.Equal(t, pd.ToCSVText(), `"1000000001","D","1985-11-17 00:00:00","\N","F","24","english","black","married","roman catholic","02140","Zip codes\Massachusetts\Cambridge\02140\","Low","","2010-11-04 10:43:00","2010-08-18 09:50:00","2010-11-04 10:43:00","DEMO","\N",`+encodedEncryptedFlag)
}

func TestConceptDimension_ToCSVText(t *testing.T) {

	csvString := `"\i2b2\Demographics\Age\>= 65 years old\100\","DEM|AGE:100"," 100 years old","","2010-09-28 11:15:00","2010-08-18 09:50:00","2010-09-28 11:40:00","DEMO","\N"`

	ac := loader.AdministrativeColumns{
		UpdateDate:     "2010-09-28 11:15:00",
		DownloadDate:   "2010-08-18 09:50:00",
		ImportDate:     "2010-09-28 11:40:00",
		SourceSystemCD: "DEMO",
		UploadID:       "\\N",
	}

	cdk := &loader.ConceptDimensionPK{
		ConceptPath: "\\i2b2\\Demographics\\Age\\>= 65 years old\\100\\",
	}

	cd := loader.ConceptDimension{
		PK:           cdk,
		ConceptCD:    "DEM|AGE:100",
		NameChar:     " 100 years old",
		ConceptBlob:  "",
		AdminColumns: ac,
	}

	assert.Equal(t, csvString, cd.ToCSVText())

	tag := libunlynx.GroupingKey("1")
	assert.Equal(t, `"\medco\tagged\concept\1\", "TAG_ID:20", "\N", "\N", "\N", "\N", "NOW()", "\N", "\N"`, loader.ConceptDimensionSensitiveToCSVText(&tag, 20))
}

func TestModifierDimension_ToCSVText(t *testing.T) {

	csvString := `"\Admit Diagnosis\","0","Admit Diagnosis","","2011-04-14 00:55:00","\N","2011-04-14 00:55:00","DEMO","\N"`

	ac := loader.AdministrativeColumns{
		UpdateDate:     "2011-04-14 00:55:00",
		DownloadDate:   "\\N",
		ImportDate:     "2011-04-14 00:55:00",
		SourceSystemCD: "DEMO",
		UploadID:       "\\N",
	}

	mdk := &loader.ModifierDimensionPK{
		ModifierPath: "\\Admit Diagnosis\\",
	}

	md := loader.ModifierDimension{
		PK:           mdk,
		ModifierCD:   "0",
		NameChar:     "Admit Diagnosis",
		ModifierBlob: "",
		AdminColumns: ac,
	}

	assert.Equal(t, csvString, md.ToCSVText())

	tag := libunlynx.GroupingKey("1")
	assert.Equal(t, `"\medco\tagged\modifier\1\", "TAG_ID:20", "\N", "\N", "\N", "\N", "NOW()", "\N", "\N"`, loader.ModifierDimensionSensitiveToCSVText(&tag, 20))
}

func TestObservationFact_ToCSVText(t *testing.T) {

	csvString := `"482232","1000000060","Affy:221610_s_at","LCS-I2B2:D000109064","2009-01-16 00:00:00","@","1","N","E","79.30000","","\N","","2009-01-16 00:00:00","@","","\N","2010-09-28 11:15:00","2010-08-18 09:50:00","2010-09-28 11:40:00","DEMO","\N","1"`

	ac := loader.AdministrativeColumns{
		UpdateDate:      "2010-09-28 11:15:00",
		DownloadDate:    "2010-08-18 09:50:00",
		ImportDate:      "2010-09-28 11:40:00",
		SourceSystemCD:  "DEMO",
		UploadID:        "\\N",
		TextSearchIndex: "1",
	}

	ofk := &loader.ObservationFactPK{
		EncounterNum: "482232",
		PatientNum:   "1000000060",
		ConceptCD:    "Affy:221610_s_at",
		ProviderID:   "LCS-I2B2:D000109064",
		StartDate:    "2009-01-16 00:00:00",
		ModifierCD:   "@",
		InstanceNum:  "1",
	}

	of := loader.ObservationFact{
		PK:              ofk,
		ValTypeCD:       "N",
		TValChar:        "E",
		NValNum:         "79.30000",
		ValueFlagCD:     "",
		QuantityNum:     "\\N",
		UnitsCD:         "",
		EndDate:         "2009-01-16 00:00:00",
		LocationCD:      "@",
		ObservationBlob: "",
		ConfidenceNum:   "\\N",
		AdminColumns:    ac,
	}

	assert.Equal(t, csvString, of.ToCSVText())
}

// ------------------------------------------------------------------------------------------------------------- //
// ---------------------------------------- FROM STRING -------------------------------------------------------- //
// ------------------------------------------------------------------------------------------------------------- //

func TestShrineOntologyFromString(t *testing.T) {
	csvString := `"0","\SHRINE\","SHRINE","N","CA ","\N","\N","","concept_cd","concept_dimension","concept_path","T","LIKE","\SHRINE\","","\N","\N","\N","\N","SHRINE","\N","@","\N"`

	ac := loader.AdministrativeColumns{
		UpdateDate:     "\\N",
		DownloadDate:   "\\N",
		ImportDate:     "\\N",
		SourceSystemCD: "SHRINE",
	}

	so := loader.ShrineOntology{
		NodeEncryptID:      -1,
		ChildrenEncryptIDs: nil,

		HLevel:           "0",
		Fullname:         "\\SHRINE\\",
		Name:             "SHRINE",
		SynonymCD:        "N",
		VisualAttributes: "CA ",
		TotalNum:         "\\N",
		BaseCode:         "\\N",
		MetadataXML:      "",
		FactTableColumn:  "concept_cd",
		Tablename:        "concept_dimension",
		ColumnName:       "concept_path",
		ColumnDataType:   "T",
		Operator:         "LIKE",
		DimCode:          "\\SHRINE\\",
		Comment:          "",
		Tooltip:          "\\N",
		AdminColumns:     ac,
		ValueTypeCD:      "\\N",
		AppliedPath:      "@",
		ExclusionCD:      "\\N",
	}

	var csvFile = strings.NewReader(csvString)
	r := csv.NewReader(csvFile)
	lines, err := r.ReadAll()
	assert.Nil(t, err, "Parsing error")

	assert.Equal(t, *loader.ShrineOntologyFromString(lines[0]), so)
}

func TestLocalOntologyFromString(t *testing.T) {
	csvString := `"4","\i2b2\Demographics\Zip codes\Arkansas\Parkdale\","Parkdale","N","FA ","\N","\N","\N","concept_cd","concept_dimension","concept_path","T","LIKE","\i2b2\Demographics\Zip codes\Arkansas\Parkdale\","\N","Demographics \ Zip codes \ Arkansas \ Parkdale","@","2007-04-10 00:00:00","2007-04-10 00:00:00","2007-04-10 00:00:00","DEMO","\N","\N","\N","\N"`

	ac := loader.AdministrativeColumns{
		UpdateDate:     "2007-04-10 00:00:00",
		DownloadDate:   "2007-04-10 00:00:00",
		ImportDate:     "2007-04-10 00:00:00",
		SourceSystemCD: "DEMO",
	}

	lo := loader.LocalOntology{
		HLevel:           "4",
		Fullname:         "\\i2b2\\Demographics\\Zip codes\\Arkansas\\Parkdale\\",
		Name:             "Parkdale",
		SynonymCD:        "N",
		VisualAttributes: "FA ",
		TotalNum:         "\\N",
		BaseCode:         "\\N",
		MetadataXML:      "\\N",
		FactTableColumn:  "concept_cd",
		Tablename:        "concept_dimension",
		ColumnName:       "concept_path",
		ColumnDataType:   "T",
		Operator:         "LIKE",
		DimCode:          "\\i2b2\\Demographics\\Zip codes\\Arkansas\\Parkdale\\",
		Comment:          "\\N",
		Tooltip:          "Demographics \\ Zip codes \\ Arkansas \\ Parkdale",
		AppliedPath:      "@",
		AdminColumns:     ac,
		ValueTypeCD:      "\\N",
		ExclusionCD:      "\\N",
		Path:             "\\N",
		Symbol:           "\\N",

		PCoriBasecode: "",
	}

	var csvFile = strings.NewReader(csvString)
	r := csv.NewReader(csvFile)
	lines, err := r.ReadAll()
	assert.Nil(t, err, "Parsing error")

	assert.Equal(t, *loader.LocalOntologyFromString(lines[0]), lo)
}

func TestPatientDimensionFromString(t *testing.T) {
	aux := [...]string{"patient_num", "vital_status_cd", "birth_date", "death_date", "sex_cd", "age_in_years_num", "language_cd", "race_cd", "marital_status_cd", "religion_cd", "zip_cd", "statecityzip_path", "income_cd", "patient_blob", "update_date", "download_date", "import_date", "sourcesystem_cd", "upload_id"}
	loader.HeaderPatientDimension = aux[:]

	ac := loader.AdministrativeColumns{
		UpdateDate:     "2010-11-04 10:43:00",
		DownloadDate:   "2010-08-18 09:50:00",
		ImportDate:     "2010-11-04 10:43:00",
		SourceSystemCD: "DEMO",
		UploadID:       "\\N",
	}

	pdk := &loader.PatientDimensionPK{
		PatientNum: "1000000001",
	}

	op := make([]loader.OptionalFields, 0)
	op = append(op, loader.OptionalFields{ValType: "sex_cd", Value: "F"})
	op = append(op, loader.OptionalFields{ValType: "age_in_years_num", Value: "24"})
	op = append(op, loader.OptionalFields{ValType: "language_cd", Value: "english"})
	op = append(op, loader.OptionalFields{ValType: "race_cd", Value: "black"})
	op = append(op, loader.OptionalFields{ValType: "marital_status_cd", Value: "married"})
	op = append(op, loader.OptionalFields{ValType: "religion_cd", Value: "roman catholic"})
	op = append(op, loader.OptionalFields{ValType: "zip_cd", Value: "02140"})
	op = append(op, loader.OptionalFields{ValType: "statecityzip_path", Value: "Zip codes\\Massachusetts\\Cambridge\\02140\\"})
	op = append(op, loader.OptionalFields{ValType: "income_cd", Value: "Low"})
	op = append(op, loader.OptionalFields{ValType: "patient_blob", Value: ""})

	_, pubKey := libunlynx.GenKey()
	enc := libunlynx.EncryptInt(pubKey, int64(2))

	pd := loader.PatientDimension{
		PK:             pdk,
		VitalStatusCD:  "D",
		BirthDate:      "1985-11-17 00:00:00",
		DeathDate:      "\\N",
		OptionalFields: op,
		AdminColumns:   ac,
		EncryptedFlag:  *enc,
	}

	csvString := `"1000000001","D","1985-11-17 00:00:00","\N","F","24","english","black","married","roman catholic","02140","Zip codes\Massachusetts\Cambridge\02140\","Low","","2010-11-04 10:43:00","2010-08-18 09:50:00","2010-11-04 10:43:00","DEMO","\N"`

	var csvFile = strings.NewReader(csvString)
	r := csv.NewReader(csvFile)
	lines, err := r.ReadAll()
	assert.Nil(t, err, "Parsing error")

	pdkExpected, pdExpected := loader.PatientDimensionFromString(lines[0], pubKey)
	assert.Equal(t, *pdkExpected, *pdk)

	// place them nil because encryption is randomized
	pdExpected.EncryptedFlag = libunlynx.CipherText{}
	pd.EncryptedFlag = libunlynx.CipherText{}

	assert.Equal(t, pdExpected, pd)
}

func TestConceptDimensionFromString(t *testing.T) {
	csvString := `"\i2b2\Demographics\Age\>= 65 years old\100\","DEM|AGE:100"," 100 years old","","2010-09-28 11:15:00","2010-08-18 09:50:00","2010-09-28 11:40:00","DEMO","\N"`

	ac := loader.AdministrativeColumns{
		UpdateDate:     "2010-09-28 11:15:00",
		DownloadDate:   "2010-08-18 09:50:00",
		ImportDate:     "2010-09-28 11:40:00",
		SourceSystemCD: "DEMO",
		UploadID:       "\\N",
	}

	cdk := &loader.ConceptDimensionPK{
		ConceptPath: "\\i2b2\\Demographics\\Age\\>= 65 years old\\100\\",
	}

	cd := loader.ConceptDimension{
		PK:           cdk,
		ConceptCD:    "DEM|AGE:100",
		NameChar:     " 100 years old",
		ConceptBlob:  "",
		AdminColumns: ac,
	}

	var csvFile = strings.NewReader(csvString)
	r := csv.NewReader(csvFile)
	lines, err := r.ReadAll()
	assert.Nil(t, err, "Parsing error")

	cdkExpected, cdExpected := loader.ConceptDimensionFromString(lines[0])

	assert.Equal(t, *cdkExpected, *cdk)
	assert.Equal(t, cdExpected, cd)
}

func TestModifierDimensionFromString(t *testing.T) {
	csvString := `"\Admit Diagnosis\","0","Admit Diagnosis","","2011-04-14 00:55:00","\N","2011-04-14 00:55:00","DEMO","\N"`

	ac := loader.AdministrativeColumns{
		UpdateDate:     "2011-04-14 00:55:00",
		DownloadDate:   "\\N",
		ImportDate:     "2011-04-14 00:55:00",
		SourceSystemCD: "DEMO",
		UploadID:       "\\N",
	}

	mdk := &loader.ModifierDimensionPK{
		ModifierPath: "\\Admit Diagnosis\\",
	}

	md := loader.ModifierDimension{
		PK:           mdk,
		ModifierCD:   "0",
		NameChar:     "Admit Diagnosis",
		ModifierBlob: "",
		AdminColumns: ac,
	}

	var csvFile = strings.NewReader(csvString)
	r := csv.NewReader(csvFile)
	lines, err := r.ReadAll()
	assert.Nil(t, err, "Parsing error")

	mdkExpected, mdExpected := loader.ModifierDimensionFromString(lines[0])

	assert.Equal(t, *mdkExpected, *mdk)
	assert.Equal(t, mdExpected, md)
}

func TestObservationFactFromString(t *testing.T) {
	csvString := `"482232","1000000060","Affy:221610_s_at","LCS-I2B2:D000109064","2009-01-16 00:00:00","@","1","N","E","79.30000","","\N","","2009-01-16 00:00:00","@","","\N","2010-09-28 11:15:00","2010-08-18 09:50:00","2010-09-28 11:40:00","DEMO","\N","1"
`
	ac := loader.AdministrativeColumns{
		UpdateDate:      "2010-09-28 11:15:00",
		DownloadDate:    "2010-08-18 09:50:00",
		ImportDate:      "2010-09-28 11:40:00",
		SourceSystemCD:  "DEMO",
		UploadID:        "\\N",
		TextSearchIndex: "1",
	}

	ofk := &loader.ObservationFactPK{
		EncounterNum: "482232",
		PatientNum:   "1000000060",
		ConceptCD:    "Affy:221610_s_at",
		ProviderID:   "LCS-I2B2:D000109064",
		StartDate:    "2009-01-16 00:00:00",
		ModifierCD:   "@",
		InstanceNum:  "1",
	}

	of := loader.ObservationFact{
		PK:              ofk,
		ValTypeCD:       "N",
		TValChar:        "E",
		NValNum:         "79.30000",
		ValueFlagCD:     "",
		QuantityNum:     "\\N",
		UnitsCD:         "",
		EndDate:         "2009-01-16 00:00:00",
		LocationCD:      "@",
		ObservationBlob: "",
		ConfidenceNum:   "\\N",
		AdminColumns:    ac,
	}

	var csvFile = strings.NewReader(csvString)
	r := csv.NewReader(csvFile)
	lines, err := r.ReadAll()
	assert.Nil(t, err, "Parsing error")

	ofkExpected, ofExpected := loader.ObservationFactFromString(lines[0])

	assert.Equal(t, *ofkExpected, *ofk)
	assert.Equal(t, ofExpected, of)
}
