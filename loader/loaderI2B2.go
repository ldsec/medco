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
		"ADAPTER_MAPPINGS":         "../data/converted/AdapterMappings.xml",
		"PATIENT_DIMENSION":        "../data/converted/patient_dimension.csv",
		"SHRINE_ONTOLOGY":          "../data/converted/shrine.csv",
		"LOCAL_ONTOLOGY_CLEAR":     "../data/converted/i2b2.csv",
		"LOCAL_ONTOLOGY_SENSITIVE": "../data/converted/i2b2.csv",
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

	// AdapterMappings maps a shrine ontology sensitive concept or modifier concept to the local ontology (we need this to know which concepts from the local ontology are sensitive)
	var am AdapterMappings

	err = xml.Unmarshal(b, &am)
	if err != nil {
		log.Fatal("Error marshaling [AdapterMappings].xml")
		return err
	}

	// convert the data in temporary maps to make it easy to traverse
	localToShrine := make(map[string][]string)
	shrineToLocal := make(map[string][]string)
	createTempMaps(am, localToShrine, shrineToLocal)

	ListSensitiveConceptsLocal = make(map[string][]string)

	// everything is fine so now we just place both the sensitive local and shrine concepts in RAM (in our global maps)
	storeSensitiveLocalConcepts(am, localToShrine, shrineToLocal)

	numElementsDel := filterSensitiveEntries(&am)
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
	err = enc.Encode(am)
	if err != nil {
		log.Fatal("Error writing converted [AdapterMappings].xml")
		return err
	}

	return nil
}

// createTempMaps simply converts the data into two different so that it's easier to traverse the data
func createTempMaps(am AdapterMappings, localToShrine map[string][]string, shrineToLocal map[string][]string) {
	for _,entry := range am.ListEntries {
		shrineKey := StripByLevel(entry.Key[1:], 1, true)
		arrValues := make([]string, 0)
		for _,value := range entry.ListLocalKeys {
			localKey := StripByLevel(value[1:], 1, true)
			arrValues = append(arrValues, localKey)

			// if the local concept is already mapped to another shrine concept
			if _,ok := localToShrine[localKey]; ok {
				localToShrine[localKey] = append(localToShrine[localKey], shrineKey)
			} else {
				aux := make([]string, 0)
				aux = append(aux, shrineKey)
				localToShrine[localKey] = aux
			}
		}
		shrineToLocal[shrineKey] = arrValues
	}
}

// storeSensitiveLocalConcepts stores the local sensitive concepts in a set that is kept in RAM during the entire loading (to make parsing the local ontology faster)
func storeSensitiveLocalConcepts(am AdapterMappings, localToShrine map[string][]string, shrineToLocal map[string][]string) {
	for shrineKey,val := range shrineToLocal {
		if containsMapString(ListSensitiveConceptsShrine, shrineKey){
			for _, localKey := range val {
				appendSensitiveConcepts(localKey,shrineKey)
				recursivelyUpdateConceptMaps(localKey, localToShrine, shrineToLocal)
			}
		}
	}
}

// appendSensitiveConcepts simply appends a shrine concept to a local concept if it has not been added before
func appendSensitiveConcepts(localKey string, shrineKey string){
	// if local concept already added
	if _,ok := ListSensitiveConceptsLocal[localKey]; ok {

		exists := false
		for _, el := range ListSensitiveConceptsLocal[localKey]{
			if shrineKey == el {
				exists = true
			}
		}
		if !exists {
			ListSensitiveConceptsLocal[localKey] = append(ListSensitiveConceptsLocal[localKey], shrineKey)
		}
	} else {
		aux := make([]string, 0)
		aux = append(aux, shrineKey)
		ListSensitiveConceptsLocal[localKey] = aux
	}
}

// recursivelyUpdateConceptMaps recursively checks if there is any shrine concept that maps to a sensitive local concept. If true all of its values (local concepts) should be set to sensitive... rinse/repeat
func recursivelyUpdateConceptMaps(localKey string,  localToShrine map[string][]string, shrineToLocal map[string][]string){
	for _, shrineKey := range localToShrine[localKey] {
		// the local concept is sensitive but its shrine key is not
		if !containsMapString(ListSensitiveConceptsShrine, shrineKey) {
			ListSensitiveConceptsShrine[shrineKey] = true
			for _, localKey := range shrineToLocal[shrineKey] {
				appendSensitiveConcepts(localKey, shrineKey)
				recursivelyUpdateConceptMaps(localKey, localToShrine, shrineToLocal)
			}
		}
	}
}

// filterSensitiveEntries filters out (removes) the <key>, <values> pair(s) that belong to sensitive concepts
func filterSensitiveEntries(am *AdapterMappings) int {
	deleted := 0
	for i := range am.ListEntries {
		j := i - deleted
		// remove the table value from the first key value like \\SHRINE
		if containsMapString(ListSensitiveConceptsShrine, StripByLevel(am.ListEntries[j].Key[1:], 1, true)) {
			am.ListEntries = am.ListEntries[:j+copy(am.ListEntries[j:], am.ListEntries[j+1:])]
			deleted++

		}
	}

	return deleted
}

func containsMapString(m map[string]bool, e string) bool {
	if _, ok := m[e]; ok {
		return true
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

		if containsMapString(ListSensitiveConceptsShrine, so.Fullname) { // if it is a sensitive concept
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
		path := so.Fullname
		for true {
			path = StripByLevel(path, 1, false)
			if path == "" {
				break
			}

			if val, ok := TableShrineOntologyEnc[path]; ok {
				val.ChildrenEncryptIDs = append(val.ChildrenEncryptIDs, so.NodeEncryptID)
			}

		}
	}

	for path, soArr := range TableShrineOntologyModifierEnc {
		for true {
			path = StripByLevel(path, 1, false)
			if path == "" {
				break
			}

			if val, ok := TableShrineOntologyModifierEnc[path]; ok {
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
	TableLocalOntologyEnc = make(map[string]*LocalOntology)

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

	//skip header
	for _, line := range lines[1:] {
		so := LocalOntologyFromString(line)

		check, mapping := true,"asdsad"

		// if we find a mapping to a Shrine Ontology term
		if check && mapping != ""{
			if containsMapString(ListSensitiveConceptsShrine, mapping) { // if it is a sensitive concept
				log.LLvl1(so.Fullname, mapping)
				TableLocalOntologyEnc[so.Fullname] = so // if a modifier already exists (it's repeated it simply gets overwritten)
			} else {
				TableLocalOntologyClear[so.Fullname] = so
			}
		} else {
			// Concept does not have a translation in the Shrine Ontology (consider as not sensitive)
			TableLocalOntologyClear[so.Fullname] = so
		}
	}

	return nil
}

// StripByLevel strips the concept path based on /. The number represents the stripping level, in other words,
// if number = 1 we strip the first element enclosed in /****/ and then on. Order means which side we start stripping: true (left-to-right),
// false (right-to-left)
func StripByLevel(conceptPath string, number int, order bool) string {
	// remove the first and last \
	conceptPath = conceptPath[1 : len(conceptPath)-1]
	pathContainer := strings.Split(conceptPath, "\\")

	if order {
		for i := 0; i < number; i++ {
			if len(pathContainer) == 0 {
				break
			}

			// reduce a 'layer' at the time -  e.g. \\Admit Diagnosis\\Leg -> \\Leg
			pathContainer = pathContainer[1:]
		}
	} else {
		for i := 0; i < number; i++ {
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
