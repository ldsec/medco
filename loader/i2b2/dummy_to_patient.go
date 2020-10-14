package loaderi2b2

import (
	"github.com/ldsec/medco/loader"
	"go.dedis.ch/onet/v3/log"
)

// tableDummyToPatient contains all dummies and the original patient that is associated with them
var tableDummyToPatient map[string]string

// parseDummyToPatient reads and parses the dummy_to_patient.csv.
func parseDummyToPatient() error {
	lines, err := readCSV(outputFilesPathsSensitive["DUMMY_TO_PATIENT"].Path)
	if err != nil {
		log.Fatal("Error in readCSV()")
		return err
	}

	tableDummyToPatient = make(map[string]string)

	/* structure of patient_dimension.csv (in order):

	"dummy",
	"patient"

	*/

	//skip header
	for _, line := range lines[1:] {
		tableDummyToPatient[line[0]] = line[1]
	}

	return nil
}

// callGenerateDummiesScript calls the python script generating the dummies
func callGenerateDummiesScript() error {
	return loader.ExecuteScript("python",
		filePythonGenerateDummies,
		outputFilesPathsSensitive["OBSERVATION_FACT_FILTERED"].Path,
		outputFilesPathsSensitive["PATIENT_DIMENSION_FILTERED"].Path,
		outputFilesPathsSensitive["OBSERVATION_FACT_WITH_DUMMIES"].Path,
		outputFilesPathsSensitive["PATIENT_DIMENSION_SCRIPT"].Path,
		outputFilesPathsSensitive["DUMMY_TO_PATIENT"].Path,
	)
}
