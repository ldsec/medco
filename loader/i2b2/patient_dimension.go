package loaderi2b2

import (
	"encoding/csv"
	libunlynx "github.com/ldsec/unlynx/lib"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/onet/v3/log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// patientDimensionRecord s a record of the patient_dimension table
type patientDimensionRecord struct {
	pk             patientDimensionPK
	vitalStatusCD  string
	birthDate      string
	deathDate      string
	optionalFields []optionalFields
	adminColumns   administrativeColumns
	encryptedFlag  libunlynx.CipherText
}

// patientDimensionPK is the primary key of the Patient_Dimension table
type patientDimensionPK struct {
	patientNum string
}

// optionalFields table contains the optional fields
type optionalFields struct {
	valType string
	value   string
}

var (
	// mapNewPatientNum keeps track of the mapping between the old patient_num and the new one
	mapNewPatientNum map[string]string
)

// tablePatientDimension is the patient_dimension table
var tablePatientDimension map[patientDimensionPK]*patientDimensionRecord

// headerPatientDimension contains the header of the Patient_Dimension table
var headerPatientDimension []string

// patientDimensionFromString generates a patientDimensionRecord struct from a parsed line of a .csv file
func patientDimensionFromString(line []string, pk kyber.Point) (patientDimensionPK, *patientDimensionRecord) {
	pdk := patientDimensionPK{
		patientNum: line[0],
	}

	pd := patientDimensionRecord{
		pk:            pdk,
		vitalStatusCD: line[1],
		birthDate:     line[2],
		deathDate:     line[3],
	}

	size := len(line)

	// optional fields
	of := make([]optionalFields, 0)

	for i := 4; i < size-5; i++ {
		of = append(of, optionalFields{valType: headerPatientDimension[i], value: line[i]})
	}

	ac := administrativeColumns{
		updateDate:     line[size-5],
		downloadDate:   line[size-4],
		importDate:     line[size-3],
		sourceSystemCD: line[size-2],
		uploadID:       line[size-1],
	}

	ef := libunlynx.EncryptInt(pk, 1)

	pd.optionalFields = of
	pd.adminColumns = ac
	pd.encryptedFlag = *ef

	return pdk, &pd
}

// toCSVText writes the patientDimensionPK struct in a way that can be added to a .csv file - "","","", etc.
func (pdk patientDimensionPK) toCSVText() string {
	return "\"" + pdk.patientNum + "\""
}

// toCSVText writes the patientDimensionRecord struct in a way that can be added to a .csv file - "","","", etc.
func (pd patientDimensionRecord) toCSVText(empty bool) string {
	encryptedFlagString, err := pd.encryptedFlag.Serialize()
	if err != nil {
		log.Error("Error during serialization:", err)
		return ""
	}
	encodedEncryptedFlag := "\"" + encryptedFlagString + "\""

	of := pd.optionalFields
	ofString := ""
	if empty == false {
		for i := 0; i < len(of); i++ {
			// +4 because there is one pk field and 3 mandatory fields
			ofString += "\"" + of[i].value + "\","
		}

		acString := "\"" + pd.adminColumns.updateDate + "\"," + "\"" + pd.adminColumns.downloadDate + "\"," + "\"" + pd.adminColumns.importDate + "\"," + "\"" + pd.adminColumns.sourceSystemCD + "\"," + "\"" + pd.adminColumns.uploadID + "\""
		finalString := pd.pk.toCSVText() + ",\"" + pd.vitalStatusCD + "\"," + "\"" + pd.birthDate + "\"," + "\"" + pd.deathDate + "\"," + ofString[:len(ofString)-1] + "," + acString + "," + encodedEncryptedFlag

		return strings.Replace(finalString, `"\N"`, "", -1)
	}

	for i := 0; i < len(of); i++ {
		// +4 because there is one pk field and 3 mandatory fields
		ofString += ","
	}

	acString := "," + "," + "," + ","
	finalString := pd.pk.toCSVText() + "," + "," + "," + "," + ofString[:len(ofString)-1] + "," + acString + "," + encodedEncryptedFlag

	return strings.Replace(finalString, `"\N"`, "", -1)
}

// filterPatientDimension creates a table containing only the patients with sensitive observations.
// This table is meant to be fed to the dummies generator script
func filterPatientDimension(pk kyber.Point) error {

	lines, err := readCSV(inputFilePaths["PATIENT_DIMENSION"])
	if err != nil {
		log.Fatal("Error in readCSV()")
		return err
	}

	csvOutputFileSensitive, err := os.Create(outputFilesPathsSensitive["PATIENT_DIMENSION_FILTERED"].Path)
	if err != nil {
		log.Fatal("Error opening sensitive [PATIENT_DIMENSION_FILTERED].csv")
		return err
	}
	defer csvOutputFileSensitive.Close()

	writerSensitive := csv.NewWriter(csvOutputFileSensitive)
	defer writerSensitive.Flush()

	writerSensitive.Write(lines[0])

	headerPatientDimension = make([]string, 0)
	for _, header := range lines[0] {
		headerPatientDimension = append(headerPatientDimension, header)
	}

	//skip header
	for _, line := range lines[1:] {
		pdk, _ := patientDimensionFromString(line, pk)
		//if the patient has sensitive observations
		if _, ok := patientsWithSensitiveObs[pdk]; ok {
			writerSensitive.Write(line)
		}

	}

	return nil

}

// parsePatientDimension reads and parses the patient_dimension.csv. This also means adding the encrypted flag.
func parsePatientDimension(pk kyber.Point) error {
	lines, err := readCSV(inputFilePaths["PATIENT_DIMENSION"])
	if err != nil {
		log.Fatal("Error in readCSV()")
		return err
	}

	tablePatientDimension = make(map[patientDimensionPK]*patientDimensionRecord)
	headerPatientDimension = make([]string, 0)
	mapNewPatientNum = make(map[string]string)

	/* structure of patient_dimension.csv (in order):

	// pk
	"patient_num",

	// MANDATORY FIELDS
	"vital_status_cd",
	"birth_date",
	"death_date",

	// OPTIONAL FIELDS
	"sex_cd","
	"age_in_years_num",
	"language_cd",
	"race_cd",
	"marital_status_cd",
	"religion_cd",
	"zip_cd",
	"statecityzip_path",
	"income_cd",
	"patient_blob",

	// ADMIN FIELDS
	"update_date",
	"download_date",
	"import_date",
	"sourcesystem_cd",
	"upload_id"

	*/

	for _, header := range lines[0] {
		headerPatientDimension = append(headerPatientDimension, header)
	}

	//skip header
	for _, line := range lines[1:] {
		pdk, pd := patientDimensionFromString(line, pk)
		tablePatientDimension[pdk] = pd
	}

	return nil
}

// convertPatientDimension converts the old patient_dimension.csv file
// If emtpy is set to true all other data except the patient_num and encrypted_dummy_flag are set to empty
func convertPatientDimension(pk kyber.Point) error {

	csvOutputFileSensitive, err := os.Create(outputFilesPathsSensitive["PATIENT_DIMENSION"].Path)
	if err != nil {
		log.Fatal("Error opening sensitive [PATIENT_DIMENSION].csv")
		return err
	}

	csvOutputFileNonSensitive, err := os.Create(outputFilesPathsNonSensitive["PATIENT_DIMENSION"].Path)
	if err != nil {
		log.Fatal("Error opening non sensitive [PATIENT_DIMENSION].csv")
		return err
	}

	defer csvOutputFileSensitive.Close()
	defer csvOutputFileNonSensitive.Close()

	headerString := ""
	for _, header := range headerPatientDimension {
		headerString += "\"" + header + "\","
	}

	// re-randomize the patient_num
	totalNbrPatients := len(tablePatientDimension) + len(tableDummyToPatient)

	rand.Seed(time.Now().UnixNano())
	perm := rand.Perm(totalNbrPatients)

	// remove the last ,
	csvOutputFileSensitive.WriteString(headerString[:len(headerString)-1] + "\n")
	csvOutputFileNonSensitive.WriteString(headerString[:len(headerString)-1] + "\n")

	i := 0
	for _, pdp := range tablePatientDimension {
		pd := *pdp
		mapNewPatientNum[pd.pk.patientNum] = strconv.FormatInt(int64(perm[i]), 10)
		pd.pk.patientNum = strconv.FormatInt(int64(perm[i]), 10)
		csvOutputFileNonSensitive.WriteString(pd.toCSVText(false) + "\n")
		// if the patient has sensitive observations, we also put him in the sensitive project
		if _, ok := patientsWithSensitiveObs[pdp.pk]; ok {
			csvOutputFileSensitive.WriteString(pd.toCSVText(true) + "\n")
		}
		i++
	}

	// add dummies
	for dummyNum, patientNum := range tableDummyToPatient {
		mapNewPatientNum[dummyNum] = strconv.FormatInt(int64(perm[i]), 10)

		patient := *tablePatientDimension[patientDimensionPK{patientNum: patientNum}]
		patient.pk.patientNum = strconv.FormatInt(int64(perm[i]), 10)
		ef := libunlynx.EncryptInt(pk, 0)
		patient.encryptedFlag = *ef

		csvOutputFileSensitive.WriteString(patient.toCSVText(true) + "\n")
		i++
	}

	// write mapNewPatientNum to csv
	csvOutputNewPatientNumFile, err := os.Create(outputFilesPathsSensitive["NEW_PATIENT_NUM"].Path)
	if err != nil {
		log.Fatal("Error opening [new_patient_num].csv")
		return err
	}
	defer csvOutputNewPatientNumFile.Close()

	csvOutputNewPatientNumFile.WriteString("\"old_patient_num\",\"new_patient_num\"\n")

	for key, value := range mapNewPatientNum {
		csvOutputNewPatientNumFile.WriteString("\"" + key + "\"," + "\"" + value + "\"\n")
	}

	return nil
}
