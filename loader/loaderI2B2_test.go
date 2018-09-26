package loader_test

import (
	"github.com/armon/go-radix"
	"github.com/dedis/kyber"
	"github.com/dedis/onet"
	"github.com/dedis/onet/app"
	"github.com/dedis/onet/log"
	"github.com/lca1/medco-loader/loader"
	"github.com/lca1/unlynx/lib"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var publicKey kyber.Point
var secretKey kyber.Scalar
var el *onet.Roster
var local *onet.LocalTest

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

	secretKey, publicKey = libunlynx.GenKey()
}

func TestStoreSensitiveLocalConcepts(t *testing.T) {
	// Test #1
	loader.ListSensitiveConceptsShrine = make(map[string]bool)
	loader.ListSensitiveConceptsShrine["a"] = true

	loader.ListSensitiveConceptsLocal = make(map[string][]string)

	localToShrine := make(map[string][]string)
	shrineToLocal := radix.New()

	shrineToLocal.Insert("a",[]string{"1", "2"})
	shrineToLocal.Insert("b",[]string{"3"})
	shrineToLocal.Insert("c",[]string{"4"})

	localToShrine["1"] = []string{"a"}
	localToShrine["2"] = []string{"a"}
	localToShrine["3"] = []string{"b"}
	localToShrine["4"] = []string{"c"}

	loader.StoreSensitiveLocalConcepts(localToShrine, shrineToLocal)

	assert.Equal(t, 1, len(loader.ListSensitiveConceptsShrine))
	assert.Equal(t, 2, len(loader.ListSensitiveConceptsLocal))

	// Test #2
	loader.ListSensitiveConceptsShrine = make(map[string]bool)
	loader.ListSensitiveConceptsShrine["a"] = true

	loader.ListSensitiveConceptsLocal = make(map[string][]string)

	localToShrine = make(map[string][]string)
	shrineToLocal = radix.New()

	shrineToLocal.Insert("a",[]string{"1", "2", "3"})
	shrineToLocal.Insert("b",[]string{"2"})
	shrineToLocal.Insert("c",[]string{"4"})

	localToShrine["1"] = []string{"a"}
	localToShrine["2"] = []string{"a", "b"}
	localToShrine["3"] = []string{"a"}
	localToShrine["4"] = []string{"c"}

	loader.StoreSensitiveLocalConcepts(localToShrine, shrineToLocal)

	assert.Equal(t, 2, len(loader.ListSensitiveConceptsShrine))
	assert.Equal(t, 3, len(loader.ListSensitiveConceptsLocal))

	// Test #3
	loader.ListSensitiveConceptsShrine = make(map[string]bool)
	loader.ListSensitiveConceptsShrine["a"] = true

	loader.ListSensitiveConceptsLocal = make(map[string][]string)

	localToShrine = make(map[string][]string)
	shrineToLocal = radix.New()

	shrineToLocal.Insert("a",[]string{"1", "2", "3"})
	shrineToLocal.Insert("b",[]string{"2", "4"})
	shrineToLocal.Insert("c",[]string{"5", "4"})
	shrineToLocal.Insert("d",[]string{"6"})

	localToShrine["1"] = []string{"a"}
	localToShrine["2"] = []string{"a", "b"}
	localToShrine["3"] = []string{"a"}
	localToShrine["4"] = []string{"c", "b"}
	localToShrine["5"] = []string{"c"}
	localToShrine["6"] = []string{"d"}

	loader.StoreSensitiveLocalConcepts(localToShrine, shrineToLocal)

	assert.Equal(t, 3, len(loader.ListSensitiveConceptsShrine))
	assert.Equal(t, 5, len(loader.ListSensitiveConceptsLocal))
}

func TestConvertAdapterMappings(t *testing.T) {
	log.SetDebugVisible(2)

	loader.ListSensitiveConceptsShrine = make(map[string]bool)
	loader.ListSensitiveConceptsShrine[`\Admit Diagnosis\`] = true
	loader.ListSensitiveConceptsShrine[`\Principal Diagnosis\`] = true
	loader.ListSensitiveConceptsShrine[`\Secondary Diagnosis\`] = true
	//loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\`] = true
	//loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\`] = true
	//loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\`] = true
	//loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.8) Benign neoplasm of short bones of lower limb\`] = true
	loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`] = true

	assert.Nil(t, loader.ParseAdapterMappings())
	assert.Nil(t, loader.ConvertAdapterMappings())
}

