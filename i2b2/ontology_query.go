package i2b2

import (
	"errors"
	"github.com/lca1/medco-connector/swagger/models"
	"github.com/lca1/medco-connector/util"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

// GetOntologyChildren makes request to browse the i2b2 ontology
func GetOntologyChildren(path string) (results []*models.SearchResultElement, err error) {

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
			util.I2b2HiveURL + "/OntologyService/getCategories",
			NewOntReqGetCategoriesMessageBody(),
			xmlResponse,
		)
		if err != nil {
			return nil, err
		}

	} else {
		err = i2b2XMLRequest(
			util.I2b2HiveURL + "/OntologyService/getChildren",
			NewOntReqGetChildrenMessageBody(convertPathToI2b2Format(path)),
			xmlResponse,
		)
		if err != nil {
			return nil, err
		}
	}

	// generate result from response
	i2b2Concepts := xmlResponse.MessageBody.(*OntRespConceptsMessageBody).Concepts
	for _, concept := range i2b2Concepts {
		results = append(results, parseI2b2Concept(concept))
	}

	return
}

func parseI2b2Concept(concept Concept) (result *models.SearchResultElement) {
	// todo: add leaf, ensure type OK
	//          type:
	//            type: "string"
	//            enum:
	//              - CONCEPT_PARENT_NODE
	//              - CONCEPT_INTERNAL_NODE
	//              - CONCEPT_LEAF
	true := true
	false := false
	result = &models.SearchResultElement{
		Name: concept.Name,
		DisplayName: concept.Name,
		Code: concept.Basecode,
		MedcoEncryption: &models.SearchResultElementMedcoEncryption{
			Encrypted: &false,
			ID: -1,
			ChildrenIds: []int64{},
		},
		Metadata: nil,
		Path: convertPathFromI2b2Format(concept.Key),
		//Type: models.SearchResultElementTypeConcept,
		//Leaf: false,
	}

	switch concept.Visualattributes[0] {
	// i2b2 leaf
	case 'L':
		result.Leaf = &true
		result.Type = models.SearchResultElementTypeConcept

	// i2b2 container
	case 'C':
		result.Leaf = &false
		result.Type = models.SearchResultElementTypeContainer

	// i2b2 folder (& default)
	default:
		fallthrough
	case 'F':
		result.Leaf = &false
		result.Type = models.SearchResultElementTypeConcept

	}

	splitCode := strings.Split(concept.Basecode, ":")

	// if clinical concept from data loader v0 (from concept code)
	if splitCode[0] == "ENC_ID" {
		result.MedcoEncryption.Encrypted = &true
		result.MedcoEncryption.ID, _ = strconv.ParseInt(splitCode[1], 10, 64)

	// if concept from loader v1 encrypted (from metadata xml)
	} else if concept.Metadataxml.ValueMetadata.EncryptedType != "" {
		result.MedcoEncryption.Encrypted = &true
		result.MedcoEncryption.ID, _ = strconv.ParseInt(concept.Metadataxml.ValueMetadata.NodeEncryptID, 10, 64)

		for _, childEncryptIDString := range strings.Split(concept.Metadataxml.ValueMetadata.ChildrenEncryptIDs, ",") {
			childEncryptID, _ := strconv.ParseInt(childEncryptIDString, 10, 64)
			result.MedcoEncryption.ChildrenIds = append(result.MedcoEncryption.ChildrenIds, childEncryptID)
		}

	// if genomic concept from data loader v0 (from concept code)
	} else if splitCode[0] == "GEN" {
		result.Name = splitCode[1]
		result.Type = models.SearchResultElementTypeGenomicAnnotation
	}

	return
}

func convertPathToI2b2Format(path string) string {
	return `\` + strings.Replace(path, "/", `\`, -1)
}

func convertPathFromI2b2Format(path string) string {
	return strings.Replace(strings.Replace(path, `\`, "/", -1), "//", "/", 1)
}
