package loadergenomic

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/csv"
	"errors"
	"github.com/dedis/onet"
	"github.com/dedis/onet/log"
	"github.com/lca1/medco-loader/loader"
	"github.com/lca1/medco-unlynx/services"
	"github.com/lca1/unlynx/lib"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// I2B2METADATA path to i2b2metadata schema
const I2B2METADATA = "i2b2metadata_i2b2."

// I2B2DEMODATA path to i2b2demodata schema
const I2B2DEMODATA = "i2b2demodata_i2b2."

// ONT path to medco_ont schema
const ONT = "medco_ont."

// ANNOTATIONS path to genomic_annotations schema
const ANNOTATIONS = "genomic_annotations."

// The different paths and handlers for all the .sql files
var (
	OutputFilePath = "../data/genomic/"

	TablenamesOntology = [...]string{ONT + "clinical_sensitive",
		ONT + "clinical_non_sensitive",
		ANNOTATIONS + "genomic_annotations",
		ONT + "sensitive_tagged"}

	TablenamesData = [...]string{I2B2DEMODATA + "concept_dimension",
		I2B2DEMODATA + "patient_mapping",
		I2B2DEMODATA + "patient_dimension",
		I2B2DEMODATA + "encounter_mapping",
		I2B2DEMODATA + "visit_dimension",
		I2B2DEMODATA + "provider_dimension",
		I2B2DEMODATA + "observation_fact"}

	FileBashPath = [...]string{"25-load-ontology.sh",
		"26-load-data.sh"}

	FilePathsOntology = [...]string{"MEDCO_ONT_CLINICAL_SENSITIVE.csv",
		"MEDCO_ONT_CLINICAL_NON_SENSITIVE.csv",
		"MEDCO_ONT_GENOMIC_ANNOTATIONS.csv",
		"MEDCO_ONT_SENSITIVE_TAGGED.csv"}

	FilePathsData = [...]string{"I2B2DEMODATA_CONCEPT_DIMENSION.csv",
		"I2B2DEMODATA_PATIENT_MAPPING.csv",
		"I2B2DEMODATA_PATIENT_DIMENSION.csv",
		"I2B2DEMODATA_ENCOUNTER_MAPPING.csv",
		"I2B2DEMODATA_VISIT_DIMENSION.csv",
		"I2B2DEMODATA_PROVIDER_DIMENSION.csv",
		"I2B2DEMODATA_OBSERVATION_FACT.csv"}
)

/*
ToIgnore: 			defines the columns to be ignored (mostly the sample and patient IDs)
TranslationDic: 	defines the translation between the fields that are present in the different datafiles and their
					'actual meaning' code-wise
AnnotationsToQuery: defines the annotations to be queried (to speed up the query)
*/
var (
	ToIgnore = map[string]struct{}{
		"PATIENT_ID":  {},
		"P_STABLE_ID": {},
		"SAMPLE_ID":   {},
		"S_STABLE_ID": {},
	}

	TranslationDic = map[string]string{
		"Tumor_Sample_Barcode": "SAMPLE_ID",
		"Chromosome":           "CHR",
		"Start_Position":       "SP",
		"Reference_Allele":     "RA",
		"Tumor_Seq_Allele1":    "TSA1",
		"Tumor_Seq_Allele2":    "TSA2",
		"PATIENT_ID":           "PATIENT_ID",
		"SAMPLE_ID":            "SAMPLE_ID",
		"CHR":                  "CHR",
		"START_POSITION":       "SP",
		"REFERENCE_ALLELE":     "RA",
		"TUMOR_SEQ_ALLELE1":    "TSA1",
		"TUMOR_SEQ_ALLELE2":    "TSA2",
	}

	AnnotationsToQuery = map[string]struct{}{
		"HUGO_GENE_SYMBOL":  {},
		"Hugo_Symbol":       {},
		"PROTEIN_CHANGE":    {},
		"MA:protein.change": {},
	}

	AllSensitive = false
)

/* NumElMap: defines an approximate size of the map (it avoids rehashing and speeds up the execution)
   NumThreads: defines the amount of go subroutines to use when parelellizing the encryption
*/
var (
	NumElMap   = int64(5e6)
	NumThreads = int(20)
)

// SensitiveIDValue contains both concept path and annotation which will be linked to a certain sensitive ID
type SensitiveIDValue struct {
	CP         ConceptPath
	Annotation string
}

// ConceptPath defines the end of the concept path tree and we use it in a map so that we do not repeat concepts
type ConceptPath struct {
	Field  string
	Record string //leaf
}

// ConceptID defines its ID (e.g., E,1 - for ENC_ID,1; C,1 - for CLEAR_ID,1; sdasdcfsx,1432 - for tagged_value,TAG_ID
type ConceptID struct {
	Identifier string
	Value      int64
}

// Support global variables
var (
	Testing         bool // testing environment
	FileHandlers    []*os.File
	OntValues       map[ConceptPath]ConceptID // stores the concept path and the correspondent ID
	TextSearchIndex int64                     // needed for the observation_fact table (counter)
)

// ReplayDataset replays the dataset x number of times
func ReplayDataset(filename string, x int) error {
	log.LLvl1("Replaying dataset", x, "times...")

	// open file to read
	fGenomic, err := os.Open(filename)
	if err != nil {
		log.Fatal("Cannot open file to read:", err)
		return err
	}

	reader := csv.NewReader(fGenomic)
	reader.Comma = '\t'

	// read all genomic file
	record, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Error in the ReplayDataset() - reading:", err)
		return err
	}

	finalResult := record[:]

	header := true
	// replay x times
	for t := 0; t < x-1; t++ {
		for _, el := range record {
			// not a comment or blank line
			if string(el[0]) == "" || string(el[0][0:1]) == "#" {
				continue
			}

			// HEADER time...
			if header == true {
				header = false
				continue
			}

			finalResult = append(finalResult, el)
		}
	}

	fGenomic.Close()

	// open file to write
	fGenomic, err = os.Create(filename)
	if err != nil {
		log.Fatal("Cannot open file to write:", err)
		return err
	}

	writer := csv.NewWriter(fGenomic)
	writer.Comma = '\t'

	err = writer.WriteAll(finalResult)
	if err != nil {
		log.Fatal("Error in the ReplayDataset() - writing:", err)
		return err
	}

	fGenomic.Close()

	return nil

}

