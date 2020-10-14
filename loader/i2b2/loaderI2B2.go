package loaderi2b2

import (
	"github.com/ldsec/medco/loader"
	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/log"
	"os"
	"strconv"
	"strings"
)

// LoadI2B2Data is the main function that performs a full conversion and loading of the i2b2 data
func LoadI2B2Data(el *onet.Roster, entryPointIdx int, directory string, files Files, allSensitive bool, mapSensitive map[string]struct{}, i2b2DB loader.DBSettings, empty bool) error {
	inputFilePaths = make(map[string]string)
	outputFilesPathsNonSensitive = make(map[string]fileInfo)
	outputFilesPathsSensitive = make(map[string]fileInfo)
	ontologyFilesPaths = make([]string, 0)

	if allSensitive {
		allSensitive = true
	} else {
		listSensitiveConcepts = mapSensitive
	}

	// change input filepaths
	changeInputFiles(files, directory)

	// change output filepaths
	changeOutputFilesNonSensitive(directory + "/" + files.OutputFolder)
	changeOutputFilesSensitive(directory + "/" + files.OutputFolder)

	log.Lvl2("--- Started v1 Data Conversion ---")

	err := parseTableAccess()
	if err != nil {
		return err
	}

	err = convertTableAccess()
	if err != nil {
		return err
	}

	log.Lvl2("--- Finished converting TABLE_ACCESS ---")

	err = convertLocalOntology(el, entryPointIdx)
	if err != nil {
		return err
	}

	log.Lvl2("--- Finished converting LOCAL_ONTOLOGY ---")

	err = generateMedCoOntology()
	if err != nil {
		return err
	}

	log.Lvl2("--- Finished generating MEDCO_ONTOLOGY ---")

	err = parseModifierDimension()
	if err != nil {
		return err
	}
	err = convertModifierDimension()
	if err != nil {
		return err
	}

	log.Lvl2("--- Finished converting MODIFIER_DIMENSION ---")

	err = parseConceptDimension()
	if err != nil {
		return err
	}
	err = convertConceptDimension()
	if err != nil {
		return err
	}

	log.Lvl2("--- Finished converting CONCEPT_DIMENSION ---")

	err = filterOldObservationFact()
	if err != nil {
		return err
	}

	log.Lvl2("--- Finished filtering OLD_OBSERVATION_FACT ---")

	err = filterPatientDimension(el.Aggregate)
	if err != nil {
		return err
	}

	log.Lvl2("--- Finished filtering PATIENT_DIMENSION ---")

	err = callGenerateDummiesScript()

	if err != nil {
		return err
	}

	log.Lvl2("--- Finished dummies generation ---")

	err = parseDummyToPatient()
	if err != nil {
		return err
	}

	err = parsePatientDimension(el.Aggregate)
	if err != nil {
		return err
	}
	err = convertPatientDimension(el.Aggregate)
	if err != nil {
		return err
	}

	log.Lvl2("--- Finished converting PATIENT_DIMENSION ---")

	err = parseVisitDimension()
	if err != nil {
		return err
	}
	err = convertVisitDimension()
	if err != nil {
		return err
	}

	log.Lvl2("--- Finished converting VISIT_DIMENSION ---")

	err = parseNonSensitiveObservationFact()
	if err != nil {
		return err
	}
	err = convertNonSensitiveObservationFact()
	if err != nil {
		return err
	}

	log.Lvl2("--- Finished converting non sensitive OBSERVATION_FACT ---")

	err = parseSensitiveObservationFact()
	if err != nil {
		return err
	}
	err = convertSensitiveObservationFact()
	if err != nil {
		return err
	}

	log.Lvl2("--- Finished converting OBSERVATION_FACT ---")

	err = generateLoadingDataScriptNonSensitive(i2b2DB)
	if err != nil {
		log.Fatal("Error while generating the loading non sensitive data .sh file", err)
		return err
	}

	err = generateLoadingDataScriptSensitive(i2b2DB)
	if err != nil {
		log.Fatal("Error while generating the loading sensitive data .sh file", err)
		return err
	}

	log.Lvl2("--- Finished generating loading scripts ---")

	err = loadDataFilesNonSensitive()
	if err != nil {
		log.Fatal("Error while loading non sensitive data", err)
		return err
	}

	err = loadDataFilesSensitive()
	if err != nil {
		log.Fatal("Error while loading sensitive data", err)
		return err
	}

	log.Lvl2("--- Finished loading data ---")

	return nil
}

