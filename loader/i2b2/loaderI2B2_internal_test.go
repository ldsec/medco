package loaderi2b2

import (
	"github.com/ldsec/medco/loader"
	"github.com/ldsec/unlynx/lib"
	"github.com/stretchr/testify/assert"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/app"
	"go.dedis.ch/onet/v3/log"
	"os"
	"testing"
)

var publicKey kyber.Point
var el *onet.Roster
var local *onet.LocalTest

func init() {
	log.SetDebugVisible(2)
	setupEncryptEnv()
	test = true
	enableModifiers = true
	allSensitive = false
}

func getRoster(groupFilePath string) (*onet.Roster, *onet.LocalTest, error) {
	// empty string: make localtest
	if len(groupFilePath) == 0 {
		log.Info("Creating local test roster")

		local := onet.NewLocalTest(libunlynx.SuiTe)
		_, el, _ := local.GenTree(3, true)
		return el, local, nil

		// generate el with group file
	} else {
		log.Info("Creating roster from group file path")

		f, err := os.Open(groupFilePath)
		if err != nil {
			log.Error("Error while opening group file", err)
			return nil, nil, err
		}
		el, err := app.ReadGroupDescToml(f)
		if err != nil {
			log.Error("Error while reading group file", err)
			return nil, nil, err
		}
		if len(el.Roster.List) <= 0 {
			log.Error("Empty or invalid group file", err)
			return nil, nil, err
		}

		return el.Roster, nil, nil
	}
}

func setupEncryptEnv() {
	elAux, localAux, err := getRoster("")
	if err != nil {
		log.Fatal("Something went wrong when creating a testing environment!")
	}
	el = elAux
	local = localAux

	_, publicKey = libunlynx.GenKey()
}

func TestConvertTableAccess(t *testing.T) {

	assert.Nil(t, parseTableAccess())
	assert.Nil(t, convertTableAccess())

}

