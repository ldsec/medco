package loaderi2b2

import (
	"go.dedis.ch/onet/v3/log"
	"os"
	"sort"
	"strconv"
	"strings"
)

// medCoOntologyRecord is a record of a MedCo ontology table
type medCoOntologyRecord struct {
	nodeEncryptID      int64
	childrenEncryptIDs []int64

	hLevel           string
	fullname         string
	name             string
	synonymCD        string
	visualAttributes string
	totalNum         string
	baseCode         string
	metadataXML      string
	factTableColumn  string
	tablename        string
	columnName       string
	columnDataType   string
	operator         string
	dimCode          string
	comment          string
	tooltip          string
	adminColumns     administrativeColumns
	valueTypeCD      string
	appliedPath      string
	exclusionCD      string
}

// medCoTableInfo stores the concepts for a specific medco ontology table
type medCoTableInfo struct {
	clear     map[string][]*medCoOntologyRecord // we need a slice because multiple modifiers can have the same path
	sensitive map[string]*medCoOntologyRecord
}

// tablesMedCoOntology stores the MedCo ontology tables
var tablesMedCoOntology map[string]medCoTableInfo // the key is the name of the table

// headerMedCoOntology contains all the headers for the medco table
var headerMedCoOntology []string

// medCoOntologyFromLocalOntologyConcept generates a medCoOntologyRecord struct from localOntologyRecord struct
func medCoOntologyFromLocalOntologyConcept(localConcept *localOntologyRecord) *medCoOntologyRecord {
	ac := administrativeColumns{
		updateDate:     localConcept.adminColumns.updateDate,
		downloadDate:   localConcept.adminColumns.downloadDate,
		importDate:     localConcept.adminColumns.importDate,
		sourceSystemCD: localConcept.adminColumns.sourceSystemCD,
	}

	so := &medCoOntologyRecord{
		nodeEncryptID:      int64(-1), //signals that this medco ontology element is not sensitive so no need for an encrypt ID
		childrenEncryptIDs: nil,       //same thing as before
		hLevel:             localConcept.hLevel,
		fullname:           localConcept.fullname,
		name:               localConcept.name,
		synonymCD:          localConcept.synonymCD,
		visualAttributes:   localConcept.visualAttributes,
		totalNum:           localConcept.totalNum,
		baseCode:           localConcept.baseCode,
		metadataXML:        strings.Replace(localConcept.metadataXML, "\"", "\"\"", -1),
		factTableColumn:    localConcept.factTableColumn,
		tablename:          localConcept.tablename,
		columnName:         localConcept.columnName,
		columnDataType:     localConcept.columnDataType,
		operator:           localConcept.operator,
		dimCode:            localConcept.dimCode,
		comment:            localConcept.comment,
		tooltip:            localConcept.tooltip,
		adminColumns:       ac,
		valueTypeCD:        localConcept.valueTypeCD,
		appliedPath:        localConcept.appliedPath,
		exclusionCD:        localConcept.exclusionCD,
	}

	return so
}

// toCSVText writes the medCoOntologyRecord object in a way that can be added to a .csv file - "","","", etc.
func (so medCoOntologyRecord) toCSVText() string {
	if so.nodeEncryptID != int64(-1) && so.visualAttributes[:1] != "M" { // sensitive
		metadata := ""

		if so.visualAttributes[:1] == "C" { // if concept_parent_node
			metadata += "<?xml version=\"\"1.0\"\"?><ValueMetadata><Version>MedCo-0.1</Version><EncryptedType>CONCEPT_PARENT_NODE</EncryptedType>"
		} else if so.visualAttributes[:1] == "F" { // else if concept_internal_node
			metadata += "<?xml version=\"\"1.0\"\"?><ValueMetadata><Version>MedCo-0.1</Version><EncryptedType>CONCEPT_INTERNAL_NODE</EncryptedType><NodeEncryptID>" + strconv.FormatInt(so.nodeEncryptID, 10) + "</NodeEncryptID>"
		} else if so.visualAttributes[:1] == "L" { // else if concept_leaf
			metadata += "<?xml version=\"\"1.0\"\"?><ValueMetadata><Version>MedCo-0.1</Version><EncryptedType>CONCEPT_LEAF</EncryptedType><NodeEncryptID>" + strconv.FormatInt(so.nodeEncryptID, 10) + "</NodeEncryptID>"
		} else {
			log.Fatal("Wrong VisualAttribute")
		}

		// only internal and parent nodes can have children ;)
		// TODO we are appending all children IDs (split by ;) in a single xml attribute. We should find a cleaner way to do this
		if len(so.childrenEncryptIDs) > 0 && so.visualAttributes[:1] != "L" && so.visualAttributes[:1] != "R" {
			metadata += "<ChildrenEncryptIDs>\""
			for _, childID := range so.childrenEncryptIDs {
				metadata += strconv.FormatInt(childID, 10) + ";"
			}
			// remove last;
			metadata = metadata[:len(metadata)-1]

			metadata += "\"</ChildrenEncryptIDs>"
		}
		so.metadataXML = metadata + "</ValueMetadata>"
	}

	acString := "\"" + so.adminColumns.updateDate + "\"," + "\"" + so.adminColumns.downloadDate + "\"," + "\"" + so.adminColumns.importDate + "\"," + "\"" + so.adminColumns.sourceSystemCD + "\""
	finalString := "\"" + so.hLevel + "\"," + "\"" + so.fullname + "\"," + "\"" + so.name + "\"," + "\"" + so.synonymCD + "\"," + "\"" + so.visualAttributes + "\"," + "\"" + so.totalNum + "\"," +
		"\"" + so.baseCode + "\",\"" + so.metadataXML + "\"," + "\"" + so.factTableColumn + "\"," + "\"" + so.tablename + "\"," + "\"" + so.columnName + "\"," + "\"" + so.columnDataType + "\"," + "\"" + so.operator + "\"," +
		"\"" + so.dimCode + "\"," + "\"" + so.comment + "\"," + "\"" + so.tooltip + "\"," + acString + "," + "\"" + so.valueTypeCD + "\"," + "\"" + so.appliedPath + "\"," + "\"" + so.exclusionCD + "\""

	return strings.Replace(finalString, `"\N"`, "", -1)
}

