package loader

import (
	"encoding/csv"
	"encoding/xml"
	"github.com/lca1/medco/services"
	"github.com/lca1/unlynx/lib"
	"gopkg.in/dedis/crypto.v0/abstract"
	"gopkg.in/dedis/onet.v1"
	"gopkg.in/dedis/onet.v1/log"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

// The different paths and handlers for all the file both for input and/or output
var (
	InputFilePaths = map[string]string{
		"ADAPTER_MAPPINGS":   "../data/original/AdapterMappings.xml",
		"SHRINE_ONTOLOGY":    "../data/original/shrine.csv",
		"LOCAL_ONTOLOGY":     "../data/original/i2b2.csv",
		"PATIENT_DIMENSION":  "../data/original/patient_dimension.csv",
		"CONCEPT_DIMENSION":  "../data/original/concept_dimension.csv",
		"MODIFIER_DIMENSION": "../data/original/modifier_dimension.csv",
		"OBSERVATION_FACT":   "../data/original/observation_fact.csv",
	}

	OutputFilePaths = map[string]string{
		"ADAPTER_MAPPINGS":         "../data/converted/AdapterMappings.xml",
		"SHRINE_ONTOLOGY":          "../data/converted/shrine.csv",
		"LOCAL_ONTOLOGY_CLEAR":     "../data/converted/i2b2.csv",
		"LOCAL_ONTOLOGY_SENSITIVE": "../data/converted/sensitive_tagged.csv",
		"PATIENT_DIMENSION":        "../data/converted/patient_dimension.csv",
		"CONCEPT_DIMENSION":        "../data/converted/concept_dimension.csv",
		"MODIFIER_DIMENSION":       "../data/converted/modifier_dimension.csv",
		"OBSERVATION_FACT":         "../data/converted/observation_fact.csv",
	}
)

// Testing defines whether we should run the DDT on test environment (locally) or using real nodes
var Testing bool // testing environment

const (
	// Header is a generic XML header suitable for use with the output of Marshal.
	// This is not automatically added to any output of this package,
	// it is provided as a convenience.
	Header = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n"
)

// ADAPTER_MAPPINGS.XML converter

// ParseAdapterMappings reads and parses the AdapterMappings.xml.
func ParseAdapterMappings() error {
	xmlInputFile, err := os.Open(InputFilePaths["ADAPTER_MAPPINGS"])
	if err != nil {
		log.Fatal("Error opening [AdapterMappings].xml")
		return err
	}
	defer xmlInputFile.Close()

	b, _ := ioutil.ReadAll(xmlInputFile)

	err = xml.Unmarshal(b, &Am)
	if err != nil {
		log.Fatal("Error marshaling [AdapterMappings].xml")
		return err
	}
	return nil
}

// ConvertAdapterMappings converts the old AdapterMappings.xml file. This file maps a shrine concept code to an i2b2 concept code
func ConvertAdapterMappings() error {

	// convert the data in temporary maps to make it easy to traverse
	localToShrine := make(map[string][]string)
	shrineToLocal := make(map[string][]string)
	createTempMaps(Am, localToShrine, shrineToLocal)

	ListSensitiveConceptsLocal = make(map[string][]string)

	// everything is fine so now we just place both the sensitive local and shrine concepts in RAM (in our global maps)
	StoreSensitiveLocalConcepts(localToShrine, shrineToLocal)

	numElementsDel := filterSensitiveEntries(&Am)
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
	err = enc.Encode(Am)
	if err != nil {
		log.Fatal("Error writing converted [AdapterMappings].xml")
		return err
	}

	return nil
}

// createTempMaps simply converts the data into two different maps so that it's easier to traverse the data
func createTempMaps(am AdapterMappings, localToShrine map[string][]string, shrineToLocal map[string][]string) {
	for _, entry := range am.ListEntries {
		// remove the first \ (in the adapter mappings file there are 2 '\\' in the beginning)
		shrineKey := StripByLevel(entry.Key[1:], 1, true)
		arrValues := make([]string, 0)
		for _, value := range entry.ListLocalKeys {
			localKey := StripByLevel(value[1:], 1, true)
			arrValues = append(arrValues, localKey)

			// if the local concept is already mapped to another shrine concept
			if _, ok := localToShrine[localKey]; ok {
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

// StoreSensitiveLocalConcepts stores the local sensitive concepts in a set that is kept in RAM during the entire loading (to make parsing the local ontology faster)
func StoreSensitiveLocalConcepts(localToShrine map[string][]string, shrineToLocal map[string][]string) {
	for shrineKey, val := range shrineToLocal {
		if containsMapString(ListSensitiveConceptsShrine, shrineKey) {
			for _, localKey := range val {
				appendSensitiveConcepts(localKey, shrineKey)
				recursivelyUpdateConceptMaps(localKey, localToShrine, shrineToLocal)
			}
		}
	}
}

// appendSensitiveConcepts simply appends/relates a shrine concept to a local concept if it has not been added before
func appendSensitiveConcepts(localKey string, shrineKey string) {
	// if local concept already added
	if _, ok := ListSensitiveConceptsLocal[localKey]; ok {
		exists := false
		for _, el := range ListSensitiveConceptsLocal[localKey] {
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
func recursivelyUpdateConceptMaps(localKey string, localToShrine map[string][]string, shrineToLocal map[string][]string) {
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
	TableShrineOntologyConceptEnc = make(map[string]*ShrineOntology)
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
				// concepts do not repeat on contrary to the modifiers
				so.NodeEncryptID = IDConcepts
				IDConcepts++
				TableShrineOntologyConceptEnc[so.Fullname] = so
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
			break
		}
		csvOutputFile.WriteString(so.ToCSVText() + "\n")
	}

	// copy the sensitive concept codes to the new csv files (it does not include the modifier concepts)
	for _, so := range TableShrineOntologyConceptEnc {
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
	for _, so := range TableShrineOntologyConceptEnc {
		path := so.Fullname
		for true {
			path = StripByLevel(path, 1, false)
			if path == "" {
				break
			}

			if val, ok := TableShrineOntologyConceptEnc[path]; ok {
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
func ParseLocalOntology(group *onet.Roster, entryPointIdx int) error {
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

	MapConceptPathToTag = make(map[string]TagAndID)
	mapConceptIDtoTagKeys := make([]string, 0)
	MapModifierPathToTag = make(map[string]TagAndID)
	mapModifierIDtoTagKeys := make([]string, 0)

	allSensitiveModifierIDs := make([]int64, 0)
	allSensitiveConceptIDs := make([]int64, 0)

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
		lo := LocalOntologyFromString(line)

		// if we find a mapping to a Shrine Ontology term
		if _, ok := ListSensitiveConceptsLocal[lo.Fullname]; ok {
			// add each shrine id (we need to replicate each concept if it matches to multiple shrine concepts)
			for _, sk := range ListSensitiveConceptsLocal[lo.Fullname] {
				if strings.ToLower(lo.FactTableColumn) == "modifier_cd" {
					shrineID := TableShrineOntologyModifierEnc[sk][0].NodeEncryptID
					// if the ID does not yet exist
					if _, ok := MapModifierPathToTag[lo.Fullname]; !ok {
						MapModifierPathToTag[lo.Fullname] = TagAndID{Tag: libunlynx.GroupingKey(-1), TagID: IDModifiers}
						IDModifiers++

						mapModifierIDtoTagKeys = append(mapModifierIDtoTagKeys, lo.Fullname)
						allSensitiveModifierIDs = append(allSensitiveModifierIDs, shrineID)
					}
				} else if strings.ToLower(lo.FactTableColumn) == "concept_cd" { // if it is a concept code
					shrineID := TableShrineOntologyConceptEnc[sk].NodeEncryptID
					// if the ID does not yet exist
					if _, ok := MapConceptPathToTag[lo.Fullname]; !ok {
						MapConceptPathToTag[lo.Fullname] = TagAndID{Tag: libunlynx.GroupingKey(-1), TagID: IDConcepts}
						IDConcepts++

						mapConceptIDtoTagKeys = append(mapConceptIDtoTagKeys, lo.Fullname)
						allSensitiveConceptIDs = append(allSensitiveConceptIDs, shrineID)
					}
				} else {
					log.Fatal("Incorrect code in the FactTable column:", strings.ToLower(lo.FactTableColumn))
				}
			}
		} else if _, ok := HasSensitiveParents(lo.Fullname); !ok {
			// Concept does not have a translation in the Shrine Ontology (consider as not sensitive)
			TableLocalOntologyClear[lo.Fullname] = lo
		}
	}

	taggedModifierValues, err := EncryptAndTag(allSensitiveModifierIDs, group, entryPointIdx)
	if err != nil {
		return err
	}

	taggedConceptValues, err := EncryptAndTag(allSensitiveConceptIDs, group, entryPointIdx)
	if err != nil {
		return err
	}

	// 'populate' map (Modifier codes)
	for i, id := range mapModifierIDtoTagKeys {
		var tmp = MapModifierPathToTag[id]
		tmp.Tag = taggedModifierValues[i]
		MapModifierPathToTag[id] = tmp
	}

	// 'populate' map (Concept codes)
	for i, id := range mapConceptIDtoTagKeys {
		var tmp = MapConceptPathToTag[id]
		tmp.Tag = taggedConceptValues[i]
		MapConceptPathToTag[id] = tmp
	}

	return nil
}

// HasSensitiveParents is a function that checks if a node whether in the LocalOntology or ConceptDimension has any siblings which are sensitive.
func HasSensitiveParents(conceptPath string) (string, bool) {
	temp := conceptPath

	isSensitive := false
	for temp != "" {
		temp = StripByLevel(temp, 1, false)
		if _, ok := ListSensitiveConceptsLocal[temp]; ok {
			isSensitive = true
			break
		}
	}
	return temp, isSensitive
}

// EncryptAndTag encrypts the elements and tags them to allow for the future comparison
func EncryptAndTag(list []int64, group *onet.Roster, entryPointIdx int) ([]libunlynx.GroupingKey, error) {

	// ENCRYPTION
	start := time.Now()
	listEncryptedElements := make(libunlynx.CipherVector, len(list))

	for i := int64(0); i < int64(len(list)); i++ {
		listEncryptedElements[i] = *libunlynx.EncryptInt(group.Aggregate, list[i])
	}
	log.LLvl1("Finished encrypting the sensitive data... (", time.Since(start), ")")

	// TAGGING
	start = time.Now()
	client := servicesmedco.NewMedCoClient(group.List[entryPointIdx], strconv.Itoa(entryPointIdx))
	_, result, tr, err := client.SendSurveyDDTRequestTerms(
		group, // Roster
		servicesmedco.SurveyID("tagging_loading_phase"), // SurveyID
		listEncryptedElements,                           // Encrypted query terms to tag
		false, // compute proofs?
		Testing,
	)

	if err != nil {
		log.Fatal("Error during DDT")
		return nil, err
	}

	totalTime := time.Since(start)

	tr.DDTRequestTimeCommunication = totalTime - tr.DDTRequestTimeExec

	log.LLvl1("DDT took: execution -", tr.DDTRequestTimeExec, "communication -", tr.DDTRequestTimeCommunication)

	log.LLvl1("Finished tagging the sensitive data... (", totalTime, ")")

	return result, nil
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
	// two new files are generated: one to store the non-sensitive data and another to store the sensitive data
	csvClearOutputFile, err := os.Create(OutputFilePaths["LOCAL_ONTOLOGY_CLEAR"])
	if err != nil {
		log.Fatal("Error opening [i2b2].csv")
		return err
	}
	defer csvClearOutputFile.Close()

	csvSensitiveOutputFile, err := os.Create(OutputFilePaths["LOCAL_ONTOLOGY_SENSITIVE"])
	if err != nil {
		log.Fatal("Error opening [sensitive_tagged].csv")
		return err
	}
	defer csvSensitiveOutputFile.Close()

	headerString := ""
	for _, header := range HeaderLocalOntology {
		headerString += "\"" + header + "\","
	}
	// remove the last ,
	csvClearOutputFile.WriteString(headerString[:len(headerString)-1] + "\n")
	csvSensitiveOutputFile.WriteString(headerString[:len(headerString)-1] + "\n")

	// non-sensitive
	for _, lo := range TableLocalOntologyClear {
		csvClearOutputFile.WriteString(lo.ToCSVText() + "\n")
	}

	// sensitive modifiers
	for _, el := range MapModifierPathToTag {
		csvSensitiveOutputFile.WriteString(LocalOntologySensitiveModifierToCSVText(&el.Tag, el.TagID) + "\n")
	}

	// sensitive concepts
	for _, el := range MapConceptPathToTag {
		csvSensitiveOutputFile.WriteString(LocalOntologySensitiveConceptToCSVText(&el.Tag, el.TagID) + "\n")
	}

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

// CONCEPT_DIMENSION.CSV converter

// ParseConceptDimension reads and parses the concept_dimension.csv.
func ParseConceptDimension() error {
	lines, err := readCSV("CONCEPT_DIMENSION")
	if err != nil {
		log.Fatal("Error in readCSV()")
		return err
	}

	TableConceptDimension = make(map[*ConceptDimensionPK]ConceptDimension)
	HeaderConceptDimension = make([]string, 0)
	MapConceptCodeToTag = make(map[string]int64)

	/* structure of concept_dimension.csv (in order):

	// PK
	"concept_path",

	// MANDATORY FIELDS
	"concept_cd",
	"name_char",
	"concept_blob",

	// ADMIN FIELDS
	"update_date",
	"download_date",
	"import_date",
	"sourcesystem_cd",
	"upload_id"
	*/

	for _, header := range lines[0] {
		HeaderConceptDimension = append(HeaderConceptDimension, header)
	}

	//skip header
	for _, line := range lines[1:] {
		cdk, cd := ConceptDimensionFromString(line)
		TableConceptDimension[cdk] = cd
	}

	return nil
}

// ConvertConceptDimension converts the old concept_dimension.csv file
func ConvertConceptDimension() error {
	csvOutputFile, err := os.Create(OutputFilePaths["CONCEPT_DIMENSION"])
	if err != nil {
		log.Fatal("Error opening [concept_dimension].csv")
		return err
	}
	defer csvOutputFile.Close()

	headerString := ""
	for _, header := range HeaderConceptDimension {
		headerString += "\"" + header + "\","
	}
	// remove the last ,
	csvOutputFile.WriteString(headerString[:len(headerString)-1] + "\n")

	for _, cd := range TableConceptDimension {
		// if the concept is non-sensitive -> keep it as it is
		if _, ok := TableLocalOntologyClear[cd.PK.ConceptPath]; ok {
			csvOutputFile.WriteString(cd.ToCSVText() + "\n")
			// if the concept is sensitive -> fetch its encrypted tag and tag_id
		} else if _, ok := MapConceptPathToTag[cd.PK.ConceptPath]; ok {
			temp := MapConceptPathToTag[cd.PK.ConceptPath].Tag
			csvOutputFile.WriteString(ConceptDimensionSensitiveToCSVText(&temp, MapConceptPathToTag[cd.PK.ConceptPath].TagID) + "\n")
			MapConceptCodeToTag[cd.ConceptCD] = MapConceptPathToTag[cd.PK.ConceptPath].TagID
			// if the concept does not exist in the LocalOntology and none of his siblings is sensitive
		} else if sensitiveParent, ok := HasSensitiveParents(cd.PK.ConceptPath); !ok {
			csvOutputFile.WriteString(cd.ToCSVText() + "\n")
			// if concept is not defined as sensitive but one of its parents is then we consider the tagID of the parent as its identifier
		} else {
			MapConceptCodeToTag[cd.ConceptCD] = MapConceptPathToTag[sensitiveParent].TagID
		}
	}

	return nil
}

// MODIFIER_DIMENSION.CSV converter

// ParseModifierDimension reads and parses the modifier_dimension.csv.
func ParseModifierDimension() error {
	lines, err := readCSV("MODIFIER_DIMENSION")
	if err != nil {
		log.Fatal("Error in readCSV()")
		return err
	}

	TableModifierDimension = make(map[*ModifierDimensionPK]ModifierDimension)
	HeaderModifierDimension = make([]string, 0)
	MapModifierCodeToTag = make(map[string]int64)

	/* structure of modifier_dimension.csv (in order):

	// PK
	"modifier_path",

	// MANDATORY FIELDS
	"modifier_cd",
	"name_char",
	"modifier_blob",

	// ADMIN FIELDS
	"update_date",
	"download_date",
	"import_date",
	"sourcesystem_cd",
	"upload_id"
	*/

	for _, header := range lines[0] {
		HeaderModifierDimension = append(HeaderModifierDimension, header)
	}

	//skip header
	for _, line := range lines[1:] {
		mdk, md := ModifierDimensionFromString(line)
		TableModifierDimension[mdk] = md
	}

	return nil
}

// ConvertModifierDimension converts the old modifier_dimension.csv file
func ConvertModifierDimension() error {
	csvOutputFile, err := os.Create(OutputFilePaths["MODIFIER_DIMENSION"])
	if err != nil {
		log.Fatal("Error opening [modifier_dimension].csv")
		return err
	}
	defer csvOutputFile.Close()

	headerString := ""
	for _, header := range HeaderModifierDimension {
		headerString += "\"" + header + "\","
	}
	// remove the last ,
	csvOutputFile.WriteString(headerString[:len(headerString)-1] + "\n")

	for _, md := range TableModifierDimension {
		// if the modifier is non-sensitive -> keep it as it is
		if _, ok := TableLocalOntologyClear[md.PK.ModifierPath]; ok {
			csvOutputFile.WriteString(md.ToCSVText() + "\n")
			// if the modifier is sensitive -> fetch its encrypted tag and tag_id
		} else if _, ok := MapModifierPathToTag[md.PK.ModifierPath]; ok {
			temp := MapModifierPathToTag[md.PK.ModifierPath].Tag
			csvOutputFile.WriteString(ModifierDimensionSensitiveToCSVText(&temp, MapModifierPathToTag[md.PK.ModifierPath].TagID) + "\n")
			MapModifierCodeToTag[md.ModifierCD] = MapModifierPathToTag[md.PK.ModifierPath].TagID
			// if the modifier does not exist in the LocalOntology and none of his siblings is sensitive
		} else if sensitiveParent, ok := HasSensitiveParents(md.PK.ModifierPath); !ok {
			csvOutputFile.WriteString(md.ToCSVText() + "\n")
			// if modifier is not defined as sensitive but one of its parents is then we consider the tagID of the parent as its identifier
		} else {
			MapModifierCodeToTag[md.ModifierCD] = MapModifierPathToTag[sensitiveParent].TagID
		}
	}

	return nil
}

// OBSERVATION_FACT.CSV converter

// ParseObservationFact reads and parses the observation_fact.csv.
func ParseObservationFact() error {
	lines, err := readCSV("OBSERVATION_FACT")
	if err != nil {
		log.Fatal("Error in readCSV()")
		return err
	}

	TableObservationFact = make(map[*ObservationFactPK]ObservationFact)
	HeaderObservationFact = make([]string, 0)

	/* structure of observation_fact.csv (in order):

	// PK
	"encounter_num",
	"patient_num",
	"concept_cd",
	"provider_id",
	"start_date",
	"modifier_cd",
	"instance_num",

	// MANDATORY FIELDS
	"valtype_cd",
	"tval_char",
	"nval_num",
	"valueflag_cd",
	"quantity_num",
	"units_cd",
	"end_date",
	"location_cd",
	"observation_blob",
	"confidence_num",

	// ADMIN FIELDS
	"update_date",
	"download_date",
	"import_date",
	"sourcesystem_cd",
	"upload_id",
	"text_search_index"
	*/

	for _, header := range lines[0] {
		HeaderObservationFact = append(HeaderObservationFact, header)
	}

	//skip header
	for _, line := range lines[1:] {
		ofk, of := ObservationFactFromString(line)
		TableObservationFact[ofk] = of
	}

	return nil
}

// ConvertObservationFact converts the old observation_fact.csv file
func ConvertObservationFact() error {
	csvOutputFile, err := os.Create(OutputFilePaths["OBSERVATION_FACT"])
	if err != nil {
		log.Fatal("Error opening [observation_fact].csv")
		return err
	}
	defer csvOutputFile.Close()

	headerString := ""
	for _, header := range HeaderObservationFact {
		headerString += "\"" + header + "\","
	}
	// remove the last ,
	csvOutputFile.WriteString(headerString[:len(headerString)-1] + "\n")

	log.LLvl1(len(MapConceptCodeToTag))
	log.LLvl1(len(MapModifierCodeToTag))

	for _, of := range TableObservationFact {
		// if the concept is sensitive we replace its code with the correspondent tag ID
		if _, ok := MapConceptCodeToTag[of.PK.ConceptCD]; ok {
			of.PK.ConceptCD = strconv.FormatInt(MapConceptCodeToTag[of.PK.ConceptCD], 10)
		}

		// if the modifier is sensitive we replace its code with the correspondent tag ID
		if _, ok := MapModifierCodeToTag[of.PK.ModifierCD]; ok {
			of.PK.ModifierCD = strconv.FormatInt(MapModifierCodeToTag[of.PK.ModifierCD], 10)
		}

		csvOutputFile.WriteString(of.ToCSVText() + "\n")
	}

	return nil
}
