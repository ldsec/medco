package i2b2

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ldsec/medco/connector/restapi/models"
	utilserver "github.com/ldsec/medco/connector/util/server"
	"github.com/sirupsen/logrus"
)

// GetOntologyConceptChildren makes request to browse the i2b2 ontology
func GetOntologyConceptChildren(path string) (results []*models.ExploreSearchResultElement, err error) {

	// craft and make request
	path = strings.TrimSpace(path)

	xmlResponse := &Response{
		MessageBody: &OntRespConceptsMessageBody{},
	}

	if len(path) == 0 {
		err = errors.New("empty path")
		logrus.Error(err)
		return

	} else if path == "/" {
		err = i2b2XMLRequest(
			utilserver.I2b2HiveURL+"/OntologyService/getCategories",
			NewOntReqGetCategoriesMessageBody(),
			xmlResponse,
		)
		if err != nil {
			return nil, err
		}

	} else {
		err = i2b2XMLRequest(
			utilserver.I2b2HiveURL+"/OntologyService/getChildren",
			NewOntReqGetChildrenMessageBody(convertPathToI2b2Format(path)),
			xmlResponse,
		)
		if err != nil {
			return nil, err
		}
	}

	// generate result from response
	i2b2Concepts := xmlResponse.MessageBody.(*OntRespConceptsMessageBody).Concepts
	results = make([]*models.ExploreSearchResultElement, 0, len(i2b2Concepts))
	for _, concept := range i2b2Concepts {
		parsedConcept, err := parseI2b2Concept(concept)
		if err != nil {
			return nil, err
		}
		results = append(results, parsedConcept)
	}

	return
}

func processModifierXMLResponse(xmlResponse *Response) (results []*models.ExploreSearchResultElement, err error) {
	// generate result from response
	i2b2Modifiers := xmlResponse.MessageBody.(*OntRespModifiersMessageBody).Modifiers
	results = make([]*models.ExploreSearchResultElement, 0, len(i2b2Modifiers))
	for _, modifier := range i2b2Modifiers {
		parsedModifier, err := parseI2b2Modifier(modifier)
		if err != nil {
			return nil, err
		}

		results = append(results, parsedModifier)
	}
	err = populateModifierTotalnum(results, i2b2Modifiers)

	if err != nil {
		logrus.Error("Cannot populate modifier's totalnum ", err)
		return
	}

	return
}

// GetOntologyModifiers retrieves the modifiers that apply to self
func GetOntologyModifiers(self string) (results []*models.ExploreSearchResultElement, err error) {

	// craft and make request
	self = convertPathToI2b2Format(strings.TrimSpace(self))

	xmlResponse := &Response{
		MessageBody: &OntRespModifiersMessageBody{},
	}

	err = i2b2XMLRequest(
		utilserver.I2b2HiveURL+"/OntologyService/getModifiers",
		NewOntReqGetModifiersMessageBody(self),
		xmlResponse,
	)
	if err != nil {
		return nil, err
	}

	return processModifierXMLResponse(xmlResponse)
}

// GetOntologyModifierChildren retrieves the children of the parent modifier which have a certain appliedPath and which apply to appliedConcept
func GetOntologyModifierChildren(parent, appliedPath, appliedConcept string) (results []*models.ExploreSearchResultElement, err error) {

	// craft and make request
	parent = convertPathToI2b2Format(strings.TrimSpace(parent))
	appliedPath = convertPathToI2b2Format(strings.TrimSpace(appliedPath))[1:]
	appliedConcept = convertPathToI2b2Format(strings.TrimSpace(appliedConcept))

	xmlResponse := &Response{
		MessageBody: &OntRespModifiersMessageBody{},
	}

	err = i2b2XMLRequest(
		utilserver.I2b2HiveURL+"/OntologyService/getModifierChildren",
		NewOntReqGetModifierChildrenMessageBody(parent, appliedPath, appliedConcept),
		xmlResponse,
	)
	if err != nil {
		return nil, err
	}

	// generate result from response
	return processModifierXMLResponse(xmlResponse)
}

const getSchemasWithATableAccessTable = "SELECT DISTINCT table_schema FROM information_schema.tables WHERE table_name='table_access'"

