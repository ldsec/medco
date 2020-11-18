package i2b2

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ldsec/medco/connector/restapi/models"
	utilserver "github.com/ldsec/medco/connector/util/server"
	"github.com/sirupsen/logrus"
)

// GetOntologyTermInfo makes request to get infromation about a node given its path
func GetOntologyTermInfo(path string) (results []*models.ExploreSearchResultElement, err error) {
	path = strings.TrimSpace(path)
	xmlResponse := &Response{
		MessageBody: &OntRespConceptsMessageBody{},
	}
	if len(path) == 0 {
		err = fmt.Errorf("empty path")
		logrus.Error(err)
		return
	}

	err = i2b2XMLRequest(
		utilserver.I2b2HiveURL+"/OntologyService/getTermInfo",
		NewOntReqGetTermInfoMessageBody(convertPathToI2b2Format(path)),
		xmlResponse,
	)
	if err != nil {
		return nil, err
	}
	i2b2TermInfo := xmlResponse.MessageBody.(*OntRespConceptsMessageBody).Concepts
	results = make([]*models.ExploreSearchResultElement, 0)
	for _, concept := range i2b2TermInfo {
		parsedConcept, err := parseI2b2Concept(concept)
		if err != nil {
			return nil, err
		}
		results = append(results, parsedConcept)
	}

	return

}

// GetOntologyChildren makes request to browse the i2b2 ontology
func GetOntologyChildren(path string) (results []*models.ExploreSearchResultElement, err error) {

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
		Metadata:    nil,
		Path:        convertPathFromI2b2Format(concept.Key),
		AppliedPath: "@",
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
	} else if concept.Metadataxml.ValueMetadata.EncryptedType != "" {
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

func convertPathToI2b2Format(path string) string {
	return `\` + strings.Replace(path, "/", `\`, -1)
}

func convertPathFromI2b2Format(path string) string {
	return strings.Replace(strings.Replace(path, `\`, "/", -1), "//", "/", 1)
}