func TestConvertOntology(t *testing.T) {

	listSensitiveConcepts = make(map[string]struct{})
	listSensitiveConcepts[`\i2b2\Diagnoses\Neoplasms (140-239)\Benign neoplasms (210-229)\(216) Benign neoplasm of skin\`] = struct{}{}

	assert.Nil(t, convertLocalOntology(el, 0))
	assert.Nil(t, generateMedCoOntology())

}

func TestConvertOntologyWithExcludedConcepts(t *testing.T) {

	// testing modifiers with m_exclusion_cd = "X"
	listSensitiveConcepts = make(map[string]struct{})
	listSensitiveConcepts[`\i2b2\Diagnoses\Neoplasms (140-239)\Malignant neoplasms (140-208)\`] = struct{}{}

	assert.Nil(t, convertLocalOntology(el, 0))
	assert.Nil(t, generateMedCoOntology())

	local.CloseAll()
}

func TestConvertModifierDimension(t *testing.T) {

	TestConvertOntology(t)

	assert.Nil(t, parseModifierDimension())
	assert.Nil(t, convertModifierDimension())

}

func TestConvertConceptDimension(t *testing.T) {

	if enableModifiers {
		TestConvertModifierDimension(t)
	} else {
		TestConvertOntology(t)
	}

	assert.Nil(t, parseConceptDimension())
	assert.Nil(t, convertConceptDimension())

}

func TestFilterOldObservationFact(t *testing.T) {

	TestConvertConceptDimension(t)

	assert.Nil(t, filterOldObservationFact())

}

func TestFilterPatientDimension(t *testing.T) {

	TestFilterOldObservationFact(t)

	assert.Nil(t, filterPatientDimension(publicKey))

}

func TestCallGenerateDummiesScript(t *testing.T) {

	TestFilterPatientDimension(t)

	assert.Nil(t, callGenerateDummiesScript())

}

func TestParseDummyToPatient(t *testing.T) {

	TestCallGenerateDummiesScript(t)

	assert.Nil(t, parseDummyToPatient())

}

func TestConvertPatientDimension(t *testing.T) {

	TestParseDummyToPatient(t)

	assert.Nil(t, parsePatientDimension(publicKey))
	assert.Nil(t, convertPatientDimension(publicKey))

	local.CloseAll()
}

func TestConvertVisitDimension(t *testing.T) {

	TestParseDummyToPatient(t)

	assert.Nil(t, parseDummyToPatient())
	assert.Nil(t, parsePatientDimension(publicKey))
	assert.Nil(t, convertPatientDimension(publicKey))

	assert.Nil(t, parseVisitDimension())
	assert.Nil(t, convertVisitDimension())

	local.CloseAll()
}

func TestUpdateChildrenEncryptIDs(t *testing.T) {
	tablesMedCoOntology = make(map[string]medCoTableInfo)
	tableMedCoOntologyConceptEnc := make(map[string]*medCoOntologyRecord)
	tablesMedCoOntology["test"] = medCoTableInfo{sensitive: tableMedCoOntologyConceptEnc}

	so0 := medCoOntologyRecord{fullname: "\\a\\", nodeEncryptID: 0}
	so1 := medCoOntologyRecord{fullname: "\\a\\b\\", nodeEncryptID: 1}
	so2 := medCoOntologyRecord{fullname: "\\a\\c\\", nodeEncryptID: 2}
	so3 := medCoOntologyRecord{fullname: "\\a\\c\\d", nodeEncryptID: 3}
	so4 := medCoOntologyRecord{fullname: "\\a\\c\\f", nodeEncryptID: 4}

	tableMedCoOntologyConceptEnc["\\a\\"] = &so0
	tableMedCoOntologyConceptEnc["\\a\\b\\"] = &so1
	tableMedCoOntologyConceptEnc["\\a\\c\\"] = &so2
	tableMedCoOntologyConceptEnc["\\a\\c\\d"] = &so3
	tableMedCoOntologyConceptEnc["\\a\\c\\f"] = &so4

	updateChildrenEncryptIDs("test")

	assert.Equal(t, 4, len(tablesMedCoOntology["test"].sensitive["\\a\\"].childrenEncryptIDs))
	assert.Equal(t, 0, len(tablesMedCoOntology["test"].sensitive["\\a\\b\\"].childrenEncryptIDs))
	assert.Equal(t, 2, len(tablesMedCoOntology["test"].sensitive["\\a\\c\\"].childrenEncryptIDs))
	assert.Equal(t, 0, len(tablesMedCoOntology["test"].sensitive["\\a\\c\\d"].childrenEncryptIDs))
	assert.Equal(t, 0, len(tablesMedCoOntology["test"].sensitive["\\a\\c\\f"].childrenEncryptIDs))
}

func TestStripByLevel(t *testing.T) {

	test := `\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`
	result := stripByLevel(test, 1, true)
	assert.Equal(t, `\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`, result)

	test = `\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`
	result = stripByLevel(test, 2, true)
	assert.Equal(t, `\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`, result)

	test = `\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`
	result = stripByLevel(test, 3, true)
	assert.Equal(t, `\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`, result)

	test = `\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`
	result = stripByLevel(test, 1, false)
	assert.Equal(t, `\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\`, result)

	test = `\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`
	result = stripByLevel(test, 2, false)
	assert.Equal(t, `\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\`, result)

	test = `\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`
	result = stripByLevel(test, 10, true)
	assert.Equal(t, "", result)

	test = `\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`
	result = stripByLevel(test, 6, true)
	assert.Equal(t, "", result)
}

func TestConvertAll(t *testing.T) {

	listSensitiveConcepts = make(map[string]struct{})
	listSensitiveConcepts[`\i2b2\Diagnoses\Neoplasms (140-239)\Benign neoplasms (210-229)\(216) Benign neoplasm of skin\`] = struct{}{}

	assert.Nil(t, convertLocalOntology(el, 0))

	log.LLvl1("--- Finished converting LOCAL_ONTOLOGY ---")

	assert.Nil(t, generateMedCoOntology())

	log.LLvl1("--- Finished generating MEDCO_ONTOLOGY ---")

	if enableModifiers {
		assert.Nil(t, parseModifierDimension())
		assert.Nil(t, convertModifierDimension())
		log.LLvl1("--- Finished converting MODIFIER_DIMENSION ---")
	}

	assert.Nil(t, parseConceptDimension())
	assert.Nil(t, convertConceptDimension())

	log.LLvl1("--- Finished converting CONCEPT_DIMENSION ---")

	assert.Nil(t, filterOldObservationFact())

	log.Lvl1("--- Finished filtering OLD_OBSERVATION_FACT ---")

	assert.Nil(t, filterPatientDimension(publicKey))

	log.Lvl1("--- Finished filtering PATIENT_DIMENSION ---")

	assert.Nil(t, callGenerateDummiesScript())

	log.Lvl1("--- Finished dummies generation ---")

	assert.Nil(t, parseDummyToPatient())

	assert.Nil(t, parsePatientDimension(publicKey))
	assert.Nil(t, convertPatientDimension(publicKey))

	log.LLvl1("--- Finished converting PATIENT_DIMENSION ---")

	assert.Nil(t, parseVisitDimension())
	assert.Nil(t, convertVisitDimension())

	log.LLvl1("--- Finished converting VISIT_DIMENSION ---")

	assert.Nil(t, parseNonSensitiveObservationFact())
	assert.Nil(t, convertNonSensitiveObservationFact())

	log.LLvl1("--- Finished converting non sensitive OBSERVATION_FACT ---")

	assert.Nil(t, parseSensitiveObservationFact())
	assert.Nil(t, convertSensitiveObservationFact())

	log.LLvl1("--- Finished converting sensitive OBSERVATION_FACT ---")

	local.CloseAll()
}

func TestGenerateLoadingDataScript(t *testing.T) {
	assert.Nil(t, generateLoadingDataScriptSensitive(loader.DBSettings{DBhost: "localhost", DBport: 5434, DBname: "i2b2medcosrv0", DBuser: "i2b2", DBpassword: "i2b2"}))
}