func createSelectTableContainingModifierQuery(schemaName string) string {
	return "SELECT DISTINCT c_table_name FROM " + schemaName + ".table_access WHERE $1 LIKE c_fullname || '%' ESCAPE '|'"
}

func closeQueryResult(conn *sql.Rows) {
	logrus.Debug("about to close conn ", conn)
	if conn == nil {
		return
	}
	closingError := conn.Close()
	if closingError != nil {
		logrus.Debug("Cannot close empty pointer")
	}
}

//This function is responsible for putting the totalnum values for the modifiers passed passed as parameter.
// By default I2B2's API does not return modifiers' "totalnum" attribute. To solve the missing patients count issue, when the back-end detects an element is a modifier,
// additional code fetches (using a mix of Golang and PostgreSQL) the patients count of that modifier.
// This additional logic is present in the function called "populateModifierTotalnum"
//  In this routine we do no make assumptions on which schema and table contains the children modifiers information.
//  Here's how we find the subjects counts associated to those modifiers
// 	_ Select all schemas in the i2b2 database which contain a "table_access" table.
// 	_ For each schema, find in the "schema.table_access" table the ontology table containing the modifiers.
// 	_ When we find the table and schema containing the modifiers, retrieve their subjects counts in the appropriate ontology table.
// Then we associate subject count with correct modifier.
func populateModifierTotalnum(modifiers []*models.ExploreSearchResultElement, unparsedModifiers []Modifier) (err error) {
	if len(modifiers) == 0 {
		return
	}

	//select all schemas in the i2b2 database which contain the table_access table.
	schemasNamesRows, err := utilserver.I2B2DBConnection.Query(getSchemasWithATableAccessTable)
	defer closeSQLConn(schemasNamesRows)

	if err != nil {
		return fmt.Errorf("Unable to fetch schemas containing a table called `table access`: %v", err)
	}

	headModifier := modifiers[0]

	var matchingTable = ""
	var matchingSchema = ""
	for schemasNamesRows.Next() {

		err = schemasNamesRows.Scan(&matchingSchema)
		if err != nil {
			return fmt.Errorf("Error while scanning the next matching schema containing a table access table: %v", err)
		}
		/*look in schema.table_access for the table name in schema that contains the modifier in its ontology*/
		i2b2ModifierParentPath := slashesToBackSlashes(headModifier.AppliedPath)

		var matchingTables *sql.Rows
		var selectTableContainingModifier string = createSelectTableContainingModifierQuery(matchingSchema)
		logrus.Debug("The select table query ", selectTableContainingModifier)
		logrus.Debug("The full name of the searched modifier ", i2b2ModifierParentPath)
		matchingTables, err = utilserver.I2B2DBConnection.Query(selectTableContainingModifier, i2b2ModifierParentPath)

		defer closeSQLConn(matchingTables)

		if err != nil {
			return fmt.Errorf("Error while executing the query fetching the table name containing the modifier: %v", err)
		}

		for matchingTables.Next() {
			matchingTables.Scan(&matchingTable)

			if matchingTable != "" {
				//we found a match for a table which contains the headModifier
				break
			}
		}

		if matchingSchema != "" && matchingTable != "" {
			//we found a match for a schema and a table that contain headModifier
			break
		}

	}

	if matchingSchema == "" || matchingTable == "" {
		return fmt.Errorf("Did not find a schema, table pair containing the modifier")
	}

	//preparing the query fetching the c_totalnum, fullname pairs

	var queryParams []interface{} = make([]interface{}, 0)

	var getTotalnumsQuery string = "SELECT c_fullname, c_totalnum FROM " + matchingSchema + "." + matchingTable + " WHERE"
	var debugStr string = ""
	for i := range modifiers {
		fullname := unparsedModifiers[i].Fullname
		queryParams = append(queryParams, fullname)

		debugStr += fullname + ", "
		//adding all conditions in the where part getTotalnumsQuery in order to create a prepared statement: OR c_totalnum = $1 OR c_totalnum = $2 ...
		argumentIndex := strconv.Itoa(i + 1)
		if i != 0 {
			getTotalnumsQuery += " OR "
		}
		getTotalnumsQuery += " c_fullname = $" + argumentIndex
	}

	logrus.Debug("Query parameters for modifiers totalnum ", debugStr)
	logrus.Debug("The query for getting fullname, totalnum ", getTotalnumsQuery)

	modifiersRows, err := utilserver.I2B2DBConnection.Query(getTotalnumsQuery, queryParams...)
	defer closeSQLConn(modifiersRows)

	if err != nil {
		return fmt.Errorf("Error while executing query to fetch totalnums: %v ", err)
	}

	for modifiersRows.Next() {
		var modifierFullname string
		var modifierTotalnum int
		modifiersRows.Scan(&modifierFullname, &modifierTotalnum)
		logrus.Debug("Query totalnum result: Found name ", modifierFullname, " found totalnum ", modifierTotalnum)

		/* modifier.Path contains the value returned by i2b2 Key attribute of a modifier.
		 * normally what is in c_fullname is a suffix of this Key attribute.
		 * This is what we check to find a matching modifier to the current sql row in the modifiers list
		 */
		for i, rawModifier := range unparsedModifiers {
			if rawModifier.Fullname == modifierFullname {
				//assigning the SQL query row totalnum result to the matching modifier.
				modifiers[i].SubjectCount = strconv.Itoa(modifierTotalnum)
				break
			}
			logrus.Debug(rawModifier.Fullname, " doesnt match ", modifierFullname)
		}
	}

	return
}

