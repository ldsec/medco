package loaderi2b2

import (
	"encoding/csv"
	"encoding/xml"
	"github.com/dedis/kyber"
	"github.com/dedis/onet"
	"github.com/dedis/onet/log"
	"github.com/lca1/medco-loader/loader"
	"github.com/lca1/medco/services"
	"github.com/lca1/unlynx/lib"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// Files is the object structure behind the files.toml
type Files struct {
	TableAccess       string
	Schemes           string
	Ontology          []string
	DummyToPatient    string
	PatientDimension  string
	VisitDimension    string
	ConceptDimension  string
	ModifierDimension string
	ObservationFact   string
	OutputFolder      string
}

// FileInfo contains the tablename where the .csv should be loaded and the output path
type FileInfo struct {
	TableName string
	Path      string
}

// The different paths and handlers for all the files both for input and/or output
var (
	InputFilePaths = map[string]string{
		"TABLE_ACCESS": "../../data/i2b2/original/table_access.csv",

		"ONTOLOGY_BIRN":        "../../data/i2b2/original/birn.csv",
		"ONTOLOGY_CUSTOM_META": "../../data/i2b2/original/custom_meta.csv",
		"ONTOLOGY_ICD10_ICD9":  "../../data/i2b2/original/icd10_icd9.csv",
		"ONTOLOGY_I2B2":        "../../data/i2b2/original/i2b2.csv",

		"DUMMY_TO_PATIENT":  "../../data/i2b2/original/dummy_to_patient.csv",
		"PATIENT_DIMENSION": "../../data/i2b2/original/patient_dimension.csv",
		"VISIT_DIMENSION":   "../../data/i2b2/original/visit_dimension.csv",
		"CONCEPT_DIMENSION": "../../data/i2b2/original/concept_dimension.csv",
		"OBSERVATION_FACT":  "../../data/i2b2/original/observation_fact.csv",
	}

	OutputFilePaths = map[string]FileInfo{
		"ADAPTER_MAPPINGS": {TableName: "", Path: "../../data/i2b2/converted/AdapterMappings.xml"},
		"SCHEMES":          {TableName: "i2b2metadata.schemes", Path: "../../data/i2b2/converted/schemes.csv"},
		"SENSITIVE_TAGGED": {TableName: "i2b2metadata.sensitive_tagged", Path: "../../data/i2b2/converted/sensitive_tagged.csv"},

		"TABLE_ACCESS_L":    {TableName: "i2b2metadata.table_access", Path: "../../data/i2b2/converted/local_table_access.csv"},
		"LOCAL_BIRN":        {TableName: "i2b2metadata.birn", Path: "../../data/i2b2/converted/local_birn.csv"},
		"LOCAL_CUSTOM_META": {TableName: "i2b2metadata.custom_meta", Path: "../../data/i2b2/converted/local_custom_meta.csv"},
		"LOCAL_ICD10_ICD9":  {TableName: "i2b2metadata.icd10_icd9", Path: "../../data/i2b2/converted/local_icd10_icd9.csv"},
		"LOCAL_I2B2":        {TableName: "i2b2metadata.i2b2", Path: "../../data/i2b2/converted/local_i2b2.csv"},

		"TABLE_ACCESS_S":     {TableName: "shrine_ont.table_access", Path: "../../data/i2b2/converted/shrine_table_access.csv"},
		"SHRINE_BIRN":        {TableName: "shrine_ont.birn", Path: "../../data/i2b2/converted/shrine_birn.csv"},
		"SHRINE_CUSTOM_META": {TableName: "shrine_ont.custom_meta", Path: "../../data/i2b2/converted/shrine_custom_meta.csv"},
		"SHRINE_ICD10_ICD9":  {TableName: "shrine_ont.icd10_icd9", Path: "../../data/i2b2/converted/shrine_icd10_icd9.csv"},
		"SHRINE_I2B2":        {TableName: "shrine_ont.i2b2", Path: "../../data/i2b2/converted/shrine_i2b2.csv"},

		"PATIENT_DIMENSION": {TableName: "i2b2demodata.patient_dimension", Path: "../../data/i2b2/converted/patient_dimension.csv"},
		"NEW_PATIENT_NUM":   {TableName: "", Path: "../../data/i2b2/converted/new_patient_num.csv"},
		"VISIT_DIMENSION":   {TableName: "i2b2demodata.visit_dimension.i2b2", Path: "../../data/i2b2/converted/visit_dimension.csv"},
		"NEW_ENCOUNTER_NUM": {TableName: "", Path: "../../data/i2b2/converted/new_encounter_num.csv"},
		"CONCEPT_DIMENSION": {TableName: "i2b2demodata.concept_dimension", Path: "../../data/i2b2/converted/concept_dimension.csv"},
		"OBSERVATION_FACT":  {TableName: "i2b2demodata.observation_fact", Path: "../../data/i2b2/converted/observation_fact.csv"},
	}

	FileBashPath = "24-load-i2b2-data.sh"
)

const (
	// Header is a generic XML header suitable for use with the output of Marshal.
	// This is not automatically added to any output of this package,
	// it is provided as a convenience.
	Header = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n"
)

// MAIN function

func generateOutputFiles(folderPath string) {
	// fixed demodata tables
	OutputFilePaths["ADAPTER_MAPPINGS"] = FileInfo{TableName: "", Path: "../../data/i2b2/converted/AdapterMappings.xml"}
	OutputFilePaths["PATIENT_DIMENSION"] = FileInfo{TableName: "i2b2demodata.patient_dimension", Path: "../../data/i2b2/converted/patient_dimension.csv"}
	OutputFilePaths["NEW_PATIENT_NUM"] = FileInfo{TableName: "", Path: "../../data/i2b2/converted/new_patient_num.csv"}
	OutputFilePaths["VISIT_DIMENSION"] = FileInfo{TableName: "i2b2demodata.visit_dimension.i2b2", Path: "../../data/i2b2/converted/visit_dimension.csv"}
	OutputFilePaths["NEW_ENCOUNTER_NUM"] = FileInfo{TableName: "", Path: "../../data/i2b2/converted/new_encounter_num.csv"}
	OutputFilePaths["CONCEPT_DIMENSION"] = FileInfo{TableName: "i2b2demodata.concept_dimension", Path: "../../data/i2b2/converted/concept_dimension.csv"}
	OutputFilePaths["OBSERVATION_FACT"] = FileInfo{TableName: "i2b2demodata.observation_fact", Path: "../../data/i2b2/converted/observation_fact.csv"}

	// fixed ontology tables
	OutputFilePaths["SENSITIVE_TAGGED"] = FileInfo{TableName: "i2b2metadata.sensitive_tagged", Path: "../../data/i2b2/converted/sensitive_tagged.csv"}
	OutputFilePaths["SCHEMES"] = FileInfo{TableName: "i2b2metadata.schemes", Path: "../../data/i2b2/converted/schemes.csv"}

	OutputFilePaths["TABLE_ACCESS_L"] = FileInfo{TableName: "i2b2metadata.table_access", Path: "../../data/i2b2/converted/local_table_access.csv"}
	OutputFilePaths["TABLE_ACCESS_S"] = FileInfo{TableName: "shrine_ont.table_access", Path: "../../data/i2b2/converted/shrine_table_access.csv"}

	for key, path := range InputFilePaths {
		if strings.HasPrefix(key, "ONTOLOGY_") {
			rawKey := strings.Split(key, "ONTOLOGY_")[1]
			tokens := strings.Split(path, "/")

			OutputFilePaths["LOCAL_"+rawKey] = FileInfo{TableName: "i2b2metadata." + strings.ToLower(rawKey), Path: folderPath + tokens[len(tokens)-1]}
			OutputFilePaths["SHRINE_"+rawKey] = FileInfo{TableName: "shrine_ont." + strings.ToLower(rawKey), Path: folderPath + tokens[len(tokens)-1]}
		}
	}
}

// ConvertI2B2 it's the main function that performs a full conversion and loading of the I2B2 data
func ConvertI2B2(el *onet.Roster, entryPointIdx int, files Files, mapSensitive map[string]struct{}, databaseS loader.DBSettings, empty bool) error {
	InputFilePaths = make(map[string]string)
	OutputFilePaths = make(map[string]FileInfo)

	ListSensitiveConcepts = mapSensitive

	// change input filepaths
	InputFilePaths["TABLE_ACCESS"] = files.TableAccess
	InputFilePaths["SCHEMES"] = files.Schemes

	if len(files.Ontology) == 0 {
		log.Fatal("No Ontology files were selected for conversion")
	}

	for _, name := range files.Ontology {
		tokens := strings.Split(name, "/")
		ontologyName := tokens[len(tokens)-1]
		InputFilePaths["ONTOLOGY_"+strings.ToUpper(strings.Split(ontologyName, ".")[0])] = name
	}
	InputFilePaths["PATIENT_DIMENSION"] = files.PatientDimension
	InputFilePaths["VISIT_DIMENSION"] = files.VisitDimension
	InputFilePaths["CONCEPT_DIMENSION"] = files.ConceptDimension
	InputFilePaths["OBSERVATION_FACT"] = files.ObservationFact
	InputFilePaths["DUMMY_TO_PATIENT"] = files.DummyToPatient

	// change output filepaths
	generateOutputFiles(files.OutputFolder)

	err := ParseTableAccess()
	if err != nil {
		return err
	}

	err = ConvertTableAccess()
	if err != nil {
		return err
	}

	log.Lvl2("--- Finished converting TABLE_ACCESS ---")

	err = ConvertLocalOntology(el, entryPointIdx)
	if err != nil {
		return err
	}

	log.Lvl2("--- Finished converting LOCAL_ONTOLOGY ---")

	err = GenerateAdapterMappings()
	if err != nil {
		return err
	}

	log.Lvl2("--- Finished generating ADAPTER_MAPPINGS ---")

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
/*func GenerateLoadingDataScript(databaseS loader.DBSettings) error {
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
}*/

// LoadDataFiles executes the loading script for the new converted data
/*func LoadDataFiles() error {
	// Display just the stderr if an error occurs
	cmd := exec.Command("/bin/sh", FileBashPath[1])
	stderr := &bytes.Buffer{} // make sure to import bytes
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		log.LLvl1("Error when running command.  Error log:", stderr.String())
		log.LLvl1("Got command status:", err.Error())
		return err
	}

	return nil
}*/

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

// TABLE_ACCESS.CSV reader

// ParseTableAccess reads and parses the table_access.csv.
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

// ConvertTableAccess generate (copy) the table_access.csv
func ConvertTableAccess() error {
	// local (we add a SENSITIVE_TAGGED entry)
	csvOutputFileL, err := os.Create(OutputFilePaths["TABLE_ACCESS_L"].Path)
	if err != nil {
		log.Fatal("Error opening [local_table_access].csv")
		return err
	}
	defer csvOutputFileL.Close()

	// shrine (copy as is)
	csvOutputFileS, err := os.Create(OutputFilePaths["TABLE_ACCESS_S"].Path)
	if err != nil {
		log.Fatal("Error opening [shrine_table_access].csv")
		return err
	}
	defer csvOutputFileS.Close()

	headerString := ""
	for _, header := range HeaderTableAccess {
		headerString += "\"" + header + "\","
	}
	// remove the last ,
	csvOutputFileL.WriteString(headerString[:len(headerString)-1] + "\n")
	csvOutputFileS.WriteString(headerString[:len(headerString)-1] + "\n")

	for _, ta := range TableAccessMap {
		csvOutputFileL.WriteString(ta.ToCSVText() + "\n")
		csvOutputFileS.WriteString(ta.ToCSVText() + "\n")
	}

	csvOutputFileL.WriteString(`"SENSITIVE_TAGGED","SENSITIVE_TAGGED","N","1","\medco\tagged\","MedCo Sensitive Tagged Ontology","N","CA ",,,,"concept_cd","concept_dimension","concept_path","T","LIKE","\medco\tagged\",,"MedCo Sensitive Tagged Ontology",,,,` + "\n")

	return nil
}

// ADAPTER_MAPPINGS.XML converter

func addNewAdapterMapping(loc string, modifier bool) {
	listLocalKeys := make([]string, 0)

	if modifier == false {
		for prefix := range TableAccessMap {
			if strings.HasPrefix(loc, prefix) {
				listLocalKeys = append(listLocalKeys, `\\`+TableAccessMap[prefix].TableCD+loc)
				continue
			}
		}
	}
	Am.ListEntries = append(Am.ListEntries, Entry{Key: listLocalKeys[0], ListLocalKeys: listLocalKeys})
}

// GenerateAdapterMappings creates a new Adapter Mappings where 1 shrine concept is unequivocally associated with 1 local concept
func GenerateAdapterMappings() error {
	xmlOutputFile, err := os.Create(OutputFilePaths["ADAPTER_MAPPINGS"].Path)
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

// SHRINE ontology converter

// GenerateNewShrineOntology generates all files for the shrine ontology (these may include multiples tables)
func GenerateNewShrineOntology() error {
	// initialize container structs and counters
	HeaderShrineOntology = []string{"c_hlevel",
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
		"update_date",
		"download_date",
		"import_date",
		"sourcesystem_cd",
		"valuetype_cd",
		"m_applied_path",
		"m_exclusion_cd"}

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

	for key := range InputFilePaths {
		if strings.HasPrefix(key, "ONTOLOGY_") {
			err := generateNewShrineTable(strings.Split(key, "ONTOLOGY_")[1])
			if err != nil {
				log.Fatal("Error generating [" + key + "].csv")
				return err
			}
		}
	}
	return nil
}

func generateNewShrineTable(rawName string) error {
	csvOutputFile, err := os.Create(OutputFilePaths["SHRINE_"+rawName].Path)
	if err != nil {
		log.Fatal("Error opening [" + rawName + "].csv")
		return err
	}
	defer csvOutputFile.Close()

	headerString := ""
	for _, header := range HeaderShrineOntology {
		headerString += "\"" + header + "\","
	}
	// remove the last ,
	csvOutputFile.WriteString(headerString[:len(headerString)-1] + "\n")

	UpdateChildrenEncryptIDs(rawName) //updates the ChildrenEncryptIDs of the internal and parent nodes

	for _, so := range TablesShrineOntology[rawName].Clear {
		csvOutputFile.WriteString(so.ToCSVText() + "\n")
	}

	// copy the sensitive concept codes to the new csv files (it does not include the modifier concepts)
	for _, so := range TablesShrineOntology[rawName].Sensitive {
		csvOutputFile.WriteString(so.ToCSVText() + "\n")
	}

	return nil
}

// UpdateChildrenEncryptIDs updates the parent and internal concept nodes with the IDs of their respective children (name identifies the name of the ontology table)
func UpdateChildrenEncryptIDs(name string) {
	for _, so := range TablesShrineOntology[name].Sensitive {
		path := so.Fullname
		for true {
			path = StripByLevel(path, 1, false)
			if path == "" {
				break
			}

			if val, ok := TablesShrineOntology[name].Sensitive[path]; ok {
				val.ChildrenEncryptIDs = append(val.ChildrenEncryptIDs, so.NodeEncryptID)
			}

		}
	}
}

// LOCAL ontology converter

// ConvertLocalOntology reads and parses all local ontology tables and generates the corresponding .csv(s) (local, shrine and adapter_mappings)
func ConvertLocalOntology(group *onet.Roster, entryPointIdx int) error {
	// initialize container structs and counters
	IDConcepts = 0
	TablesShrineOntology = make(map[string]ShrineTableInfo)
	MapConceptPathToTag = make(map[string]TagAndID)

	for key := range InputFilePaths {
		if strings.HasPrefix(key, "ONTOLOGY_") {
			rawName := strings.Split(key, "ONTOLOGY_")[1]
			err := ParseLocalTable(group, entryPointIdx, key)
			if err != nil {
				log.Fatal("Error parsing [" + strings.ToLower(rawName) + "].csv")
				return err
			}
			err = ConvertClearLocalTable(rawName)
			if err != nil {
				log.Fatal("Error converting [" + strings.ToLower(rawName) + "].csv")
				return err
			}
		}
	}

	err := ConvertSensitiveLocalTable()
	if err != nil {
		log.Fatal("Error converting [sensitive_tagged].csv")
		return err
	}

	return nil
}

// ParseLocalTable reads and parses the xxxx.csv (part of the local ontology)
// The adapter_mappings and shrine ontology are also generated based on the local ontology. Each local table (in i2b2metadata is replicated in the shrine_ont, with some minor changes)
func ParseLocalTable(group *onet.Roster, entryPointIdx int, name string) error {
	lines, err := readCSV(name)
	if err != nil {
		log.Fatal("Error in readCSV()")
		return err
	}
	rawName := strings.Split(name, "ONTOLOGY_")[1]

	HeaderLocalOntology = make([]string, 0)
	TableLocalOntologyClear = make(map[string]*LocalOntology)

	listConceptCD := make([]string, 0)
	allSensitiveConceptIDs := make([]int64, 0)
	listEntries := make([]Entry, 0)

	Am = AdapterMappings{ListEntries: listEntries}
	Am.Hostname = "cw-ptrevvett.MED.HARVARD.EDU"
	Am.TimeStamp = "2014-02-26T15:48:27.286-05:00"

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

			// if the flag AllSensitive is 'active'
			sensitive := false
			if AllSensitive == true {
				sensitive = true
			} else {
				// check if local ontology concept is a child of any of the sensitive concepts selected by the client
				_, sensitive = HasSensitiveParents(lo.Fullname)
			}

			// if it is sensitive or has a sensitive parent
			if sensitive {
				ListSensitiveConcepts[lo.Fullname] = struct{}{}

				//TODO for now we remove all modifiers
				if strings.ToLower(so.FactTableColumn) == "concept_cd" {
					so.NodeEncryptID = IDConcepts

					if _, ok := TablesShrineOntology[rawName]; ok {
						TablesShrineOntology[rawName].Sensitive[so.Fullname] = so
					} else {
						sensitive := make(map[string]*ShrineOntology)
						clear := make(map[string]*ShrineOntology)
						sensitive[so.Fullname] = so
						TablesShrineOntology[rawName] = ShrineTableInfo{Clear: clear, Sensitive: sensitive}
					}

					// if the ID does not yet exist
					if _, ok := MapConceptPathToTag[lo.Fullname]; !ok {
						MapConceptPathToTag[lo.Fullname] = TagAndID{Tag: libunlynx.GroupingKey(-1), TagID: -1}
						listConceptCD = append(listConceptCD, lo.Fullname)
						allSensitiveConceptIDs = append(allSensitiveConceptIDs, IDConcepts)
					}

					IDConcepts++
				}
			} else {

				//TODO for now we remove all modifiers
				if strings.ToLower(so.FactTableColumn) == "concept_cd" {
					if lo.HLevel != "0" {
						// add a new entry for the AdapterMappings (1-1 mapping between shrine and local concept)
						addNewAdapterMapping(lo.Fullname, false)
					}
					// add a new entry to the local ontology table
					TableLocalOntologyClear[lo.Fullname] = lo
					// add a new entry to the shrine ontology table
					if _, ok := TablesShrineOntology[rawName]; ok {
						TablesShrineOntology[rawName].Clear[so.Fullname] = so
					} else {
						sensitive := make(map[string]*ShrineOntology)
						clear := make(map[string]*ShrineOntology)
						clear[so.Fullname] = so
						TablesShrineOntology[rawName] = ShrineTableInfo{Clear: clear, Sensitive: sensitive}
					}

					TablesShrineOntology[rawName].Clear[so.Fullname] = so
				}
			}
		}
	}

	// if there are sensitive concepts
	if len(allSensitiveConceptIDs) > 0 {
		taggedConceptValues, err := EncryptAndTag(allSensitiveConceptIDs, group, entryPointIdx)
		if err != nil {
			return err
		}

		// re-randomize TAG_IDs
		rand.Seed(time.Now().UnixNano())
		perm := rand.Perm(len(MapConceptPathToTag))

		// 'populate' map (Concept codes)
		for i, concept := range listConceptCD {
			var tmp = MapConceptPathToTag[concept]
			tmp.TagID = int64(perm[i])
			tmp.Tag = taggedConceptValues[i]
			MapConceptPathToTag[concept] = tmp
		}
	}

	return nil
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

// ConvertClearLocalTable converts the old xxxx.csv local ontology file
func ConvertClearLocalTable(rawName string) error {
	// two new files are generated: one to store the non-sensitive data and another to store the sensitive data
	csvClearOutputFile, err := os.Create(OutputFilePaths["LOCAL_"+rawName].Path)
	if err != nil {
		log.Fatal("Error opening [" + strings.ToLower(rawName) + "].csv")
		return err
	}
	defer csvClearOutputFile.Close()

	headerString := ""
	for _, header := range HeaderLocalOntology {
		headerString += "\"" + header + "\","
	}
	// remove the last ,
	csvClearOutputFile.WriteString(headerString[:len(headerString)-1] + "\n")

	// non-sensitive
	for _, lo := range TableLocalOntologyClear {
		csvClearOutputFile.WriteString(lo.ToCSVText() + "\n")
	}

	return nil
}

// ConvertSensitiveLocalTable generates the sensitive_tagged file
func ConvertSensitiveLocalTable() error {
	csvSensitiveOutputFile, err := os.Create(OutputFilePaths["SENSITIVE_TAGGED"].Path)
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
	csvSensitiveOutputFile.WriteString(headerString[:len(headerString)-1] + "\n")

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
	csvOutputFile, err := os.Create(OutputFilePaths["PATIENT_DIMENSION"].Path)
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
	csvOutputNewPatientNumFile, err := os.Create(OutputFilePaths["NEW_PATIENT_NUM"].Path)
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

	csvOutputFile, err := os.Create(OutputFilePaths["VISIT_DIMENSION"].Path)
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
	csvOutputNewEncounterNumFile, err := os.Create(OutputFilePaths["NEW_ENCOUNTER_NUM"].Path)
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
	csvOutputFile, err := os.Create(OutputFilePaths["CONCEPT_DIMENSION"].Path)
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

	csvOutputFile, err := os.Create(OutputFilePaths["OBSERVATION_FACT"].Path)
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
			copyObs.PK.ConceptCD = "TAG_ID:" + strconv.FormatInt(MapConceptCodeToTag[copyObs.PK.ConceptCD], 10)
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
