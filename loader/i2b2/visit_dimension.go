package loaderi2b2

import (
	"go.dedis.ch/onet/v3/log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// visitDimension table represents a visit in the database
type visitDimension struct {
	pk             visitDimensionPK
	activeStatusCD string
	startDate      string
	endDate        string
	optionalFields []optionalFields
	adminColumns   administrativeColumns
}

// visitDimensionPK is the primary key of the visit_dimension table
type visitDimensionPK struct {
	encounterNum string
	patientNum   string
}

// tableVisitDimension is the visit_dimension table
var tableVisitDimension map[visitDimensionPK]*visitDimension

// headerVisitDimension contains the header of the visit_dimension table
var headerVisitDimension []string

var (
	// mapNewEncounterNum maps [old_patient_num old_encounter_num] -> [new_patient_num new_encounter_num].
	// For the dummies the [old_patient_num old_encounter_num] refers to the original values
	mapNewEncounterNum map[visitDimensionPK]visitDimensionPK
	// mapPatientVisits maps a patient_num to all its encounter_nums
	mapPatientVisits map[string][]string
	// maxVisits keeps track of the maximum number of visits of all the patients
	maxVisits int
)

// visitDimensionFromString generates a visitDimension struct from a parsed line of a .csv file
func visitDimensionFromString(line []string) (visitDimensionPK, *visitDimension) {
	vdk := visitDimensionPK{
		encounterNum: line[0],
		patientNum:   line[1],
	}

	vd := visitDimension{
		pk:             vdk,
		activeStatusCD: line[2],
		startDate:      line[3],
		endDate:        line[4],
	}

	size := len(line)

	// optional fields
	of := make([]optionalFields, 0)

	for i := 5; i < size-5; i++ {
		of = append(of, optionalFields{valType: headerPatientDimension[i], value: line[i]})
	}

	ac := administrativeColumns{
		updateDate:     line[size-5],
		downloadDate:   line[size-4],
		importDate:     line[size-3],
		sourceSystemCD: line[size-2],
		uploadID:       line[size-1],
	}

	vd.optionalFields = of
	vd.adminColumns = ac

	return vdk, &vd
}

// toCSVText writes the visitDimensionPK struct in a way that can be added to a .csv file - "","","", etc.
func (vdk visitDimensionPK) toCSVText() string {
	return "\"" + vdk.encounterNum + "\"," + "\"" + vdk.patientNum + "\""
}

// toCSVText writes the visitDimension struct in a way that can be added to a .csv file - "","","", etc.
func (vd visitDimension) toCSVText(empty bool) string {
	of := vd.optionalFields
	ofString := ""
	if empty == false {
		for i := 0; i < len(of); i++ {
			// +4 because there is two pk field and 3 mandatory fields
			ofString += "\"" + of[i].value + "\","
		}

		acString := "\"" + vd.adminColumns.updateDate + "\"," + "\"" + vd.adminColumns.downloadDate + "\"," + "\"" + vd.adminColumns.importDate + "\"," + "\"" + vd.adminColumns.sourceSystemCD + "\"," + "\"" + vd.adminColumns.uploadID + "\""
		finalString := vd.pk.toCSVText() + ",\"" + vd.activeStatusCD + "\"," + "\"" + vd.startDate + "\"," + "\"" + vd.endDate + "\"," + ofString[:len(ofString)-1] + "," + acString

		return strings.Replace(finalString, `"\N"`, "", -1)
	}

	for i := 0; i < len(of); i++ {
		// +4 because there is on pk field and 3 mandatory fields
		ofString += ","
	}

	acString := "," + "," + "," + ","
	finalString := vd.pk.toCSVText() + "," + "," + "," + "," + ofString[:len(ofString)-1] + "," + acString

	return strings.Replace(finalString, `"\N"`, "", -1)
}

// parseVisitDimension reads and parses the visit_dimension.csv.
func parseVisitDimension() error {
	lines, err := readCSV(inputFilePaths["VISIT_DIMENSION"])
	if err != nil {
		log.Fatal("Error in readCSV()")
		return err
	}

	tableVisitDimension = make(map[visitDimensionPK]*visitDimension)
	headerVisitDimension = make([]string, 0)
	mapNewEncounterNum = make(map[visitDimensionPK]visitDimensionPK)
	mapPatientVisits = make(map[string][]string)
	maxVisits = 0

	/* structure of visit_dimension.csv (in order):

	// pk
	"encounter_num",
	"patient_num",

	// MANDATORY FIELDS
	"active_status_cd",
	"start_date",
	"end_date",

	// OPTIONAL FIELDS
	"inout_cd",
	"location_cd",
	"location_path",
	"length_of_stay",
	"visit_blob",

	// ADMIN FIELDS
	"update_date",
	"download_date",
	"import_date",
	"sourcesystem_cd",
	"upload_id"

	*/

	for _, header := range lines[0] {
		headerVisitDimension = append(headerVisitDimension, header)
	}

	//skip header
	for _, line := range lines[1:] {
		vdk, vd := visitDimensionFromString(line)
		tableVisitDimension[vdk] = vd

		// if patient does not exist
		if _, ok := mapPatientVisits[vdk.patientNum]; !ok {
			// create array and add the encounter
			tmp := make([]string, 0)
			tmp = append(tmp, vdk.encounterNum)
			mapPatientVisits[vdk.patientNum] = tmp
		} else {
			// append encounter to array
			mapPatientVisits[vdk.patientNum] = append(mapPatientVisits[vdk.patientNum], vdk.encounterNum)
		}

		if maxVisits < len(mapPatientVisits[vdk.patientNum]) {
			maxVisits = len(mapPatientVisits[vdk.patientNum])
		}
	}

	return nil
}

// convertVisitDimension converts the old visit_dimension.csv file. The means re-randomizing the encounter_num.
// If emtpy is set to true all other data except the patient_num and encounter_num are set to empty
func convertVisitDimension() error {

	csvOutputFile, err := os.Create(outputFilesPathsNonSensitive["VISIT_DIMENSION"].Path)
	if err != nil {
		log.Fatal("Error opening [visit_dimension].csv")
		return err
	}

	defer csvOutputFile.Close()

	headerString := ""
	for _, header := range headerVisitDimension {
		headerString += "\"" + header + "\","
	}

	// re-randomize the encounter_num
	totalNbrVisits := len(tableVisitDimension)
	rand.Seed(time.Now().UnixNano())
	perm := rand.Perm(totalNbrVisits)

	// remove the last ,
	csvOutputFile.WriteString(headerString[:len(headerString)-1] + "\n")

	i := 0
	for vdk, vdp := range tableVisitDimension {
		if _, ok := visitsWithNonSensitiveObs[vdk]; ok {
			vd := *vdp
			mapNewEncounterNum[vdk] = visitDimensionPK{encounterNum: strconv.FormatInt(int64(perm[i]), 10), patientNum: mapNewPatientNum[vd.pk.patientNum]}
			vd.pk.encounterNum = strconv.FormatInt(int64(perm[i]), 10)
			vd.pk.patientNum = mapNewPatientNum[vd.pk.patientNum]
			csvOutputFile.WriteString(vd.toCSVText(false) + "\n")
			i++
		}
	}

	return nil
}
