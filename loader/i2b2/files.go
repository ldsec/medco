package loaderi2b2

import (
	"encoding/csv"
	"go.dedis.ch/onet/v3/log"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	dpath := os.Getenv("DEFAULT_DATA_PATH")
	if dpath == "" {
		defaultDataPath = "../../test/data/"
		log.Warn("Couldn't parse DEFAULT_DATA_PATH, using default value: ", defaultDataPath)
	} else {
		defaultDataPath = dpath
	}

	addDataPathToFiles()
	createDirectories()
}

// Files is the object structure behind the files.toml
type Files struct {
	TableAccess       string
	Ontology          []string
	DummyToPatient    string
	PatientDimension  string
	VisitDimension    string
	ConceptDimension  string
	ModifierDimension string
	ObservationFact   string
	OutputFolder      string
}

// fileInfo contains the path of the .csv output file and the name of the table where it should be loaded
type fileInfo struct {
	TableName string
	Path      string
}

const (
	// pathNonSensitive path of non sensitive data files
	pathNonSensitive = "i2b2/converted/non_sensitive/"

	// pathSensitive path of sensitive data files
	pathSensitive = "i2b2/converted/sensitive/"

	// i2b2MetadataNonSensitive path to i2b2metadata_i2b2_non_sensitive schema
	i2b2MetadataNonSensitive = "i2b2metadata_i2b2_non_sensitive."

	// i2b2DemoDataNonSensitive path to i2b2demodata_i2b2_non_sensitive schema
	i2b2DemoDataNonSensitive = "i2b2demodata_i2b2_non_sensitive."

	// i2b2DemoDataSensitive path to i2b2demodata_i2b2_sensitive schema
	i2b2DemoDataSensitive = "i2b2demodata_i2b2_sensitive."

	// ontNonSensitive path to medco_ont_non_sensitive schema
	ontNonSensitive = "medco_ont_non_sensitive."

	// ontSensitive path to medco_ont_sensitive schema
	ontSensitive = "medco_ont_sensitive."

	// Header is a generic XML header suitable for use with the output of Marshal.
	// This is not automatically added to any output of this package,
	// it is provided as a convenience.
	Header = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n"
)

