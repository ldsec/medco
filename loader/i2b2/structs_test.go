package loaderi2b2_test

import (
	"encoding/base64"
	"encoding/csv"
	"github.com/lca1/medco-loader/loader/i2b2"
	"github.com/lca1/unlynx/lib"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

// ----------------------------------------------------------------------------------------------------------- //
// ---------------------------------------- TO STRING -------------------------------------------------------- //
// ----------------------------------------------------------------------------------------------------------- //

func TestSchemes_ToCSVText(t *testing.T) {
	sk := &loaderi2b2.SchemesPK{
		Key: "NDC:",
	}

	s := loaderi2b2.Schemes{
		PK:          sk,
		Name:        "NDC",
		Description: "National Drug Code",
	}

	assert.Equal(t, s.ToCSVText(), `"NDC:","NDC","National Drug Code"`)

}

func TestTableAccess_ToCSVText(t *testing.T) {
	ta := loaderi2b2.TableAccess{
		TableCD:          "i2b2_DEMO",
		TableName:        "I2B2",
		ProtectedAccess:  "N",
		Hlevel:           "1",
		Fullname:         "\\i2b2\\Demographics\\",
		Name:             "Demographics",
		SynonymCD:        "N",
		Visualattributes: "CA ",
		Totalnum:         "\\N",
		Basecode:         "\\N",
		Metadataxml:      "\\N",
		Facttablecolumn:  "concept_cd",
		Dimtablename:     "concept_dimension",
		Columnname:       "concept_path",
		Columndatatype:   "T",
		Operator:         "LIKE",
		Dimcode:          "\\i2b2\\Demographics\\",
		Comment:          "\\N",
		Tooltip:          "Demographics",
		EntryDate:        "\\N",
		ChangeDate:       "\\N",
		StatusCD:         "\\N",
		ValuetypeCD:      "\\N",
	}

	assert.Equal(t, ta.ToCSVText(), `"i2b2_DEMO","I2B2","N","1","\i2b2\Demographics\","Demographics","N","CA ",,,,"concept_cd","concept_dimension","concept_path","T","LIKE","\i2b2\Demographics\",,"Demographics",,,,`)

}

func TestShrineOntology_ToCSVText(t *testing.T) {

	ac := loaderi2b2.AdministrativeColumns{
		UpdateDate:     "\\N",
		DownloadDate:   "\\N",
		ImportDate:     "\\N",
		SourceSystemCD: "SHRINE",
	}

	so := loaderi2b2.ShrineOntology{
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
	assert.Equal(t, so.ToCSVText(), `"0","\SHRINE\","SHRINE","N","CA ",,,"","concept_cd","concept_dimension","concept_path","T","LIKE","\SHRINE\","",,,,,"SHRINE",,"@",`)

	so.NodeEncryptID = 1
	assert.Equal(t, so.ToCSVText(), `"0","\SHRINE\","SHRINE","N","CA ",,,"<?xml version=""1.0""?><ValueMetadata><Version>MedCo-0.1</Version><EncryptedType>CONCEPT_PARENT_NODE</EncryptedType></ValueMetadata>","concept_cd","concept_dimension","concept_path","T","LIKE","\SHRINE\","",,,,,"SHRINE",,"@",`)

	so.VisualAttributes = "LA "
	assert.Equal(t, so.ToCSVText(), `"0","\SHRINE\","SHRINE","N","LA ",,,"<?xml version=""1.0""?><ValueMetadata><Version>MedCo-0.1</Version><EncryptedType>CONCEPT_LEAF</EncryptedType><NodeEncryptID>1</NodeEncryptID></ValueMetadata>","concept_cd","concept_dimension","concept_path","T","LIKE","\SHRINE\","",,,,,"SHRINE",,"@",`)

	so.ChildrenEncryptIDs = append(so.ChildrenEncryptIDs, 2)
	so.ChildrenEncryptIDs = append(so.ChildrenEncryptIDs, 3)
	assert.Equal(t, so.ToCSVText(), `"0","\SHRINE\","SHRINE","N","LA ",,,"<?xml version=""1.0""?><ValueMetadata><Version>MedCo-0.1</Version><EncryptedType>CONCEPT_LEAF</EncryptedType><NodeEncryptID>1</NodeEncryptID></ValueMetadata>","concept_cd","concept_dimension","concept_path","T","LIKE","\SHRINE\","",,,,,"SHRINE",,"@",`)

	so.VisualAttributes = "FA "
	assert.Equal(t, so.ToCSVText(), `"0","\SHRINE\","SHRINE","N","FA ",,,"<?xml version=""1.0""?><ValueMetadata><Version>MedCo-0.1</Version><EncryptedType>CONCEPT_INTERNAL_NODE</EncryptedType><NodeEncryptID>1</NodeEncryptID><ChildrenEncryptIDs>"2;3"</ChildrenEncryptIDs></ValueMetadata>","concept_cd","concept_dimension","concept_path","T","LIKE","\SHRINE\","",,,,,"SHRINE",,"@",`)

	so.VisualAttributes = "M "
	assert.Equal(t, so.ToCSVText(), `"0","\SHRINE\","SHRINE","N","M ",,,"","concept_cd","concept_dimension","concept_path","T","LIKE","\SHRINE\","",,,,,"SHRINE",,"@",`)

}

func TestLocalOntology_ToCSVText(t *testing.T) {

	ac := loaderi2b2.AdministrativeColumns{
		UpdateDate:     "2007-04-10 00:00:00",
		DownloadDate:   "2007-04-10 00:00:00",
		ImportDate:     "2007-04-10 00:00:00",
		SourceSystemCD: "DEMO",
	}

	lo := loaderi2b2.LocalOntology{
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

	assert.Equal(t, lo.ToCSVText(), `"4","\i2b2\Demographics\Zip codes\Arkansas\Parkdale\","Parkdale","N","FA ",,,,"concept_cd","concept_dimension","concept_path","T","LIKE","\i2b2\Demographics\Zip codes\Arkansas\Parkdale\",,"Demographics \ Zip codes \ Arkansas \ Parkdale","@","2007-04-10 00:00:00","2007-04-10 00:00:00","2007-04-10 00:00:00","DEMO",,,,`)

	tag := libunlynx.GroupingKey("1")
	assert.Equal(t, loaderi2b2.LocalOntologySensitiveConceptToCSVText(&tag, 20), `"3","\medco\tagged\1\","","N","LA ",,"TAG_ID:20",,"concept_cd","concept_dimension","concept_path","T","LIKE","\medco\tagged\concept\1\",,,"NOW()",,,,"TAG_ID","@",,,,`)

}

func TestPatientDimension_ToCSVText(t *testing.T) {

	ac := loaderi2b2.AdministrativeColumns{
		UpdateDate:     "2010-11-04 10:43:00",
		DownloadDate:   "2010-08-18 09:50:00",
		ImportDate:     "2010-11-04 10:43:00",
		SourceSystemCD: "DEMO",
		UploadID:       "\\N",
	}

	pdk := loaderi2b2.PatientDimensionPK{
		PatientNum: "1000000001",
	}

	op := make([]loaderi2b2.OptionalFields, 0)
	op = append(op, loaderi2b2.OptionalFields{ValType: "sex_cd", Value: "F"})
	op = append(op, loaderi2b2.OptionalFields{ValType: "age_in_years_num", Value: "24"})
	op = append(op, loaderi2b2.OptionalFields{ValType: "language_cd", Value: "english"})
	op = append(op, loaderi2b2.OptionalFields{ValType: "race_cd", Value: "black"})
	op = append(op, loaderi2b2.OptionalFields{ValType: "marital_status_cd", Value: "married"})
	op = append(op, loaderi2b2.OptionalFields{ValType: "religion_cd", Value: "roman catholic"})
	op = append(op, loaderi2b2.OptionalFields{ValType: "zip_cd", Value: "02140"})
	op = append(op, loaderi2b2.OptionalFields{ValType: "statecityzip_path", Value: "Zip codes\\Massachusetts\\Cambridge\\02140\\"})
	op = append(op, loaderi2b2.OptionalFields{ValType: "income_cd", Value: "Low"})
	op = append(op, loaderi2b2.OptionalFields{ValType: "patient_blob", Value: ""})

	_, pubKey := libunlynx.GenKey()
	enc := libunlynx.EncryptInt(pubKey, int64(2))

	pd := loaderi2b2.PatientDimension{
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

	assert.Equal(t, pd.ToCSVText(false), `"1000000001","D","1985-11-17 00:00:00",,"F","24","english","black","married","roman catholic","02140","Zip codes\Massachusetts\Cambridge\02140\","Low","","2010-11-04 10:43:00","2010-08-18 09:50:00","2010-11-04 10:43:00","DEMO",,`+encodedEncryptedFlag)
	assert.Equal(t, pd.ToCSVText(true), `"1000000001",,,,,,,,,,,,,,,,,,,`+encodedEncryptedFlag)

}

func TestVisitDimension_ToCSVText(t *testing.T) {
	ac := loaderi2b2.AdministrativeColumns{
		UpdateDate:     "2010-11-04 10:43:00",
		DownloadDate:   "2010-08-18 09:50:00",
		ImportDate:     "2010-11-04 10:43:00",
		SourceSystemCD: "DEMO",
		UploadID:       "\\N",
	}

	vdk := loaderi2b2.VisitDimensionPK{
		EncounterNum: "471185",
		PatientNum:   "1000000101",
	}

	op := make([]loaderi2b2.OptionalFields, 0)
	op = append(op, loaderi2b2.OptionalFields{ValType: "inout_cd", Value: "O"})
	op = append(op, loaderi2b2.OptionalFields{ValType: "location_cd", Value: ""})
	op = append(op, loaderi2b2.OptionalFields{ValType: "location_path", Value: ""})
	op = append(op, loaderi2b2.OptionalFields{ValType: "length_of_stay", Value: "\\N"})
	op = append(op, loaderi2b2.OptionalFields{ValType: "visit_blob", Value: ""})

	vd := loaderi2b2.VisitDimension{
		PK:             vdk,
		ActiveStatusCD: "U",
		StartDate:      "1997-01-02 00:00:00",
		EndDate:        "\\N",
		OptionalFields: op,
		AdminColumns:   ac,
	}

	assert.Equal(t, vd.ToCSVText(false), `"471185","1000000101","U","1997-01-02 00:00:00",,"O","","",,"","2010-11-04 10:43:00","2010-08-18 09:50:00","2010-11-04 10:43:00","DEMO",`)
	assert.Equal(t, vd.ToCSVText(true), `"471185","1000000101",,,,,,,,,,,,,`)
}

func TestConceptDimension_ToCSVText(t *testing.T) {

	csvString := `"\i2b2\Demographics\Age\>= 65 years old\100\","DEM|AGE:100"," 100 years old","","2010-09-28 11:15:00","2010-08-18 09:50:00","2010-09-28 11:40:00","DEMO",`

	ac := loaderi2b2.AdministrativeColumns{
		UpdateDate:     "2010-09-28 11:15:00",
		DownloadDate:   "2010-08-18 09:50:00",
		ImportDate:     "2010-09-28 11:40:00",
		SourceSystemCD: "DEMO",
		UploadID:       "\\N",
	}

	cdk := &loaderi2b2.ConceptDimensionPK{
		ConceptPath: "\\i2b2\\Demographics\\Age\\>= 65 years old\\100\\",
	}

	cd := loaderi2b2.ConceptDimension{
		PK:           cdk,
		ConceptCD:    "DEM|AGE:100",
		NameChar:     " 100 years old",
		ConceptBlob:  "",
		AdminColumns: ac,
	}

	assert.Equal(t, csvString, cd.ToCSVText())

	tag := libunlynx.GroupingKey("1")
	assert.Equal(t, `"\medco\tagged\concept\1\","TAG_ID:20",,,,,"NOW()",,`, loaderi2b2.ConceptDimensionSensitiveToCSVText(&tag, 20))
}

func TestObservationFact_ToCSVText(t *testing.T) {

	csvString := `"482232","1000000060","Affy:221610_s_at","LCS-I2B2:D000109064","2009-01-16 00:00:00","@","1","N","E","79.30000","",,"","2009-01-16 00:00:00","@","",,"2010-09-28 11:15:00","2010-08-18 09:50:00","2010-09-28 11:40:00","DEMO",,"1"`

	ac := loaderi2b2.AdministrativeColumns{
		UpdateDate:      "2010-09-28 11:15:00",
		DownloadDate:    "2010-08-18 09:50:00",
		ImportDate:      "2010-09-28 11:40:00",
		SourceSystemCD:  "DEMO",
		UploadID:        "\\N",
		TextSearchIndex: "1",
	}

	ofk := &loaderi2b2.ObservationFactPK{
		EncounterNum: "482232",
		PatientNum:   "1000000060",
		ConceptCD:    "Affy:221610_s_at",
		ProviderID:   "LCS-I2B2:D000109064",
		StartDate:    "2009-01-16 00:00:00",
		ModifierCD:   "@",
		InstanceNum:  "1",
	}

	of := loaderi2b2.ObservationFact{
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

func TestTableAccessFromString(t *testing.T) {
	csvString := `"i2b2_DEMO","I2B2","N","1","\i2b2\Demographics\","Demographics","N","CA ","\N","\N","\N","concept_cd","concept_dimension","concept_path","T","LIKE","\i2b2\Demographics\","\N","Demographics","\N","\N","\N","\N"`

	ta := loaderi2b2.TableAccess{
		TableCD:          "i2b2_DEMO",
		TableName:        "I2B2",
		ProtectedAccess:  "N",
		Hlevel:           "1",
		Fullname:         `\i2b2\Demographics\`,
		Name:             "Demographics",
		SynonymCD:        "N",
		Visualattributes: "CA ",
		Totalnum:         "\\N",
		Basecode:         "\\N",
		Metadataxml:      "\\N",
		Facttablecolumn:  "concept_cd",
		Dimtablename:     "concept_dimension",
		Columnname:       "concept_path",
		Columndatatype:   "T",
		Operator:         "LIKE",
		Dimcode:          `\i2b2\Demographics\`,
		Comment:          "\\N",
		Tooltip:          "Demographics",
		EntryDate:        "\\N",
		ChangeDate:       "\\N",
		StatusCD:         "\\N",
		ValuetypeCD:      "\\N",
	}

	var csvFile = strings.NewReader(csvString)
	r := csv.NewReader(csvFile)
	lines, err := r.ReadAll()
	assert.Nil(t, err, "Parsing error")

	assert.Equal(t, *loaderi2b2.TableAccessFromString(lines[0]), ta)

}

func TestLocalOntologyFromString(t *testing.T) {
	csvString := `"4","\i2b2\Demographics\Zip codes\Arkansas\Parkdale\","Parkdale","N","FA ","\N","\N","\N","concept_cd","concept_dimension","concept_path","T","LIKE","\i2b2\Demographics\Zip codes\Arkansas\Parkdale\","\N","Demographics \ Zip codes \ Arkansas \ Parkdale","@","2007-04-10 00:00:00","2007-04-10 00:00:00","2007-04-10 00:00:00","DEMO","\N","\N","\N","\N","\N"`

	ac := loaderi2b2.AdministrativeColumns{
		UpdateDate:     "2007-04-10 00:00:00",
		DownloadDate:   "2007-04-10 00:00:00",
		ImportDate:     "2007-04-10 00:00:00",
		SourceSystemCD: "DEMO",
	}

	lo := loaderi2b2.LocalOntology{
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

		PlainCode: "",
	}

	var csvFile = strings.NewReader(csvString)
	r := csv.NewReader(csvFile)
	lines, err := r.ReadAll()
	assert.Nil(t, err, "Parsing error")

	assert.Equal(t, *loaderi2b2.LocalOntologyFromString(lines[0], false), lo)
	lo.PlainCode = "\\N"
	assert.Equal(t, *loaderi2b2.LocalOntologyFromString(lines[0], true), lo)
}

func TestPatientDimensionFromString(t *testing.T) {
	aux := [...]string{"patient_num", "vital_status_cd", "birth_date", "death_date", "sex_cd", "age_in_years_num", "language_cd", "race_cd", "marital_status_cd", "religion_cd", "zip_cd", "statecityzip_path", "income_cd", "patient_blob", "update_date", "download_date", "import_date", "sourcesystem_cd", "upload_id"}
	loaderi2b2.HeaderPatientDimension = aux[:]

	ac := loaderi2b2.AdministrativeColumns{
		UpdateDate:     "2010-11-04 10:43:00",
		DownloadDate:   "2010-08-18 09:50:00",
		ImportDate:     "2010-11-04 10:43:00",
		SourceSystemCD: "DEMO",
		UploadID:       "\\N",
	}

	pdk := loaderi2b2.PatientDimensionPK{
		PatientNum: "1000000001",
	}

	op := make([]loaderi2b2.OptionalFields, 0)
	op = append(op, loaderi2b2.OptionalFields{ValType: "sex_cd", Value: "F"})
	op = append(op, loaderi2b2.OptionalFields{ValType: "age_in_years_num", Value: "24"})
	op = append(op, loaderi2b2.OptionalFields{ValType: "language_cd", Value: "english"})
	op = append(op, loaderi2b2.OptionalFields{ValType: "race_cd", Value: "black"})
	op = append(op, loaderi2b2.OptionalFields{ValType: "marital_status_cd", Value: "married"})
	op = append(op, loaderi2b2.OptionalFields{ValType: "religion_cd", Value: "roman catholic"})
	op = append(op, loaderi2b2.OptionalFields{ValType: "zip_cd", Value: "02140"})
	op = append(op, loaderi2b2.OptionalFields{ValType: "statecityzip_path", Value: "Zip codes\\Massachusetts\\Cambridge\\02140\\"})
	op = append(op, loaderi2b2.OptionalFields{ValType: "income_cd", Value: "Low"})
	op = append(op, loaderi2b2.OptionalFields{ValType: "patient_blob", Value: ""})

	_, pubKey := libunlynx.GenKey()
	enc := libunlynx.EncryptInt(pubKey, int64(2))

	pd := loaderi2b2.PatientDimension{
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

	pdkExpected, pdExpected := loaderi2b2.PatientDimensionFromString(lines[0], pubKey)
	assert.Equal(t, pdkExpected, pdk)

	// place them nil because encryption is randomized
	pdExpected.EncryptedFlag = libunlynx.CipherText{}
	pd.EncryptedFlag = libunlynx.CipherText{}

	assert.Equal(t, pdExpected, pd)
}

func TestVisitDimensionFromString(t *testing.T) {
	aux := [...]string{"encounter_num", "patient_num", "active_status_cd", "start_date", "end_date", "inout_cd", "location_cd", "location_path", "length_of_stay", "visit_blob", "update_date", "download_date", "import_date", "sourcesystem_cd", "upload_id"}
	loaderi2b2.HeaderPatientDimension = aux[:]

	ac := loaderi2b2.AdministrativeColumns{
		UpdateDate:     "2010-11-04 10:43:00",
		DownloadDate:   "2010-08-18 09:50:00",
		ImportDate:     "2010-11-04 10:43:00",
		SourceSystemCD: "DEMO",
		UploadID:       "\\N",
	}

	vdk := loaderi2b2.VisitDimensionPK{
		EncounterNum: "471185",
		PatientNum:   "1000000101",
	}

	op := make([]loaderi2b2.OptionalFields, 0)
	op = append(op, loaderi2b2.OptionalFields{ValType: "inout_cd", Value: "O"})
	op = append(op, loaderi2b2.OptionalFields{ValType: "location_cd", Value: ""})
	op = append(op, loaderi2b2.OptionalFields{ValType: "location_path", Value: ""})
	op = append(op, loaderi2b2.OptionalFields{ValType: "length_of_stay", Value: "\\N"})
	op = append(op, loaderi2b2.OptionalFields{ValType: "visit_blob", Value: ""})

	vd := loaderi2b2.VisitDimension{
		PK:             vdk,
		ActiveStatusCD: "U",
		StartDate:      "1997-01-02 00:00:00",
		EndDate:        "\\N",
		OptionalFields: op,
		AdminColumns:   ac,
	}

	csvString := `"471185","1000000101","U","1997-01-02 00:00:00","\N","O","","","\N","","2010-11-04 10:43:00","2010-08-18 09:50:00","2010-11-04 10:43:00","DEMO","\N"`

	var csvFile = strings.NewReader(csvString)
	r := csv.NewReader(csvFile)
	lines, err := r.ReadAll()
	assert.Nil(t, err, "Parsing error")

	vdkExpected, vdExpected := loaderi2b2.VisitDimensionFromString(lines[0])
	assert.Equal(t, vdkExpected, vdk)

	assert.Equal(t, vdExpected, vd)
}

func TestConceptDimensionFromString(t *testing.T) {
	csvString := `"\i2b2\Demographics\Age\>= 65 years old\100\","DEM|AGE:100"," 100 years old","","2010-09-28 11:15:00","2010-08-18 09:50:00","2010-09-28 11:40:00","DEMO","\N"`

	ac := loaderi2b2.AdministrativeColumns{
		UpdateDate:     "2010-09-28 11:15:00",
		DownloadDate:   "2010-08-18 09:50:00",
		ImportDate:     "2010-09-28 11:40:00",
		SourceSystemCD: "DEMO",
		UploadID:       "\\N",
	}

	cdk := &loaderi2b2.ConceptDimensionPK{
		ConceptPath: "\\i2b2\\Demographics\\Age\\>= 65 years old\\100\\",
	}

	cd := loaderi2b2.ConceptDimension{
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

	cdkExpected, cdExpected := loaderi2b2.ConceptDimensionFromString(lines[0])

	assert.Equal(t, *cdkExpected, *cdk)
	assert.Equal(t, cdExpected, cd)
}

func TestObservationFactFromString(t *testing.T) {
	csvString := `"482232","1000000060","Affy:221610_s_at","LCS-I2B2:D000109064","2009-01-16 00:00:00","@","1","N","E","79.30000","","\N","","2009-01-16 00:00:00","@","","\N","2010-09-28 11:15:00","2010-08-18 09:50:00","2010-09-28 11:40:00","DEMO","\N","1"
`
	loaderi2b2.TextSearchIndex = 0

	ac := loaderi2b2.AdministrativeColumns{
		UpdateDate:      "2010-09-28 11:15:00",
		DownloadDate:    "2010-08-18 09:50:00",
		ImportDate:      "2010-09-28 11:40:00",
		SourceSystemCD:  "DEMO",
		UploadID:        "\\N",
		TextSearchIndex: "0",
	}

	ofk := &loaderi2b2.ObservationFactPK{
		EncounterNum: "482232",
		PatientNum:   "1000000060",
		ConceptCD:    "Affy:221610_s_at",
		ProviderID:   "LCS-I2B2:D000109064",
		StartDate:    "2009-01-16 00:00:00",
		ModifierCD:   "@",
		InstanceNum:  "1",
	}

	of := loaderi2b2.ObservationFact{
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

	ofkExpected, ofExpected := loaderi2b2.ObservationFactFromString(lines[0])

	assert.Equal(t, ofkExpected, ofk)
	assert.Equal(t, ofExpected, of)
}
