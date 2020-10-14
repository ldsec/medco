package loaderi2b2

import (
	"go.dedis.ch/onet/v3/log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// observationFactRecord is a record of the observation_fact table
type observationFactRecord struct {
	pk              observationFactPK
	valTypeCD       string
	tValChar        string
	nValNum         string
	valueFlagCD     string
	quantityNum     string
	unitsCD         string
	endDate         string
	locationCD      string
	observationBlob string
	confidenceNum   string
	adminColumns    administrativeColumns
}

// observationFactPK is the primary key of observationFactRecord
type observationFactPK struct {
	encounterNum string
	patientNum   string
	conceptCD    string
	providerID   string
	startDate    string
	modifierCD   string
	instanceNum  string
}

var (
	// mapPatientObs maps patients to their observations
	mapPatientObs map[string][]*observationFactPK
	// mapDummyObs maps dummies to the observation of the original patients they are related too
	mapDummyObs map[string][]*observationFactPK
	// patientsWithSensitiveObs contains patients with sensitive observations
	patientsWithSensitiveObs map[patientDimensionPK]struct{}
	// visitsWithNonSensitiveObs contains the visits with at least one non sensitive observation
	visitsWithNonSensitiveObs map[visitDimensionPK]struct{}
)

// tableObservationFactNonSensitive is the observation_fact table of the non sensitive project
var tableObservationFactNonSensitive map[observationFactPK]*observationFactRecord

// tableObservationFactSensitive is the observation_fact table of the sensitive project
var tableObservationFactSensitive map[observationFactPK]*observationFactRecord

// headerObservationFact contains the header of the old_observation_fact table
var headerObservationFact []string

// headerObservationFactNonSensitive contains the header of the non sensitive observation_fact table
var headerObservationFactNonSensitive []string

// headerObservationFactSensitive contains the header of the sensitive observation_fact table
var headerObservationFactSensitive []string

// textSearchIndex counter used to fill up the last column of observation_fact tables
var textSearchIndex int64

// observationFactFromString generates a observationFactRecord struct from a parsed line of a .csv file
func observationFactFromString(line []string) (observationFactPK, *observationFactRecord) {
	ofk := observationFactPK{
		encounterNum: line[0],
		patientNum:   line[1],
		conceptCD:    line[2],
		providerID:   line[3],
		startDate:    line[4],
		modifierCD:   line[5],
		instanceNum:  line[6],
	}

	of := observationFactRecord{
		pk:              ofk,
		valTypeCD:       line[7],
		tValChar:        line[8],
		nValNum:         line[9],
		valueFlagCD:     line[10],
		quantityNum:     line[11],
		unitsCD:         line[12],
		endDate:         line[13],
		locationCD:      line[14],
		observationBlob: line[15],
		confidenceNum:   line[16],
	}

	ac := administrativeColumns{
		updateDate:      line[17],
		downloadDate:    line[18],
		importDate:      line[19],
		sourceSystemCD:  line[20],
		uploadID:        line[21],
		textSearchIndex: strconv.FormatInt(textSearchIndex, 10),
	}
	textSearchIndex++

	of.adminColumns = ac

	return ofk, &of
}

// toCSVText writes the observationFactRecord object in a way that can be added to a .csv file - "","","", etc.
func (of observationFactRecord) toCSVText() string {
	acString := "\"" + of.adminColumns.updateDate + "\"," + "\"" + of.adminColumns.downloadDate + "\"," + "\"" + of.adminColumns.importDate + "\"," + "\"" + of.adminColumns.sourceSystemCD + "\"," + "\"" + of.adminColumns.uploadID + "\"," + "\"" + of.adminColumns.textSearchIndex + "\""
	finalString := "\"" + of.pk.encounterNum + "\"," + "\"" + of.pk.patientNum + "\"," + "\"" + of.pk.conceptCD + "\"," + "\"" + of.pk.providerID + "\"," + "\"" + of.pk.startDate + "\"," + "\"" + of.pk.modifierCD + "\"," +
		"\"" + of.pk.instanceNum + "\"," + "\"" + of.valTypeCD + "\"," + "\"" + of.tValChar + "\"," + "\"" + of.nValNum + "\"," + "\"" + of.valueFlagCD + "\"," + "\"" + of.quantityNum + "\"," + "\"" + of.unitsCD + "\"," +
		"\"" + of.endDate + "\"," + "\"" + of.locationCD + "\"," + "\"" + of.observationBlob + "\"," + "\"" + of.confidenceNum + "\"," + acString

	return strings.Replace(finalString, `"\N"`, "", -1)
}

// filterOldObservationFact splits the observation_fact_old table in two tables
// containing sensitive and non sensitive observations
func filterOldObservationFact() error {

	patientsWithSensitiveObs = make(map[patientDimensionPK]struct{})
	visitsWithNonSensitiveObs = make(map[visitDimensionPK]struct{})

	lines, err := readCSV(inputFilePaths["OBSERVATION_FACT_OLD"])
	if err != nil {
		log.Fatal("Error in readCSV()")
		return err
	}

	csvOutputFileNonSensitive, err := os.Create(outputFilesPathsNonSensitive["OBSERVATION_FACT_FILTERED"].Path)
	if err != nil {
		log.Fatal("Error opening non sensitive [OBSERVATION_FACT_FILTERED].csv")
		return err
	}
	defer csvOutputFileNonSensitive.Close()

	csvOutputFileSensitive, err := os.Create(outputFilesPathsSensitive["OBSERVATION_FACT_FILTERED"].Path)
	if err != nil {
		log.Fatal("Error opening sensitive [OBSERVATION_FACT_FILTERED].csv")
		return err
	}
	defer csvOutputFileSensitive.Close()

	for _, header := range lines[0] {
		headerObservationFact = append(headerObservationFact, header)
	}

	headerString := ""
	for _, header := range headerObservationFact {
		headerString += "\"" + header + "\","
	}

	// remove the last ,
	csvOutputFileNonSensitive.WriteString(headerString[:len(headerString)-1] + "\n")
	csvOutputFileSensitive.WriteString(headerString[:len(headerString)-1] + "\n")

	//skip header
	for _, line := range lines[1:] {
		ofk, of := observationFactFromString(line)
		if ofk.modifierCD == "@" {
			//if the concept is sensitive
			if _, ok := mapConceptCodeToTag[ofk.conceptCD]; ok {
				csvOutputFileSensitive.WriteString(of.toCSVText() + "\n")
				patientsWithSensitiveObs[patientDimensionPK{ofk.patientNum}] = struct{}{}
			} else {
				csvOutputFileNonSensitive.WriteString(of.toCSVText() + "\n")
				visitsWithNonSensitiveObs[visitDimensionPK{
					encounterNum: ofk.encounterNum,
					patientNum:   ofk.patientNum,
				}] = struct{}{}
			}
		} else {
			conceptCD := ofk.conceptCD + "\\" + ofk.modifierCD
			//if the concept is sensitive
			if _, ok := mapConceptCodeToTag[conceptCD]; ok {
				ofk.conceptCD = conceptCD
				ofk.modifierCD = "@"
				csvOutputFileSensitive.WriteString(of.toCSVText() + "\n")
				patientsWithSensitiveObs[patientDimensionPK{ofk.patientNum}] = struct{}{}
			} else {
				csvOutputFileNonSensitive.WriteString(of.toCSVText() + "\n")
				visitsWithNonSensitiveObs[visitDimensionPK{
					encounterNum: ofk.encounterNum,
					patientNum:   ofk.patientNum,
				}] = struct{}{}
			}
		}

	}

	return nil
}

// parseNonSensitiveObservationFact reads and parses the non sensitive observation_fact_filtered.csv file
func parseNonSensitiveObservationFact() error {

	lines, err := readCSV(outputFilesPathsNonSensitive["OBSERVATION_FACT_FILTERED"].Path)
	if err != nil {
		log.Fatal("Error in readCSV()")
		return err
	}

	tableObservationFactNonSensitive = make(map[observationFactPK]*observationFactRecord)
	headerObservationFactNonSensitive = make([]string, 0)
	mapPatientObs = make(map[string][]*observationFactPK)
	textSearchIndex = 0

	/* structure of observation_fact_old.csv (in order):

	// pk
	"encounter_num",
	"patient_num",
	"concept_cd",
	"provider_id",
	"start_date",
	"modifier_cd",
	"instance_num",

	// MANDATORY FIELDS
	"valtype_cd",
	"tval_char",
	"nval_num",
	"valueflag_cd",
	"quantity_num",
	"units_cd",
	"end_date",
	"location_cd",
	"observation_blob",
	"confidence_num",

	// ADMIN FIELDS
	"update_date",
	"download_date",
	"import_date",
	"sourcesystem_cd",
	"upload_id",
	"text_search_index"
	*/

	for _, header := range lines[0] {
		headerObservationFactNonSensitive = append(headerObservationFactNonSensitive, header)
	}

	//skip header
	for _, line := range lines[1:] {
		ofk, of := observationFactFromString(line)

		//TODO do not consider observations where the concept is not mapped in the ontology
		if _, ok := listConceptsToIgnore[ofk.conceptCD]; !ok {
			tableObservationFactNonSensitive[ofk] = of

			// if patient does not exist
			if _, ok := mapPatientObs[ofk.patientNum]; !ok {
				// create array and add the observation
				tmp := make([]*observationFactPK, 0)
				tmp = append(tmp, &ofk)
				mapPatientObs[ofk.patientNum] = tmp
			} else {
				// append encounter to array
				mapPatientObs[ofk.patientNum] = append(mapPatientObs[ofk.patientNum], &ofk)
			}
		}
	}

	return nil

}

// convertNonSensitiveObservationFact converts the non sensitive observation_fact_filtered.csv file
func convertNonSensitiveObservationFact() error {
	rand.Seed(time.Now().UnixNano())

	csvOutputFile, err := os.Create(outputFilesPathsNonSensitive["OBSERVATION_FACT"].Path)
	if err != nil {
		log.Fatal("Error opening non sensitive [observation_fact].csv")
		return err
	}
	defer csvOutputFile.Close()

	headerString := ""
	for _, header := range headerObservationFactNonSensitive {
		headerString += "\"" + header + "\","
	}
	// remove the last ,
	csvOutputFile.WriteString(headerString[:len(headerString)-1] + "\n")

	for _, of := range tableObservationFactNonSensitive {
		copyObs := *of

		// change patient_num and encounter_num
		tmp := mapNewEncounterNum[visitDimensionPK{encounterNum: of.pk.encounterNum, patientNum: of.pk.patientNum}]
		copyObs.pk = regenerateObservationPK(&copyObs.pk, tmp.patientNum, tmp.encounterNum)

		// TODO: find out why this can be 0 (the generation should not allow this
		if copyObs.pk.encounterNum != "" {
			csvOutputFile.WriteString(copyObs.toCSVText() + "\n")

		}
	}

	return nil
}

// parseSensitiveObservationFact reads and parses the observation_fact_with_dummies.csv file
// It has to be called after convertNonSensitiveObservationFact
func parseSensitiveObservationFact() error {
	lines, err := readCSV(outputFilesPathsSensitive["OBSERVATION_FACT_WITH_DUMMIES"].Path)
	if err != nil {
		log.Fatal("Error in readCSV()")
		return err
	}

	tableObservationFactSensitive = make(map[observationFactPK]*observationFactRecord)
	headerObservationFactSensitive = make([]string, 0)
	mapDummyObs = make(map[string][]*observationFactPK)
	textSearchIndex = 0

	/* structure of observation_fact_old.csv (in order):

	// pk
	"encounter_num",
	"patient_num",
	"concept_cd",
	"provider_id",
	"start_date",
	"modifier_cd",
	"instance_num",

	// MANDATORY FIELDS
	"valtype_cd",
	"tval_char",
	"nval_num",
	"valueflag_cd",
	"quantity_num",
	"units_cd",
	"end_date",
	"location_cd",
	"observation_blob",
	"confidence_num",

	// ADMIN FIELDS
	"update_date",
	"download_date",
	"import_date",
	"sourcesystem_cd",
	"upload_id",
	"text_search_index"

	// EXTRA FIELDS (added during dummy generation)
	"cluster_label"
	*/

	for _, header := range lines[0] {
		headerObservationFactSensitive = append(headerObservationFactSensitive, header)
	}
	// remove "cluster_label"
	headerObservationFactSensitive = headerObservationFactSensitive[:len(headerObservationFactSensitive)-1]

	//skip header
	for _, line := range lines[1:] {
		ofk, of := observationFactFromString(line)

		//TODO do not consider observations where the concept is not mapped in the ontology
		if _, ok := listConceptsToIgnore[ofk.conceptCD]; !ok {
			tableObservationFactSensitive[ofk] = of

			// if patient does not exist
			if _, ok := mapPatientObs[ofk.patientNum]; !ok {
				// create array and add the observation
				tmp := make([]*observationFactPK, 0)
				tmp = append(tmp, &ofk)
				mapPatientObs[ofk.patientNum] = tmp
			} else {
				// append encounter to array
				mapPatientObs[ofk.patientNum] = append(mapPatientObs[ofk.patientNum], &ofk)
			}

			// if dummy
			if originalPatient, ok := tableDummyToPatient[ofk.patientNum]; ok {
				mapDummyObs[ofk.patientNum] = mapPatientObs[originalPatient]
			}
		}
	}

	return nil
}

// convertSensitiveObservationFact converts the observation_fact_with_dummies.csv file
func convertSensitiveObservationFact() error {
	rand.Seed(time.Now().UnixNano())

	csvOutputFile, err := os.Create(outputFilesPathsSensitive["OBSERVATION_FACT"].Path)
	if err != nil {
		log.Fatal("Error opening sensitive [observation_fact].csv")
		return err
	}
	defer csvOutputFile.Close()

	headerString := ""
	for _, header := range headerObservationFactSensitive {
		headerString += "\"" + header + "\","
	}
	// remove the last ,
	csvOutputFile.WriteString(headerString[:len(headerString)-1] + "\n")

	for _, of := range tableObservationFactSensitive {
		copyObs := *of

		// if dummy observation
		if _, ok := tableDummyToPatient[of.pk.patientNum]; ok {
			// 1. choose a random observation from the original patient
			// 2. copy the data
			// 3. change patient_num and encounter_num
			listObs := mapDummyObs[of.pk.patientNum]

			// TODO: find out why this can be 0 (the generation should not allow this
			if len(listObs) == 0 {
				continue
			}
			index := rand.Intn(len(listObs))

			obs, ok := tableObservationFactNonSensitive[*(listObs[index])]
			if !ok {
				obs = tableObservationFactSensitive[*(listObs[index])]
			}
			copyObs = *obs

			// change patient_num and encounter_num
			tmp := mapNewEncounterNum[visitDimensionPK{encounterNum: copyObs.pk.encounterNum, patientNum: of.pk.patientNum}]
			copyObs.pk = regenerateObservationPK(&copyObs.pk, tmp.patientNum, tmp.encounterNum)
			// keep the same concept (and text_search_index) that was already there
			copyObs.pk.conceptCD = of.pk.conceptCD
			copyObs.adminColumns.textSearchIndex = of.adminColumns.textSearchIndex

			// delete observation from the list (so we don't choose it again)
			listObs[index] = listObs[len(listObs)-1]
			listObs = listObs[:len(listObs)-1]
			mapDummyObs[of.pk.patientNum] = listObs

		} else { // if real observation
			// change patient_num and encounter_num
			tmp := mapNewEncounterNum[visitDimensionPK{encounterNum: of.pk.encounterNum, patientNum: of.pk.patientNum}]
			copyObs.pk = regenerateObservationPK(&copyObs.pk, tmp.patientNum, tmp.encounterNum)
		}

		// we replace the code of the sensitive concept with the correspondent tag ID
		copyObs.pk.conceptCD = "TAG_ID:" + strconv.FormatInt(mapConceptCodeToTag[copyObs.pk.conceptCD], 10)

		// TODO: connected with the previous TODO
		if copyObs.pk.encounterNum != "" {
			csvOutputFile.WriteString(copyObs.toCSVText() + "\n")
		}
	}

	return nil
}

func regenerateObservationPK(ofk *observationFactPK, patientNum, encounterNum string) observationFactPK {
	ofkNew := observationFactPK{
		encounterNum: encounterNum,
		patientNum:   patientNum,
		conceptCD:    ofk.conceptCD,
		providerID:   ofk.providerID,
		startDate:    ofk.startDate,
		modifierCD:   ofk.modifierCD,
		instanceNum:  ofk.instanceNum,
	}
	return ofkNew
}