// The paths for all input and output files
var (
	// defaultDataPath is the default path for the data folder
	defaultDataPath string

	ontologyFilesPaths = []string{
		"ONTOLOGY_BIRN",
		"ONTOLOGY_CUSTOM_META",
		"ONTOLOGY_ICD10_ICD9",
		"ONTOLOGY_I2B2",
	}

	inputFilePaths = map[string]string{
		"ONTOLOGY_BIRN":        "i2b2/original/birn.csv",
		"ONTOLOGY_CUSTOM_META": "i2b2/original/custom_meta.csv",
		"ONTOLOGY_ICD10_ICD9":  "i2b2/original/icd10_icd9.csv",
		"ONTOLOGY_I2B2":        "i2b2/original/i2b2.csv",

		"TABLE_ACCESS":         "i2b2/original/table_access.csv",
		"PATIENT_DIMENSION":    "i2b2/original/patient_dimension.csv",
		"VISIT_DIMENSION":      "i2b2/original/visit_dimension.csv",
		"CONCEPT_DIMENSION":    "i2b2/original/concept_dimension.csv",
		"MODIFIER_DIMENSION":   "i2b2/original/modifier_dimension.csv",
		"OBSERVATION_FACT_OLD": "i2b2/original/observation_fact_old.csv",
	}

	outputFilesPathsNonSensitive = map[string]fileInfo{
		"TABLE_ACCESS": {TableName: ontNonSensitive + "table_access", Path: pathNonSensitive + "table_access.csv"},

		"LOCAL_BIRN":        {TableName: i2b2MetadataNonSensitive + "birn", Path: pathNonSensitive + "local_birn.csv"},
		"LOCAL_CUSTOM_META": {TableName: i2b2MetadataNonSensitive + "custom_meta", Path: pathNonSensitive + "local_custom_meta.csv"},
		"LOCAL_ICD10_ICD9":  {TableName: i2b2MetadataNonSensitive + "icd10_icd9", Path: pathNonSensitive + "local_icd10_icd9.csv"},
		"LOCAL_I2B2":        {TableName: i2b2MetadataNonSensitive + "i2b2", Path: pathNonSensitive + "local_i2b2.csv"},

		"MEDCO_BIRN":        {TableName: ontNonSensitive + "birn", Path: pathNonSensitive + "medco_birn.csv"},
		"MEDCO_CUSTOM_META": {TableName: ontNonSensitive + "custom_meta", Path: pathNonSensitive + "medco_custom_meta.csv"},
		"MEDCO_ICD10_ICD9":  {TableName: ontNonSensitive + "icd10_icd9", Path: pathNonSensitive + "medco_icd10_icd9.csv"},
		"MEDCO_I2B2":        {TableName: ontNonSensitive + "i2b2", Path: pathNonSensitive + "medco_i2b2.csv"},

		"PATIENT_DIMENSION":         {TableName: i2b2DemoDataNonSensitive + "patient_dimension", Path: pathNonSensitive + "patient_dimension.csv"},
		"VISIT_DIMENSION":           {TableName: i2b2DemoDataNonSensitive + "visit_dimension", Path: pathNonSensitive + "visit_dimension.csv"},
		"CONCEPT_DIMENSION":         {TableName: i2b2DemoDataNonSensitive + "concept_dimension", Path: pathNonSensitive + "concept_dimension.csv"},
		"MODIFIER_DIMENSION":        {TableName: i2b2DemoDataNonSensitive + "modifier_dimension", Path: pathNonSensitive + "modifier_dimension.csv"},
		"OBSERVATION_FACT_FILTERED": {TableName: "", Path: pathNonSensitive + "observation_fact_filtered.csv"},
		"OBSERVATION_FACT":          {TableName: i2b2DemoDataNonSensitive + "observation_fact", Path: pathNonSensitive + "observation_fact.csv"},
	}

	outputFilesPathsSensitive = map[string]fileInfo{
		"TABLE_ACCESS": {TableName: ontSensitive + "table_access", Path: pathSensitive + "table_access.csv"},

		"SENSITIVE_TAGGED": {TableName: ontSensitive + "sensitive_tagged", Path: pathSensitive + "sensitive_tagged.csv"},

		"MEDCO_BIRN":        {TableName: ontSensitive + "birn", Path: pathSensitive + "medco_birn.csv"},
		"MEDCO_CUSTOM_META": {TableName: ontSensitive + "custom_meta", Path: pathSensitive + "medco_custom_meta.csv"},
		"MEDCO_ICD10_ICD9":  {TableName: ontSensitive + "icd10_icd9", Path: pathSensitive + "medco_icd10_icd9.csv"},
		"MEDCO_I2B2":        {TableName: ontSensitive + "i2b2", Path: pathSensitive + "medco_i2b2.csv"},

		"PATIENT_DIMENSION_FILTERED":    {TableName: "", Path: pathSensitive + "patient_dimension_filtered.csv"},
		"PATIENT_DIMENSION":             {TableName: i2b2DemoDataSensitive + "patient_dimension", Path: pathSensitive + "patient_dimension.csv"},
		"NEW_PATIENT_NUM":               {TableName: "", Path: pathSensitive + "new_patient_num.csv"},
		"CONCEPT_DIMENSION":             {TableName: i2b2DemoDataSensitive + "concept_dimension", Path: pathSensitive + "concept_dimension.csv"},
		"OBSERVATION_FACT_FILTERED":     {TableName: "", Path: pathSensitive + "observation_fact_filtered.csv"},
		"OBSERVATION_FACT_WITH_DUMMIES": {TableName: "", Path: pathSensitive + "observation_fact_with_dummies.csv"},
		"DUMMY_TO_PATIENT":              {TableName: "", Path: pathSensitive + "dummy_to_patient.csv"},
		"OBSERVATION_FACT":              {TableName: i2b2DemoDataSensitive + "observation_fact", Path: pathSensitive + "observation_fact.csv"},

		// we do not used this file, it is needed by the script generating the dummies
		"PATIENT_DIMENSION_SCRIPT": {TableName: "", Path: pathSensitive + "patient_dimension_script.csv"},
	}

	fileBashPathNonSensitive = "24-load-non-sensitive-i2b2-data.sh"
	fileBashPathSensitive    = "24-load-sensitive-i2b2-data.sh"

	filePythonGenerateDummies = "../import-tool/using_clustering.py"
)