func TestDummyToPatient(t *testing.T) {
	log.SetDebugVisible(2)

	assert.Nil(t, loader.ParseDummyToPatient())
}

func TestConvertPatientDimension(t *testing.T) {
	log.SetDebugVisible(2)
	setupEncryptEnv()

	loader.ParseDummyToPatient()

	assert.Nil(t, loader.ParsePatientDimension(publicKey))
	assert.Nil(t, loader.ConvertPatientDimension(publicKey, true))

	local.CloseAll()
}

func TestConvertVisitDimension(t *testing.T) {
	log.SetDebugVisible(2)
	setupEncryptEnv()

	loader.ParseDummyToPatient()

	loader.ParsePatientDimension(publicKey)
	loader.ConvertPatientDimension(publicKey, true)

	assert.Nil(t, loader.ParseVisitDimension())
	assert.Nil(t, loader.ConvertVisitDimension(true))

	local.CloseAll()
}

func TestUpdateChildrenEncryptIDs(t *testing.T) {
	loader.TableShrineOntologyConceptEnc = make(map[string]*loader.ShrineOntology)
	loader.TableShrineOntologyModifierEnc = make(map[string][]*loader.ShrineOntology)

	so0 := loader.ShrineOntology{Fullname: "\\a\\", NodeEncryptID: 0}
	so1 := loader.ShrineOntology{Fullname: "\\a\\b\\", NodeEncryptID: 1}
	so2 := loader.ShrineOntology{Fullname: "\\a\\c\\", NodeEncryptID: 2}
	so3 := loader.ShrineOntology{Fullname: "\\a\\c\\d", NodeEncryptID: 3}
	so4 := loader.ShrineOntology{Fullname: "\\a\\c\\f", NodeEncryptID: 4}

	loader.TableShrineOntologyConceptEnc["\\a\\"] = &so0
	loader.TableShrineOntologyConceptEnc["\\a\\b\\"] = &so1
	loader.TableShrineOntologyConceptEnc["\\a\\c\\"] = &so2
	loader.TableShrineOntologyConceptEnc["\\a\\c\\d"] = &so3
	loader.TableShrineOntologyConceptEnc["\\a\\c\\f"] = &so4

	soM0 := loader.ShrineOntology{Fullname: "\\a\\", NodeEncryptID: 0}
	soM1 := loader.ShrineOntology{Fullname: "\\a\\", NodeEncryptID: 0}
	soM2 := loader.ShrineOntology{Fullname: "\\a\\b\\", NodeEncryptID: 1}
	soM3 := loader.ShrineOntology{Fullname: "\\a\\b\\", NodeEncryptID: 1}

	loader.TableShrineOntologyModifierEnc["\\a\\"] = []*loader.ShrineOntology{&soM0, &soM1}
	loader.TableShrineOntologyModifierEnc["\\a\\b\\"] = []*loader.ShrineOntology{&soM2, &soM3}

	loader.UpdateChildrenEncryptIDs()

	assert.Equal(t, 4, len(loader.TableShrineOntologyConceptEnc["\\a\\"].ChildrenEncryptIDs))
	assert.Equal(t, 0, len(loader.TableShrineOntologyConceptEnc["\\a\\b\\"].ChildrenEncryptIDs))
	assert.Equal(t, 2, len(loader.TableShrineOntologyConceptEnc["\\a\\c\\"].ChildrenEncryptIDs))
	assert.Equal(t, 0, len(loader.TableShrineOntologyConceptEnc["\\a\\c\\d"].ChildrenEncryptIDs))
	assert.Equal(t, 0, len(loader.TableShrineOntologyConceptEnc["\\a\\c\\f"].ChildrenEncryptIDs))

	assert.Equal(t, []int64{1}, loader.TableShrineOntologyModifierEnc["\\a\\"][0].ChildrenEncryptIDs)
	assert.Equal(t, []int64{1}, loader.TableShrineOntologyModifierEnc["\\a\\"][1].ChildrenEncryptIDs)
}

