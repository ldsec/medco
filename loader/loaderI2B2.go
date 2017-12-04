package loader

import (
	"encoding/xml"
	"os"
	"gopkg.in/dedis/onet.v1/log"
	"io/ioutil"
	"io"
	"encoding/csv"
	"gopkg.in/dedis/crypto.v0/abstract"
	"gopkg.in/dedis/crypto.v0/base64"
	"strings"
)

// ListConceptsPaths list all the sensitive concepts (paths)
var ListConceptsPaths []string

// The different paths and handlers for all the file both for input and/or output
var (
	InputFilePaths = map[string]string{
		"ADAPTER_MAPPINGS"	: "../data/original/AdapterMappings.xml",
		"PATIENT_DIMENSION"	: "../data/original/patient_dimension.csv",
	}

	OutputFilePaths = map[string]string{
		"ADAPTER_MAPPINGS": "../data/converted/AdapterMappings.xml",
		"PATIENT_DIMENSION"	: "../data/converted/patient_dimension.csv",
	}
)

const (
	// A generic XML header suitable for use with the output of Marshal.
	// This is not automatically added to any output of this package,
	// it is provided as a convenience.
	Header = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n"
)


// ADAPTER_MAPPINGS.XML converter

// ConvertAdapterMappings converts the old AdapterMappings.xml file. This file maps a shrine concept code to an i2b2 concept code
func ConvertAdapterMappings() error{
	xmlInputFile, err := os.Open(InputFilePaths["ADAPTER_MAPPINGS"])
	if err != nil {
		log.Fatal("Error opening [AdapterMappings].xml")
		return err
	}
	defer xmlInputFile.Close()

	b, _ := ioutil.ReadAll(xmlInputFile)

	var am AdapterMappings

	err = xml.Unmarshal(b, &am)
	if err != nil {
		log.Fatal("Error marshaling [AdapterMappings].xml")
		return err
	}

	// filter out sensitive entries
	numElementsDel := FilterSensitiveEntries(&am)
	log.Lvl2(numElementsDel,"entries deleted")

	xmlOutputFile, err := os.Create(OutputFilePaths["ADAPTER_MAPPINGS"])
	if err != nil {
		log.Fatal("Error creating converted [AdapterMappings].xml")
		return err
	}
	xmlOutputFile.Write([]byte(Header))

	xmlWriter := io.Writer(xmlOutputFile)

	enc := xml.NewEncoder(xmlWriter)
	enc.Indent("", "\t")
	err = enc.Encode(am)
	if err != nil {
		log.Fatal("Error writing converted [AdapterMappings].xml")
		return err
	}
	return nil
}


// FilterSensitiveEntries filters out (removes) the <key>, <values> pair(s) that belong to sensitive concepts
func FilterSensitiveEntries(am *AdapterMappings) int{
	m := am.ListEntries

	deleted := 0
	for i := range m {
		j := i - deleted
		if containsArrayString(ListConceptsPaths, m[j].Key){
			m = m[:j+copy(m[j:], m[j+1:])]
			deleted++
		}
	}

	return deleted
}

func containsArrayString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func readCSV(filename string) ([][]string, error){
	csvInputFile , err := os.Open(InputFilePaths[filename])
	if err != nil {
		log.Fatal("Error opening [" + strings.ToLower(filename) + "].csv")
		return nil, err
	}
	defer csvInputFile.Close()

	reader := csv.NewReader(csvInputFile)
	reader.Comma = ','

	lines, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Error reading [" + strings.ToLower(filename) + "].csv")
		return nil, err
	}

	return lines, nil
}


// PATIENT_DIMENSION.CSV converter

// ParsePatientDimension reads and parses the patient_dimension.csv. This also means adding the encrypted flag.
func ParsePatientDimension(pk abstract.Point) error{
	lines, err := readCSV("PATIENT_DIMENSION")
	if err != nil {
		log.Fatal("Error in readCSV()")
		return err
	}

	TablePatientDimension = make(map[PatientDimensionPK]PatientDimension)
	HeaderPatientDimension = make([]string,0)

	/* structure of patient_dimension.csv (in order):

	// PK
	"patient_num",

	// MANDATORY FIELDS
	"vital_status_cd",
	"birth_date",
	death_date",

	// OPTIONAL FIELDS
	"sex_cd","
	age_in_years_num",
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
		HeaderPatientDimension = append(HeaderPatientDimension,header)
	}

	// the encrypted_flag term
	HeaderPatientDimension = append(HeaderPatientDimension,"enc_dummy_flag_cd")

	//skip header
	for _, line := range lines[1:] {
		PatientDimensionFromString(line,pk)
	}

	return nil
}

// ConvertPatientDimension converts the old patient_dimension.csv file
func ConvertPatientDimension() error {
	csvOutputFile , err := os.Create(OutputFilePaths["PATIENT_DIMENSION"])
	if err != nil {
		log.Fatal("Error opening [patient_dimension].csv")
		return err
	}
	defer csvOutputFile.Close()

	headerString := ""
	for _,header := range HeaderPatientDimension {
		headerString += "\"" + header + "\","
	}
	// remove the last ,
	csvOutputFile.WriteString(headerString[:len(headerString)-1]+"\n")

	for pdk,pd := range TablePatientDimension {
		lineString := "\"" + pdk.PatientNum + "\"," + "\"" + pd.VitalStatusCD + "\"," + "\"" + pd.BirthDate + "\"," + "\"" + pd.DeathDate + "\","

		for i:=0; i<len(pd.OptionalFields); i++{
			// +4 because there is on pk field and 3 mandatory fields
			lineString += "\"" + pd.OptionalFields[HeaderPatientDimension[i+4]] + "\","
		}

		lineString += pd.AdminColumns.ToCSVText()

		b := pd.EncryptedFlag.ToBytes()
		lineString += ",\"" + base64.StdEncoding.EncodeToString(b) + "\""

		csvOutputFile.WriteString(lineString+"\n")
	}

	return nil
}