func addDataPathToFiles() {
	for k, v := range inputFilePaths {
		inputFilePaths[k] = defaultDataPath + v
	}
	for k, v := range outputFilesPathsNonSensitive {
		tmp := outputFilesPathsNonSensitive[k]
		tmp.Path = defaultDataPath + v.Path
		outputFilesPathsNonSensitive[k] = tmp
	}
	for k, v := range outputFilesPathsSensitive {
		tmp := outputFilesPathsSensitive[k]
		tmp.Path = defaultDataPath + v.Path
		outputFilesPathsSensitive[k] = tmp
	}
}

func createDirectories() {

	os.MkdirAll(filepath.Join(defaultDataPath, pathNonSensitive), os.ModePerm)
	os.MkdirAll(filepath.Join(defaultDataPath, pathSensitive), os.ModePerm)

}

func changeInputFiles(files Files, directory string) {
	if len(files.Ontology) == 0 {
		log.Fatal("No Ontology files were selected for conversion")
	}

	for _, name := range files.Ontology {
		tokens := strings.Split(name, "/")
		ontologyName := "ONTOLOGY_" + strings.ToUpper(strings.Split(tokens[len(tokens)-1], ".")[0])
		inputFilePaths[ontologyName] = directory + "/" + name
		ontologyFilesPaths = append(ontologyFilesPaths, ontologyName)
	}
	inputFilePaths["TABLE_ACCESS"] = directory + "/" + files.TableAccess
	inputFilePaths["PATIENT_DIMENSION"] = directory + "/" + files.PatientDimension
	inputFilePaths["VISIT_DIMENSION"] = directory + "/" + files.VisitDimension
	inputFilePaths["CONCEPT_DIMENSION"] = directory + "/" + files.ConceptDimension
	inputFilePaths["MODIFIER_DIMENSION"] = directory + "/" + files.ModifierDimension
	inputFilePaths["OBSERVATION_FACT"] = directory + "/" + files.ObservationFact
}

func changeOutputFilesNonSensitive(folderPath string) {
	// fixed demodata tables
	outputFilesPathsNonSensitive["PATIENT_DIMENSION"] = fileInfo{TableName: i2b2DemoDataNonSensitive + "patient_dimension", Path: folderPath + "non_sensitive/patient_dimension.csv"}
	outputFilesPathsNonSensitive["VISIT_DIMENSION"] = fileInfo{TableName: i2b2DemoDataNonSensitive + "visit_dimension", Path: folderPath + "non_sensitive/visit_dimension.csv"}
	outputFilesPathsNonSensitive["CONCEPT_DIMENSION"] = fileInfo{TableName: i2b2DemoDataNonSensitive + "concept_dimension", Path: folderPath + "non_sensitive/concept_dimension.csv"}
	outputFilesPathsNonSensitive["MODIFIER_DIMENSION"] = fileInfo{TableName: i2b2DemoDataNonSensitive + "modifier_dimension", Path: folderPath + "non_sensitive/modifier_dimension.csv"}
	outputFilesPathsNonSensitive["OBSERVATION_FACT_FILTERED"] = fileInfo{TableName: "", Path: folderPath + "non_sensitive/observation_fact_filtered.csv"}
	outputFilesPathsNonSensitive["OBSERVATION_FACT"] = fileInfo{TableName: i2b2DemoDataNonSensitive + "observation_fact", Path: folderPath + "non_sensitive/observation_fact.csv"}

	// fixed ontology tables
	outputFilesPathsNonSensitive["TABLE_ACCESS"] = fileInfo{TableName: ontNonSensitive + "table_access", Path: folderPath + "non_sensitive/table_access.csv"}

	for key, path := range inputFilePaths {
		if strings.HasPrefix(key, "ONTOLOGY_") {
			rawKey := strings.Split(key, "ONTOLOGY_")[1]
			tokens := strings.Split(path, "/")

			outputFilesPathsNonSensitive["LOCAL_"+rawKey] = fileInfo{TableName: i2b2MetadataNonSensitive + strings.ToLower(rawKey), Path: folderPath + "non_sensitive/local_" + tokens[len(tokens)-1]}
			outputFilesPathsNonSensitive["MEDCO_"+rawKey] = fileInfo{TableName: ontNonSensitive + strings.ToLower(rawKey), Path: folderPath + "non_sensitive/medco_" + tokens[len(tokens)-1]}

		}
	}
}