// generateMedCoOntology generates all files for the medco ontology (these may include multiples tables)
func generateMedCoOntology() error {
	// initialize container structs and counters
	headerMedCoOntology = []string{"c_hlevel",
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

	for _, key := range ontologyFilesPaths {
		err := generateNewMedCoTableNonSensitive(strings.Split(key, "ONTOLOGY_")[1])
		if err != nil {
			log.Fatal("Error generating non sensitive [" + key + "].csv")
			return err
		}
		err = generateNewMedCoTableSensitive(strings.Split(key, "ONTOLOGY_")[1])
		if err != nil {
			log.Fatal("Error generating sensitive [" + key + "].csv")
			return err
		}
	}
	return nil
}

func generateNewMedCoTableNonSensitive(rawName string) error {
	csvOutputFile, err := os.Create(outputFilesPathsNonSensitive["MEDCO_"+rawName].Path)
	if err != nil {
		log.Fatal("Error opening [" + rawName + "].csv")
		return err
	}
	defer csvOutputFile.Close()

	headerString := ""
	for _, header := range headerMedCoOntology {
		headerString += "\"" + header + "\","
	}
	// remove the last ,
	csvOutputFile.WriteString(headerString[:len(headerString)-1] + "\n")

	keys := make([]string, 0, len(tablesMedCoOntology[rawName].clear))
	for key := range tablesMedCoOntology[rawName].clear {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		sos := tablesMedCoOntology[rawName].clear[key]
		for _, so := range sos {
			csvOutputFile.WriteString(so.toCSVText() + "\n")
		}
	}

	return nil
}

func generateNewMedCoTableSensitive(rawName string) error {
	csvOutputFile, err := os.Create(outputFilesPathsSensitive["MEDCO_"+rawName].Path)
	if err != nil {
		log.Fatal("Error opening [" + rawName + "].csv")
		return err
	}
	defer csvOutputFile.Close()

	headerString := ""
	for _, header := range headerMedCoOntology {
		headerString += "\"" + header + "\","
	}
	// remove the last ,
	csvOutputFile.WriteString(headerString[:len(headerString)-1] + "\n")

	updateChildrenEncryptIDs(rawName) //updates the childrenEncryptIDs of the internal and parent nodes

	keys := make([]string, 0, len(tablesMedCoOntology[rawName].sensitive))
	for key := range tablesMedCoOntology[rawName].sensitive {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		so := tablesMedCoOntology[rawName].sensitive[key]
		csvOutputFile.WriteString(so.toCSVText() + "\n")
	}

	return nil
}

// updateChildrenEncryptIDs updates the parent and internal concept nodes with the IDs of their respective children (name identifies the name of the ontology table)
func updateChildrenEncryptIDs(name string) {
	for _, so := range tablesMedCoOntology[name].sensitive {
		path := so.fullname
		for true {
			path = stripByLevel(path, 1, false)
			if path == "" {
				break
			}

			if val, ok := tablesMedCoOntology[name].sensitive[path]; ok {
				val.childrenEncryptIDs = append(val.childrenEncryptIDs, so.nodeEncryptID)
			}

		}
	}
}
