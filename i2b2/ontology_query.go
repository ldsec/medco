package i2b2

import (
	"errors"
	"github.com/lca1/medco-connector/swagger/models"
	"github.com/lca1/medco-connector/util"
	"log"
	"strconv"
	"strings"
)

// make request to browse the i2b2 ontology
func GetOntologyChildren(path string) (results []*models.SearchResultElement, err error) {

	// craft and make request
	path = strings.TrimSpace(path)

	xmlResponse := &Response{
		MessageBody: &OntRespConceptsMessageBody{},
	}

	if len(path) == 0 {
		err = errors.New("empty path")
		log.Print(err)
		return

	} else if path == "/" {
		err = i2b2XMLRequest(
			util.I2b2ONTCellURL() + "/getCategories",
			NewOntReqGetCategoriesRequest(),
			xmlResponse,
		)
		if err != nil {
			return nil, err
		}

	} else {
		err = i2b2XMLRequest(
			util.I2b2ONTCellURL() + "/getChildren",
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

// parse the i2b2 concept into the result model
func parseI2b2Concept(concept Concept) (result *models.SearchResultElement) {

	result = &models.SearchResultElement{
		Name: concept.Name,
		DisplayName: concept.Name,
		Code: concept.Basecode,
		MedcoEncryption: &models.SearchResultElementMedcoEncryption{
			Encrypted: false,
			ID: -1,
			ChildrenIds: []int64{},
			Type: "",
		},
		Metadata: nil,
		Path: convertPathFromI2b2Format(concept.Key),
		ValueType: models.SearchResultElementValueTypeNone,
	}

	splitCode := strings.Split(concept.Basecode, ":")

	// if clinical concept from data loader v0 (from concept code)
	if splitCode[0] == "ENC_ID" {
		result.MedcoEncryption.Encrypted = true
		result.MedcoEncryption.Type = models.SearchResultElementMedcoEncryptionTypeCONCEPTLEAF
		result.MedcoEncryption.ID, _ = strconv.ParseInt(splitCode[1], 10, 64)

	// if concept from loader v1 encrypted (from metadata xml)
	} else if concept.Metadataxml.ValueMetadata.EncryptedType != "" {
		result.MedcoEncryption.Encrypted = true
		result.MedcoEncryption.Type = concept.Metadataxml.ValueMetadata.EncryptedType
		result.MedcoEncryption.ID, _ = strconv.ParseInt(concept.Metadataxml.ValueMetadata.NodeEncryptID, 10, 64)

		for _, childEncryptIdString := range strings.Split(concept.Metadataxml.ValueMetadata.ChildrenEncryptIDs, ",") {
			childEncryptId, _ := strconv.ParseInt(childEncryptIdString, 10, 64)
			result.MedcoEncryption.ChildrenIds = append(result.MedcoEncryption.ChildrenIds, childEncryptId)
		}

		// if genomic concept from data loader v0 (from concept code)
	} else if splitCode[0] == "GEN" {
		result.ValueType = models.SearchResultElementValueTypeGenomicAnnotation
	}

	return
}

func convertPathToI2b2Format(path string) string {
	return "\\" + strings.Replace(path, "/", "\\", -1)
}

func convertPathFromI2b2Format(path string) string {
	return strings.Replace(strings.Replace(path, "\\", "/", -1), "//", "/", 1)
}
