package loaderi2b2

import (
	"go.dedis.ch/onet/v3/log"
	"os"
	"sort"
	"strings"
)

// modifierDimensionRecord is a record of the modifier_dimension table
type modifierDimensionRecord struct {
	pk           modifierDimensionPK
	modifierCD   string
	nameChar     string
	modifierBlob string
	adminColumns administrativeColumns
}

// modifierDimensionPK is the primary key of the modifier_dimension table
type modifierDimensionPK struct {
	modifierPath string
}

// tableModifierDimension is the modifier_dimension table
var tableModifierDimension map[modifierDimensionPK]*modifierDimensionRecord

// headerModifierDimension contains the header of the modifier_dimension table
var headerModifierDimension []string

// modifierDimensionFromString generates a modifierDimensionRecord struct from a parsed line of a .csv file
func modifierDimensionFromString(line []string) (modifierDimensionPK, *modifierDimensionRecord) {
	mdk := modifierDimensionPK{
		modifierPath: strings.Replace(line[0], "\"", "\"\"", -1),
	}

	md := modifierDimensionRecord{
		pk:           mdk,
		modifierCD:   line[1],
		nameChar:     strings.Replace(line[2], "\"", "\"\"", -1),
		modifierBlob: line[3],
	}

	ac := administrativeColumns{
		updateDate:     line[4],
		downloadDate:   line[5],
		importDate:     line[6],
		sourceSystemCD: line[7],
		uploadID:       line[8],
	}

	md.adminColumns = ac

	return mdk, &md
}

// modifierDimensionToConceptDimension generates a conceptDimensionRecord instance starting from a modifierDimensionRecord instance
func modifierDimensionToConceptDimension(mdk *modifierDimensionRecord) *conceptDimensionRecord {

	cdk := conceptDimensionPK{
		conceptPath: mdk.pk.modifierPath,
	}

	cd := &conceptDimensionRecord{
		pk:          cdk,
		conceptCD:   mdk.modifierCD,
		nameChar:    mdk.nameChar,
		conceptBlob: mdk.modifierBlob,
	}

	cd.adminColumns = mdk.adminColumns

	return cd

}

// toCSVText writes the modifierDimensionRecord object in a way that can be added to a .csv file - "","","", etc.
func (md modifierDimensionRecord) toCSVText() string {
	acString := "\"" + md.adminColumns.updateDate + "\"," + "\"" + md.adminColumns.downloadDate + "\"," + "\"" + md.adminColumns.importDate + "\"," + "\"" + md.adminColumns.sourceSystemCD + "\"," + "\"" + md.adminColumns.uploadID + "\""
	finalString := "\"" + md.pk.modifierPath + "\"," + "\"" + md.modifierCD + "\"," + "\"" + md.nameChar + "\"," + "\"" + md.modifierBlob + "\"," + acString

	return strings.Replace(finalString, `"\N"`, "", -1)
}

// parseModifierDimension reads and parses the modifier_dimension.csv file
func parseModifierDimension() error {
	lines, err := readCSV(inputFilePaths["MODIFIER_DIMENSION"])
	if err != nil {
		log.Fatal("Error in readCSV()")
		return err
	}

	tableModifierDimension = make(map[modifierDimensionPK]*modifierDimensionRecord)
	headerModifierDimension = make([]string, 0)

	/* structure of modifier_dimension.csv (in order):

	// pk
	"modifier_path",

	// MANDATORY FIELDS
	"modifier_cd",
	"name_char",
	"modifier_blob",

	// ADMIN FIELDS
	"update_date",
	"download_date",
	"import_date",
	"sourcesystem_cd",
	"upload_id"
	*/

	for _, header := range lines[0] {
		headerModifierDimension = append(headerModifierDimension, header)
	}

	//skip header
	for _, line := range lines[1:] {
		mdk, md := modifierDimensionFromString(line)
		tableModifierDimension[mdk] = md
	}

	return nil

}

// convertModifierDimension converts the modifier_dimension.csv file
func convertModifierDimension() error {

	csvOutputFileNonSensitive, err := os.Create(outputFilesPathsNonSensitive["MODIFIER_DIMENSION"].Path)
	if err != nil {
		log.Fatal("Error opening [modifier_dimension].csv")
		return err
	}
	defer csvOutputFileNonSensitive.Close()

	headerString := ""
	for _, header := range headerModifierDimension {
		headerString += "\"" + header + "\","
	}
	// remove the last ,
	csvOutputFileNonSensitive.WriteString(headerString[:len(headerString)-1] + "\n")

	mdks := make([]string, 0, len(tableModifierDimension))
	for mdk := range tableModifierDimension {
		mdks = append(mdks, mdk.modifierPath)
	}
	sort.Strings(mdks)

	for _, mdk := range mdks {
		// if the modifier can be applied to non sensitive concepts
		if _, ok := tableLocalOntologyClear[mdk]; ok {
			md := tableModifierDimension[modifierDimensionPK{modifierPath: mdk}]
			csvOutputFileNonSensitive.WriteString(md.toCSVText() + "\n")
		}
	}

	return nil
}