func TestConvertShrineOntology(t *testing.T) {
	log.SetDebugVisible(2)

	loader.ListSensitiveConceptsShrine = make(map[string]bool)
	loader.ListSensitiveConceptsShrine[`\Admit Diagnosis\`] = true
	loader.ListSensitiveConceptsShrine[`\Principal Diagnosis\`] = true
	loader.ListSensitiveConceptsShrine[`\Secondary Diagnosis\`] = true
	//loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\`] = true
	//loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\`] = true
	//loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\`] = true
	//loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.8) Benign neoplasm of short bones of lower limb\`] = true
	loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`] = true

	assert.Nil(t, loader.ParseAdapterMappings())
	assert.Nil(t, loader.ConvertAdapterMappings())

	assert.Nil(t, loader.ParseShrineOntology())
	assert.Nil(t, loader.ConvertShrineOntology())
}

func TestStripByLevel(t *testing.T) {

	test := `\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`
	result := loader.StripByLevel(test, 1, true)
	assert.Equal(t, `\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`, result)

	test = `\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`
	result = loader.StripByLevel(test, 2, true)
	assert.Equal(t, `\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`, result)

	test = `\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`
	result = loader.StripByLevel(test, 3, true)
	assert.Equal(t, `\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`, result)

	test = `\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`
	result = loader.StripByLevel(test, 1, false)
	assert.Equal(t, `\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\`, result)

	test = `\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`
	result = loader.StripByLevel(test, 2, false)
	assert.Equal(t, `\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\`, result)

	test = `\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`
	result = loader.StripByLevel(test, 10, true)
	assert.Equal(t, "", result)

	test = `\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`
	result = loader.StripByLevel(test, 6, true)
	assert.Equal(t, "", result)
}

func TestConvertLocalOntology(t *testing.T) {
	log.SetDebugVisible(2)
	setupEncryptEnv()
	loader.Testing = true

	loader.ListSensitiveConceptsShrine = make(map[string]bool)
	loader.ListSensitiveConceptsShrine[`\Admit Diagnosis\`] = true
	loader.ListSensitiveConceptsShrine[`\Principal Diagnosis\`] = true
	loader.ListSensitiveConceptsShrine[`\Secondary Diagnosis\`] = true
	//loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\`] = true
	//loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\`] = true
	//loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\`] = true
	//loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.8) Benign neoplasm of short bones of lower limb\`] = true
	loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`] = true

	assert.Nil(t, loader.ParseAdapterMappings())
	assert.Nil(t, loader.ConvertAdapterMappings())

	assert.Nil(t, loader.ParseShrineOntology())
	assert.Nil(t, loader.ConvertShrineOntology())

	assert.Nil(t, loader.ParseLocalOntology(el, 0))
	assert.Nil(t, loader.ConvertLocalOntology())

	local.CloseAll()
}

func TestConvertConceptDimension(t *testing.T) {
	log.SetDebugVisible(2)
	setupEncryptEnv()
	loader.Testing = true

	loader.ListSensitiveConceptsShrine = make(map[string]bool)
	loader.ListSensitiveConceptsShrine[`\Admit Diagnosis\`] = true
	loader.ListSensitiveConceptsShrine[`\Principal Diagnosis\`] = true
	loader.ListSensitiveConceptsShrine[`\Secondary Diagnosis\`] = true
	//loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\`] = true
	//loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\`] = true
	//loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\`] = true
	//loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.8) Benign neoplasm of short bones of lower limb\`] = true
	loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`] = true

	assert.Nil(t, loader.ParseAdapterMappings())
	assert.Nil(t, loader.ConvertAdapterMappings())

	assert.Nil(t, loader.ParseShrineOntology())
	assert.Nil(t, loader.ConvertShrineOntology())

	assert.Nil(t, loader.ParseLocalOntology(el, 0))
	assert.Nil(t, loader.ConvertLocalOntology())

	assert.Nil(t, loader.ParseConceptDimension())
	assert.Nil(t, loader.ConvertConceptDimension())

	local.CloseAll()

}