// LoadGenomicData initiates the loading process
func LoadGenomicData(el *onet.Roster, entryPointIdx int, fOntClinical, fOntGenomic, fClinical, fGenomic *os.File, outputPath string, allSensitive bool, mapSensitive map[string]struct{}, databaseS loader.DBSettings, testing bool) error {
	start := time.Now()

	// init global variables
	FileHandlers = make([]*os.File, 0)
	OntValues = make(map[ConceptPath]ConceptID)
	Testing = testing
	TextSearchIndex = int64(1) // needed for the observation_fact table (counter)
	OutputFilePath = outputPath
	AllSensitive = allSensitive

	for i := range FilePathsOntology {
		FilePathsOntology[i] = OutputFilePath + FilePathsOntology[i]
		fp, err := os.Create(FilePathsOntology[i])
		if err != nil {
			log.Fatal("Error while creating", FilePathsOntology[i])
			return err
		}
		FileHandlers = append(FileHandlers, fp)
	}

	for i := range FilePathsData {
		FilePathsData[i] = OutputFilePath + FilePathsData[i]
		fp, err := os.Create(FilePathsData[i])
		if err != nil {
			log.Fatal("Error while creating", FilePathsData[i])
			return err
		}
		FileHandlers = append(FileHandlers, fp)
	}

	err := GenerateOntologyFiles(el, entryPointIdx, fOntClinical, fOntGenomic, mapSensitive)
	if err != nil {
		log.Fatal("Error while generating the ontology .csv files", err)
		return err
	}

	// to free

	err = GenerateDataFiles(el, fClinical, fGenomic)
	if err != nil {
		log.Fatal("Error while generating the data .csv files", err)
		return err
	}

	fClinical.Close()
	fGenomic.Close()

	startLoadingOntology := time.Now()

	err = GenerateLoadingOntologyScript(databaseS)
	if err != nil {
		log.Fatal("Error while generating the loading ontology .sh file", err)
		return err
	}

	err = LoadOntologyFiles()
	if err != nil {
		log.Fatal("Error while loading ontology .sql file", err)
		return err
	}

	loadTime := time.Since(startLoadingOntology)
	log.LLvl1("Loading ontology took:", loadTime)

	startLoadingData := time.Now()

	err = GenerateLoadingDataScript(databaseS)
	if err != nil {
		log.Fatal("Error while generating the loading dataset .sh file", err)
		return err
	}

	err = LoadDataFiles()
	if err != nil {
		log.Fatal("Error while loading dataset .sql file", err)
		return err
	}

	loadTime = time.Since(startLoadingData)
	log.LLvl1("Loading dataset took:", loadTime)

	fOntClinical.Close()
	fOntGenomic.Close()

	for _, fp := range FileHandlers {
		fp.Close()
	}

	// to free memory
	OntValues = make(map[ConceptPath]ConceptID)
	FileHandlers = make([]*os.File, 0)

	etlTime := time.Since(start)
	log.LLvl1("The ETL took:", etlTime)

	return nil
}

// GenerateLoadingOntologyScript creates a load ontology .sql script
func GenerateLoadingOntologyScript(databaseS loader.DBSettings) error {
	fp, err := os.Create(FileBashPath[0])
	if err != nil {
		return err
	}

	loading := `#!/usr/bin/env bash` + "\n" + "\n" + `PGPASSWORD=` + databaseS.DBpassword + ` psql -v ON_ERROR_STOP=1 -h "` + databaseS.DBhost +
		`" -U "` + databaseS.DBuser + `" -p ` + strconv.FormatInt(int64(databaseS.DBport), 10) + ` -d "` + databaseS.DBname + `" <<-EOSQL` + "\n"

	loading += "BEGIN;\n"

	//update table access
	loading += `INSERT INTO medco_ont.table_access (c_table_cd, c_table_name, c_protected_access, c_hlevel, c_fullname, c_name,
				c_synonym_cd, c_visualattributes, c_facttablecolumn, c_dimtablename,
        		c_columnname, c_columndatatype, c_operator, c_dimcode, c_tooltip) VALUES
        		('CLINICAL_SENSITIVE', 'CLINICAL_SENSITIVE', 'N', 2, '\medco\clinical\sensitive\', 'MedCo Clinical Sensitive Ontology',
        		'N', 'CA', 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\medco\clinical\sensitive\', 'MedCo Clinical Sensitive Ontology');
    			INSERT INTO medco_ont.table_access (c_table_cd, c_table_name, c_protected_access, c_hlevel, c_fullname, c_name,
        		c_synonym_cd, c_visualattributes, c_facttablecolumn, c_dimtablename,
        		c_columnname, c_columndatatype, c_operator, c_dimcode, c_tooltip) VALUES
        		('CLINICAL_NON_SENSITIVE', 'CLINICAL_NON_SENSITIVE', 'N', 2, '\medco\clinical\nonsensitive\', 'MedCo Clinical Non-Sensitive Ontology',
        		'N', 'CA', 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\medco\clinical\nonsensitive\', 'MedCo Clinical Non-Sensitive Ontology');
    			INSERT INTO medco_ont.table_access (c_table_cd, c_table_name, c_protected_access, c_hlevel, c_fullname, c_name,
        		c_synonym_cd, c_visualattributes, c_facttablecolumn, c_dimtablename,
        		c_columnname, c_columndatatype, c_operator, c_dimcode, c_tooltip) VALUES
				('GENOMIC', 'GENOMIC', 'N', 1, '\medco\genomic\', 'MedCo Genomic Ontology',
        		'N', 'CA', 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\medco\genomic\', 'MedCo Genomic Ontology');` + "\n"

	loading += `CREATE TABLE medco_ont.clinical_sensitive(
        		c_hlevel numeric(22,0) not null,
        		c_fullname character varying(900) not null,
        		c_name character varying(2000) not null,
        		c_synonym_cd character(1) not null,
        		c_visualattributes character(3) not null,
        		c_totalnum numeric(22,0),
        		c_basecode character varying(450),
        		c_metadataxml text,
        		c_facttablecolumn character varying(50) not null,
        		c_tablename character varying(50) not null,
        		c_columnname character varying(50) not null,
        		c_columndatatype character varying(50) not null,
        		c_operator character varying(10) not null,
        		c_dimcode character varying(900) not null,
        		c_comment text,
        		c_tooltip character varying(900),
        		update_date date not null,
        		download_date date,
        		import_date date,
        		sourcesystem_cd character varying(50),
        		valuetype_cd character varying(50),
        		m_applied_path character varying(900) not null,
        		m_exclusion_cd character varying(900),
        		c_path character varying(700),
        		c_symbol character varying(50),
        		pcori_basecode character varying(50));
    			ALTER TABLE ONLY medco_ont.clinical_sensitive ADD CONSTRAINT fullname_pk_20 PRIMARY KEY (c_fullname);
    			ALTER TABLE ONLY medco_ont.clinical_sensitive ADD CONSTRAINT basecode_un_20 UNIQUE (c_basecode);
	
	   			CREATE TABLE medco_ont.clinical_non_sensitive(
        		c_hlevel numeric(22,0) not null,
        		c_fullname character varying(900) not null,
        		c_name character varying(2000) not null,
        		c_synonym_cd character(1) not null,
        		c_visualattributes character(3) not null,
				c_totalnum numeric(22,0),
        		c_basecode character varying(450),
        		c_metadataxml text,
        		c_facttablecolumn character varying(50) not null,
        		c_tablename character varying(50) not null,
        		c_columnname character varying(50) not null,
        		c_columndatatype character varying(50) not null,
        		c_operator character varying(10) not null,
        		c_dimcode character varying(900) not null,
        		c_comment text,
        		c_tooltip character varying(900),
        		update_date date not null,
				download_date date,
        		import_date date,
        		sourcesystem_cd character varying(50),
        		valuetype_cd character varying(50),
        		m_applied_path character varying(900) not null,
        		m_exclusion_cd character varying(900),
        		c_path character varying(700),
        		c_symbol character varying(50),
        		pcori_basecode character varying(50));
    			ALTER TABLE ONLY medco_ont.clinical_non_sensitive ADD CONSTRAINT fullname_pk_21 PRIMARY KEY (c_fullname);
    			ALTER TABLE ONLY medco_ont.clinical_non_sensitive ADD CONSTRAINT basecode_un_21 UNIQUE (c_basecode);
	
	
    			CREATE TABLE medco_ont.genomic(
        		c_hlevel numeric(22,0) not null,
        		c_fullname character varying(900) not null,
        		c_name character varying(2000) not null,
        		c_synonym_cd character(1) not null,
        		c_visualattributes character(3) not null,
        		c_totalnum numeric(22,0),
        		c_basecode character varying(450),
        		c_metadataxml text,
        		c_facttablecolumn character varying(50) not null,
        		c_tablename character varying(50) not null,
        		c_columnname character varying(50) not null,
        		c_columndatatype character varying(50) not null,
        		c_operator character varying(10) not null,
        		c_dimcode character varying(900) not null,
        		c_comment text,
        		c_tooltip character varying(900),
        		update_date date not null,
        		download_date date,
        		import_date date,
        		sourcesystem_cd character varying(50),
        		valuetype_cd character varying(50),
        		m_applied_path character varying(900) not null,
        		m_exclusion_cd character varying(900),
        		c_path character varying(700),
        		c_symbol character varying(50),
        		pcori_basecode character varying(50));
    			ALTER TABLE ONLY medco_ont.genomic ADD CONSTRAINT fullname_pk_22 PRIMARY KEY (c_fullname);
    			ALTER TABLE ONLY medco_ont.genomic ADD CONSTRAINT basecode_un_22 UNIQUE (c_basecode);` + "\n"

	loading += `INSERT INTO medco_ont.genomic (c_hlevel, c_fullname, c_name, c_synonym_cd, c_visualattributes, c_totalnum,
        		c_facttablecolumn, c_tablename, c_columnname, c_columndatatype, c_operator, c_dimcode, c_comment, c_tooltip, update_date,
        		download_date, import_date, valuetype_cd, m_applied_path) values
        		('1', '\medco\genomic\', 'MedCo Genomic Ontology', 'N', 'CA', '0', 'concept_cd', 'concept_dimension', 'concept_path',
        		'T', 'LIKE', '\medco\genomic\', 'MedCo Genomic Ontology', '\medco\genomic\',
        		'NOW()', 'NOW()', 'NOW()', 'GEN', '@');
            	INSERT INTO medco_ont.genomic (c_hlevel, c_fullname, c_name, c_synonym_cd, c_visualattributes, c_totalnum, c_basecode,
        		c_facttablecolumn, c_tablename, c_columnname, c_columndatatype, c_operator, c_dimcode, c_comment, c_tooltip, update_date,
        		download_date, import_date, valuetype_cd, m_applied_path) values
        		('2', '\medco\genomic\annotations_Hugo_Symbol\', 'Gene Name', 'N', 'LA', '0', 'GEN:Hugo_Gene_Symbol', 'concept_cd', 'concept_dimension', 'concept_path',
        		'T', 'LIKE', '\medco\genomic\annotations_Hugo_Symbol\', 'Gene Name', '\medco\genomic\annotations_Hugo_Symbol\',
        		'NOW()', 'NOW()', 'NOW()', 'GEN', '@');
    			INSERT INTO medco_ont.genomic (c_hlevel, c_fullname, c_name, c_synonym_cd, c_visualattributes, c_totalnum, c_basecode,
        		c_facttablecolumn, c_tablename, c_columnname, c_columndatatype, c_operator, c_dimcode, c_comment, c_tooltip, update_date,
        		download_date, import_date, valuetype_cd, m_applied_path) values
        		('2', '\medco\genomic\annotations_Protein_position\', 'Protein Position', 'N', 'LA', '0', 'GEN:Protein_Change', 'concept_cd', 'concept_dimension', 'concept_path',
        		'T', 'LIKE', '\medco\genomic\annotations_Protein_position\', 'Protein Position', '\medco\genomic\annotations_Protein_position\',
        		'NOW()', 'NOW()', 'NOW()', 'GEN', '@');
    			INSERT INTO medco_ont.genomic (c_hlevel, c_fullname, c_name, c_synonym_cd, c_visualattributes, c_totalnum, c_basecode,
				c_facttablecolumn, c_tablename, c_columnname, c_columndatatype, c_operator, c_dimcode, c_comment, c_tooltip, update_date,
        		download_date, import_date, valuetype_cd, m_applied_path) values
        		('2', '\medco\genomic\variant\', 'Variant Name', 'N', 'LA', '0', 'GEN:variant_name', 'concept_cd', 'concept_dimension', 'concept_path',
        		'T', 'LIKE', '\medco\genomic\variant\', 'Variant Name', '\medco\genomic\variant\',
        		'NOW()', 'NOW()', 'NOW()', 'GEN', '@');` + "\n"

	for i := 0; i < len(TablenamesOntology); i++ {

		//TODO: Delete this please
		if TablenamesOntology[i] != ONT+"non_sensitive_clear" {
			loading += "TRUNCATE " + TablenamesOntology[i] + ";\n"
			loading += `\copy ` + TablenamesOntology[i] + ` FROM '` + FilePathsOntology[i] + `' ESCAPE '"' DELIMITER ',' CSV;` + "\n"
		}
	}
	loading += "\n"
	// create annotations table
	loading += `DROP TABLE IF EXISTS genomic_annotations.hugo_gene_symbol;` + "\n"
	loading += `CREATE TABLE genomic_annotations.hugo_gene_symbol as select distinct hugo_gene_symbol as annotation_value from genomic_annotations.genomic_annotations;` + "\n"

	loading += `DROP TABLE IF EXISTS genomic_annotations.protein_change;` + "\n"
	loading += `CREATE TABLE genomic_annotations.protein_change as select distinct protein_change as annotation_value from genomic_annotations.genomic_annotations;` + "\n"

	loading += `DROP TABLE IF EXISTS genomic_annotations.variant_name;` + "\n"
	loading += `CREATE TABLE genomic_annotations.variant_name as select distinct variant_name as annotation_value from genomic_annotations.genomic_annotations;` + "\n"

	loading += `ALTER TABLE medco_ont.genomic OWNER TO i2b2;
    			ALTER TABLE medco_ont.clinical_sensitive OWNER TO i2b2;
    			ALTER TABLE medco_ont.clinical_non_sensitive OWNER TO i2b2;` + "\n"

	loading += "COMMIT;\n"
	loading += "EOSQL"

	_, err = fp.WriteString(loading)
	if err != nil {
		return err
	}

	fp.Close()

	return nil

}