func changeOutputFilesSensitive(folderPath string) {

	// fixed demodata tables$
	outputFilesPathsSensitive["PATIENT_DIMENSION_FILTERED"] = fileInfo{TableName: "", Path: folderPath + "sensitive/patient_dimension_filtered.csv"}
	outputFilesPathsSensitive["PATIENT_DIMENSION"] = fileInfo{TableName: i2b2DemoDataSensitive + "patient_dimension", Path: folderPath + "sensitive/patient_dimension.csv"}
	outputFilesPathsSensitive["NEW_PATIENT_NUM"] = fileInfo{TableName: "", Path: folderPath + "sensitive/new_patient_num.csv"}
	outputFilesPathsSensitive["CONCEPT_DIMENSION"] = fileInfo{TableName: i2b2DemoDataSensitive + "concept_dimension", Path: folderPath + "sensitive/concept_dimension.csv"}
	outputFilesPathsSensitive["OBSERVATION_FACT_FILTERED"] = fileInfo{TableName: "", Path: folderPath + "sensitive/observation_fact_filtered.csv"}
	outputFilesPathsSensitive["OBSERVATION_FACT_WITH_DUMMIES"] = fileInfo{TableName: "", Path: folderPath + "sensitive/observation_fact_with_dummies.csv"}
	outputFilesPathsSensitive["OBSERVATION_FACT"] = fileInfo{TableName: i2b2DemoDataSensitive + "observation_fact", Path: folderPath + "sensitive/observation_fact.csv"}

	// fixed ontology tables
	outputFilesPathsSensitive["TABLE_ACCESS"] = fileInfo{TableName: ontSensitive + "table_access", Path: folderPath + "sensitive/table_access.csv"}
	outputFilesPathsSensitive["SENSITIVE_TAGGED"] = fileInfo{TableName: ontSensitive + "sensitive_tagged", Path: folderPath + "sensitive/sensitive_tagged.csv"}

	for key, path := range inputFilePaths {
		if strings.HasPrefix(key, "ONTOLOGY_") {
			rawKey := strings.Split(key, "ONTOLOGY_")[1]
			tokens := strings.Split(path, "/")

			outputFilesPathsSensitive["MEDCO_"+rawKey] = fileInfo{TableName: ontSensitive + strings.ToLower(rawKey), Path: folderPath + "sensitive/medco_" + tokens[len(tokens)-1]}
		}
	}
}

func readCSV(filename string) ([][]string, error) {
	csvInputFile, err := os.Open(filename)
	if err != nil {
		log.Fatal("Error opening " + filename)
		return nil, err
	}
	defer csvInputFile.Close()

	reader := csv.NewReader(csvInputFile)
	reader.Comma = ','

	lines, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Error reading "+filename, err)
		return nil, err
	}

	return lines, nil
}
