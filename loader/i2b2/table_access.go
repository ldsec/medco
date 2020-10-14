package loaderi2b2

import (
	"go.dedis.ch/onet/v3/log"
	"os"
	"strings"
)

// tableAccessRecord is a record of the table_access table
type tableAccessRecord struct {
	tableCD          string
	tableName        string
	protectedAccess  string
	hlevel           string
	fullname         string
	name             string
	synonymCD        string
	visualAttributes string
	totalNum         string
	baseCode         string
	metadataXML      string
	factTableColumn  string
	dimTableName     string
	columnName       string
	columnDataType   string
	operator         string
	dimCode          string
	comment          string
	tooltip          string
	entryData        string
	changeData       string
	statusCD         string
	valueType        string
}

// headerTableAccess contains the header of the table_access table
var headerTableAccess []string

// tableTableAccess is the table that contains all data from the table_access table
var tableTableAccess []tableAccessRecord

// tableAccessFromString generates a tableAccessRecord struct from a parsed line of a .csv file
func tableAccessFromString(line []string) tableAccessRecord {
	ta := tableAccessRecord{
		tableCD:          line[0],
		tableName:        line[1],
		protectedAccess:  line[2],
		hlevel:           line[3],
		fullname:         line[4],
		name:             line[5],
		synonymCD:        line[6],
		visualAttributes: line[7],
		totalNum:         line[8],
		baseCode:         line[9],
		metadataXML:      line[10],
		factTableColumn:  line[11],
		dimTableName:     line[12],
		columnName:       line[13],
		columnDataType:   line[14],
		operator:         line[15],
		dimCode:          line[16],
		comment:          line[17],
		tooltip:          line[18],
		entryData:        line[19],
		changeData:       line[20],
		statusCD:         line[21],
		valueType:        line[22],
	}
	return ta
}

// toCSVText writes the tableAccessRecord object in a way that can be added to a .csv file - "","","", etc.
func (ta tableAccessRecord) toCSVText() string {
	finalString := "\"" + ta.tableCD + "\"," + "\"" + ta.tableName + "\"," + "\"" + ta.protectedAccess + "\"," + "\"" + ta.hlevel + "\"," + "\"" + ta.fullname + "\"," + "\"" + ta.name + "\"," +
		"\"" + ta.synonymCD + "\",\"" + ta.visualAttributes + "\"," + "\"" + ta.totalNum + "\"," + "\"" + ta.baseCode + "\"," + "\"" + ta.metadataXML + "\"," + "\"" + ta.factTableColumn + "\"," + "\"" + ta.dimTableName + "\"," +
		"\"" + ta.columnName + "\"," + "\"" + ta.columnDataType + "\"," + "\"" + ta.operator + "\"," + "\"" + ta.dimCode + "\"," + "\"" + ta.comment + "\"," + "\"" + ta.tooltip + "\"," + "\"" + ta.entryData + "\"," + "\"" +
		ta.changeData + "\"," + "\"" + ta.statusCD + "\"," + "\"" + ta.valueType + "\""

	return strings.Replace(finalString, `"\N"`, "", -1)
}

// parseTableAccess reads and parses the table_access.csv file
func parseTableAccess() error {
	lines, err := readCSV(inputFilePaths["TABLE_ACCESS"])
	if err != nil {
		log.Fatal("Error in readCSV()")
		return err
	}

	headerTableAccess = make([]string, 0)
	tableTableAccess = make([]tableAccessRecord, 0)

	/* structure of table_access.csv (in order):

	"c_table_cd",
	"c_table_name",
	"c_protected_access",
	"c_hlevel",
	"c_fullname",
	"c_name",
	"c_synonym_cd",
	"c_visualattributes",
	"c_totalnum",
	"c_basecode",
	"c_metadataxml",
	"c_facttablecolumn",
	"c_dimtablename",
	"c_columnname",
	"c_columndatatype",
	"c_operator",
	"c_dimcode",
	"c_comment",
	"c_tooltip",
	"c_entry_date",
	"c_change_date",
	"c_status_cd",
	"valuetype_cd"

	*/

	for _, header := range lines[0] {
		headerTableAccess = append(headerTableAccess, header)
	}

	//skip header
	for _, line := range lines[1:] {
		tableTableAccess = append(tableTableAccess, tableAccessFromString(line))
	}

	return nil
}

// convertTableAccess converts the table_access.csv file
func convertTableAccess() error {

	// two new files are generated: one to store the non-sensitive data and another to store the sensitive data
	err := convertTableAccessLogic(outputFilesPathsNonSensitive["TABLE_ACCESS"].Path)
	if err != nil {
		return err
	}

	err = convertTableAccessLogic(outputFilesPathsSensitive["TABLE_ACCESS"].Path)

	return err

}

func convertTableAccessLogic(tableAccessPath string) error {

	csvOutputFile, err := os.Create(tableAccessPath)
	if err != nil {
		log.Fatal("Error opening " + tableAccessPath)
		return err
	}
	defer csvOutputFile.Close()

	headerString := ""
	for _, header := range headerTableAccess {
		headerString += "\"" + header + "\","
	}
	// remove the last ,
	csvOutputFile.WriteString(headerString[:len(headerString)-1] + "\n")

	for _, ta := range tableTableAccess {
		csvOutputFile.WriteString(ta.toCSVText() + "\n")
	}

	return nil

}
