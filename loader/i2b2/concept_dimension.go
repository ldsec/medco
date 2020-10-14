package loaderi2b2

import (
	libunlynx "github.com/ldsec/unlynx/lib"
	"go.dedis.ch/onet/v3/log"
	"os"
	"sort"
	"strconv"
	"strings"
)

// conceptDimensionRecord is a record of the concept_dimension table
type conceptDimensionRecord struct {
	pk           conceptDimensionPK
	conceptCD    string
	nameChar     string
	conceptBlob  string
	adminColumns administrativeColumns
}

// conceptDimensionPK is the primary key of the concept_dimension table
type conceptDimensionPK struct {
	conceptPath string
}

var (
	// listConceptsToIgnore lists concepts that appear in the concept_dimension and not in the ontology (which is kind of strange)
	listConceptsToIgnore map[string]struct{}
)

// tableConceptDimension is the concept_dimension table
var tableConceptDimension map[conceptDimensionPK]*conceptDimensionRecord

// headerConceptDimension contains the header of the concept_dimension table
var headerConceptDimension []string

// mapConceptCodeToTag maps the concept code (in the concept dimension) to the tag ID value (for the sensitive terms)
var mapConceptCodeToTag map[string]int64

// conceptDimensionFromString generates a conceptDimensionRecord struct from a parsed line of a .csv file
func conceptDimensionFromString(line []string) (conceptDimensionPK, *conceptDimensionRecord) {
	cdk := conceptDimensionPK{
		conceptPath: strings.Replace(line[0], "\"", "\"\"", -1),
	}

	cd := conceptDimensionRecord{
		pk:          cdk,
		conceptCD:   line[1],
		nameChar:    strings.Replace(line[2], "\"", "\"\"", -1),
		conceptBlob: line[3],
	}

	ac := administrativeColumns{
		updateDate:     line[4],
		downloadDate:   line[5],
		importDate:     line[6],
		sourceSystemCD: line[7],
		uploadID:       line[8],
	}

	cd.adminColumns = ac

	return cdk, &cd
}

// toCSVText writes the conceptDimensionRecord object in a way that can be added to a .csv file - "","","", etc.
func (cd conceptDimensionRecord) toCSVText() string {
	acString := "\"" + cd.adminColumns.updateDate + "\"," + "\"" + cd.adminColumns.downloadDate + "\"," + "\"" + cd.adminColumns.importDate + "\"," + "\"" + cd.adminColumns.sourceSystemCD + "\"," + "\"" + cd.adminColumns.uploadID + "\""
	finalString := "\"" + cd.pk.conceptPath + "\"," + "\"" + cd.conceptCD + "\"," + "\"" + cd.nameChar + "\"," + "\"" + cd.conceptBlob + "\"," + acString

	return strings.Replace(finalString, `"\N"`, "", -1)
}

// conceptDimensionSensitiveToCSVText writes the tagging information of a concept of the concept_dimension table in a way that can be added to a .csv file - "","","", etc.
func conceptDimensionSensitiveToCSVText(tag *libunlynx.GroupingKey, tagID int64) string {
	finalString := `"\medco\tagged\concept\` + string(*tag) + `\","TAG_ID:` + strconv.FormatInt(tagID, 10) + `","\N","\N","\N","\N","NOW()","\N","\N"`

	return strings.Replace(finalString, `"\N"`, "", -1)
}

