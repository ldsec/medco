package loaderi2b2

import (
	"fmt"
	servicesmedco "github.com/ldsec/medco/unlynx/services"
	libunlynx "github.com/ldsec/unlynx/lib"
	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// localOntologyRecord is a record of a local i2b2 ontology table
type localOntologyRecord struct {
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
	appliedPath      string
	adminColumns     administrativeColumns
	valueTypeCD      string
	exclusionCD      string
	path             string
	symbol           string

	// this only exists in the sensitive tagged
	pCoriBasecode string

	// only exists in some strange tables (like icd10_icd9)
	plainCode string
}

// administrativeColumns are a set of columns that exist in every i2b2 table
type administrativeColumns struct {
	updateDate      string
	downloadDate    string
	importDate      string
	sourceSystemCD  string
	uploadID        string
	textSearchIndex string
}

// tagAndID is a struct that contains both tag and tagID for a concept or modifier
type tagAndID struct {
	tag   libunlynx.GroupingKey
	tagID int64
}

type modifiersInfo struct {
	applied  []*localOntologyRecord
	excluded map[string]*localOntologyRecord //the key is the fullname of the modifier
}

var (
	// test defines whether we should run the DDT on testing environment (locally) or using real nodes
	test bool // testing environment
	// allSensitive is a flag that defines whether all concepts are to be considered sensitive or not (-allSens flag)
	allSensitive = false
	// listSensitiveConcepts list all sensitive concepts (paths) - MedCo and LOCAL (the bool is for nothing)
	listSensitiveConcepts map[string]struct{}
	// enableModifiers defines whether we are using modifiers or not
	enableModifiers bool
	// mapGeneratedConcepts contains the concepts generated from the modifiers
	mapGeneratedConcepts map[string][]*medCoOntologyRecord // the key is the path of the generator concept
	// idConcepts used to assign IDs (NodeEncryptIDs) to be encrypted to the different concepts
	idConcepts int64
	// tagIDConceptsUsed used to keep track of the number of TAG_IDs that have been used
	tagIDConceptsUsed int64
	// mapConceptPathToTag maps a sensitive concept path to its respective tag and tag_id
	mapConceptPathToTag map[string]tagAndID
	// mapModifiers contains all the modifiers contained in the ontology tables indexed by their m_applied_path
	mapModifiers map[string]*modifiersInfo
)

// tableLocalOntologyClear is the local i2b2 ontology table containing only the non sensitive concepts (it includes modifiers with non sensitive m_applied_path)
var tableLocalOntologyClear map[string][]*localOntologyRecord // we need a slice because multiple modifiers can have the same path

// headerLocalOntology contains the headers for the local i2b2 ontology table(s)
var headerLocalOntology []string

// localOntologyFromString generates a localOntologyRecord struct from a parsed line of a local i2b2 ontology file
func localOntologyFromString(line []string, plainCode bool) *localOntologyRecord {
	ac := administrativeColumns{
		updateDate:     line[17],
		downloadDate:   line[18],
		importDate:     line[19],
		sourceSystemCD: line[20],
	}

	lo := &localOntologyRecord{
		hLevel:           line[0],
		fullname:         strings.Replace(line[1], "\"", "\"\"", -1),
		name:             strings.Replace(line[2], "\"", "\"\"", -1),
		synonymCD:        line[3],
		visualAttributes: line[4],
		totalNum:         line[5],
		baseCode:         line[6],
		metadataXML:      strings.Replace(line[7], "\"", "\"\"", -1),
		factTableColumn:  line[8],
		tablename:        line[9],
		columnName:       line[10],
		columnDataType:   line[11],
		operator:         line[12],
		dimCode:          strings.Replace(line[13], "\"", "\"\"", -1),
		comment:          line[14],
		tooltip:          strings.Replace(line[15], "\"", "\"\"", -1),
		appliedPath:      line[16],
		adminColumns:     ac,
		valueTypeCD:      line[21],
		exclusionCD:      line[22],
		path:             strings.Replace(line[23], "\"", "\"\"", -1),
		symbol:           strings.Replace(line[24], "\"", "\"\"", -1),
	}

	if plainCode {
		lo.plainCode = line[25]
	}

	return lo

}

// toCSVText writes the localOntologyRecord object in a way that can be added to a .csv file - "","","", etc.
func (lo localOntologyRecord) toCSVText() string {
	acString := "\"" + lo.adminColumns.updateDate + "\"," + "\"" + lo.adminColumns.downloadDate + "\"," + "\"" + lo.adminColumns.importDate + "\"," + "\"" + lo.adminColumns.sourceSystemCD + "\""
	finalString := "\"" + lo.hLevel + "\"," + "\"" + lo.fullname + "\"," + "\"" + lo.name + "\"," + "\"" + lo.synonymCD + "\"," + "\"" + lo.visualAttributes + "\"," + "\"" + lo.totalNum + "\"," +
		"\"" + lo.baseCode + "\"," + "\"" + lo.metadataXML + "\"," + "\"" + lo.factTableColumn + "\"," + "\"" + lo.tablename + "\"," + "\"" + lo.columnName + "\"," + "\"" + lo.columnDataType + "\"," + "\"" + lo.operator + "\"," +
		"\"" + lo.dimCode + "\"," + "\"" + lo.comment + "\"," + "\"" + lo.tooltip + "\"," + "\"" + lo.appliedPath + "\"," + acString + "," + "\"" + lo.valueTypeCD + "\"," + "\"" + lo.exclusionCD + "\"," +
		"\"" + lo.path + "\"," + "\"" + lo.symbol + "\""

	if lo.plainCode != "" {
		finalString += ",\"" + lo.plainCode + "\""
	}

	return strings.Replace(finalString, `"\N"`, "", -1)
}

// localOntologySensitiveConceptToCSVText writes the tagging information of a concept of the local ontology in a way that can be added to a .csv file - "","","", etc.
func localOntologySensitiveConceptToCSVText(tag *libunlynx.GroupingKey, tagID int64) string {
	finalString := `"3","\medco\tagged\` + string(*tag) + `\","","N","LA ","\N","TAG_ID:` + strconv.FormatInt(tagID, 10) + `","\N","concept_cd","concept_dimension","concept_path","T","LIKE","\medco\tagged\concept\` + string(*tag) +
		`\","\N","\N","NOW()","\N","\N","\N","TAG_ID","@","\N","\N","\N","\N"`

	return strings.Replace(finalString, `"\N"`, "", -1)
}

// convertLocalOntology reads and parses all local i2b2 ontology tables and generates the corresponding .csv(s) (local, medco and adapter_mappings)
func convertLocalOntology(group *onet.Roster, entryPointIdx int) (err error) {
	// initialize container structs and counters
	idConcepts = 0
	tagIDConceptsUsed = 0
	tablesMedCoOntology = make(map[string]medCoTableInfo)
	mapConceptPathToTag = make(map[string]tagAndID)

	if enableModifiers {
		err = mapModifiersToAppliedPaths()
		if err != nil {
			return err
		}
		mapGeneratedConcepts = make(map[string][]*medCoOntologyRecord)
	}

	for _, key := range ontologyFilesPaths {
		rawName := strings.Split(key, "ONTOLOGY_")[1]
		err := parseLocalOntologyTable(group, entryPointIdx, key)
		if err != nil {
			log.Fatal("Error parsing [" + strings.ToLower(rawName) + "].csv")
			return err
		}
		err = convertClearLocalOntologyTable(rawName)
		if err != nil {
			log.Fatal("Error converting [" + strings.ToLower(rawName) + "].csv")
			return err
		}
	}

	err = convertSensitiveLocalOntologyTable()
	if err != nil {
		log.Fatal("Error converting [sensitive_tagged].csv")
		return err
	}

	return nil
}

// mapModifiersToAppliedPaths looks for the modifiers in the local i2b2 ontology tables and store them in mapModifiers
func mapModifiersToAppliedPaths() (err error) {

	mapModifiers = make(map[string]*modifiersInfo)

	for _, name := range ontologyFilesPaths {

		lines, err := readCSV(inputFilePaths[name])
		if err != nil {
			log.Fatal("Error in readCSV()")
			return err
		}

		headerLocalOntology = make([]string, 0)

		for _, header := range lines[0] {
			headerLocalOntology = append(headerLocalOntology, header)
		}

		plainCode := false
		if headerLocalOntology[len(headerLocalOntology)-1] == "plain_code" {
			plainCode = true
		}

		// the pcori_basecode
		headerLocalOntology = append(headerLocalOntology, "pcori_basecode")

		//skip header
		for _, line := range lines[1:] {
			lo := localOntologyFromString(line, plainCode)

			if strings.ToLower(lo.factTableColumn) == "modifier_cd" {

				if mapModifiers[lo.appliedPath] == nil {
					mapModifiers[lo.appliedPath] = new(modifiersInfo)
				}

				if lo.exclusionCD != "X" {
					mapModifiers[lo.appliedPath].applied = append(mapModifiers[lo.appliedPath].applied, lo)
				} else { //m_applied_path contains an exclusion path
					if mapModifiers[lo.appliedPath].excluded == nil {
						mapModifiers[lo.appliedPath].excluded = make(map[string]*localOntologyRecord, 0)
					}
					mapModifiers[lo.appliedPath].excluded[lo.fullname] = lo
				}
			}
		}
	}

	return

}

// parseLocalOntologyTable reads and parses the local i2b2 ontology table(s) (i2b2.csv, icd10_icd9.csv, etc.) and creates the MedCo ontology
func parseLocalOntologyTable(group *onet.Roster, entryPointIdx int, name string) error {
	lines, err := readCSV(inputFilePaths[name])
	if err != nil {
		log.Fatal("Error in readCSV()")
		return err
	}
	rawName := strings.Split(name, "ONTOLOGY_")[1]

	headerLocalOntology = make([]string, 0)
	tableLocalOntologyClear = make(map[string][]*localOntologyRecord)

	listConceptCD := make([]string, 0)
	allSensitiveConceptIDs := make([]int64, 0)

	/* structure of the i2b2 ontology tables (in order):

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
		headerLocalOntology = append(headerLocalOntology, header)
	}

	plainCode := false
	if headerLocalOntology[len(headerLocalOntology)-1] == "plain_code" {
		plainCode = true
	}

	// the pcori_basecode
	headerLocalOntology = append(headerLocalOntology, "pcori_basecode")

	//skip header
	for _, line := range lines[1:] {
		lo := localOntologyFromString(line, plainCode)

		// TODO for now we remove all synonyms from the i2b2 local ontology
		// if it is the original concept (N = original, Y = synonym)
		if strings.ToLower(lo.synonymCD) == "n" || strings.ToLower(lo.synonymCD) == "" {
			// create entry for medco ontology (direct copy)
			so := medCoOntologyFromLocalOntologyConcept(lo)

			_, sensitive := hasSensitiveParents(lo.fullname)

			// if it is sensitive or has a sensitive parent
			if sensitive {

				// add the original concept to the concepts to parse
				concepts := []*medCoOntologyRecord{so}

				if enableModifiers {
					generatedConcepts, err := generateConceptsFromModifiers(so)
					if err != nil {
						return err
					}
					concepts = append(concepts, generatedConcepts...)
					mapGeneratedConcepts[so.fullname] = generatedConcepts
				}
				parseSensitiveMedCoOntology(rawName, concepts, &listConceptCD, &allSensitiveConceptIDs)

			} else {

				if strings.ToLower(so.factTableColumn) != "modifier_cd" || enableModifiers {

					if strings.ToLower(so.factTableColumn) == "modifier_cd" {
						// if the modifier applies to sensitive concepts only, exclude it from the clear ontology
						_, sensitive := hasSensitiveParents(strings.TrimSuffix(lo.appliedPath, "%"))
						if sensitive {
							continue
						}
					}

					// add a new entry to the local ontology table
					tableLocalOntologyClear[lo.fullname] = append(tableLocalOntologyClear[lo.fullname], lo)
					// add a new entry to the medco ontology table
					if _, ok := tablesMedCoOntology[rawName]; ok {
						tablesMedCoOntology[rawName].clear[so.fullname] = append(tablesMedCoOntology[rawName].clear[so.fullname], so)
					} else {
						sensitive := make(map[string]*medCoOntologyRecord)
						clear := make(map[string][]*medCoOntologyRecord)
						clear[so.fullname] = append(clear[so.fullname], so)
						tablesMedCoOntology[rawName] = medCoTableInfo{clear: clear, sensitive: sensitive}
					}
				}
			}
		}
	}

	// if there are sensitive concepts
	if len(allSensitiveConceptIDs) > 0 {
		taggedConceptValues, err := encryptAndTag(allSensitiveConceptIDs, group, entryPointIdx)
		if err != nil {
			return err
		}

		// re-randomize TAG_IDs
		rand.Seed(time.Now().UnixNano())
		perm := rand.Perm(len(mapConceptPathToTag))

		// 'populate' map (Concept codes)
		// we create a permutation of [0, n] and then add #concepts_already_parsed
		for i, concept := range listConceptCD {
			var tmp = mapConceptPathToTag[concept]
			tmp.tagID = tagIDConceptsUsed + int64(perm[i])
			tmp.tag = taggedConceptValues[i]
			mapConceptPathToTag[concept] = tmp
		}

		tagIDConceptsUsed += int64(len(mapConceptPathToTag))
	}

	return nil
}

// hasSensitiveParents is a function that checks if a node whether in the LocalOntology or ConceptDimension has any siblings which are sensitive.
func hasSensitiveParents(conceptPath string) (string, bool) {
	if allSensitive == true {
		return "", true
	}

	temp := conceptPath

	isSensitive := false
	for temp != "" {
		if _, ok := listSensitiveConcepts[temp]; ok {
			isSensitive = true
			listSensitiveConcepts[conceptPath] = struct{}{}
			break
		}
		temp = stripByLevel(temp, 1, false)
	}
	return temp, isSensitive
}

// stripByLevel strips the concept path based on /. The number represents the stripping level, in other words,
// if number = 1 we strip the first element enclosed in /****/ and then on. Order means which side we start stripping: true (left-to-right),
// false (right-to-left)
func stripByLevel(conceptPath string, number int, order bool) string {
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

// generateConceptsFromModifiers looks for the modifiers that apply to the so concept
// and turn them into concepts by concatenating the modifiers' fullnames to the concept's
func generateConceptsFromModifiers(so *medCoOntologyRecord) (concepts []*medCoOntologyRecord, err error) {

	concepts = make([]*medCoOntologyRecord, 0)

	loModifiersAppliedTmp := make([]*localOntologyRecord, 0)
	loModifiersApplied := make([]*localOntologyRecord, 0)
	loModifiersExcluded := make(map[string]*localOntologyRecord, 0)

	// look for the modifiers that apply only to the so concept (m_applied_path with no trailing %)
	if m, ok := mapModifiers[so.fullname]; ok {
		for _, ma := range m.applied {
			loModifiersApplied = append(loModifiersAppliedTmp, ma)
		}
	}

	// look for all the other modifiers
	soFullname := so.fullname
	for soFullname != "" {
		if m, ok := mapModifiers[soFullname+"%"]; ok {
			for _, ma := range m.applied {
				loModifiersAppliedTmp = append(loModifiersAppliedTmp, ma)
			}
			for fullname, me := range m.excluded {
				loModifiersExcluded[fullname] = me
			}
		}
		soFullname = stripByLevel(soFullname, 1, false)
	}

	for _, loModifier := range loModifiersAppliedTmp {
		// if loModifier is not in the excluded path for the so concept
		if _, ok := loModifiersExcluded[loModifier.fullname]; !ok {
			loModifiersApplied = append(loModifiersApplied, loModifier)
		}
	}

	// modify the modifiers' columns to turn them into concepts
	for _, loModifier := range loModifiersApplied {

		concept := medCoOntologyFromLocalOntologyConcept(loModifier)

		// sum the c_hlevels
		conceptHLevel, err := strconv.Atoi(concept.hLevel)
		if err != nil {
			return nil, err
		}
		soHLevel, err := strconv.Atoi(so.hLevel)
		if err != nil {
			return nil, err
		}
		concept.hLevel = strconv.Itoa(conceptHLevel + soHLevel)

		// concatenate the fullnames
		concept.fullname = strings.TrimSuffix(so.fullname, "\\") + concept.fullname

		// concatenate the names
		concept.name = so.name + "\\" + concept.name

		switch concept.visualAttributes[0] {
		// modifier container
		case 'O':
			concept.visualAttributes = "C" + concept.visualAttributes[1:]
		// modifier folder
		case 'D':
			concept.visualAttributes = "F" + concept.visualAttributes[1:]
		// modifier leaf
		case 'R':
			concept.visualAttributes = "L" + concept.visualAttributes[1:]
		default:
			return nil, fmt.Errorf("wrong visual attribute: %s", concept.visualAttributes)
		}

		// concatenate basecodes
		concept.baseCode = so.baseCode + "\\" + concept.baseCode

		concept.factTableColumn = "concept_cd"
		concept.tablename = "concept_dimension"
		concept.columnName = "concept_path"
		concept.dimCode = concept.fullname

		// concatenate tooltips
		concept.tooltip = so.tooltip + " \\ " + concept.tooltip

		concept.appliedPath = "@"
		concept.exclusionCD = "\\N"

		concepts = append(concepts, concept)
	}

	return
}

// parseSensitiveMedCoOntology parses and stores sensitive MedCo ontology concepts
func parseSensitiveMedCoOntology(rawName string, concepts []*medCoOntologyRecord, listConceptCD *[]string, allSensitiveConceptIDs *[]int64) {

	for _, so := range concepts {

		so.nodeEncryptID = idConcepts

		if _, ok := tablesMedCoOntology[rawName]; ok {
			tablesMedCoOntology[rawName].sensitive[so.fullname] = so
		} else {
			sensitive := make(map[string]*medCoOntologyRecord)
			clear := make(map[string][]*medCoOntologyRecord)
			sensitive[so.fullname] = so
			tablesMedCoOntology[rawName] = medCoTableInfo{clear: clear, sensitive: sensitive}
		}

		// if the ID does not yet exist
		if _, ok := mapConceptPathToTag[so.fullname]; !ok {
			mapConceptPathToTag[so.fullname] = tagAndID{tag: libunlynx.GroupingKey(-1), tagID: -1}
			*listConceptCD = append(*listConceptCD, so.fullname)
			*allSensitiveConceptIDs = append(*allSensitiveConceptIDs, idConcepts)
		}

		idConcepts++
	}

	return

}

// encryptAndTag encrypts the elements and tags them to allow for future comparison
func encryptAndTag(list []int64, group *onet.Roster, entryPointIdx int) ([]libunlynx.GroupingKey, error) {

	// ENCRYPTION
	start := time.Now()
	listEncryptedElements := make(libunlynx.CipherVector, len(list))

	for i := int64(0); i < int64(len(list)); i++ {
		listEncryptedElements[i] = *libunlynx.EncryptInt(group.Aggregate, list[i])
	}
	log.Lvl2("Finished encrypting the sensitive data... ["+strconv.FormatInt(int64(len(listEncryptedElements)), 10)+"] (", time.Since(start), ")")

	// TAGGING
	start = time.Now()
	client := servicesmedco.NewMedCoClient(group.List[entryPointIdx], strconv.Itoa(entryPointIdx))

	_, result, tr, err := client.SendSurveyDDTRequestTerms(
		group, // Roster
		servicesmedco.SurveyID("tagging_loading_phase"), // SurveyID
		listEncryptedElements,                           // Encrypted query terms to tag
		false,                                           // compute proofs?
		test,
	)

	if err != nil {
		log.Fatal("Error during DDT:", err)
		return nil, err
	}

	totalTime := time.Since(start)

	log.Lvl2("DDT took: execution -", tr.MapTR[servicesmedco.TaggingTimeExec], "communication -", tr.MapTR[servicesmedco.TaggingTimeCommunication])

	log.Lvl2("Finished tagging the sensitive data... (", totalTime, ")")

	return result, nil
}

// convertClearLocalOntologyTable converts the local i2b2 ontology table(s) (i2b2.csv, icd10_icd9.csv, etc.)
func convertClearLocalOntologyTable(rawName string) error {
	// two new files are generated: one to store the non-sensitive data and another to store the sensitive data
	csvClearOutputFile, err := os.Create(outputFilesPathsNonSensitive["LOCAL_"+rawName].Path)
	if err != nil {
		log.Fatal("Error opening [" + strings.ToLower(rawName) + "].csv")
		return err
	}
	defer csvClearOutputFile.Close()

	headerString := ""
	for _, header := range headerLocalOntology {
		headerString += "\"" + header + "\","
	}
	// remove the last ,
	csvClearOutputFile.WriteString(headerString[:len(headerString)-1] + "\n")

	// non-sensitive
	keys := make([]string, 0, len(tableLocalOntologyClear))
	for key := range tableLocalOntologyClear {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		los := tableLocalOntologyClear[key]
		for _, lo := range los {
			csvClearOutputFile.WriteString(lo.toCSVText() + "\n")
		}
	}

	return nil
}

// convertSensitiveLocalOntologyTable generates the sensitive_tagged.csv file
func convertSensitiveLocalOntologyTable() error {
	csvSensitiveOutputFile, err := os.Create(outputFilesPathsSensitive["SENSITIVE_TAGGED"].Path)
	if err != nil {
		log.Fatal("Error opening [sensitive_tagged].csv")
		return err
	}
	defer csvSensitiveOutputFile.Close()

	headerString := ""
	for _, header := range headerLocalOntology {
		headerString += "\"" + header + "\","
	}
	// remove the last ,
	csvSensitiveOutputFile.WriteString(headerString[:len(headerString)-1] + "\n")

	// sensitive concepts
	for _, el := range mapConceptPathToTag {
		csvSensitiveOutputFile.WriteString(localOntologySensitiveConceptToCSVText(&el.tag, el.tagID) + "\n")
	}

	return nil
}