// GenerateLoadingDataScript creates a load dataset .sql script
func GenerateLoadingDataScript(databaseS loader.DBSettings) error {
	fp, err := os.Create(FileBashPath[1])
	if err != nil {
		return err
	}

	loading := `#!/usr/bin/env bash` + "\n" + "\n" + `PGPASSWORD=` + databaseS.DBpassword + ` psql -v ON_ERROR_STOP=1 -h "` + databaseS.DBhost +
		`" -U "` + databaseS.DBuser + `" -p ` + strconv.FormatInt(int64(databaseS.DBport), 10) + ` -d "` + databaseS.DBname + `" <<-EOSQL` + "\n"

	loading += "BEGIN;\n"
	for i := 0; i < len(TablenamesData); i++ {
		loading += "TRUNCATE " + TablenamesData[i] + ";\n"
		loading += `\copy ` + TablenamesData[i] + ` FROM '` + FilePathsData[i] + `' ESCAPE '"' DELIMITER ',' CSV;` + "\n"
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

// LoadOntologyFiles executes the loading script for the ontology
func LoadOntologyFiles() error {
	// Display just the stderr if an error occurs
	cmd := exec.Command("/bin/sh", FileBashPath[0])
	stderr := &bytes.Buffer{} // make sure to import bytes
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		log.LLvl1("Error when running command.  Error log:", stderr.String())
		log.LLvl1("Got command status:", err.Error())
		return err
	}

	return nil
}

// LoadDataFiles executes the loading script for the dataset
func LoadDataFiles() error {
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
}

// GenerateOntologyFiles generates the .csv files that 'belong' to the whole ontology (metadata & medco)
func GenerateOntologyFiles(group *onet.Roster, entryPointIdx int, fOntClinical, fOntGenomic *os.File, mapSensitive map[string]struct{}) error {
	parsingTime := time.Duration(0)
	startParsing := time.Now()

	writeMedCoOntologyClearHeader()
	writeMedCoOntologyEncHeader()
	writeMedCoSensitiveTaggedHeader()

	allSensitiveIDs := make(map[int64]SensitiveIDValue, NumElMap) // maps the EncID(s) to the concept path
	toTraverseIndex := make([]int, 0)                             // the indexes of the columns that matter

	encID := int64(1)   // clinical sensitive IDs
	clearID := int64(1) // clinical non-sensitive IDs

	// load clinical ontology
	reader := csv.NewReader(fOntClinical)
	reader.Comma = '\t'

	first := true
	headerClinical := make([]string, 0)
	for {
		// read just one record, but we could ReadAll() as well
		record, err := reader.Read()
		// end-of-file is fitted into err
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		// if it is not a commented line
		if len(record) > 0 && string(record[0]) != "" && string(record[0][0:1]) != "#" {
			// the HEADER
			if first == true {
				for i, rec := range record {
					// skip SampleID and PatientID and other similar fields
					if _, ok := ToIgnore[rec]; !ok {
						// sensitive
						if _, ok := mapSensitive[rec]; ok || AllSensitive == true {
							if err := writeMedCoOntologyEnc(rec); err != nil {
								return err
							}
							// we don't generate the MetadataOntologyEnc because we will do this afterwards (so that we only perform 1 DDT with all sensitive elements)
						} else {
							if err := writeMedCoOntologyClear(rec); err != nil {
								return err
							}
						}
						headerClinical = append(headerClinical, rec)
						toTraverseIndex = append(toTraverseIndex, i)
					}
				}
				first = false
				// the RECORDS
			} else {

				j := 0
				for _, i := range toTraverseIndex {

					// uncomment if you want to include the <empty> fields as part of the ontology
					/*if record[i] == "" || record[i] == "NA" {
						record[i] = "<empty>"
					}*/

					// skip empty fields
					if record[i] == "" || record[i] == "NA" {
						j++
						continue
					}

					// sensitive
					if _, ok := mapSensitive[headerClinical[j]]; ok || AllSensitive == true {
						// if concept path does not exist
						if _, ok := OntValues[ConceptPath{Field: headerClinical[j], Record: record[i]}]; ok == false {
							if err := writeMedCoOntologyLeafEnc(headerClinical[j], record[i], encID); err != nil {
								return err
							}
							// we don't generate the MetadataOntologyLeafEnc because we will do this afterwards (so that we only perform 1 DDT with all sensitive elements)
							allSensitiveIDs[encID] = SensitiveIDValue{CP: ConceptPath{Field: headerClinical[j], Record: record[i]}, Annotation: "NA"}
							OntValues[ConceptPath{Field: headerClinical[j], Record: record[i]}] = ConceptID{Identifier: "E", Value: encID}
							encID++
						}
						// non-sensitive
					} else {
						// if concept path does not exist
						if _, ok := OntValues[ConceptPath{Field: headerClinical[j], Record: record[i]}]; ok == false {
							if err := writeMedCoOntologyLeafClear(headerClinical[j], record[i], clearID); err != nil {
								return err
							}

							OntValues[ConceptPath{Field: headerClinical[j], Record: record[i]}] = ConceptID{Identifier: "C", Value: clearID}
							clearID++
						}

					}
					j++
				}

			}
		}
	}

	fOntClinical.Close()

	log.LLvl1("Finished parsing the clinical ontology... (", len(allSensitiveIDs), ")")

	// load genomic
	reader = csv.NewReader(fOntGenomic)
	reader.Comma = '\t'

	first = true
	headerGenomic := make([]string, 0)
	// this arrays stores the indexes of the fields we need to use to generate a genomic id
	indexGenVariant := make(map[string]int)

	progress := int64(0)
	for {
		// read just one record, but we could ReadAll() as well
		record, err := reader.Read()
		progress++

		// end-of-file is fitted into err
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		// for every 100,000 rows parsed print a message
		if progress%100000 == 0 {
			log.LLvl1("Continuing parsing the genomic ontology... (", progress, ")")
		}

		// if it is not a commented line
		if len(record) > 0 && string(record[0]) != "" && string(record[0][0:1]) != "#" {

			// the HEADER
			if first == true {
				for i, el := range record {
					// the fields we need to generate the genomic id
					if val, ok := TranslationDic[el]; ok {
						indexGenVariant[val] = i
					}
					headerGenomic = append(headerGenomic, el)

				}
				first = false
			} else {
				// the number of genomic ids does not match the number of distinct mutation because if the RA is too big we discard the mutation
				genomicID, err := generateGenomicID(indexGenVariant, record)

				// if genomic id already exist we don't need to add it to the medco_ont.genomic_annotations
				if _, ok := allSensitiveIDs[genomicID]; ok == false && err == nil {
					allSensitiveIDs[genomicID] = SensitiveIDValue{CP: ConceptPath{Field: strconv.FormatInt(genomicID, 10), Record: ""}, Annotation: generateMedCoOntologyGenomicAnnotation(headerGenomic, record)}
				}
			}

		}

	}

	fOntGenomic.Close()

	log.LLvl1("Finished parsing the genomic ontology... (", len(allSensitiveIDs), ")")

	// convert the map of sensitive IDs to a slice (this is what the DDT service/protocol gets)
	listSensitiveIDs := make([]int64, 0)
	annotations := make([]string, 0)
	keyForSensitiveIDs := make([]ConceptPath, 0)
	for k, v := range allSensitiveIDs {
		listSensitiveIDs = append(listSensitiveIDs, k)
		annotations = append(annotations, v.Annotation)
		keyForSensitiveIDs = append(keyForSensitiveIDs, v.CP)
	}

	parsingTime += time.Since(startParsing)

	// encrypt sensitive ids
	listEncryptedElements := EncryptElements(listSensitiveIDs, group)
	if err := writeMedCoOntologyGenomicAnnotations(listEncryptedElements, annotations); err != nil {
		return err
	}

	// write the tagged values
	taggedValues, err := TagElements(listEncryptedElements, group, entryPointIdx)
	if err != nil {
		return err
	}

	startParsing = time.Now()
	err = writeMedCoSensitiveTagged(taggedValues, keyForSensitiveIDs)
	parsingTime += time.Since(startParsing)

	log.LLvl1("Parsing all ontology files took (", parsingTime, ")")

	return err
}

// GenerateDataFiles generates the .csv files that 'belong' to the dataset (demodata)
func GenerateDataFiles(group *onet.Roster, fClinical, fGenomic *os.File) error {
	parsingTime := time.Duration(0)
	startParsing := time.Now()

	// patient_id column index
	pidIndex := 0
	// encounter_id (sample_id) column index
	eidIndex := 0
	// patient_id counter
	pid := int64(1)
	// encounter_id counter
	eid := int64(1)

	ontValuesSmallCopy := make(map[ConceptPath]bool) // reduced set of ontology data to ensure that no repeated elements are added to the concept dimension table
	visitMapping := make(map[string]int64)           // map a sample ID to a numeric ID
	patientMapping := make(map[string]int64)         // map a patient ID to a numeric ID
	toTraverseIndex := make([]int, 0)                // the indexes of the columns that matter

	if err := writeDemodataProviderDimension(); err != nil {
		return err
	}

	// load clinical
	reader := csv.NewReader(fClinical)
	reader.Comma = '\t'

	first := true
	headerClinical := make([]string, 0)
	for {
		// read just one record, but we could ReadAll() as well
		record, err := reader.Read()
		// end-of-file is fitted into err
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		// if it is not a commented line
		if len(record) > 0 && string(record[0]) != "" && string(record[0][0:1]) != "#" {

			// the HEADER
			if first == true {

				for i, rec := range record {
					// skip SampleID and PatientID and other similar fields
					if _, ok := ToIgnore[rec]; !ok {
						headerClinical = append(headerClinical, record[i])
						toTraverseIndex = append(toTraverseIndex, i)
					} else {
						// if no keep track of the index of the patient_id and encounter_id (sample_id)
						if TranslationDic[rec] == "PATIENT_ID" {
							pidIndex = i
						} else if TranslationDic[rec] == "SAMPLE_ID" {
							eidIndex = i
						}
					}
				}
				first = false
			} else {
				// patient not yet exists
				if _, ok := patientMapping[record[pidIndex]]; ok == false {
					patientMapping[record[pidIndex]] = pid

					if err := writeDemodataPatientMapping(record[pidIndex], patientMapping[record[pidIndex]]); err != nil {
						return err
					}
					if err := writeDemodataPatientDimension(group, patientMapping[record[pidIndex]]); err != nil {
						return err
					}

					pid++
				}

				// sample not yet exists
				if _, ok := visitMapping[record[eidIndex]]; ok == false {
					visitMapping[record[eidIndex]] = eid

					if err := writeDemodataEncounterMapping(record[eidIndex], record[pidIndex], visitMapping[record[eidIndex]]); err != nil {
						return err
					}
					if err := writeDemodataVisitDimension(visitMapping[record[eidIndex]], patientMapping[record[pidIndex]]); err != nil {
						return err
					}

					eid++
				}

				j := 0
				for _, i := range toTraverseIndex {

					// uncomment if you want to include the <empty> fields as part of the ontology
					/*if record[i] == "" || record[i] == "NA" {
						record[i] = "<empty>"
					}*/

					// skip empty fields
					if record[i] == "" || record[i] == "NA" {
						j++
						continue
					}

					// check if it exists in the ontology
					if _, ok := OntValues[ConceptPath{Field: headerClinical[j], Record: record[i]}]; ok == true {
						// sensitive
						if OntValues[ConceptPath{Field: headerClinical[j], Record: record[i]}].Identifier != "C" {
							// if concept path does not exist
							if _, ok := ontValuesSmallCopy[ConceptPath{Field: headerClinical[j], Record: record[i]}]; ok == false {
								if err := writeDemodataConceptDimensionTaggedConcepts(headerClinical[j], record[i]); err != nil {
									return err
								}
								ontValuesSmallCopy[ConceptPath{Field: headerClinical[j], Record: record[i]}] = true
							}

							if err := writeDemodataObservationFactEnc(OntValues[ConceptPath{Field: headerClinical[j], Record: record[i]}].Value,
								patientMapping[record[pidIndex]],
								visitMapping[record[eidIndex]]); err != nil {
								return err
							}
							// non-sensitive
						} else {
							// if concept path does not exist
							if _, ok := ontValuesSmallCopy[ConceptPath{Field: headerClinical[j], Record: record[i]}]; ok == false {
								if err := writeDemodataConceptDimensionCleartextConcepts(headerClinical[j], record[i]); err != nil {
									return err
								}
								ontValuesSmallCopy[ConceptPath{Field: headerClinical[j], Record: record[i]}] = true
							}

							if err := writeDemodataObservationFactClear(OntValues[ConceptPath{Field: headerClinical[j], Record: record[i]}].Value,
								patientMapping[record[pidIndex]],
								visitMapping[record[eidIndex]]); err != nil {
								return err
							}
						}
					} else {
						log.Fatal("There are elements in the dataset that do not belong to the existing ontology")
						return err
					}
					j++
				}

			}
		}
	}
	fClinical.Close()

	log.LLvl1("Finished parsing the clinical dataset...")

	// load genomic
	reader = csv.NewReader(fGenomic)
	reader.Comma = '\t'

	first = true
	headerGenomic := make([]string, 0)
	// this arrays stores the indexes of the fields we need to use to generate a genomic id
	indexGenVariant := make(map[string]int)
	for {
		// read just one record, but we could ReadAll() as well
		record, err := reader.Read()

		// end-of-file is fitted into err
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		// if it is not a commented line
		if len(record) > 0 && string(record[0]) != "" && string(record[0][0:1]) != "#" {

			// the HEADER
			if first == true {
				for i, el := range record {
					if val, ok := TranslationDic[el]; ok {
						indexGenVariant[val] = i
					}

					// if no keep track of the index of the patient_id and encounter_id (sample_id)
					if TranslationDic[el] == "PATIENT_ID" {
						pidIndex = i
					} else if TranslationDic[el] == "SAMPLE_ID" {
						eidIndex = i
					}

					headerGenomic = append(headerGenomic, el)

				}
				first = false
			} else {
				genomicID, err := generateGenomicID(indexGenVariant, record)

				if err == nil {

					// check if it exists in the ontology
					if _, ok := OntValues[ConceptPath{Field: strconv.FormatInt(genomicID, 10), Record: ""}]; ok == true {
						// if concept path does not exist
						if _, ok := ontValuesSmallCopy[ConceptPath{Field: strconv.FormatInt(genomicID, 10), Record: ""}]; ok == false {
							if err := writeDemodataConceptDimensionTaggedConcepts(strconv.FormatInt(genomicID, 10), ""); err != nil {
								return err
							}
							ontValuesSmallCopy[ConceptPath{Field: strconv.FormatInt(genomicID, 10), Record: ""}] = true
						}

						if err := writeDemodataObservationFactEnc(OntValues[ConceptPath{Field: strconv.FormatInt(genomicID, 10), Record: ""}].Value,
							patientMapping[record[pidIndex]],
							visitMapping[record[eidIndex]]); err != nil {
							return err
						}
					} else {
						log.Fatal("There are elements in the dataset that do not belong to the existing ontology")
						return err
					}
				}
			}

		}
	}

	fGenomic.Close()

	parsingTime += time.Since(startParsing)
	log.LLvl1("Finished parsing the genomic dataset...")
	log.LLvl1("Parsing all dataset files took (", parsingTime, ")")

	log.LLvl1("The End. Only loading left...")

	return nil
}

func writeMedCoOntologyEncHeader() error {
	clinicalSensitive := `"2","\medco\clinical\sensitive\","MedCo Clinical Sensitive Ontology","N","CA","0",,,"concept_cd","concept_dimension","concept_path","T","LIKE","\medco\clinical\sensitive\","MedCo Clinical Sensitive Ontology","\medco\clinical\sensitive\","NOW()","NOW()","NOW()",,"ENC_ID","@",,,,` + "\n"

	_, err := FileHandlers[0].WriteString(clinicalSensitive)

	if err != nil {
		log.Fatal("Error in the writeMedCoOntologyEnc():", err)
		return err
	}

	return nil
}

func writeMedCoOntologyEnc(el string) error {
	el = SanitizeHeader(el)

	/*clinicalSensitive := `INSERT INTO medco_ont.clinical_sensitive VALUES (3, '\medco\clinical\sensitive\` + el + `\', '` + el + `', 'N', 'CA', NULL, NULL, NULL, 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE',
	  '\medco\clinical\sensitive\` + el + `\', 'Sensitive field encrypted by Unlynx', '\medco\clinical\sensitive\` + el + `\',
	   'NOW()', NULL, NULL, NULL, 'ENC_ID', '@', NULL, NULL, NULL, NULL);` + "\n"*/

	clinicalSensitive := `"3","\medco\clinical\sensitive\` + el + `\","` + el + `","N","CA",,,,"concept_cd","concept_dimension","concept_path","T","LIKE","\medco\clinical\sensitive\` + el + `\","Sensitive field encrypted by Unlynx","\medco\clinical\sensitive\` + el + `\","NOW()",,,,"ENC_ID","@",,,,` + "\n"

	_, err := FileHandlers[0].WriteString(clinicalSensitive)

	if err != nil {
		log.Fatal("Error in the writeMedCoOntologyEnc():", err)
		return err
	}

	return nil
}

func writeMedCoOntologyLeafEnc(field, el string, id int64) error {
	field = SanitizeHeader(field)

	/*clinicalSensitive := `INSERT INTO medco_ont.clinical_sensitive VALUES (4, '\medco\clinical\sensitive\` + field + `\` + el + `\', '` + el + `', 'N', 'LA', NULL, 'ENC_ID:` + strconv.FormatInt(id, 10) + `', NULL, 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE',
	  '\medco\clinical\sensitive\` + field + `\` + el + `\', 'Sensitive value encrypted by Unlynx',  '\medco\clinical\sensitive\` + field + `\` + el + `\',
	   'NOW()', NULL, NULL, NULL, 'ENC_ID', '@', NULL, NULL, NULL, NULL);` + "\n"*/

	clinicalSensitive := `"4","\medco\clinical\sensitive\` + field + `\` + el + `\","` + el + `","N","LA",,"ENC_ID:` + strconv.FormatInt(id, 10) + `",,"concept_cd","concept_dimension","concept_path","T","LIKE","\medco\clinical\sensitive\` + field + `\` + el + `\","Sensitive value encrypted by Unlynx","\medco\clinical\sensitive\` + field + `\` + el + `\","NOW()",,,,"ENC_ID","@",,,,` + "\n"

	_, err := FileHandlers[0].WriteString(clinicalSensitive)

	if err != nil {
		log.Fatal("Error in the writeMedCoOntologyLeafEnc():", err)
		return err
	}

	return nil
}

func writeMedCoOntologyClearHeader() error {
	clinical := `"2","\medco\clinical\nonsensitive\","MedCo Clinical Non-Sensitive Ontology","N","CA","0",,,"concept_cd","concept_dimension","concept_path","T","LIKE","\medco\clinical\nonsensitive\","MedCo Clinical Non-Sensitive Ontology","\medco\clinical\nonsensitive\","NOW()","NOW()","NOW()",,"CLEAR","@",,,,` + "\n"

	_, err := FileHandlers[1].WriteString(clinical)

	if err != nil {
		log.Fatal("Error in the writeMedCoOntologyClear():", err)
		return err
	}

	return nil
}

func writeMedCoOntologyClear(el string) error {
	el = SanitizeHeader(el)

	/*clinical := `INSERT INTO medco_ont.clinical_non_sensitive VALUES (3, '\medco\clinical\nonsensitive\` + el + `\', '` + el + `', 'N', 'CA', NULL, NULL, NULL, 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE',
	  '\medco\clinical\nonsensitive\` + el + `\', 'Non-sensitive field', '\medco\clinical\nonsensitive\` + el + `\',
	   'NOW()', NULL, NULL, NULL, 'CLEAR', '@', NULL, NULL, NULL, NULL);` + "\n"*/

	clinical := `"3","\medco\clinical\nonsensitive\` + el + `\","` + el + `","N","CA",,,,"concept_cd","concept_dimension","concept_path","T","LIKE","\medco\clinical\nonsensitive\` + el + `\","Non-sensitive field","\medco\clinical\nonsensitive\` + el + `\","NOW()",,,,"CLEAR","@",,,,` + "\n"

	_, err := FileHandlers[1].WriteString(clinical)

	if err != nil {
		log.Fatal("Error in the writeMedCoOntologyClear():", err)
		return err
	}

	return nil
}

func writeMedCoOntologyLeafClear(field, el string, id int64) error {
	field = SanitizeHeader(field)

	/*clinical := `INSERT INTO medco_ont.clinical_non_sensitive VALUES (4, '\medco\clinical\nonsensitive\` + field + `\` + el + `\', '` + el + `', 'N', 'LA', NULL, 'CLEAR:` + strconv.FormatInt(id, 10) + `', NULL, 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE',
	  '\medco\clinical\nonsensitive\` + field + `\` + el + `\', 'Non-sensitive value',  '\medco\clinical\sensitive\` + field + `\` + el + `\',
	   'NOW()', NULL, NULL, NULL, 'CLEAR', '@', NULL, NULL, NULL, NULL);` + "\n"*/

	clinical := `"4","\medco\clinical\nonsensitive\` + field + `\` + el + `\","` + el + `","N","LA",,"CLEAR:` + strconv.FormatInt(id, 10) + `",,"concept_cd","concept_dimension","concept_path","T","LIKE","\medco\clinical\nonsensitive\` + field + `\` + el + `\","Non-sensitive value","\medco\clinical\sensitive\` + field + `\` + el + `\","NOW()",,,,"CLEAR","@",,,,` + "\n"

	_, err := FileHandlers[1].WriteString(clinical)

	if err != nil {
		log.Fatal("Error in the writeMedCoOntologyLeafClear():", err)
		return err
	}

	return nil
}

func generateGenomicID(indexGenVariant map[string]int, record []string) (int64, error) {

	// if the ref and alt are too big ignore them (for now....)
	if len(record[indexGenVariant["RA"]]) > 6 || len(record[indexGenVariant["TSA1"]]) > 6 {
		return int64(-1), errors.New("reference and/or Alternate base size is bigger than the maximum allowed")
	}

	// generate id
	aux, err := strconv.ParseInt(record[indexGenVariant["SP"]], 10, 64)
	if err != nil {
		return int64(-1), err
	}

	id, err := loader.GetVariantID(record[indexGenVariant["CHR"]], aux, record[indexGenVariant["RA"]], record[indexGenVariant["TSA1"]])
	if err != nil {
		return int64(-1), err
	}

	return id, nil

}

func generateMedCoOntologyGenomicAnnotation(fields []string, record []string) string {
	// genomic info
	chr, sp, ra, tsa1, tsa2 := "?", "?", "?", "?", "?"

	// annotations that are to be queried
	queryFields := ""
	// annotations that are NOT to be queried (at least in a fast way)
	otherFields := ""

	for i, el := range record {
		// if element is CHR, SP, RA, TSA1
		if val, ok := TranslationDic[fields[i]]; ok == true {
			if val == "CHR" && el != "" {
				chr = el
			} else if val == "SP" && el != "" {
				sp = el
			} else if val == "RA" && el != "" && el != "-" {
				ra = el
			} else if val == "TSA1" && el != "" && el != "-" {
				tsa1 = el
			} else if val == "TSA2" && el != "" && el != "-" {
				tsa2 = el
			}
			// if element is selected to be queried
		} else if _, ok := AnnotationsToQuery[fields[i]]; ok == true {
			queryFields += `"` + el + `",`
			// if element is not to be ignored
		} else if _, ok := ToIgnore[fields[i]]; ok == false {
			field := SanitizeHeader(fields[i])
			otherFields += field + "=" + el + ";"
		}
	}
	// remove the last ", " and "; "
	queryFields = queryFields[:len(queryFields)-1]
	otherFields = otherFields[:len(otherFields)-2]

	// tsa1  tsa2
	// nil   nil     Unknown
	// A     nil     Unknown
	// nil   B       Unknown
	// A     B       Heterozygous
	// A     A       Homozygous
	zigosity := ""
	alt := ""
	if tsa1 == "?" && tsa2 == "?" {
		zigosity = "Unknown"
		alt = "?"
	} else if tsa1 == "?" && tsa2 != "?" {
		zigosity = "Unknown"
		alt = tsa2
	} else if tsa1 != "?" && tsa2 == "?" {
		zigosity = "Unknown"
		alt = tsa1
	} else if tsa1 == tsa2 {
		zigosity = "Homozygous"
		alt = tsa1
	} else {
		zigosity = "Heterozygous"
		alt = tsa1
	}

	annotation := `"` + chr + `:` + sp + `:` + ra + `>` + alt + `",` + queryFields + `,"` + zigosity + `;` + otherFields + `"` + "\n"
	return annotation
}

func writeMedCoOntologyGenomicAnnotations(listEncryptedElements *libunlynx.CipherVector, annotations []string) error {
	for i, annotation := range annotations {
		if annotation != "NA" && annotation != "" {
			ciphertextStr := (*listEncryptedElements)[i].Serialize()
			_, err := FileHandlers[2].WriteString(`"` + ciphertextStr + `",` + annotation)

			if err != nil {
				log.Fatal("Error in the writeMedCoOntologyGenomicAnnotations():", err)
				return err
			}
		}
	}
	return nil
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// EncryptElements encrypts the genomic ids
func EncryptElements(list []int64, group *onet.Roster) *libunlynx.CipherVector {
	// ENCRYPTION
	start := time.Now()
	listEncryptedElements := make(libunlynx.CipherVector, len(list))

	// parallelize the encryption (we need this because this is so slow)
	blockSize := int64(len(list)) / int64(NumThreads)
	wg := libunlynx.StartParallelize(NumThreads)
	for i := 0; i < NumThreads; i++ {
		blockEnd := int64(i)*blockSize + blockSize
		// the last goroutine gets the remaining content of the array
		if i == NumThreads-1 {
			blockEnd = int64(len(list))
		}

		go func(init int64, length int64) {
			defer wg.Done()
			for j := init; j < length; j++ {
				listEncryptedElements[j] = *libunlynx.EncryptInt(group.Aggregate, list[j])
			}
			log.LLvl1("Encrypted (", length-init, ")elements")
		}(int64(i)*blockSize, blockEnd)
	}
	libunlynx.EndParallelize(wg)

	log.LLvl1("Finished encrypting the sensitive data... (", time.Since(start), ")")

	return &listEncryptedElements

}

// TagElements tags the genomic ids to allow for the comparison
func TagElements(listEncryptedElements *libunlynx.CipherVector, group *onet.Roster, entryPointIdx int) ([]libunlynx.GroupingKey, error) {
	// TAGGING
	start := time.Now()
	client := servicesmedco.NewMedCoClient(group.List[entryPointIdx], strconv.Itoa(entryPointIdx))
	_, result, tr, err := client.SendSurveyDDTRequestTerms(
		group, // Roster
		servicesmedco.SurveyID("tagging_loading_phase"), // SurveyID
		*listEncryptedElements,                          // Encrypted query terms to tag
		false,                                           // compute proofs?
		Testing,
	)

	if err != nil {
		log.Fatal("Error during DDT:", err)
		return nil, err
	}

	totalTime := time.Since(start)

	tr.DDTRequestTimeCommunication = totalTime - tr.DDTRequestTimeExec

	log.LLvl1("DDT took: exec -", tr.DDTRequestTimeExec, "commun -", tr.DDTRequestTimeCommunication)

	log.LLvl1("Finished tagging the sensitive data... (", totalTime, ")")

	return result, nil
}

func writeMedCoSensitiveTaggedHeader() error {
	sensitive := `"1","\medco\tagged\","MedCo Sensitive Tagged Ontology","N","CA","0",,,"concept_cd","concept_dimension","concept_path","T","LIKE","\medco\tagged\","MedCo Sensitive Tagged Ontology","\medco\tagged\","NOW()","NOW()","NOW()",,"TAG_ID","@",,,,` + "\n"

	_, err := FileHandlers[3].WriteString(sensitive)

	if err != nil {
		log.Fatal("Error in the writeMedCoSensitiveTagged():", err)
		return err
	}

	return nil
}

func writeMedCoSensitiveTagged(list []libunlynx.GroupingKey, keyForSensitiveIDs []ConceptPath) error {

	if len(list) != len(keyForSensitiveIDs) {
		log.Fatal("The number of sensitive elements does not match the number of 'KeyForSensitiveID's.")
		return errors.New("")
	}

	tagIDs := make(map[int64]bool)

	for i, el := range list {
		// generate a tagID with 32bits (cannot be repeated)
		ok := false
		var tagID uint32

		// while random tag is not unique
		for ok == false {
			b, err := GenerateRandomBytes(4)

			if err != nil {
				log.Fatal("Error while generating random number", err)
				return err
			}

			tagID = binary.BigEndian.Uint32(b)

			// if random tag does not exist yet
			if _, okTagID := tagIDs[int64(tagID)]; okTagID == false {
				tagIDs[int64(tagID)] = true
				ok = true
			}
		}

		/*sensitive := `INSERT INTO medco_ont.sensitive_tagged VALUES (2, '\medco\tagged\` + string(el) + `\', '', 'N', 'LA ', NULL, 'TAG_ID:` + strconv.FormatUint(int64(tagID), 10) + `', NULL, 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE',
		'\medco\tagged\` + string(el) + `\', NULL, NULL, 'NOW()', NULL, NULL, NULL, 'TAG_ID', '@', NULL, NULL, NULL, NULL);` + "\n"*/

		sensitive := `"2","\medco\tagged\` + string(el) + `\","""","N","LA",,"TAG_ID:` + strconv.FormatInt(int64(tagID), 10) + `",,"concept_cd","concept_dimension","concept_path","T","LIKE","\medco\tagged\` + string(el) + `\",,,"NOW()",,,,"TAG_ID","@",,,,` + "\n"

		_, err := FileHandlers[3].WriteString(sensitive)

		if err != nil {
			log.Fatal("Error in the writeMedCoSensitiveTagged():", err)
			return err
		}

		OntValues[keyForSensitiveIDs[i]] = ConceptID{Identifier: string(el), Value: int64(tagID)}
	}
	return nil
}

func writeDemodataConceptDimensionCleartextConcepts(field, el string) error {
	/*cleartextConcepts := `INSERT INTO i2b2demodata.concept_dimension VALUES ('\medco\clinical\nonsensitive\` + field + `\` + record + `\', 'CLEAR:` + strconv.FormatInt(OntValues[ConceptPath{Field: field, Record: record}].Value, 10) + `', '` + record + `', NULL, NULL, NULL, 'NOW()', NULL, NULL);` + "\n"*/

	cleartextConcepts := `"\medco\clinical\nonsensitive\` + SanitizeHeader(field) + `\` + el + `\","CLEAR:` + strconv.FormatInt(OntValues[ConceptPath{Field: field, Record: el}].Value, 10) + `","` + el + `",,,,"NOW()",,` + "\n"

	_, err := FileHandlers[4].WriteString(cleartextConcepts)

	if err != nil {
		log.Fatal("Error in the writeDemodataConceptDimensionCleartextConcepts():", err)
		return err
	}

	return nil

}

func writeDemodataConceptDimensionTaggedConcepts(field string, el string) error {

	/*taggedConcepts := `INSERT INTO i2b2demodata.concept_dimension VALUES ('\medco\tagged\` + OntValues[ConceptPath{Field: field, Record: el}].Identifier + `\', 'TAG_ID:` + strconv.FormatInt(OntValues[ConceptPath{Field: field, Record: el}].Value, 10) + `', NULL, NULL, NULL, NULL, 'NOW()', NULL, NULL);` + "\n"*/

	taggedConcepts := `"\medco\tagged\` + OntValues[ConceptPath{Field: field, Record: el}].Identifier + `\","TAG_ID:` + strconv.FormatInt(OntValues[ConceptPath{Field: field, Record: el}].Value, 10) + `",,,,,"NOW()",,` + "\n"

	_, err := FileHandlers[4].WriteString(taggedConcepts)

	if err != nil {
		log.Fatal("Error in the writeDemodataConceptDimensionTaggedConcepts():", err)
		return err
	}

	return nil
}

func writeDemodataPatientMapping(el string, id int64) error {

	/*chuv := `INSERT INTO i2b2demodata.patient_mapping VALUES ('` + el + `', 'chuv', ` + strconv.FormatInt(id, 10) + `, NULL, 'Demo', NULL, NULL, NULL, 'NOW()', NULL, 1);` + "\n"*/

	chuv := `"` + el + `","chuv","` + strconv.FormatInt(id, 10) + `",,"Demo",,,,"NOW()",,"1"` + "\n"

	_, err := FileHandlers[5].WriteString(chuv)

	if err != nil {
		log.Fatal("Error in the writeDemodataPatientMapping()-Chuv:", err)
		return err
	}

	/*hive := `INSERT INTO i2b2demodata.patient_mapping VALUES ('` + strconv.FormatInt(id, 10) + `', 'HIVE', ` + strconv.FormatInt(id, 10) + `, 'A', 'HIVE', NULL, 'NOW()', 'NOW()', 'NOW()', 'edu.harvard.i2b2.crc', 1);` + "\n"*/

	hive := `"` + strconv.FormatInt(id, 10) + `","HIVE","` + strconv.FormatInt(id, 10) + `","A","HIVE",,"NOW()","NOW()","NOW()","edu.harvard.i2b2.crc","1"` + "\n"

	_, err = FileHandlers[5].WriteString(hive)

	if err != nil {
		log.Fatal("Error in the writeDemodataPatientMapping()-Hive:", err)
		return err
	}

	return nil

}

// TODO: No dummy data. Basically all flags are
func writeDemodataPatientDimension(group *onet.Roster, id int64) error {

	encryptedFlag := libunlynx.EncryptInt(group.Aggregate, 1)
	b := encryptedFlag.ToBytes()

	/*patientDimension := `INSERT INTO i2b2demodata.patient_dimension VALUES (` + strconv.FormatInt(id, 10) + `, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 'NOW()', NULL, 1, '` + base64.StdEncoding.EncodeToString(b) + `');` + "\n"*/

	patientDimension := `"` + strconv.FormatInt(id, 10) + `",,,,,,,,,,,,,,,,"NOW()",,"1","` + base64.StdEncoding.EncodeToString(b) + `"` + "\n"

	_, err := FileHandlers[6].WriteString(patientDimension)

	if err != nil {
		log.Fatal("Error in the writeDemodataPatientDimension()-Hive:", err)
		return err
	}

	return nil
}

func writeDemodataEncounterMapping(sampleID, patientID string, id int64) error {

	/*encounterChuv := `INSERT INTO i2b2demodata.encounter_mapping VALUES ('` + sampleID + `', 'chuv', 'Demo', ` + strconv.FormatInt(id, 10) + `, '` + patientID + `', 'chuv', NULL, NULL, NULL, NULL, 'NOW()', NULL, 1);` + "\n"*/

	encounterChuv := `"` + sampleID + `","chuv","Demo","` + strconv.FormatInt(id, 10) + `","` + patientID + `","chuv",,,,,"NOW()",,"1"` + "\n"

	_, err := FileHandlers[7].WriteString(encounterChuv)

	if err != nil {
		log.Fatal("Error in the writeDemodataEncounterMapping()-Chuv:", err)
		return err
	}

	/*encounterHive := `INSERT INTO i2b2demodata.encounter_mapping VALUES ('` + strconv.FormatInt(id, 10) + `', 'HIVE', 'HIVE', ` + strconv.FormatInt(id, 10) + `, '` + sampleID + `', 'chuv', 'A', NULL, 'NOW()', 'NOW()', 'NOW()', 'edu.harvard.i2b2.crc', 1);` + "\n"*/

	encounterHive := `"` + strconv.FormatInt(id, 10) + `","HIVE","HIVE","` + strconv.FormatInt(id, 10) + `","` + sampleID + `","chuv","A",,"NOW()","NOW()","NOW()","edu.harvard.i2b2.crc","1"` + "\n"

	_, err = FileHandlers[7].WriteString(encounterHive)

	if err != nil {
		log.Fatal("Error in the writeDemodataEncounterMapping()-Chuv:", err)
		return err
	}

	return nil
}

func writeDemodataVisitDimension(idV, idP int64) error {

	/*visit := `INSERT INTO i2b2demodata.visit_dimension VALUES (` + strconv.FormatInt(idV, 10) + `, ` + strconv.FormatInt(idP, 10) + `, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 'NOW()', 'chuv', 1);` + "\n"*/

	visit := `"` + strconv.FormatInt(idV, 10) + `","` + strconv.FormatInt(idP, 10) + `",,,,,,,,,,,"NOW()","chuv","1"` + "\n"

	_, err := FileHandlers[8].WriteString(visit)

	if err != nil {
		log.Fatal("Error in the writeDemodataVisitDimension():", err)
		return err
	}

	return nil
}

func writeDemodataProviderDimension() error {

	/*provider := `INSERT INTO i2b2demodata.provider_dimension VALUES ('chuv', '\medco\institutions\chuv\', 'chuv', NULL, NULL, NULL, 'NOW()', NULL, 1);` + "\n"*/

	provider := `"chuv","\medco\institutions\chuv\","chuv",,,,"NOW()",,"1"` + "\n"

	_, err := FileHandlers[9].WriteString(provider)

	if err != nil {
		log.Fatal("Error in the writeDemodateProviderDimension():", err)
		return err
	}

	return nil
}

func writeDemodataObservationFactClear(el, idP, idV int64) error {

	/*clear := `INSERT INTO i2b2demodata.observation_fact VALUES (` + strconv.FormatInt(idP, 10) + `, ` + strconv.FormatInt(idV, 10), 10) + `,
	'CLEAR:` + strconv.FormatInt(el, 10) + `', 'chuv', 'NOW()', '@', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL,
	'chuv', NULL, NULL, NULL, NULL, 'NOW()', NULL, 1, ` + strconv.FormatInt(TextSearchIndex, 10) + `);` + "\n"*/

	clear := `"` + strconv.FormatInt(idP, 10) + `","` + strconv.FormatInt(idV, 10) + `","CLEAR:` + strconv.FormatInt(el, 10) + `","chuv","NOW()","@","1",,,,,,,,"chuv",,,,,"NOW()",,"1","` + strconv.FormatInt(TextSearchIndex, 10) + `"` + "\n"

	_, err := FileHandlers[10].WriteString(clear)

	if err != nil {
		log.Fatal("Error in the writeDemodataObservationFactClear():", err)
		return err
	}

	TextSearchIndex++

	return nil
}

func writeDemodataObservationFactEnc(el int64, idP, idV int64) error {

	/*encrypted := `INSERT INTO i2b2demodata.observation_fact VALUES (` + strconv.FormatInt(idP, 10) + `, ` + strconv.FormatInt(idV, 10) + `, 'TAG_ID:` + strconv.FormatInt(el, 10) + `',
	'chuv', 'NOW()', '@', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 'chuv', NULL, NULL, NULL, NULL, 'NOW()', NULL, 1, ` + strconv.FormatInt(TextSearchIndex, 10) + `);` + "\n"*/

	encrypted := `"` + strconv.FormatInt(idP, 10) + `","` + strconv.FormatInt(idV, 10) + `","TAG_ID:` + strconv.FormatInt(el, 10) + `","chuv","NOW()","@","` + strconv.FormatInt(TextSearchIndex, 10) + `",,,,,,,,"chuv",,,,,"NOW()",,"1","` + strconv.FormatInt(TextSearchIndex, 10) + `"` + "\n"

	_, err := FileHandlers[10].WriteString(encrypted)

	if err != nil {
		log.Fatal("Error in the writeDemodataObservationFactEnc():", err)
		return err
	}

	TextSearchIndex++

	return nil

}

// SanitizeHeader gets and header name and transforms it in the form Xxx Yyy Zzz
func SanitizeHeader(header string) string {
	tokens := strings.Split(header, "_")
	for i, token := range tokens {
		tokens[i] = strings.Title(strings.ToLower(token))
	}

	res := ""
	for _, token := range tokens {
		res += token + " "
	}

	return res[:len(res)-1]
}