// parseConceptDimension reads and parses the concept_dimension.csv file
// It has to be called after convertModifierDimension()
func parseConceptDimension() error {
	lines, err := readCSV(inputFilePaths["CONCEPT_DIMENSION"])
	if err != nil {
		log.Fatal("Error in readCSV()")
		return err
	}

	listConceptsToIgnore = make(map[string]struct{})
	tableConceptDimension = make(map[conceptDimensionPK]*conceptDimensionRecord)
	headerConceptDimension = make([]string, 0)
	mapConceptCodeToTag = make(map[string]int64)

	/* structure of concept_dimension.csv (in order):

	// pk
	"concept_path",

	// MANDATORY FIELDS
	"concept_cd",
	"name_char",
	"concept_blob",

	// ADMIN FIELDS
	"update_date",
	"download_date",
	"import_date",
	"sourcesystem_cd",
	"upload_id"
	*/

	for _, header := range lines[0] {
		headerConceptDimension = append(headerConceptDimension, header)
	}

	//skip header
	for _, line := range lines[1:] {
		cdk, cd := conceptDimensionFromString(line)
		tableConceptDimension[cdk] = cd
	}

	if enableModifiers {
		// add to the table the concepts generated from modifiers
		for generatorConceptPath, generatedConcepts := range mapGeneratedConcepts {
			cd, ok := tableConceptDimension[conceptDimensionPK{conceptPath: generatorConceptPath}]
			if ok {
				for _, generatedConcept := range generatedConcepts {
					modifierPath := "\\" + strings.TrimPrefix(generatedConcept.fullname, generatorConceptPath)
					mdk := modifierDimensionPK{modifierPath: modifierPath}
					md, ok := tableModifierDimension[mdk]
					if ok {
						newConcept := modifierDimensionToConceptDimension(md)
						newConcept.pk.conceptPath = strings.TrimSuffix(cd.pk.conceptPath, "\\") + newConcept.pk.conceptPath
						newConcept.conceptCD = cd.conceptCD + "\\" + newConcept.conceptCD
						newConcept.nameChar = cd.nameChar + "\\" + newConcept.nameChar
						tableConceptDimension[conceptDimensionPK{conceptPath: generatedConcept.fullname}] = newConcept
					}
				}
			}
		}
	}

	return nil

}

// convertConceptDimension converts the concept_dimension.csv file
func convertConceptDimension() error {

	csvOutputFileNonSensitive, err := os.Create(outputFilesPathsNonSensitive["CONCEPT_DIMENSION"].Path)
	if err != nil {
		log.Fatal("Error opening [concept_dimension].csv")
		return err
	}
	csvOutputFileSensitive, err := os.Create(outputFilesPathsSensitive["CONCEPT_DIMENSION"].Path)
	if err != nil {
		log.Fatal("Error opening [concept_dimension].csv")
		return err
	}
	defer csvOutputFileNonSensitive.Close()
	defer csvOutputFileSensitive.Close()

	headerString := ""
	for _, header := range headerConceptDimension {
		headerString += "\"" + header + "\","
	}
	// remove the last ,
	csvOutputFileNonSensitive.WriteString(headerString[:len(headerString)-1] + "\n")
	csvOutputFileSensitive.WriteString(headerString[:len(headerString)-1] + "\n")

	cdks := make([]string, 0, len(tableConceptDimension))
	for cdk := range tableConceptDimension {
		cdks = append(cdks, cdk.conceptPath)
	}
	sort.Strings(cdks)

	for _, cdk := range cdks {
		cd := tableConceptDimension[conceptDimensionPK{conceptPath: cdk}]
		// if the concept is non-sensitive -> keep it as it is
		if _, ok := tableLocalOntologyClear[cd.pk.conceptPath]; ok {
			csvOutputFileNonSensitive.WriteString(cd.toCSVText() + "\n")
			// if the concept is sensitive -> fetch its encrypted tag and tag_id
		} else if _, ok := mapConceptPathToTag[cd.pk.conceptPath]; ok {
			temp := mapConceptPathToTag[cd.pk.conceptPath].tag
			csvOutputFileSensitive.WriteString(conceptDimensionSensitiveToCSVText(&temp, mapConceptPathToTag[cd.pk.conceptPath].tagID) + "\n")
			mapConceptCodeToTag[cd.conceptCD] = mapConceptPathToTag[cd.pk.conceptPath].tagID
			// if the concept does not exist in the LocalOntology and none of his siblings is sensitive
		} else if _, ok := hasSensitiveParents(cd.pk.conceptPath); !ok {
			csvOutputFileNonSensitive.WriteString(cd.toCSVText() + "\n")
		} else {
			listConceptsToIgnore[cd.conceptCD] = struct{}{}
		}
	}

	return nil
}