// GetOntologyConceptInfo makes request to get information about a node given its path
func GetOntologyConceptInfo(path string) (results []*models.ExploreSearchResultElement, err error) {

	if path == "/" {
		return
	}

	// craft and make request
	path = convertPathToI2b2Format(strings.TrimSpace(path))

	xmlResponse := &Response{
		MessageBody: &OntRespConceptsMessageBody{},
	}

	err = i2b2XMLRequest(
		utilserver.I2b2HiveURL+"/OntologyService/getTermInfo",
		NewOntReqGetTermInfoMessageBody(path),
		xmlResponse,
	)

	if err != nil {
		return nil, err
	}

	// generate result from response
	i2b2ConceptsInfo := xmlResponse.MessageBody.(*OntRespConceptsMessageBody).Concepts
	results = make([]*models.ExploreSearchResultElement, 0)
	for _, concept := range i2b2ConceptsInfo {
		parsedConcept, err := parseI2b2Concept(concept)
		if err != nil {
			return nil, err
		}
		results = append(results, parsedConcept)
	}

	return
}

// GetOntologyModifierInfo retrieves the info of the modifier identified by path and appliedPath
func GetOntologyModifierInfo(path, appliedPath string) (results []*models.ExploreSearchResultElement, err error) {

	// craft and make request
	path = convertPathToI2b2Format(strings.TrimSpace(path))
	appliedPath = convertPathToI2b2Format(strings.TrimSpace(appliedPath))[1:]

	xmlResponse := &Response{
		MessageBody: &OntRespModifiersMessageBody{},
	}

	err = i2b2XMLRequest(
		utilserver.I2b2HiveURL+"/OntologyService/getModifierInfo",
		NewOntReqGetModifierInfoMessageBody(path, appliedPath),
		xmlResponse,
	)
	if err != nil {
		return nil, err
	}

	// generate result from response
	i2b2Modifiers := xmlResponse.MessageBody.(*OntRespModifiersMessageBody).Modifiers
	results = make([]*models.ExploreSearchResultElement, 0, len(i2b2Modifiers))
	for _, modifier := range i2b2Modifiers {
		parsedModifier, err := parseI2b2Modifier(modifier)
		if err != nil {
			return nil, err
		}
		results = append(results, parsedModifier)
	}

	return
}

