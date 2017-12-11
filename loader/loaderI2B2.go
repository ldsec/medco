package loader

import (
	"encoding/csv"
	"encoding/xml"
	"gopkg.in/dedis/crypto.v0/abstract"
	"gopkg.in/dedis/onet.v1/log"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// The different paths and handlers for all the file both for input and/or output
var (
	InputFilePaths = map[string]string{
		"ADAPTER_MAPPINGS":  "../data/original/AdapterMappings.xml",
		"PATIENT_DIMENSION": "../data/original/patient_dimension.csv",
		"SHRINE_ONTOLOGY":   "../data/original/shrine.csv",
		"LOCAL_ONTOLOGY":    "../data/original/i2b2.csv",
	}

	OutputFilePaths = map[string]string{
		"ADAPTER_MAPPINGS":  		"../data/converted/AdapterMappings.xml",
		"PATIENT_DIMENSION": 		"../data/converted/patient_dimension.csv",
		"SHRINE_ONTOLOGY":   		"../data/converted/shrine.csv",
		"LOCAL_ONTOLOGY_CLEAR":		"../data/converted/i2b2.csv",
		"LOCAL_ONTOLOGY_SENSITIVE":	"../data/converted/i2b2.csv",
	}
)

const (
	// Header is a generic XML header suitable for use with the output of Marshal.
	// This is not automatically added to any output of this package,
	// it is provided as a convenience.
	Header = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n"
)

// ADAPTER_MAPPINGS.XML converter

// ConvertAdapterMappings converts the old AdapterMappings.xml file. This file maps a shrine concept code to an i2b2 concept code
func ConvertAdapterMappings() error {
	xmlInputFile, err := os.Open(InputFilePaths["ADAPTER_MAPPINGS"])
	if err != nil {
		log.Fatal("Error opening [AdapterMappings].xml")
		return err
	}
	defer xmlInputFile.Close()

	b, _ := ioutil.ReadAll(xmlInputFile)

	err = xml.Unmarshal(b,&am)
	if err != nil {
		log.Fatal("Error marshaling [AdapterMappings].xml")
		return err
	}

	// filter out sensitive entries
	var amCopy AdapterMappings = am

	log.LLvl1(len(amCopy.ListEntries))

	numElementsDel := FilterSensitiveEntries(&amCopy)
	log.Lvl2(numElementsDel, "entries deleted")

	xmlOutputFile, err := os.Create(OutputFilePaths["ADAPTER_MAPPINGS"])
	if err != nil {
		log.Fatal("Error creating converted [AdapterMappings].xml")
		return err
	}
	xmlOutputFile.Write([]byte(Header))

	xmlWriter := io.Writer(xmlOutputFile)

	enc := xml.NewEncoder(xmlWriter)
	enc.Indent("", "\t")
	err = enc.Encode(amCopy)
	if err != nil {
		log.Fatal("Error writing converted [AdapterMappings].xml")
		return err
	}

	return nil
}

// FilterSensitiveEntries filters out (removes) the <key>, <values> pair(s) that belong to sensitive concepts
func FilterSensitiveEntries(am *AdapterMappings) int {
	deleted := 0
	for i := range am.ListEntries {
		j := i - deleted
		// remove the table value from the key value like \\SHRINE or \\i2b2_DEMO
		if containsArrayString(ListSensitiveConcepts, "\\"+strings.SplitN((am.ListEntries[j].Key)[2:], "\\", 2)[1]) {
			am.ListEntries = am.ListEntries[:j+copy(am.ListEntries[j:], am.ListEntries[j+1:])]
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

func readCSV(filename string) ([][]string, error) {
	csvInputFile, err := os.Open(InputFilePaths[filename])
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

// SHRINE.CSV converter (shrine ontology)

// ParseShrineOntology reads and parses the shrine.csv.
func ParseShrineOntology() error {
	lines, err := readCSV("SHRINE_ONTOLOGY")
	if err != nil {
		log.Fatal("Error in readCSV()")
		return err
	}

	// initialize container structs and counters
	IDModifiers = 0
	IDConcepts = 0
	TableShrineOntologyClear = make(map[string]*ShrineOntology)
	TableShrineOntologyEnc = make(map[string]*ShrineOntology)
	TableShrineOntologyModifierEnc = make(map[string][]*ShrineOntology)
	HeaderShrineOntology = make([]string, 0)

	/* structure of patient_dimension.csv (in order):

	// MANDATORY FIELDS
	"c_hlevel",
	"c_fullname",
	"c_name",
	"c_synonym_cd",
	"c_visualattributes",
	"c_totalnum",
	"c_basecode",
	"c_metadataxml",
	"c_facttablecolumn",
	"c_tablename",
	"c_columnname",
	"c_columndatatype",
	"c_operator",
	"c_dimcode",
	"c_comment",
	"c_tooltip",

	// ADMIN FIELDS
	"update_date",
	"download_date",
	"import_date",
	"sourcesystem_cd",

	// MANDATORY FIELDS
	"valuetype_cd",
	"m_applied_path",
	"m_exclusion_cd"

	*/

	for _, header := range lines[0] {
		HeaderShrineOntology = append(HeaderShrineOntology, header)
	}

	//skip header
	for _, line := range lines[1:] {
		so := ShrineOntologyFromString(line)

		if containsArrayString(ListSensitiveConcepts, so.Fullname) { // if it is a sensitive concept
			so.ChildrenEncryptIDs = make([]int64, 0)

			// if it is a modifier
			if strings.ToLower(so.FactTableColumn) == "modifier_cd" {
				// if value already present in the map
				if val, ok := TableShrineOntologyModifierEnc[so.Fullname]; ok {
					so.NodeEncryptID = val[0].NodeEncryptID
					TableShrineOntologyModifierEnc[so.Fullname] = append(TableShrineOntologyModifierEnc[so.Fullname], so)
				} else {
					so.NodeEncryptID = IDModifiers
					IDModifiers++
					TableShrineOntologyModifierEnc[so.Fullname] = make([]*ShrineOntology, 0)
					TableShrineOntologyModifierEnc[so.Fullname] = append(TableShrineOntologyModifierEnc[so.Fullname], so)
				}
			} else if strings.ToLower(so.FactTableColumn) == "concept_cd" { // if it is a concept code
				so.NodeEncryptID = IDConcepts
				IDConcepts++
				TableShrineOntologyEnc[so.Fullname] = so
			} else {
				log.Fatal("Incorrect code in the FactTable column:", strings.ToLower(so.FactTableColumn))
			}
		} else {
			TableShrineOntologyClear[so.Fullname] = so
		}
	}

	return nil
}

// ConvertShrineOntology converts the old shrine.csv file (the shrine ontology)
func ConvertShrineOntology() error {
	csvOutputFile, err := os.Create(OutputFilePaths["SHRINE_ONTOLOGY"])
	if err != nil {
		log.Fatal("Error opening [shrine].csv")
		return err
	}
	defer csvOutputFile.Close()

	headerString := ""
	for _, header := range HeaderShrineOntology {
		headerString += "\"" + header + "\","
	}
	// remove the last ,
	csvOutputFile.WriteString(headerString[:len(headerString)-1] + "\n")

	UpdateChildrenEncryptIDs() //updates the ChildrenEncryptIDs of the internal and parent nodes

	// copy the non-sensitive concept codes to the new csv file and change the name of the ONTOLOGYVERSION to blabla_convert
	prefix := "\\SHRINE\\ONTOLOGYVERSION\\"

	for _, so := range TableShrineOntologyClear {
		// search the \SHRINE\ONTOLOGYVERSION\blabla and change the name to blabla_Converted
		if strings.HasPrefix(so.Fullname, prefix) && len(so.Fullname) > len(prefix) {
			newName := so.Fullname[:len(so.Fullname)-1] + "_Converted\\"
			so.Fullname = newName
			so.Name = newName
			so.DimCode = newName
			so.Tooltip = newName
		}
		csvOutputFile.WriteString(so.ToCSVText() + "\n")
	}

	// copy the sensitive concept codes to the new csv files (it does not include the modifier concepts)
	for _, so := range TableShrineOntologyEnc {
		//log.LLvl1(so.Fullname, so.NodeEncryptID, so.ChildrenEncryptIDs, so.VisualAttributes)
		csvOutputFile.WriteString(so.ToCSVText() + "\n")
	}

	// copy the sensitive modifier concept codes to the new csv files
	for _, soArr := range TableShrineOntologyModifierEnc {
		for _, so := range soArr {
			//log.LLvl1(so.Fullname, so.NodeEncryptID, so.ChildrenEncryptIDs, so.VisualAttributes)
			csvOutputFile.WriteString(so.ToCSVText() + "\n")
		}
	}

	return nil
}

// UpdateChildrenEncryptIDs updates the parent and internal concept nodes with the IDs of their respective children
func UpdateChildrenEncryptIDs() {
	for _, so := range TableShrineOntologyEnc {

		StripByLevel(so.Fullname, 1, false)


		path := so.Fullname[1 : len(so.Fullname)-1] // remove the first and last \
		pathContainer := strings.Split(path, "\\")

		for len(pathContainer) > 0 {
			// reduce a 'layer' at the time -  e.g. \\SHRINE\\Diagnosis\\Haematite\\Leg -> \\SHRINE\\Diagnosis\\Haematite
			pathContainer = pathContainer[:len(pathContainer)-1]
			conceptPath := strings.Join(pathContainer, "\\")

			// if we remove the first and last \ in the beginning when comparing we need add them again
			if val, ok := TableShrineOntologyEnc["\\"+conceptPath+"\\"]; ok {
				val.ChildrenEncryptIDs = append(val.ChildrenEncryptIDs, so.NodeEncryptID)
			}
		}
	}

	for path, soArr := range TableShrineOntologyModifierEnc {
		// remove the first and last \
		path = path[1 : len(path)-1]
		pathContainer := strings.Split(path, "\\")

		for len(pathContainer) > 0 {
			// reduce a 'layer' at the time -  e.g. \\Admit Diagnosis\\Leg -> \\Admit Diagnosis
			pathContainer = pathContainer[:len(pathContainer)-1]
			conceptPath := strings.Join(pathContainer, "\\")

			// if we remove the first and last \ in the beginning when comparing we need add them again
			if val, ok := TableShrineOntologyModifierEnc["\\"+conceptPath+"\\"]; ok {
				for _, el := range val {
					// no matter the element in the array they all have the same NodeEncryptID
					el.ChildrenEncryptIDs = append(el.ChildrenEncryptIDs, soArr[0].NodeEncryptID)
				}
			}

		}
	}
}

// I2B2.CSV converter (local ontology)

// ParseLocalOntology reads and parses the i2b2.csv.
func ParseLocalOntology() error {
	lines, err := readCSV("LOCAL_ONTOLOGY")
	if err != nil {
		log.Fatal("Error in readCSV()")
		return err
	}

	// initialize container structs and counters
	IDModifiers = 0
	IDConcepts = 0
	HeaderLocalOntology = make([]string, 0)
	TableLocalOntologyClear = make(map[string]*LocalOntology)

	/* structure of i2b2.csv (in order):

	// MANDATORY FIELDS
	"c_hlevel",
	"c_fullname",
	"c_name",
	"c_synonym_cd",
	"c_visualattributes",
	"c_totalnum",
	"c_basecode",
	"c_metadataxml",
	"c_facttablecolumn",
	"c_tablename",
	"c_columnname",
	"c_columndatatype",
	"c_operator",
	"c_dimcode",
	"c_comment",
	"c_tooltip",
	"m_applied_path",

	// ADMIN FIELDS
	"update_date",
	"download_date",
	"import_date",
	"sourcesystem_cd",

	// MANDATORY FIELDS
	"valuetype_cd",
	"m_exclusion_cd",
	"c_path",
	"c_symbol"
	"pcori_basecode" (only exists in the sensitive tagged output csv file)
	*/

	for _, header := range lines[0] {
		HeaderLocalOntology = append(HeaderLocalOntology, header)
	}

	// the pcori_basecode
	HeaderPatientDimension = append(HeaderPatientDimension, "pcori_basecode")

	// this the identifier that represents the root of tree (by default it's i2b2) with HLevel = "0"
	// TODO: search this dynamically
	rootIdentifier := "\\i2b2\\"

	//skip header
	for _, line := range lines[1:] {
		so := LocalOntologyFromString(line)

		pathConcept := strings.Trim(so.Fullname, rootIdentifier)

		// the root so add it to the
		if pathConcept == "" {
			TableLocalOntologyClear[so.Fullname] = so
		}
	}

	return nil
}

// FindLocalConceptAdapterMapping checks if a local concept is present in the adapter mappings and what its translation in shrine ontology terms
func FindLocalConceptAdapterMapping(conceptPath string) (bool, string) {
	for _,entry := range am.ListEntries {
		for _,value := range entry.ListLocalKeys {
			// remove the first / out of the two that exist //
			// strip the path
			if StripByLevel(value[1:], 1, true) == conceptPath {
				return true, StripByLevel(entry.Key[1:], 1, true)
			}
		}
	}
	return false, ""
}

// StripFirstLevel strips the concept path based on /. The number represents the stripping level, in other words,
// if number = 1 we strip the first element enclosed in /****/ and then on. Order means which side we start stripping: true (left-to-right),
// false (right-to-left)
func StripByLevel(conceptPath string, number int, order bool) string {
	// remove the first and last \
	conceptPath = conceptPath[1 : len(conceptPath)-1]
	pathContainer := strings.Split(conceptPath, "\\")

	if order {
		for i:=0; i<number; i++ {
			if len(pathContainer) == 0 {
				break
			}

			// reduce a 'layer' at the time -  e.g. \\Admit Diagnosis\\Leg -> \\Leg
			pathContainer = pathContainer[1:]
		}
	} else {
		for i:=0; i<number; i++ {
			if len(pathContainer) == 0 {
				break
			}

			// reduce a 'layer' at the time -  e.g. \\Admit Diagnosis\\Leg -> \\Admit Diagnosis
			pathContainer = pathContainer[:len(pathContainer)-1]
		}
	}
	conceptPathFinal := strings.Join(pathContainer, "\\")

	if conceptPathFinal == "" {
		return conceptPathFinal
	}

	// if not empty we remove the first and last \ in the beginning when comparing we need add them again
	return "\\" + conceptPathFinal + "\\"
}

// ConvertLocalOntology converts the old i2b2.csv file
func ConvertLocalOntology() error {


	return nil
}

// PATIENT_DIMENSION.CSV converter

// ParsePatientDimension reads and parses the patient_dimension.csv. This also means adding the encrypted flag.
func ParsePatientDimension(pk abstract.Point) error {
	lines, err := readCSV("PATIENT_DIMENSION")
	if err != nil {
		log.Fatal("Error in readCSV()")
		return err
	}

	TablePatientDimension = make(map[*PatientDimensionPK]PatientDimension)
	HeaderPatientDimension = make([]string, 0)

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
		HeaderPatientDimension = append(HeaderPatientDimension, header)
	}

	// the encrypted_flag term
	HeaderPatientDimension = append(HeaderPatientDimension, "enc_dummy_flag_cd")

	//skip header
	for _, line := range lines[1:] {
		pdk, pd := PatientDimensionFromString(line, pk)
		TablePatientDimension[pdk] = pd
	}

	return nil
}

// ConvertPatientDimension converts the old patient_dimension.csv file
func ConvertPatientDimension() error {
	csvOutputFile, err := os.Create(OutputFilePaths["PATIENT_DIMENSION"])
	if err != nil {
		log.Fatal("Error opening [patient_dimension].csv")
		return err
	}
	defer csvOutputFile.Close()

	headerString := ""
	for _, header := range HeaderPatientDimension {
		headerString += "\"" + header + "\","
	}
	// remove the last ,
	csvOutputFile.WriteString(headerString[:len(headerString)-1] + "\n")

	for _, pd := range TablePatientDimension {
		csvOutputFile.WriteString(pd.ToCSVText() + "\n")
	}

	return nil
}
