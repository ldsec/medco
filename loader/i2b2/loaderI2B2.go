package loaderi2b2

import (
	"encoding/csv"
	"encoding/xml"
	"github.com/armon/go-radix"
	"github.com/dedis/kyber"
	"github.com/dedis/onet"
	"github.com/dedis/onet/log"
	"github.com/lca1/medco-loader/loader"
	"github.com/lca1/medco/services"
	"github.com/lca1/unlynx/lib"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// Files is the object structure behind the files.toml
type Files struct {
	TableAccess 	  string
	AdapterMappings   string
	I2B2              string
	SHRINE            string
	DummyToPatient    string
	PatientDimension  string
	VisitDimension    string
	ConceptDimension  string
	ModifierDimension string
	ObservationFact   string
	OutputFolder      string
}

// The different paths and handlers for all the file both for input and/or output
var (
	InputFilePaths = map[string]string{
		"TABLE_ACCESS":		  "../../data/i2b2/original/table_access.csv",
		"ADAPTER_MAPPINGS":   "../../data/i2b2/original/AdapterMappings.xml",
		"SHRINE_ONTOLOGY":    "../../data/i2b2/original/shrine.csv",
		"LOCAL_ONTOLOGY":     "../../data/i2b2/original/i2b2.csv",
		"PATIENT_DIMENSION":  "../../data/i2b2/original/patient_dimension.csv",
		"VISIT_DIMENSION":    "../../data/i2b2/original/visit_dimension.csv",
		"CONCEPT_DIMENSION":  "../../data/i2b2/original/concept_dimension.csv",
		"MODIFIER_DIMENSION": "../../data/i2b2/original/modifier_dimension.csv",
		"OBSERVATION_FACT":   "../../data/i2b2/original/observation_fact.csv",
		"DUMMY_TO_PATIENT":   "../../data/i2b2/original/dummy_to_patient.csv",
	}

	OutputFilePaths = map[string]string{
		"TABLE_ACCESS":		  		"../../data/i2b2/converted/table_access.csv",
		"ADAPTER_MAPPINGS":         "../../data/i2b2/converted/AdapterMappings.xml",
		"SHRINE_ONTOLOGY":          "../../data/i2b2/converted/shrine.csv",
		"LOCAL_ONTOLOGY_CLEAR":     "../../data/i2b2/converted/i2b2.csv",
		"LOCAL_ONTOLOGY_SENSITIVE": "../../data/i2b2/converted/sensitive_tagged.csv",
		"PATIENT_DIMENSION":        "../../data/i2b2/converted/patient_dimension.csv",
		"NEW_PATIENT_NUM":          "../../data/i2b2/converted/new_patient_num.csv",
		"VISIT_DIMENSION":          "../../data/i2b2/converted/visit_dimension.csv",
		"NEW_ENCOUNTER_NUM":        "../../data/i2b2/converted/new_encounter_num.csv",
		"CONCEPT_DIMENSION":        "../../data/i2b2/converted/concept_dimension.csv",
		"MODIFIER_DIMENSION":       "../../data/i2b2/converted/modifier_dimension.csv",
		"OBSERVATION_FACT":         "../../data/i2b2/converted/observation_fact.csv",
	}

	FileBashPath = "24-load-i2b2-data.sh"

	FilePathsData = [...]string{
		"SHRINE_ONTOLOGY",
		"LOCAL_ONTOLOGY_CLEAR",
		"LOCAL_ONTOLOGY_SENSITIVE",
		"PATIENT_DIMENSION",
		"VISIT_DIMENSION",
		"CONCEPT_DIMENSION",
		"MODIFIER_DIMENSION",
		"OBSERVATION_FACT",

	}

	TablenamesData = [...]string{
		"shrine_ont.shrine",
		"i2b2metadata.i2b2",
		"i2b2metadata.sensitive_tagged",
		"i2b2demodata.patient_dimension",
		"i2b2demodata.visit_dimension",
		"i2b2demodata.concept_dimension",
		"i2b2demodata.modifier_dimension",
		"i2b2demodata.observation_fact",

	}
)

const (
	// Header is a generic XML header suitable for use with the output of Marshal.
	// This is not automatically added to any output of this package,
	// it is provided as a convenience.
	Header = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n"
)

// MAIN function

func replaceOutputFolder(folderPath string) {
	tokens := strings.Split(OutputFilePaths["TABLE_ACCESS"], "/")
	OutputFilePaths["TABLE_ACCESS"] = folderPath + tokens[len(tokens)-1]

	tokens = strings.Split(OutputFilePaths["ADAPTER_MAPPINGS"], "/")
	OutputFilePaths["ADAPTER_MAPPINGS"] = folderPath + tokens[len(tokens)-1]

	tokens = strings.Split(OutputFilePaths["SHRINE_ONTOLOGY"], "/")
	OutputFilePaths["SHRINE_ONTOLOGY"] = folderPath + tokens[len(tokens)-1]

	tokens = strings.Split(OutputFilePaths["LOCAL_ONTOLOGY_CLEAR"], "/")
	OutputFilePaths["LOCAL_ONTOLOGY_CLEAR"] = folderPath + tokens[len(tokens)-1]

	tokens = strings.Split(OutputFilePaths["LOCAL_ONTOLOGY_SENSITIVE"], "/")
	OutputFilePaths["LOCAL_ONTOLOGY_SENSITIVE"] = folderPath + tokens[len(tokens)-1]

	tokens = strings.Split(OutputFilePaths["PATIENT_DIMENSION"], "/")
	OutputFilePaths["PATIENT_DIMENSION"] = folderPath + tokens[len(tokens)-1]

	tokens = strings.Split(OutputFilePaths["NEW_PATIENT_NUM"], "/")
	OutputFilePaths["NEW_PATIENT_NUM"] = folderPath + tokens[len(tokens)-1]

	tokens = strings.Split(OutputFilePaths["VISIT_DIMENSION"], "/")
	OutputFilePaths["VISIT_DIMENSION"] = folderPath + tokens[len(tokens)-1]

	tokens = strings.Split(OutputFilePaths["NEW_ENCOUNTER_NUM"], "/")
	OutputFilePaths["NEW_ENCOUNTER_NUM"] = folderPath + tokens[len(tokens)-1]

	tokens = strings.Split(OutputFilePaths["CONCEPT_DIMENSION"], "/")
	OutputFilePaths["CONCEPT_DIMENSION"] = folderPath + tokens[len(tokens)-1]

	tokens = strings.Split(OutputFilePaths["MODIFIER_DIMENSION"], "/")
	OutputFilePaths["MODIFIER_DIMENSION"] = folderPath + tokens[len(tokens)-1]

	tokens = strings.Split(OutputFilePaths["OBSERVATION_FACT"], "/")
	OutputFilePaths["OBSERVATION_FACT"] = folderPath + tokens[len(tokens)-1]
}

// ConvertI2B2 it's the main function that performs a full conversion and loading of the I2B2 data
func ConvertI2B2(el *onet.Roster, entryPointIdx int, files Files, mapSensitive map[string]struct{}, databaseS loader.DBSettings, empty bool) error {

	ListSensitiveConcepts = mapSensitive

	// change input filepaths
	InputFilePaths["TABLE_ACCESS"] = files.TableAccess
	InputFilePaths["ADAPTER_MAPPINGS"] = files.AdapterMappings
	InputFilePaths["SHRINE_ONTOLOGY"] = files.SHRINE
	InputFilePaths["LOCAL_ONTOLOGY"] = files.I2B2
	InputFilePaths["PATIENT_DIMENSION"] = files.PatientDimension
	InputFilePaths["VISIT_DIMENSION"] = files.VisitDimension
	InputFilePaths["CONCEPT_DIMENSION"] = files.ConceptDimension
	InputFilePaths["MODIFIER_DIMENSION"] = files.ModifierDimension
	InputFilePaths["OBSERVATION_FACT"] = files.ObservationFact
	InputFilePaths["DUMMY_TO_PATIENT"] = files.DummyToPatient

	// change output filepaths
	replaceOutputFolder(files.OutputFolder)

	err := ParseTableAccess()
	if err != nil {
		return err
	}

	err = GenerateNewTableAccess()
	if err != nil {
		return err
	}

	log.Lvl2("--- Finished generating TABLE_ACCESS ---")

	err = ParseLocalOntology(el, entryPointIdx)
	if err != nil {
		return err
	}
	err = ConvertLocalOntology()
	if err != nil {
		return err
	}

	log.Lvl2("--- Finished converting LOCAL_ONTOLOGY ---")

	err = GenerateNewAdapterMappings()
	if err != nil {
		return err
	}

	log.Lvl2("--- Finished generating ADAPTER_MAPPINGS ---")

	err = ParseShrineOntologyHeader()
	if err != nil {
		return err
	}
	err = GenerateNewShrineOntology()
	if err != nil {
		return err
	}

	log.Lvl2("--- Finished generating SHRINE_ONTOLOGY ---")

	err = ParseDummyToPatient()
	if err != nil {
		return err
	}

	err = ParsePatientDimension(el.Aggregate)
	if err != nil {
		return err
	}
	err = ConvertPatientDimension(el.Aggregate, empty)
	if err != nil {
		return err
	}

	log.Lvl2("--- Finished converting PATIENT_DIMENSION ---")

	err = ParseVisitDimension()
	if err != nil {
		return err
	}
	err = ConvertVisitDimension(empty)
	if err != nil {
		return err
	}

	log.Lvl2("--- Finished converting VISIT_DIMENSION ---")

	err = ParseConceptDimension()
	if err != nil {
		return err
	}
	err = ConvertConceptDimension()
	if err != nil {
		return err
	}

	log.Lvl2("--- Finished converting CONCEPT_DIMENSION ---")

	err = ParseObservationFact()
	if err != nil {
		return err
	}
	err = ConvertObservationFact()
	if err != nil {
		return err
	}

	log.Lvl2("--- Finished converting OBSERVATION_FACT ---")

	/*err = GenerateLoadingDataScript(databaseS)
	if err != nil {
		log.Fatal("Error while generating the loading data .sh file", err)
		return err
	}

	err = LoadDataFiles()
	if err != nil {
		log.Fatal("Error while loading ontology .sql file", err)
		return err
	}*/

	return nil
}

// GenerateLoadingDataScript creates a load dataset .sql script (deletes the data in the corresponding tables and reloads the new 'protected' data)
func GenerateLoadingDataScript(databaseS loader.DBSettings) error {
	fp, err := os.Create(FileBashPath)
	if err != nil {
		return err
	}

	loading := `#!/usr/bin/env bash` + "\n" + "\n" + `PGPASSWORD=` + databaseS.DBpassword + ` psql -v ON_ERROR_STOP=1 -h "` + databaseS.DBhost +
		`" -U "` + databaseS.DBuser + `" -p ` + strconv.FormatInt(int64(databaseS.DBport), 10) + ` -d "` + databaseS.DBname + `" <<-EOSQL` + "\n"

	loading += "BEGIN;\n"

	for i := 0; i < len(TablenamesData); i++ {
		tokens := strings.Split(FilePathsData[i], "/")
		loading += `\copy ` + TablenamesData[i] + ` FROM 'files/` + tokens[1] + `' ESCAPE '"' DELIMITER ',' CSV;` + "\n"
	}
	loading += "COMMIT;\n"
	loading += "EOSQL"

	_, err = fp.WriteString(loading)
	if err != nil {
		return err
	}

	fp.Close()

	return nil
}

// LoadDataFiles executes the loading script for the new converted data
func LoadDataFiles() error {
	// Display just the stderr if an error occurs
	/*cmd := exec.Command("/bin/sh", FileBashPath[1])
	stderr := &bytes.Buffer{} // make sure to import bytes
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		log.LLvl1("Error when running command.  Error log:", stderr.String())
		log.LLvl1("Got command status:", err.Error())
		return err
	}*/

	return nil
}

// TABLE_ACCESS.CSV reader

func ParseTableAccess() error {
	lines, err := readCSV("TABLE_ACCESS")
	if err != nil {
		log.Fatal("Error in readCSV()")
		return err
	}

	// initialize container structs and counters
	TableAccessMap = make(map[string]*TableAccess)
	HeaderTableAccess = make([]string, 0)

	/* structure of table_access.csv (in order):

	// FIELDS
	"c_table_cd",
	"c_table_name",
	"c_protected_access",
	"c_hlevel",
	"c_fullname",
	"c_name",
	"c_synonym_cd",
	"c_visualattributes",
	"c_totalnum",
	"c_basecode",
	"c_metadataxml",
	"c_facttablecolumn",
	"c_dimtablename",
	"c_columnname",
	"c_columndatatype",
	"c_operator",
	"c_dimcode",
	"c_comment",
	"c_tooltip",
	"c_entry_date",
	"c_change_date",
	"c_status_cd",
	"valuetype_cd"

	*/

	for _, header := range lines[0] {
		HeaderTableAccess = append(HeaderTableAccess, header)
	}

	//skip header
	for _, line := range lines[1:] {
		ta := TableAccessFromString(line)
		TableAccessMap[ta.Fullname] = ta
	}

	return nil
}

// GenerateNewAdapterMappings generate (copy) the table_access.csv
func GenerateNewTableAccess() error {
	csvOutputFile, err := os.Create(OutputFilePaths["TABLE_ACCESS"])
	if err != nil {
		log.Fatal("Error opening [table_access].csv")
		return err
	}
	defer csvOutputFile.Close()

	headerString := ""
	for _, header := range HeaderTableAccess {
		headerString += "\"" + header + "\","
	}
	// remove the last ,
	csvOutputFile.WriteString(headerString[:len(headerString)-1] + "\n")


	for _, ta := range TableAccessMap {
		csvOutputFile.WriteString(ta.ToCSVText() + "\n")
	}

	return nil

}

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
	shrineToLocal := radix.New()
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

func AddNewAdapterMapping(soc, loc string, modifier bool) {
	listLocalKeys := make([]string,0)

	if modifier == false {
		tokens := strings.Split(loc, `\`)
		// skip things like \i2b2\
		if len(tokens) <= 3 {
			return
		}
		fullname := `\` + tokens[1] + `\` + tokens[2] + `\`
		listLocalKeys = append(listLocalKeys, `\\` + TableAccessMap[fullname].TableCD + loc)
	} else { //TODO
		// if it is a modifier we have to create a new concept and append the modifier
		// e.g. \Admit Diagnosis\
	}

	Am.ListEntries = append(Am.ListEntries, Entry{Key: listLocalKeys[0], ListLocalKeys: listLocalKeys})
}

// GenerateNewAdapterMappings creates a new Adapter Mappings where 1 shrine concept is unequivocally is associated with 1 local concept
func GenerateNewAdapterMappings() error {
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

// createTempMaps simply converts the data into one map and one radix tree so that it's easier to traverse and manage the data
func createTempMaps(am AdapterMappings, localToShrine map[string][]string, shrineToLocal *radix.Tree) {
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
		shrineToLocal.Insert(shrineKey, arrValues)
	}
}

// StoreSensitiveLocalConcepts stores the local sensitive concepts in a set that is kept in RAM during the entire loading (to make parsing the local ontology faster)
func StoreSensitiveLocalConcepts(localToShrine map[string][]string, shrineToLocal *radix.Tree) {
	// STEP 1: Pick sensitive shrine concept
	for shrineKey := range ListSensitiveConceptsShrine {
		// STEP 2: Traverse prefix and add other related sensitive concepts
		shrineToLocal.WalkPrefix(shrineKey, func(s string, v interface{}) bool {
			for _, localKey := range v.([]string) {
				appendSensitiveConcepts(localKey, shrineKey)
				recursivelyUpdateConceptMaps(localKey, localToShrine, shrineToLocal)
			}
			return false
		})
		// STEP 3: Delete prefix subtree
		shrineToLocal.DeletePrefix(shrineKey)
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
func recursivelyUpdateConceptMaps(localKey string, localToShrine map[string][]string, shrineToLocal *radix.Tree) {
	for _, shrineKey := range localToShrine[localKey] {
		// the local concept is sensitive but its shrine key is not
		if !containsMapString(ListSensitiveConceptsShrine, shrineKey) {
			ListSensitiveConceptsShrine[shrineKey] = true
			shrineToLocal.WalkPrefix(shrineKey, func(s string, v interface{}) bool {
				for _, localKey := range v.([]string) {
					appendSensitiveConcepts(localKey, shrineKey)
					recursivelyUpdateConceptMaps(localKey, localToShrine, shrineToLocal)
				}
				return false
			})
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

// DUMMY_TO_PATIENT.csv parser

// ParseDummyToPatient reads and parses the dummy_to_patient.csv.
func ParseDummyToPatient() error {
	lines, err := readCSV("DUMMY_TO_PATIENT")
	if err != nil {
		log.Fatal("Error in readCSV()")
		return err
	}

	TableDummyToPatient = make(map[string]string)

	/* structure of patient_dimension.csv (in order):

	"dummy",
	"patient"

	*/

	//skip header
	for _, line := range lines[1:] {
		TableDummyToPatient[line[0]] = line[1]
	}

	return nil
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

	/* structure of shrine.csv (in order):

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

func ParseShrineOntologyHeader() error {
	lines, err := readCSV("SHRINE_ONTOLOGY")
	if err != nil {
		log.Fatal("Error in readCSV()")
		return err
	}

	// initialize container structs and counters
	HeaderShrineOntology = make([]string, 0)

	/* structure of shrine.csv (in order):

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

	return nil
}

func GenerateNewShrineOntology() error {
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

	// manually add first two shrine rows
	csvOutputFile.WriteString(`"1","\SHRINE\ONTOLOGYVERSION\","ONTOLOGYVERSION","N","FH ",,,"","concept_cd","concept_dimension","concept_path","T","LIKE","\SHRINE\ONTOLOGYVERSION\","","ONTOLOGYVERSION\",,,,"SHRINE",,"@",` + "\n")
	csvOutputFile.WriteString(`"2","\SHRINE\ONTOLOGYVERSION\SHRINE_DEMO-Download_ONTOLOGY_V-1.18_1-20-14_Converted\","SHRINE_DEMO-Download_ONTOLOGY_V-1.18_1-20-14_Converted","N","LH ",,,"","concept_cd","concept_dimension","concept_path","T","LIKE","\SHRINE\ONTOLOGYVERSION\SHRINE_DEMO-Download_ONTOLOGY_V-1.18_1-20-14_Converted\","","ONTOLOGYVERSION\SHRINE_DEMO-Download_ONTOLOGY_V-1.18_1-20-14_Converted\",,,,"SHRINE",,"@",` + "\n")

	for _, so := range TableShrineOntologyClear {
		csvOutputFile.WriteString(so.ToCSVText() + "\n")
	}

	// copy the sensitive concept codes to the new csv files (it does not include the modifier concepts)
	for _, so := range TableShrineOntologyConceptEnc {
		csvOutputFile.WriteString(so.ToCSVText() + "\n")
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

// ParseLocalOntology reads and parses the i2b2.csv. and generates the shrine.csv and adapter_mappings.xml
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
	TableShrineOntologyClear = make(map[string]*ShrineOntology)
	TableShrineOntologyConceptEnc = make(map[string]*ShrineOntology)

	MapConceptPathToTag = make(map[string]TagAndID)
	mapConceptIDtoTagKeys := make([]string, 0)
	allSensitiveConceptIDs := make([]int64, 0)

	listEntries := make([]Entry,0)
	Am = AdapterMappings{ListEntries: listEntries}
	Am.Hostname = "cw-ptrevvett.MED.HARVARD.EDU"
	Am.TimeStamp = "2014-02-26T15:48:27.286-05:00"

	MapConceptPathToTag = make(map[string]TagAndID)
	//mapConceptIDtoTagKeys := make([]string, 0)

	MapModifiers = make(map[string]struct{})

	//allSensitiveConceptIDs := make([]int64, 0)

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

		// TODO for now we remove all synonyms from the i2b2 local ontology
		// if it is the original concept (N = original, Y = synonym)
		if strings.ToLower(lo.SynonymCD) == "n" || strings.ToLower(lo.SynonymCD) == "" {
			// create entry for shrine ontology (direct copy)
			so := ShrineOntologyFromLocalConcept(lo)

			// check if local ontology concept is a child of any of the sensitive concepts selected by the client
			_, sensitive := HasSensitiveParents(lo.Fullname)
			// if it is sensitive or has a sensitive parent
			if sensitive {
				ListSensitiveConcepts[lo.Fullname] = struct{}{}

				//TODO for now we remove all modifiers
				if strings.ToLower(so.FactTableColumn) == "modifier_cd" {
					//if _, ok := MapModifiers[so.Fullname]; !ok {
					//AddNewAdapterMapping(so.Fullname, lo.Fullname, true)
					//}
				} else {
					// concepts do not repeat on contrary to the modifiers
					so.NodeEncryptID = IDConcepts
					TableShrineOntologyConceptEnc[so.Fullname] = so

					// if the ID does not yet exist
					if _, ok := MapConceptPathToTag[lo.Fullname]; !ok {
						MapConceptPathToTag[lo.Fullname] = TagAndID{Tag: libunlynx.GroupingKey(-1), TagID: IDConcepts}
						mapConceptIDtoTagKeys = append(mapConceptIDtoTagKeys, lo.Fullname)
						allSensitiveConceptIDs = append(allSensitiveConceptIDs, IDConcepts)
					}
					IDConcepts++
				}
			} else {

				//TODO for now we remove all modifiers
				if strings.ToLower(so.FactTableColumn) == "modifier_cd" {
					//if _, ok := MapModifiers[so.Fullname]; !ok {
					//AddNewAdapterMapping(so.Fullname, lo.Fullname, true)
					//}
				} else {
					// add a new entry for the AdapterMappings (1-1 mapping between shrine and local concept)
					AddNewAdapterMapping(so.Fullname, lo.Fullname, false)
					// add a new entry i2b2.csv (same as before)
					TableLocalOntologyClear[lo.Fullname] = lo
					// add a new entry shrine.csv
					TableShrineOntologyClear[so.Fullname] = so
				}
			}
		}
	}

	taggedConceptValues, err := EncryptAndTag(allSensitiveConceptIDs, group, entryPointIdx)
	if err != nil {
		return err
	}

	// 'populate' map (Concept codes)
	for i, id := range mapConceptIDtoTagKeys {
		var tmp = MapConceptPathToTag[id]
		tmp.Tag = taggedConceptValues[i]
		MapConceptPathToTag[id] = tmp
	}

		/*// if we find a mapping to a Shrine Ontology term
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

						log.LLvl1(lo.Fullname)
						log.LLvl1(TableShrineOntologyConceptEnc[sk].Fullname, shrineID)

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
	}*/

	return nil
}

// HasSensitiveParents is a function that checks if a node whether in the LocalOntology or ConceptDimension has any siblings which are sensitive.
func HasSensitiveParents(conceptPath string) (string, bool) {
	temp := conceptPath

	isSensitive := false
	for temp != "" {
		if _, ok := ListSensitiveConcepts[temp]; ok {
			isSensitive = true
			break
		}
		temp = StripByLevel(temp, 1, false)
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
	log.Lvl2("Finished encrypting the sensitive data... (", time.Since(start), ")")

	log.LLvl1("LI", list)
	log.LLvl1("LIST", listEncryptedElements)
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

	log.Lvl2("DDT took: execution -", tr.DDTRequestTimeExec, "communication -", tr.DDTRequestTimeCommunication)

	log.Lvl2("Finished tagging the sensitive data... (", totalTime, ")")

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
	/*for _, el := range MapModifierPathToTag {
		csvSensitiveOutputFile.WriteString(LocalOntologySensitiveModifierToCSVText(&el.Tag, el.TagID) + "\n")
	}*/

	// sensitive concepts
	for _, el := range MapConceptPathToTag {
		csvSensitiveOutputFile.WriteString(LocalOntologySensitiveConceptToCSVText(&el.Tag, el.TagID) + "\n")
	}

	return nil

}

// PATIENT_DIMENSION.CSV converter

// ParsePatientDimension reads and parses the patient_dimension.csv. This also means adding the encrypted flag.
func ParsePatientDimension(pk kyber.Point) error {
	lines, err := readCSV("PATIENT_DIMENSION")
	if err != nil {
		log.Fatal("Error in readCSV()")
		return err
	}

	TablePatientDimension = make(map[PatientDimensionPK]PatientDimension)
	HeaderPatientDimension = make([]string, 0)
	MapNewPatientNum = make(map[string]string)

	/* structure of patient_dimension.csv (in order):

	// PK
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
		HeaderPatientDimension = append(HeaderPatientDimension, header)
	}

	//skip header
	for _, line := range lines[1:] {
		pdk, pd := PatientDimensionFromString(line, pk)
		TablePatientDimension[pdk] = pd
	}

	return nil
}

// ConvertPatientDimension converts the old patient_dimension.csv file
// If emtpy is set to true all other data except the patient_num and encrypted_dummy_flag are set to empty
func ConvertPatientDimension(pk kyber.Point, empty bool) error {
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

	// re-randomize the patient_num
	totalNbrPatients := len(TablePatientDimension) + len(TableDummyToPatient)
	rand.Seed(time.Now().UnixNano())
	perm := rand.Perm(totalNbrPatients)


	// remove the last ,
	csvOutputFile.WriteString(headerString[:len(headerString)-1] + "\n")

	i := 0
	for _, pd := range TablePatientDimension {
		MapNewPatientNum[pd.PK.PatientNum] = strconv.FormatInt(int64(perm[i]), 10)
		pd.PK.PatientNum = strconv.FormatInt(int64(perm[i]), 10)
		csvOutputFile.WriteString(pd.ToCSVText(empty) + "\n")
		i++
	}

	// add dummies
	for dummyNum, patientNum := range TableDummyToPatient {
		MapNewPatientNum[dummyNum] = strconv.FormatInt(int64(perm[i]), 10)

		patient := TablePatientDimension[PatientDimensionPK{PatientNum: patientNum}]
		patient.PK.PatientNum = strconv.FormatInt(int64(perm[i]), 10)
		ef := libunlynx.EncryptInt(pk, 0)
		patient.EncryptedFlag = *ef

		csvOutputFile.WriteString(patient.ToCSVText(empty) + "\n")
		i++
	}

	// write MapNewPatientNum to csv
	csvOutputNewPatientNumFile, err := os.Create(OutputFilePaths["NEW_PATIENT_NUM"])
	if err != nil {
		log.Fatal("Error opening [new_patient_num].csv")
		return err
	}
	defer csvOutputNewPatientNumFile.Close()

	csvOutputNewPatientNumFile.WriteString("\"old_patient_num\",\"new_patient_num\"\n")

	for key, value := range MapNewPatientNum {
		csvOutputNewPatientNumFile.WriteString("\"" + key + "\"," + "\"" + value + "\"\n")
	}

	return nil
}

// VISIT_DIMENSION.CSV converter

// ParseVisitDimension reads and parses the visit_dimension.csv.
func ParseVisitDimension() error {
	lines, err := readCSV("VISIT_DIMENSION")
	if err != nil {
		log.Fatal("Error in readCSV()")
		return err
	}

	TableVisitDimension = make(map[VisitDimensionPK]VisitDimension)
	HeaderVisitDimension = make([]string, 0)
	MapNewEncounterNum = make(map[VisitDimensionPK]VisitDimensionPK)
	MapPatientVisits = make(map[string][]string)
	MaxVisits = 0

	/* structure of visit_dimension.csv (in order):

	// PK
	"encounter_num",
	"patient_num",

	// MANDATORY FIELDS
	"active_status_cd",
	"start_date",
	"end_date",

	// OPTIONAL FIELDS
	"inout_cd",
	"location_cd",
	"location_path",
	"length_of_stay",
	"visit_blob",

	// ADMIN FIELDS
	"update_date",
	"download_date",
	"import_date",
	"sourcesystem_cd",
	"upload_id"

	*/

	for _, header := range lines[0] {
		HeaderVisitDimension = append(HeaderVisitDimension, header)
	}

	//skip header
	for _, line := range lines[1:] {
		vdk, vd := VisitDimensionFromString(line)
		TableVisitDimension[vdk] = vd

		// if patient does not exist
		if _, ok := MapPatientVisits[vdk.PatientNum]; !ok {
			// create array and add the encounter
			tmp := make([]string, 0)
			tmp = append(tmp, vdk.EncounterNum)
			MapPatientVisits[vdk.PatientNum] = tmp
		} else {
			// append encounter to array
			MapPatientVisits[vdk.PatientNum] = append(MapPatientVisits[vdk.PatientNum], vdk.EncounterNum)
		}

		if MaxVisits < len(MapPatientVisits[vdk.PatientNum]) {
			MaxVisits = len(MapPatientVisits[vdk.PatientNum])
		}
	}

	return nil
}

// ConvertVisitDimension converts the old visit_dimension.csv file. The means re-randomizing the encounter_num.
// If emtpy is set to true all other data except the patient_num and encounter_num are set to empty
func ConvertVisitDimension(empty bool) error {

	csvOutputFile, err := os.Create(OutputFilePaths["VISIT_DIMENSION"])
	if err != nil {
		log.Fatal("Error opening [visit_dimension].csv")
		return err
	}
	defer csvOutputFile.Close()

	headerString := ""
	for _, header := range HeaderVisitDimension {
		headerString += "\"" + header + "\","
	}

	// re-randomize the encounter_num
	totalNbrVisits := len(TableVisitDimension) + len(TableDummyToPatient)*MaxVisits
	rand.Seed(time.Now().UnixNano())
	perm := rand.Perm(totalNbrVisits)

	// remove the last ,
	csvOutputFile.WriteString(headerString[:len(headerString)-1] + "\n")

	i := 0
	for _, vd := range TableVisitDimension {
		MapNewEncounterNum[VisitDimensionPK{EncounterNum: vd.PK.EncounterNum, PatientNum: vd.PK.PatientNum}] = VisitDimensionPK{EncounterNum: strconv.FormatInt(int64(perm[i]), 10), PatientNum: MapNewPatientNum[vd.PK.PatientNum]}
		vd.PK.EncounterNum = strconv.FormatInt(int64(perm[i]), 10)
		vd.PK.PatientNum = MapNewPatientNum[vd.PK.PatientNum]
		csvOutputFile.WriteString(vd.ToCSVText(empty) + "\n")
		i++
	}

	// add dummies
	for dummyNum, patientNum := range TableDummyToPatient {
		listPatientVisits := MapPatientVisits[patientNum]

		for _, el := range listPatientVisits {
			MapNewEncounterNum[VisitDimensionPK{EncounterNum: el, PatientNum: dummyNum}] = VisitDimensionPK{EncounterNum: strconv.FormatInt(int64(perm[i]), 10), PatientNum: MapNewPatientNum[dummyNum]}
			visit := TableVisitDimension[VisitDimensionPK{EncounterNum: el, PatientNum: patientNum}]
			visit.PK.EncounterNum = strconv.FormatInt(int64(perm[i]), 10)
			visit.PK.PatientNum = MapNewPatientNum[dummyNum]
			csvOutputFile.WriteString(visit.ToCSVText(empty) + "\n")
			i++
		}
	}

	// write MapNewEncounterNum to csv
	csvOutputNewEncounterNumFile, err := os.Create(OutputFilePaths["NEW_ENCOUNTER_NUM"])
	if err != nil {
		log.Fatal("Error opening [new_encounter_num].csv")
		return err
	}
	defer csvOutputNewEncounterNumFile.Close()

	csvOutputNewEncounterNumFile.WriteString("\"old_encounter_num\",\"old_patient_num\",\"new_encounter_num\",\"new_patient_num\"\n")

	for key, value := range MapNewEncounterNum {
		csvOutputNewEncounterNumFile.WriteString("\"" + key.EncounterNum + "\"," + "\"" + key.PatientNum + "\"," + "\"" + value.EncounterNum + "\"," + "\"" + value.PatientNum + "\"\n")
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

// ParseObservationFact reads and parses the observation_fact_old.csv.
func ParseObservationFact() error {
	lines, err := readCSV("OBSERVATION_FACT")
	if err != nil {
		log.Fatal("Error in readCSV()")
		return err
	}

	TableObservationFact = make(map[*ObservationFactPK]ObservationFact)
	HeaderObservationFact = make([]string, 0)
	MapPatientObs = make(map[string][]*ObservationFactPK)
	MapDummyObs = make(map[string][]*ObservationFactPK)
	TextSearchIndex = 0

	/* structure of observation_fact_old.csv (in order):

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

	// EXTRA FIELDS (added during dummy generation)
	"cluster_label"
	*/

	for _, header := range lines[0] {
		HeaderObservationFact = append(HeaderObservationFact, header)
	}
	// remove "cluster_label"
	HeaderObservationFact = HeaderObservationFact[:len(HeaderObservationFact)-1]

	//skip header
	for _, line := range lines[1:] {
		ofk, of := ObservationFactFromString(line)
		TableObservationFact[ofk] = of

		// if patient does not exist
		if _, ok := MapPatientObs[ofk.PatientNum]; !ok {
			// create array and add the observation
			tmp := make([]*ObservationFactPK, 0)
			tmp = append(tmp, ofk)
			MapPatientObs[ofk.PatientNum] = tmp
		} else {
			// append encounter to array
			MapPatientObs[ofk.PatientNum] = append(MapPatientObs[ofk.PatientNum], ofk)
		}

		// if dummy
		if originalPatient, ok := TableDummyToPatient[ofk.PatientNum]; ok {
			MapDummyObs[ofk.PatientNum] = MapPatientObs[originalPatient]
		}
	}

	return nil
}

// ConvertObservationFact converts the old observation_fact_old.csv file
func ConvertObservationFact() error {
	rand.Seed(time.Now().UnixNano())

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

	for _, of := range TableObservationFact {
		copyObs := of

		// if dummy observation
		if _, ok := TableDummyToPatient[of.PK.PatientNum]; ok {
			// 1. choose a random observation from the original patient
			// 2. copy the data
			// 3. change patient_num and encounter_num
			listObs := MapDummyObs[of.PK.PatientNum]
			index := rand.Intn(len(listObs))

			copyObs = TableObservationFact[listObs[index]]
			// change patient_num and encounter_num
			tmp := MapNewEncounterNum[VisitDimensionPK{EncounterNum: copyObs.PK.EncounterNum, PatientNum: of.PK.PatientNum}]
			copyObs.PK = regenerateObservationPK(copyObs.PK, tmp.PatientNum, tmp.EncounterNum)
			// keep the same concept (and text_search_index) that was already there
			copyObs.PK.ConceptCD = of.PK.ConceptCD
			copyObs.AdminColumns.TextSearchIndex = of.AdminColumns.TextSearchIndex

			// delete observation from the list (so we don't choose it again)
			listObs[index] = listObs[len(listObs)-1]
			listObs = listObs[:len(listObs)-1]
			MapDummyObs[of.PK.PatientNum] = listObs

		} else { // if real observation
			// change patient_num and encounter_num
			tmp := MapNewEncounterNum[VisitDimensionPK{EncounterNum: of.PK.EncounterNum, PatientNum: of.PK.PatientNum}]
			copyObs.PK = regenerateObservationPK(copyObs.PK, tmp.PatientNum, tmp.EncounterNum)
		}

		// if the concept is sensitive we replace its code with the correspondent tag ID
		if _, ok := MapConceptCodeToTag[copyObs.PK.ConceptCD]; ok {
			copyObs.PK.ConceptCD = strconv.FormatInt(MapConceptCodeToTag[copyObs.PK.ConceptCD], 10)
		}

		csvOutputFile.WriteString(copyObs.ToCSVText() + "\n")
	}

	return nil
}

func regenerateObservationPK(ofk *ObservationFactPK, patientNum, encounterNum string) *ObservationFactPK {
	ofkNew := &ObservationFactPK{
		EncounterNum: encounterNum,
		PatientNum:   patientNum,
		ConceptCD:    ofk.ConceptCD,
		ProviderID:   ofk.ProviderID,
		StartDate:    ofk.StartDate,
		ModifierCD:   ofk.ModifierCD,
		InstanceNum:  ofk.InstanceNum,
	}
	return ofkNew
}