func parseI2b2Concept(concept Concept) (result *models.ExploreSearchResultElement, err error) {
	// todo: add leaf, ensure type OK
	//          type:
	//            type: "string"
	//            enum:
	//              - CONCEPT_PARENT_NODE
	//              - CONCEPT_INTERNAL_NODE
	//              - CONCEPT_LEAF
	true := true
	false := false

	result = &models.ExploreSearchResultElement{
		Name:        concept.Name,
		DisplayName: concept.Name,
		Code:        concept.Basecode,
		MedcoEncryption: &models.ExploreSearchResultElementMedcoEncryption{
			Encrypted:   &false,
			ChildrenIds: []int64{},
		},
		Metadata:     concept.Metadataxml,
		Path:         convertPathFromI2b2Format(concept.Key),
		AppliedPath:  "@",
		SubjectCount: concept.Totalnum,
		//Type: models.SearchResultElementTypeConcept,
		//Leaf: false,
	}

	switch concept.Visualattributes[0] {
	// i2b2 leaf
	case 'L':
		result.Leaf = &true
		result.Type = models.ExploreSearchResultElementTypeConcept
	case 'R':
		result.Leaf = &true
		result.Type = models.ExploreSearchResultElementTypeModifier

	// i2b2 container
	case 'C':
		result.Leaf = &false
		result.Type = models.ExploreSearchResultElementTypeConceptContainer
	case 'O':
		result.Leaf = &false
		result.Type = models.ExploreSearchResultElementTypeModifierContainer

	// i2b2 folder (& default)
	default:
		fallthrough
	case 'F':
		result.Leaf = &false
		result.Type = models.ExploreSearchResultElementTypeConceptFolder
	case 'D':
		result.Leaf = &false
		result.Type = models.ExploreSearchResultElementTypeModifierFolder
	}

	splitCode := strings.Split(concept.Basecode, ":")

	// if clinical concept from data loader v0 (from concept code)
	if splitCode[0] == "ENC_ID" {
		result.MedcoEncryption.Encrypted = &true

		if parsedCode, parseErr := strconv.ParseInt(splitCode[1], 10, 64); parseErr != nil {
			logrus.Error("Malformed concept could not be parsed: ", concept.Basecode, "error: ", parseErr)
			return nil, parseErr
		} else if len(splitCode) != 2 {
			err = errors.New("Malformed concept: " + concept.Basecode)
			logrus.Error(err)
			return
		} else {
			result.MedcoEncryption.ID = &parsedCode
		}

		// if concept from loader v1 encrypted (from metadata xml)
	} else if concept.Metadataxml != nil && concept.Metadataxml.ValueMetadata.EncryptedType != "" {
		result.MedcoEncryption.Encrypted = &true

		// node ID
		if concept.Metadataxml.ValueMetadata.NodeEncryptID != "" {
			parsedNodeID, parseErr := strconv.ParseInt(concept.Metadataxml.ValueMetadata.NodeEncryptID, 10, 64)
			if parseErr != nil {
				logrus.Error("Malformed ID could not be parsed: ", concept.Metadataxml.ValueMetadata.NodeEncryptID, "error: ", parseErr)
				// if error: not interrupting
			}
			result.MedcoEncryption.ID = &parsedNodeID
		}

		// children nodes IDs
		if concept.Metadataxml.ValueMetadata.ChildrenEncryptIDs != "" {
			for _, childEncryptIDString := range strings.Split(concept.Metadataxml.ValueMetadata.ChildrenEncryptIDs, ";") {

				childEncryptID, parseErr := strconv.ParseInt(childEncryptIDString, 10, 64)
				if parseErr != nil {
					logrus.Error("Malformed ID could not be parsed: ", childEncryptIDString, "error: ", parseErr)
					// if error: not interrupting
				}
				result.MedcoEncryption.ChildrenIds = append(result.MedcoEncryption.ChildrenIds, childEncryptID)
			}
		}

		// if genomic concept from data loader v0 (from concept code)
	} else if splitCode[0] == "GEN" {
		result.Type = models.ExploreSearchResultElementTypeGenomicAnnotation

		if len(splitCode) != 2 {
			err = errors.New("Malformed concept: " + concept.Basecode)
			logrus.Error(err)
			return
		}
		result.Name = splitCode[1]
	}

	return
}

func parseI2b2Modifier(modifier Modifier) (result *models.ExploreSearchResultElement, err error) {

	result, err = parseI2b2Concept(modifier.Concept)
	if err != nil {
		return
	}
	result.AppliedPath = convertPathFromI2b2Format(modifier.AppliedPath)
	return

}

func slashesToBackSlashes(path string) string {
	return strings.Replace(path, "/", `\`, -1)
}

func convertPathToI2b2Format(path string) string {
	return `\` + slashesToBackSlashes(path)
}

func convertPathFromI2b2Format(path string) string {
	return strings.Replace(strings.Replace(path, `\`, "/", -1), "//", "/", 1)
}
