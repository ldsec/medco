package i2b2

import (
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ldsec/medco/connector/restapi/models"
	utilserver "github.com/ldsec/medco/connector/util/server"
	"github.com/sirupsen/logrus"
)

// GetOntologyElements retrieves the info about the ontology elements.
func GetOntologyElements(path string, limit int64) (results []*models.ExploreSearchResultElement, err error) {

	results = make([]*models.ExploreSearchResultElement, 0)

	// transform to i2b2 format
	path = strings.Replace(strings.TrimSpace(path), "/", `\\`, -1)

	if len(path) == 0 {
		err = fmt.Errorf("empty path")
		logrus.Error(err)
		return
	}

	row, err := utilserver.I2B2DBConnection.Query("SELECT * FROM medco_ont.get_ontology_elements($1,$2) ORDER BY id, fullpath DESC", path, limit)
	if err != nil {
		return nil, fmt.Errorf("while calling i2b2 database for retrieving ontology elements: %v", err)
	}

	var fullName, name, visualAttributes, baseCode, metaDataXML, comment, appliedPath, fullPath sql.NullString
	var id int
	currentID := 0
	var currentElement *models.ExploreSearchResultElement

	for row.Next() {
		err = row.Scan(&fullName, &name, &visualAttributes, &baseCode, &metaDataXML, &comment, &appliedPath, &id, &fullPath)
		if err != nil {
			return nil, fmt.Errorf("while reading database record stream for retrieving ontology elements: %v", err)
		}

		var metadataXML Metadataxml
		var metadataXMLtoJSON *models.Metadataxml
		if metaDataXML.Valid {

			err = xml.Unmarshal([]byte(metaDataXML.String), &metadataXML.ValueMetadata)
			if err != nil {
				return nil, fmt.Errorf("while unmarshalling xml metadata of ontology element %s: %v", fullName.String, err)
			}

			var unitValues []*models.UnitValues
			for _, unitValue := range metadataXML.ValueMetadata.UnitValues {

				var unitValuesConvertingUnits []*models.UnitValuesConvertingUnitsItems0
				for _, unitValuesConvertingUnit := range unitValue.ConvertingUnits {
					unitValuesConvertingUnits = append(unitValuesConvertingUnits, &models.UnitValuesConvertingUnitsItems0{
						MultiplyingFactor: unitValuesConvertingUnit.MultiplyingFactor,
						Units:             unitValuesConvertingUnit.Units,
					})
				}

				unitValues = append(unitValues, &models.UnitValues{
					ConvertingUnits: unitValuesConvertingUnits,
					EqualUnits:      unitValue.EqualUnits,
					ExcludingUnits:  unitValue.ExcludingUnits,
					NormalUnits:     unitValue.NormalUnits,
				})
			}

			metadataXMLtoJSON = &models.Metadataxml{
				ValueMetadata: &models.MetadataxmlValueMetadata{
					ChildrenEncryptIDs: metadataXML.ValueMetadata.ChildrenEncryptIDs,
					CreationDateTime:   metadataXML.ValueMetadata.CreationDateTime,
					DataType:           metadataXML.ValueMetadata.DataType,
					EncryptedType:      metadataXML.ValueMetadata.EncryptedType,
					EnumValues:         metadataXML.ValueMetadata.EnumValues,
					Flagstouse:         metadataXML.ValueMetadata.Flagstouse,
					NodeEncryptID:      metadataXML.ValueMetadata.NodeEncryptID,
					Oktousevalues:      metadataXML.ValueMetadata.Oktousevalues,
					TestID:             metadataXML.ValueMetadata.TestID,
					TestName:           metadataXML.ValueMetadata.TestName,
					UnitValues:         unitValues,
					Version:            metadataXML.ValueMetadata.Version,
				},
			}

		}

		kind, leaf := parseVisualAttributes(visualAttributes.String)
		ontologyElement := &models.ExploreSearchResultElement{
			Path:        convertPathFromI2b2Format(fullName.String),
			AppliedPath: appliedPath.String,
			Type:        kind,
			Name:        name.String,
			DisplayName: name.String,
			Code:        baseCode.String,
			Metadata:    metadataXMLtoJSON,
			Comment:     comment.String,
			MedcoEncryption: &models.ExploreSearchResultElementMedcoEncryption{
				Encrypted:   func() *bool { b := false; return &b }(),
				ChildrenIds: []int64{},
			},
			Leaf: leaf,
		}

		if ontologyElement.AppliedPath != "@" {
			ontologyElement.AppliedPath = convertPathFromI2b2Format(ontologyElement.AppliedPath)
		}

		// a found element and its ancestors have the same id
		// the result of the query is ordered, i.e. element at position i have its father at position i+1
		if id != currentID {
			currentID = id
			results = append(results, ontologyElement)
		} else {
			currentElement.Parent = ontologyElement
		}

		currentElement = ontologyElement
	}

	return
}

