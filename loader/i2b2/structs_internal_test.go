package loaderi2b2

import (
	"encoding/csv"
	"github.com/ldsec/unlynx/lib"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

// ----------------------------------------------------------------------------------------------------------- //
// ---------------------------------------- TO STRING -------------------------------------------------------- //
// ----------------------------------------------------------------------------------------------------------- //

func TestTableAccess_ToCSVText(t *testing.T) {
	ta := tableAccessRecord{
		tableCD:          "BIRN",
		tableName:        "BIRN",
		protectedAccess:  "N",
		hlevel:           "0",
		fullname:         "\\BIRN\\",
		name:             "Clinical Trials",
		synonymCD:        "N",
		visualAttributes: "CA ",
		totalNum:         "",
		baseCode:         "",
		metadataXML:      "",
		factTableColumn:  "concept_cd",
		dimTableName:     "concept_dimension",
		columnName:       "concept_path",
		columnDataType:   "T",
		operator:         "LIKE",
		dimCode:          "\\BIRN\\",
		comment:          "",
		tooltip:          "Clinical Trials",
		entryData:        "",
		changeData:       "",
		statusCD:         "",
		valueType:        "",
	}

	assert.Equal(t, ta.toCSVText(), `"BIRN","BIRN","N","0","\BIRN\","Clinical Trials","N","CA ","","","","concept_cd","concept_dimension","concept_path","T","LIKE","\BIRN\","","Clinical Trials","","","",""`)
}

func TestMedCoOntology_ToCSVText(t *testing.T) {

	ac := administrativeColumns{
		updateDate:     "\\N",
		downloadDate:   "\\N",
		importDate:     "\\N",
		sourceSystemCD: "SHRINE",
	}

	so := medCoOntologyRecord{
		nodeEncryptID:      -1,
		childrenEncryptIDs: nil,

		hLevel:           "0",
		fullname:         "\\SHRINE\\",
		name:             "SHRINE",
		synonymCD:        "N",
		visualAttributes: "CA ",
		totalNum:         "\\N",
		baseCode:         "\\N",
		metadataXML:      "",
		factTableColumn:  "concept_cd",
		tablename:        "concept_dimension",
		columnName:       "concept_path",
		columnDataType:   "T",
		operator:         "LIKE",
		dimCode:          "\\SHRINE\\",
		comment:          "",
		tooltip:          "\\N",
		adminColumns:     ac,
		valueTypeCD:      "\\N",
		appliedPath:      "@",
		exclusionCD:      "\\N",
	}
	assert.Equal(t, so.toCSVText(), `"0","\SHRINE\","SHRINE","N","CA ",,,"","concept_cd","concept_dimension","concept_path","T","LIKE","\SHRINE\","",,,,,"SHRINE",,"@",`)

	so.nodeEncryptID = 1
	assert.Equal(t, so.toCSVText(), `"0","\SHRINE\","SHRINE","N","CA ",,,"<?xml version=""1.0""?><ValueMetadata><Version>MedCo-0.1</Version><EncryptedType>CONCEPT_PARENT_NODE</EncryptedType></ValueMetadata>","concept_cd","concept_dimension","concept_path","T","LIKE","\SHRINE\","",,,,,"SHRINE",,"@",`)

	so.visualAttributes = "LA "
	assert.Equal(t, so.toCSVText(), `"0","\SHRINE\","SHRINE","N","LA ",,,"<?xml version=""1.0""?><ValueMetadata><Version>MedCo-0.1</Version><EncryptedType>CONCEPT_LEAF</EncryptedType><NodeEncryptID>1</NodeEncryptID></ValueMetadata>","concept_cd","concept_dimension","concept_path","T","LIKE","\SHRINE\","",,,,,"SHRINE",,"@",`)

	so.childrenEncryptIDs = append(so.childrenEncryptIDs, 2)
	so.childrenEncryptIDs = append(so.childrenEncryptIDs, 3)
	assert.Equal(t, so.toCSVText(), `"0","\SHRINE\","SHRINE","N","LA ",,,"<?xml version=""1.0""?><ValueMetadata><Version>MedCo-0.1</Version><EncryptedType>CONCEPT_LEAF</EncryptedType><NodeEncryptID>1</NodeEncryptID></ValueMetadata>","concept_cd","concept_dimension","concept_path","T","LIKE","\SHRINE\","",,,,,"SHRINE",,"@",`)

	so.visualAttributes = "FA "
	assert.Equal(t, so.toCSVText(), `"0","\SHRINE\","SHRINE","N","FA ",,,"<?xml version=""1.0""?><ValueMetadata><Version>MedCo-0.1</Version><EncryptedType>CONCEPT_INTERNAL_NODE</EncryptedType><NodeEncryptID>1</NodeEncryptID><ChildrenEncryptIDs>"2;3"</ChildrenEncryptIDs></ValueMetadata>","concept_cd","concept_dimension","concept_path","T","LIKE","\SHRINE\","",,,,,"SHRINE",,"@",`)

	so.visualAttributes = "M "
	assert.Equal(t, so.toCSVText(), `"0","\SHRINE\","SHRINE","N","M ",,,"","concept_cd","concept_dimension","concept_path","T","LIKE","\SHRINE\","",,,,,"SHRINE",,"@",`)

}

func TestLocalOntology_ToCSVText(t *testing.T) {

	ac := administrativeColumns{
		updateDate:     "2007-04-10 00:00:00",
		downloadDate:   "2007-04-10 00:00:00",
		importDate:     "2007-04-10 00:00:00",
		sourceSystemCD: "DEMO",
	}

	lo := localOntologyRecord{
		hLevel:           "4",
		fullname:         "\\i2b2\\Demographics\\Zip codes\\Arkansas\\Parkdale\\",
		name:             "Parkdale",
		synonymCD:        "N",
		visualAttributes: "FA ",
		totalNum:         "\\N",
		baseCode:         "\\N",
		metadataXML:      "\\N",
		factTableColumn:  "concept_cd",
		tablename:        "concept_dimension",
		columnName:       "concept_path",
		columnDataType:   "T",
		operator:         "LIKE",
		dimCode:          "\\i2b2\\Demographics\\Zip codes\\Arkansas\\Parkdale\\",
		comment:          "\\N",
		tooltip:          "Demographics \\ Zip codes \\ Arkansas \\ Parkdale",
		appliedPath:      "@",
		adminColumns:     ac,
		valueTypeCD:      "\\N",
		exclusionCD:      "\\N",
		path:             "\\N",
		symbol:           "\\N",

		pCoriBasecode: "\\N",
	}

	assert.Equal(t, lo.toCSVText(), `"4","\i2b2\Demographics\Zip codes\Arkansas\Parkdale\","Parkdale","N","FA ",,,,"concept_cd","concept_dimension","concept_path","T","LIKE","\i2b2\Demographics\Zip codes\Arkansas\Parkdale\",,"Demographics \ Zip codes \ Arkansas \ Parkdale","@","2007-04-10 00:00:00","2007-04-10 00:00:00","2007-04-10 00:00:00","DEMO",,,,`)

	tag := libunlynx.GroupingKey("1")
	assert.Equal(t, localOntologySensitiveConceptToCSVText(&tag, 20), `"3","\medco\tagged\1\","","N","LA ",,"TAG_ID:20",,"concept_cd","concept_dimension","concept_path","T","LIKE","\medco\tagged\concept\1\",,,"NOW()",,,,"TAG_ID","@",,,,`)

}

func TestPatientDimension_ToCSVText(t *testing.T) {

	ac := administrativeColumns{
		updateDate:     "2010-11-04 10:43:00",
		downloadDate:   "2010-08-18 09:50:00",
		importDate:     "2010-11-04 10:43:00",
		sourceSystemCD: "DEMO",
		uploadID:       "\\N",
	}

	pdk := patientDimensionPK{
		patientNum: "1000000001",
	}

	op := make([]optionalFields, 0)
	op = append(op, optionalFields{valType: "sex_cd", value: "F"})
	op = append(op, optionalFields{valType: "age_in_years_num", value: "24"})
	op = append(op, optionalFields{valType: "language_cd", value: "english"})
	op = append(op, optionalFields{valType: "race_cd", value: "black"})
	op = append(op, optionalFields{valType: "marital_status_cd", value: "married"})
	op = append(op, optionalFields{valType: "religion_cd", value: "roman catholic"})
	op = append(op, optionalFields{valType: "zip_cd", value: "02140"})
	op = append(op, optionalFields{valType: "statecityzip_path", value: "Zip codes\\Massachusetts\\Cambridge\\02140\\"})
	op = append(op, optionalFields{valType: "income_cd", value: "Low"})
	op = append(op, optionalFields{valType: "patient_blob", value: ""})

	_, pubKey := libunlynx.GenKey()
	enc := libunlynx.EncryptInt(pubKey, int64(2))

	pd := patientDimensionRecord{
		pk:             pdk,
		vitalStatusCD:  "D",
		birthDate:      "1985-11-17 00:00:00",
		deathDate:      "\\N",
		optionalFields: op,
		adminColumns:   ac,
		encryptedFlag:  *enc,
	}

	encryptedFlagString, err := pd.encryptedFlag.Serialize()
	assert.NoError(t, err)
	encodedEncryptedFlag := "\"" + encryptedFlagString + "\""

	assert.Equal(t, pd.toCSVText(false), `"1000000001","D","1985-11-17 00:00:00",,"F","24","english","black","married","roman catholic","02140","Zip codes\Massachusetts\Cambridge\02140\","Low","","2010-11-04 10:43:00","2010-08-18 09:50:00","2010-11-04 10:43:00","DEMO",,`+encodedEncryptedFlag)
	assert.Equal(t, pd.toCSVText(true), `"1000000001",,,,,,,,,,,,,,,,,,,`+encodedEncryptedFlag)

}

func TestVisitDimension_ToCSVText(t *testing.T) {
	ac := administrativeColumns{
		updateDate:     "2010-11-04 10:43:00",
		downloadDate:   "2010-08-18 09:50:00",
		importDate:     "2010-11-04 10:43:00",
		sourceSystemCD: "DEMO",
		uploadID:       "\\N",
	}

	vdk := visitDimensionPK{
		encounterNum: "471185",
		patientNum:   "1000000101",
	}

	op := make([]optionalFields, 0)
	op = append(op, optionalFields{valType: "inout_cd", value: "O"})
	op = append(op, optionalFields{valType: "location_cd", value: ""})
	op = append(op, optionalFields{valType: "location_path", value: ""})
	op = append(op, optionalFields{valType: "length_of_stay", value: "\\N"})
	op = append(op, optionalFields{valType: "visit_blob", value: ""})

	vd := visitDimension{
		pk:             vdk,
		activeStatusCD: "U",
		startDate:      "1997-01-02 00:00:00",
		endDate:        "\\N",
		optionalFields: op,
		adminColumns:   ac,
	}

	assert.Equal(t, vd.toCSVText(false), `"471185","1000000101","U","1997-01-02 00:00:00",,"O","","",,"","2010-11-04 10:43:00","2010-08-18 09:50:00","2010-11-04 10:43:00","DEMO",`)
	assert.Equal(t, vd.toCSVText(true), `"471185","1000000101",,,,,,,,,,,,,`)
}

func TestConceptDimension_ToCSVText(t *testing.T) {

	csvString := `"\i2b2\Demographics\Age\>= 65 years old\100\","DEM|AGE:100"," 100 years old","","2010-09-28 11:15:00","2010-08-18 09:50:00","2010-09-28 11:40:00","DEMO",`

	ac := administrativeColumns{
		updateDate:     "2010-09-28 11:15:00",
		downloadDate:   "2010-08-18 09:50:00",
		importDate:     "2010-09-28 11:40:00",
		sourceSystemCD: "DEMO",
		uploadID:       "\\N",
	}

	cdk := conceptDimensionPK{
		conceptPath: "\\i2b2\\Demographics\\Age\\>= 65 years old\\100\\",
	}

	cd := conceptDimensionRecord{
		pk:           cdk,
		conceptCD:    "DEM|AGE:100",
		nameChar:     " 100 years old",
		conceptBlob:  "",
		adminColumns: ac,
	}

	assert.Equal(t, csvString, cd.toCSVText())

	tag := libunlynx.GroupingKey("1")
	assert.Equal(t, `"\medco\tagged\concept\1\","TAG_ID:20",,,,,"NOW()",,`, conceptDimensionSensitiveToCSVText(&tag, 20))
}

func TestObservationFact_ToCSVText(t *testing.T) {

	csvString := `"482232","1000000060","Affy:221610_s_at","LCS-I2B2:D000109064","2009-01-16 00:00:00","@","1","N","E","79.30000","",,"","2009-01-16 00:00:00","@","",,"2010-09-28 11:15:00","2010-08-18 09:50:00","2010-09-28 11:40:00","DEMO",,"1"`

	ac := administrativeColumns{
		updateDate:      "2010-09-28 11:15:00",
		downloadDate:    "2010-08-18 09:50:00",
		importDate:      "2010-09-28 11:40:00",
		sourceSystemCD:  "DEMO",
		uploadID:        "\\N",
		textSearchIndex: "1",
	}

	ofk := observationFactPK{
		encounterNum: "482232",
		patientNum:   "1000000060",
		conceptCD:    "Affy:221610_s_at",
		providerID:   "LCS-I2B2:D000109064",
		startDate:    "2009-01-16 00:00:00",
		modifierCD:   "@",
		instanceNum:  "1",
	}

	of := observationFactRecord{
		pk:              ofk,
		valTypeCD:       "N",
		tValChar:        "E",
		nValNum:         "79.30000",
		valueFlagCD:     "",
		quantityNum:     "\\N",
		unitsCD:         "",
		endDate:         "2009-01-16 00:00:00",
		locationCD:      "@",
		observationBlob: "",
		confidenceNum:   "\\N",
		adminColumns:    ac,
	}

	assert.Equal(t, csvString, of.toCSVText())
}

// ------------------------------------------------------------------------------------------------------------- //
// ---------------------------------------- FROM STRING -------------------------------------------------------- //
// ------------------------------------------------------------------------------------------------------------- //

func TestTableAccessFromString(t *testing.T) {
	csvString := `"BIRN","BIRN","N","0","\BIRN\","Clinical Trials","N","CA ","","","","concept_cd","concept_dimension","concept_path","T","LIKE","\BIRN\","","Clinical Trials","","","",""`

	ta := tableAccessRecord{
		tableCD:          "BIRN",
		tableName:        "BIRN",
		protectedAccess:  "N",
		hlevel:           "0",
		fullname:         "\\BIRN\\",
		name:             "Clinical Trials",
		synonymCD:        "N",
		visualAttributes: "CA ",
		totalNum:         "",
		baseCode:         "",
		metadataXML:      "",
		factTableColumn:  "concept_cd",
		dimTableName:     "concept_dimension",
		columnName:       "concept_path",
		columnDataType:   "T",
		operator:         "LIKE",
		dimCode:          "\\BIRN\\",
		comment:          "",
		tooltip:          "Clinical Trials",
		entryData:        "",
		changeData:       "",
		statusCD:         "",
		valueType:        "",
	}

	var csvFile = strings.NewReader(csvString)
	r := csv.NewReader(csvFile)
	lines, err := r.ReadAll()
	assert.Nil(t, err, "Parsing error")

	assert.Equal(t, tableAccessFromString(lines[0]), ta)
}

func TestLocalOntologyFromString(t *testing.T) {
	csvString := `"4","\i2b2\Demographics\Zip codes\Arkansas\Parkdale\","Parkdale","N","FA ","\N","\N","\N","concept_cd","concept_dimension","concept_path","T","LIKE","\i2b2\Demographics\Zip codes\Arkansas\Parkdale\","\N","Demographics \ Zip codes \ Arkansas \ Parkdale","@","2007-04-10 00:00:00","2007-04-10 00:00:00","2007-04-10 00:00:00","DEMO","\N","\N","\N","\N","\N"`

	ac := administrativeColumns{
		updateDate:     "2007-04-10 00:00:00",
		downloadDate:   "2007-04-10 00:00:00",
		importDate:     "2007-04-10 00:00:00",
		sourceSystemCD: "DEMO",
	}

	lo := localOntologyRecord{
		hLevel:           "4",
		fullname:         "\\i2b2\\Demographics\\Zip codes\\Arkansas\\Parkdale\\",
		name:             "Parkdale",
		synonymCD:        "N",
		visualAttributes: "FA ",
		totalNum:         "\\N",
		baseCode:         "\\N",
		metadataXML:      "\\N",
		factTableColumn:  "concept_cd",
		tablename:        "concept_dimension",
		columnName:       "concept_path",
		columnDataType:   "T",
		operator:         "LIKE",
		dimCode:          "\\i2b2\\Demographics\\Zip codes\\Arkansas\\Parkdale\\",
		comment:          "\\N",
		tooltip:          "Demographics \\ Zip codes \\ Arkansas \\ Parkdale",
		appliedPath:      "@",
		adminColumns:     ac,
		valueTypeCD:      "\\N",
		exclusionCD:      "\\N",
		path:             "\\N",
		symbol:           "\\N",

		pCoriBasecode: "",

		plainCode: "",
	}

	var csvFile = strings.NewReader(csvString)
	r := csv.NewReader(csvFile)
	lines, err := r.ReadAll()
	assert.Nil(t, err, "Parsing error")

	assert.Equal(t, *localOntologyFromString(lines[0], false), lo)
	lo.plainCode = "\\N"
	assert.Equal(t, *localOntologyFromString(lines[0], true), lo)
}

func TestPatientDimensionFromString(t *testing.T) {
	aux := [...]string{"patient_num", "vital_status_cd", "birth_date", "death_date", "sex_cd", "age_in_years_num", "language_cd", "race_cd", "marital_status_cd", "religion_cd", "zip_cd", "statecityzip_path", "income_cd", "patient_blob", "update_date", "download_date", "import_date", "sourcesystem_cd", "upload_id"}
	headerPatientDimension = aux[:]

	ac := administrativeColumns{
		updateDate:     "2010-11-04 10:43:00",
		downloadDate:   "2010-08-18 09:50:00",
		importDate:     "2010-11-04 10:43:00",
		sourceSystemCD: "DEMO",
		uploadID:       "\\N",
	}

	pdk := patientDimensionPK{
		patientNum: "1000000001",
	}

	op := make([]optionalFields, 0)
	op = append(op, optionalFields{valType: "sex_cd", value: "F"})
	op = append(op, optionalFields{valType: "age_in_years_num", value: "24"})
	op = append(op, optionalFields{valType: "language_cd", value: "english"})
	op = append(op, optionalFields{valType: "race_cd", value: "black"})
	op = append(op, optionalFields{valType: "marital_status_cd", value: "married"})
	op = append(op, optionalFields{valType: "religion_cd", value: "roman catholic"})
	op = append(op, optionalFields{valType: "zip_cd", value: "02140"})
	op = append(op, optionalFields{valType: "statecityzip_path", value: "Zip codes\\Massachusetts\\Cambridge\\02140\\"})
	op = append(op, optionalFields{valType: "income_cd", value: "Low"})
	op = append(op, optionalFields{valType: "patient_blob", value: ""})

	_, pubKey := libunlynx.GenKey()
	enc := libunlynx.EncryptInt(pubKey, int64(2))

	pd := patientDimensionRecord{
		pk:             pdk,
		vitalStatusCD:  "D",
		birthDate:      "1985-11-17 00:00:00",
		deathDate:      "\\N",
		optionalFields: op,
		adminColumns:   ac,
		encryptedFlag:  *enc,
	}

	csvString := `"1000000001","D","1985-11-17 00:00:00","\N","F","24","english","black","married","roman catholic","02140","Zip codes\Massachusetts\Cambridge\02140\","Low","","2010-11-04 10:43:00","2010-08-18 09:50:00","2010-11-04 10:43:00","DEMO","\N"`

	var csvFile = strings.NewReader(csvString)
	r := csv.NewReader(csvFile)
	lines, err := r.ReadAll()
	assert.Nil(t, err, "Parsing error")

	pdkExpected, pdExpected := patientDimensionFromString(lines[0], pubKey)
	assert.Equal(t, pdkExpected, pdk)

	// place them nil because encryption is randomized
	pdExpected.encryptedFlag = libunlynx.CipherText{}
	pd.encryptedFlag = libunlynx.CipherText{}

	assert.Equal(t, pdExpected, &pd)
}

func TestVisitDimensionFromString(t *testing.T) {
	aux := [...]string{"encounter_num", "patient_num", "active_status_cd", "start_date", "end_date", "inout_cd", "location_cd", "location_path", "length_of_stay", "visit_blob", "update_date", "download_date", "import_date", "sourcesystem_cd", "upload_id"}
	headerPatientDimension = aux[:]

	ac := administrativeColumns{
		updateDate:     "2010-11-04 10:43:00",
		downloadDate:   "2010-08-18 09:50:00",
		importDate:     "2010-11-04 10:43:00",
		sourceSystemCD: "DEMO",
		uploadID:       "\\N",
	}

	vdk := visitDimensionPK{
		encounterNum: "471185",
		patientNum:   "1000000101",
	}

	op := make([]optionalFields, 0)
	op = append(op, optionalFields{valType: "inout_cd", value: "O"})
	op = append(op, optionalFields{valType: "location_cd", value: ""})
	op = append(op, optionalFields{valType: "location_path", value: ""})
	op = append(op, optionalFields{valType: "length_of_stay", value: "\\N"})
	op = append(op, optionalFields{valType: "visit_blob", value: ""})

	vd := visitDimension{
		pk:             vdk,
		activeStatusCD: "U",
		startDate:      "1997-01-02 00:00:00",
		endDate:        "\\N",
		optionalFields: op,
		adminColumns:   ac,
	}

	csvString := `"471185","1000000101","U","1997-01-02 00:00:00","\N","O","","","\N","","2010-11-04 10:43:00","2010-08-18 09:50:00","2010-11-04 10:43:00","DEMO","\N"`

	var csvFile = strings.NewReader(csvString)
	r := csv.NewReader(csvFile)
	lines, err := r.ReadAll()
	assert.Nil(t, err, "Parsing error")

	vdkExpected, vdExpected := visitDimensionFromString(lines[0])
	assert.Equal(t, vdkExpected, vdk)

	assert.Equal(t, vdExpected, &vd)
}

func TestConceptDimensionFromString(t *testing.T) {
	csvString := `"\i2b2\Demographics\Age\>= 65 years old\100\","DEM|AGE:100"," 100 years old","","2010-09-28 11:15:00","2010-08-18 09:50:00","2010-09-28 11:40:00","DEMO","\N"`

	ac := administrativeColumns{
		updateDate:     "2010-09-28 11:15:00",
		downloadDate:   "2010-08-18 09:50:00",
		importDate:     "2010-09-28 11:40:00",
		sourceSystemCD: "DEMO",
		uploadID:       "\\N",
	}

	cdk := conceptDimensionPK{
		conceptPath: "\\i2b2\\Demographics\\Age\\>= 65 years old\\100\\",
	}

	cd := &conceptDimensionRecord{
		pk:           cdk,
		conceptCD:    "DEM|AGE:100",
		nameChar:     " 100 years old",
		conceptBlob:  "",
		adminColumns: ac,
	}

	var csvFile = strings.NewReader(csvString)
	r := csv.NewReader(csvFile)
	lines, err := r.ReadAll()
	assert.Nil(t, err, "Parsing error")

	cdkExpected, cdExpected := conceptDimensionFromString(lines[0])

	assert.Equal(t, cdkExpected, cdk)
	assert.Equal(t, cdExpected, cd)
}

func TestObservationFactFromString(t *testing.T) {
	csvString := `"482232","1000000060","Affy:221610_s_at","LCS-I2B2:D000109064","2009-01-16 00:00:00","@","1","N","E","79.30000","","\N","","2009-01-16 00:00:00","@","","\N","2010-09-28 11:15:00","2010-08-18 09:50:00","2010-09-28 11:40:00","DEMO","\N","1"
`
	textSearchIndex = 0

	ac := administrativeColumns{
		updateDate:      "2010-09-28 11:15:00",
		downloadDate:    "2010-08-18 09:50:00",
		importDate:      "2010-09-28 11:40:00",
		sourceSystemCD:  "DEMO",
		uploadID:        "\\N",
		textSearchIndex: "0",
	}

	ofk := observationFactPK{
		encounterNum: "482232",
		patientNum:   "1000000060",
		conceptCD:    "Affy:221610_s_at",
		providerID:   "LCS-I2B2:D000109064",
		startDate:    "2009-01-16 00:00:00",
		modifierCD:   "@",
		instanceNum:  "1",
	}

	of := observationFactRecord{
		pk:              ofk,
		valTypeCD:       "N",
		tValChar:        "E",
		nValNum:         "79.30000",
		valueFlagCD:     "",
		quantityNum:     "\\N",
		unitsCD:         "",
		endDate:         "2009-01-16 00:00:00",
		locationCD:      "@",
		observationBlob: "",
		confidenceNum:   "\\N",
		adminColumns:    ac,
	}

	var csvFile = strings.NewReader(csvString)
	r := csv.NewReader(csvFile)
	lines, err := r.ReadAll()
	assert.Nil(t, err, "Parsing error")

	ofkExpected, ofExpected := observationFactFromString(lines[0])

	assert.Equal(t, ofkExpected, ofk)
	assert.Equal(t, ofExpected, &of)
}