func TestConvertModifierDimension(t *testing.T) {
	log.SetDebugVisible(2)
	setupEncryptEnv()
	loader.Testing = true

	loader.ListSensitiveConceptsShrine = make(map[string]bool)
	loader.ListSensitiveConceptsShrine[`\Admit Diagnosis\`] = true
	loader.ListSensitiveConceptsShrine[`\Principal Diagnosis\`] = true
	loader.ListSensitiveConceptsShrine[`\Secondary Diagnosis\`] = true
	//loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\`] = true
	//loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\`] = true
	//loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\`] = true
	//loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.8) Benign neoplasm of short bones of lower limb\`] = true
	loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`] = true

	assert.Nil(t, loader.ParseAdapterMappings())
	assert.Nil(t, loader.ConvertAdapterMappings())

	assert.Nil(t, loader.ParseShrineOntology())
	assert.Nil(t, loader.ConvertShrineOntology())

	assert.Nil(t, loader.ParseLocalOntology(el, 0))
	assert.Nil(t, loader.ConvertLocalOntology())

	assert.Nil(t, loader.ParseModifierDimension())
	assert.Nil(t, loader.ConvertModifierDimension())

	local.CloseAll()

}

func TestConvertObservationFact(t *testing.T) {
	log.SetDebugVisible(2)
	setupEncryptEnv()
	loader.Testing = true

	loader.ListSensitiveConceptsShrine = make(map[string]bool)
	loader.ListSensitiveConceptsShrine[`\Admit Diagnosis\`] = true
	loader.ListSensitiveConceptsShrine[`\Principal Diagnosis\`] = true
	loader.ListSensitiveConceptsShrine[`\Secondary Diagnosis\`] = true
	//loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\`] = true
	//loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\`] = true
	//loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\`] = true
	//loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.8) Benign neoplasm of short bones of lower limb\`] = true
	loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`] = true

	assert.Nil(t, loader.ParseAdapterMappings())
	assert.Nil(t, loader.ConvertAdapterMappings())

	log.LLvl1("\n\n\n--- Finished converting ADAPTER_MAPPINGS ---\n\n")

	assert.Nil(t, loader.ParseShrineOntology())
	assert.Nil(t, loader.ConvertShrineOntology())

	log.LLvl1("\n\n\n--- Finished converting SHRINE_ONTOLOGY ---\n\n")

	assert.Nil(t, loader.ParseLocalOntology(el, 0))
	assert.Nil(t, loader.ConvertLocalOntology())

	log.LLvl1("\n\n\n--- Finished converting LOCAL_ONTOLOGY ---\n\n")

	assert.Nil(t, loader.ParseDummyToPatient())

	assert.Nil(t, loader.ParsePatientDimension(publicKey))
	assert.Nil(t, loader.ConvertPatientDimension(publicKey, true))

	log.LLvl1("\n\n\n--- Finished converting PATIENT_DIMENSION ---\n\n")

	assert.Nil(t, loader.ParseVisitDimension())
	assert.Nil(t, loader.ConvertVisitDimension(true))

	log.LLvl1("\n\n\n--- Finished converting VISIT_DIMENSION ---\n\n")

	assert.Nil(t, loader.ParseConceptDimension())
	assert.Nil(t, loader.ConvertConceptDimension())

	log.LLvl1("\n\n\n--- Finished converting CONCEPT_DIMENSION ---\n\n")

	assert.Nil(t, loader.ParseModifierDimension())
	assert.Nil(t, loader.ConvertModifierDimension())

	log.LLvl1("\n\n\n--- Finished converting MODIFIER_DIMENSION ---\n\n")

	assert.Nil(t, loader.ParseObservationFact())
	assert.Nil(t, loader.ConvertObservationFact())

	log.LLvl1("\n\n\n--- Finished converting OBSERVATION_FACT ---\n\n")

	local.CloseAll()
}