// generateLoadingDataScriptNonSensitive creates a load dataset .sql script for non sensitive data (deletes the data in the corresponding tables and reloads the new 'protected' data)
func generateLoadingDataScriptNonSensitive(i2b2DB loader.DBSettings) error {
	fp, err := os.Create(fileBashPathNonSensitive)
	if err != nil {
		return err
	}

	loading := `#!/usr/bin/env bash` + "\n" + "\n" + `PGPASSWORD=` + i2b2DB.DBpassword + ` psql -v ON_ERROR_STOP=1 -h "` + i2b2DB.DBhost +
		`" -U "` + i2b2DB.DBuser + `" -p ` + strconv.FormatInt(int64(i2b2DB.DBport), 10) + ` -d "` + i2b2DB.DBname + `" <<-EOSQL` + "\n"

	loading += "BEGIN;\n"

	loading += "TRUNCATE TABLE " + i2b2DemoDataNonSensitive + "patient_mapping;\n" +
		"TRUNCATE TABLE " + i2b2DemoDataNonSensitive + "encounter_mapping;\n" +
		"TRUNCATE TABLE " + i2b2DemoDataNonSensitive + "concept_dimension;\n" +
		"TRUNCATE TABLE " + i2b2DemoDataNonSensitive + "modifier_dimension;\n" +
		"TRUNCATE TABLE " + i2b2DemoDataNonSensitive + "patient_dimension;\n" +
		"TRUNCATE TABLE " + i2b2DemoDataNonSensitive + "visit_dimension;\n" +
		"TRUNCATE TABLE " + i2b2DemoDataNonSensitive + "observation_fact;\n"

	loading += `\copy ` + outputFilesPathsNonSensitive["CONCEPT_DIMENSION"].TableName + ` FROM '` + outputFilesPathsNonSensitive["CONCEPT_DIMENSION"].Path + `' ESCAPE '"' DELIMITER ',' CSV HEADER;` + "\n" +
		`\copy ` + outputFilesPathsNonSensitive["MODIFIER_DIMENSION"].TableName + ` FROM '` + outputFilesPathsNonSensitive["MODIFIER_DIMENSION"].Path + `' ESCAPE '"' DELIMITER ',' CSV HEADER;` + "\n" +
		`\copy ` + outputFilesPathsNonSensitive["PATIENT_DIMENSION"].TableName + ` FROM '` + outputFilesPathsNonSensitive["PATIENT_DIMENSION"].Path + `' ESCAPE '"' DELIMITER ',' CSV HEADER;` + "\n" +
		`\copy ` + outputFilesPathsNonSensitive["VISIT_DIMENSION"].TableName + ` FROM '` + outputFilesPathsNonSensitive["VISIT_DIMENSION"].Path + `' ESCAPE '"' DELIMITER ',' CSV HEADER;` + "\n" +
		`\copy ` + outputFilesPathsNonSensitive["OBSERVATION_FACT"].TableName + ` FROM '` + outputFilesPathsNonSensitive["OBSERVATION_FACT"].Path + `' ESCAPE '"' DELIMITER ',' CSV HEADER;` + "\n"

	loading += "\n"

	for file, fI := range outputFilesPathsNonSensitive {
		if strings.HasPrefix(file, "LOCAL_") {
			loading += "TRUNCATE TABLE " + fI.TableName + ";\n"
			loading += `\copy ` + fI.TableName + ` FROM '` + fI.Path + `' ESCAPE '"' DELIMITER ',' CSV HEADER;` + "\n"
		}
	}

	loading += `\copy ` + outputFilesPathsNonSensitive["TABLE_ACCESS"].TableName + ` FROM '` + outputFilesPathsNonSensitive["TABLE_ACCESS"].Path + `' ESCAPE '"' DELIMITER ',' CSV HEADER;` + "\n"
	loading += "\n"

	// Create MedCo Table
	loading += createMedcoTable(outputFilesPathsNonSensitive)

	loading += "COMMIT;\n"
	loading += "EOSQL"

	_, err = fp.WriteString(loading)
	if err != nil {
		return err
	}

	fp.Close()

	return nil
}