// GetOntologyConceptChildren retrieves the children of the concept identified by path.
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

// GetOntologyModifiers retrieves the modifiers that apply to the concept identified by path.
func GetOntologyModifiers(path string) (results []*models.ExploreSearchResultElement, err error) {

	// craft and make request
	path = convertPathToI2b2Format(strings.TrimSpace(path))

	xmlResponse := &Response{
		MessageBody: &OntRespModifiersMessageBody{},
	}

	err = i2b2XMLRequest(
		utilserver.I2b2HiveURL+"/OntologyService/getModifiers",
		NewOntReqGetModifiersMessageBody(path),
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

// GetOntologyModifierChildren retrieves the children of the modifier identified by path and appliedPath, and that apply to appliedConcept.
func GetOntologyModifierChildren(path, appliedPath, appliedConcept string) (results []*models.ExploreSearchResultElement, err error) {

	// craft and make request
	path = convertPathToI2b2Format(strings.TrimSpace(path))
	appliedPath = convertPathToI2b2Format(strings.TrimSpace(appliedPath))[1:]
	appliedConcept = convertPathToI2b2Format(strings.TrimSpace(appliedConcept))

	xmlResponse := &Response{
		MessageBody: &OntRespModifiersMessageBody{},
	}

	err = i2b2XMLRequest(
		utilserver.I2b2HiveURL+"/OntologyService/getModifierChildren",
		NewOntReqGetModifierChildrenMessageBody(path, appliedPath, appliedConcept),
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

// GetOntologyConceptInfo retrieves the info of the concept identified by path.
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

// GetOntologyModifierInfo retrieves the info of the modifier identified by path and appliedPath.
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
		Metadata:    concept.Metadataxml,
		Path:        convertPathFromI2b2Format(concept.Key),
		AppliedPath: "@",
		Comment:     concept.Comment,
	}

	result.Type, result.Leaf = parseVisualAttributes(concept.Visualattributes)

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

func parseVisualAttributes(visualAttributes string) (kind string, leaf *bool) {

	false := false
	true := true

	switch visualAttributes[0] {
	// i2b2 leaf
	case 'L':
		leaf = &true
		kind = models.ExploreSearchResultElementTypeConcept
	case 'R':
		leaf = &true
		kind = models.ExploreSearchResultElementTypeModifier

	// i2b2 container
	case 'C':
		leaf = &false
		kind = models.ExploreSearchResultElementTypeConceptContainer
	case 'O':
		leaf = &false
		kind = models.ExploreSearchResultElementTypeModifierContainer

	// i2b2 folder (& default)
	default:
		fallthrough
	case 'F':
		leaf = &false
		kind = models.ExploreSearchResultElementTypeConceptFolder
	case 'D':
		leaf = &false
		kind = models.ExploreSearchResultElementTypeModifierFolder
	}

	return
}

func convertPathToI2b2Format(path string) string {
	return `\` + strings.Replace(path, "/", `\`, -1)
}

func convertPathFromI2b2Format(path string) string {
	return strings.Replace(strings.Replace(path, `\`, "/", -1), "//", "/", 1)
}