// generateLoadingDataScriptSensitive creates a load dataset .sql script for sensitive data (deletes the data in the corresponding tables and reloads the new 'protected' data)
func generateLoadingDataScriptSensitive(i2b2DB loader.DBSettings) error {

	fp, err := os.Create(fileBashPathSensitive)
	if err != nil {
		return err
	}

	loading := `#!/usr/bin/env bash` + "\n" + "\n" + `PGPASSWORD=` + i2b2DB.DBpassword + ` psql -v ON_ERROR_STOP=1 -h "` + i2b2DB.DBhost +
		`" -U "` + i2b2DB.DBuser + `" -p ` + strconv.FormatInt(int64(i2b2DB.DBport), 10) + ` -d "` + i2b2DB.DBname + `" <<-EOSQL` + "\n"

	loading += "BEGIN;\n"

	loading += "TRUNCATE TABLE " + i2b2DemoDataSensitive + "patient_mapping;\n" +
		"TRUNCATE TABLE " + i2b2DemoDataSensitive + "encounter_mapping;\n" +
		"TRUNCATE TABLE " + i2b2DemoDataSensitive + "concept_dimension;\n" +
		"TRUNCATE TABLE " + i2b2DemoDataSensitive + "patient_dimension;\n" +
		"TRUNCATE TABLE " + i2b2DemoDataSensitive + "observation_fact;\n"

	loading += `\copy ` + outputFilesPathsSensitive["CONCEPT_DIMENSION"].TableName + ` FROM '` + outputFilesPathsSensitive["CONCEPT_DIMENSION"].Path + `' ESCAPE '"' DELIMITER ',' CSV HEADER;` + "\n" +
		`\copy ` + outputFilesPathsSensitive["PATIENT_DIMENSION"].TableName + ` FROM '` + outputFilesPathsSensitive["PATIENT_DIMENSION"].Path + `' ESCAPE '"' DELIMITER ',' CSV HEADER;` + "\n" +
		`\copy ` + outputFilesPathsSensitive["OBSERVATION_FACT"].TableName + ` FROM '` + outputFilesPathsSensitive["OBSERVATION_FACT"].Path + `' ESCAPE '"' DELIMITER ',' CSV HEADER;` + "\n"

	loading += "\n"

	loading += `\copy ` + outputFilesPathsSensitive["TABLE_ACCESS"].TableName + ` FROM '` + outputFilesPathsSensitive["TABLE_ACCESS"].Path + `' ESCAPE '"' DELIMITER ',' CSV HEADER;` + "\n"
	loading += "TRUNCATE TABLE " + outputFilesPathsSensitive["SENSITIVE_TAGGED"].TableName + ";\n"
	loading += `\copy ` + outputFilesPathsSensitive["SENSITIVE_TAGGED"].TableName + ` FROM '` + outputFilesPathsSensitive["SENSITIVE_TAGGED"].Path + `' ESCAPE '"' DELIMITER ',' CSV HEADER;` + "\n"
	loading += "\n"

	// Create MedCo Table
	loading += createMedcoTable(outputFilesPathsSensitive)

	loading += "COMMIT;\n"
	loading += "EOSQL"

	_, err = fp.WriteString(loading)
	if err != nil {
		return err
	}

	fp.Close()

	return nil
}

func createMedcoTable(outputFilesPaths map[string]fileInfo) (loading string) {
	for file, fI := range outputFilesPaths {
		if strings.HasPrefix(file, "MEDCO_") {
			loading = `CREATE TABLE IF NOT EXISTS ` + fI.TableName + ` (
        				C_HLEVEL NUMERIC(22,0),
        				C_FULLNAME VARCHAR(900),
        				C_NAME VARCHAR(2000),
        				C_SYNONYM_CD CHAR(1),
        				C_VISUALATTRIBUTES CHAR(3),
        				C_TOTALNUM NUMERIC(22,0),
        				C_BASECODE VARCHAR(450),
        				C_METADATAXML TEXT,
        				C_FACTTABLECOLUMN VARCHAR(50),
        				C_TABLENAME VARCHAR(50),
						C_COLUMNNAME VARCHAR(50),
        				C_COLUMNDATATYPE VARCHAR(50),
        				C_OPERATOR VARCHAR(10),
        				C_DIMCODE VARCHAR(900),
        				C_COMMENT TEXT,
        				C_TOOLTIP VARCHAR(900),
        				UPDATE_DATE DATE,
						DOWNLOAD_DATE DATE,
        				IMPORT_DATE DATE,
        				SOURCESYSTEM_CD VARCHAR(50),
        				VALUETYPE_CD VARCHAR(50),
        				M_APPLIED_PATH VARCHAR(900),
        				M_EXCLUSION_CD VARCHAR(900));
        				
						ALTER TABLE ` + fI.TableName + ` OWNER TO i2b2;` + "\n"

			loading += "TRUNCATE TABLE " + fI.TableName + ";\n"
			loading += `\copy ` + fI.TableName + ` FROM '` + fI.Path + `' ESCAPE '"' DELIMITER ',' CSV HEADER;` + "\n"
		}
	}

	return

}

// loadDataFilesNonSensitive executes the loading script for the new converted non sensitive data
func loadDataFilesNonSensitive() error {
	return loader.ExecuteScript("/bin/sh", fileBashPathNonSensitive)
}

// loadDataFilesNonSensitive executes the loading script for the new converted sensitive data
func loadDataFilesSensitive() error {
	return loader.ExecuteScript("/bin/sh", fileBashPathSensitive)
}
