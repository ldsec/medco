package loaderi2b2_test

import (
	"github.com/armon/go-radix"
	"github.com/dedis/kyber"
	"github.com/dedis/onet"
	"github.com/dedis/onet/app"
	"github.com/dedis/onet/log"
	"github.com/lca1/medco-loader/loader/i2b2"
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
	loaderi2b2.ListSensitiveConceptsShrine = make(map[string]bool)
	loaderi2b2.ListSensitiveConceptsShrine["a"] = true

	loaderi2b2.ListSensitiveConceptsLocal = make(map[string][]string)

	localToShrine := make(map[string][]string)
	shrineToLocal := radix.New()

	shrineToLocal.Insert("a", []string{"1", "2"})
	shrineToLocal.Insert("b", []string{"3"})
	shrineToLocal.Insert("c", []string{"4"})

	localToShrine["1"] = []string{"a"}
	localToShrine["2"] = []string{"a"}
	localToShrine["3"] = []string{"b"}
	localToShrine["4"] = []string{"c"}

	loaderi2b2.StoreSensitiveLocalConcepts(localToShrine, shrineToLocal)

	assert.Equal(t, 1, len(loaderi2b2.ListSensitiveConceptsShrine))
	assert.Equal(t, 2, len(loaderi2b2.ListSensitiveConceptsLocal))

	// Test #2
	loaderi2b2.ListSensitiveConceptsShrine = make(map[string]bool)
	loaderi2b2.ListSensitiveConceptsShrine["a"] = true

	loaderi2b2.ListSensitiveConceptsLocal = make(map[string][]string)

	localToShrine = make(map[string][]string)
	shrineToLocal = radix.New()

	shrineToLocal.Insert("a", []string{"1", "2", "3"})
	shrineToLocal.Insert("b", []string{"2"})
	shrineToLocal.Insert("c", []string{"4"})

	localToShrine["1"] = []string{"a"}
	localToShrine["2"] = []string{"a", "b"}
	localToShrine["3"] = []string{"a"}
	localToShrine["4"] = []string{"c"}

	loaderi2b2.StoreSensitiveLocalConcepts(localToShrine, shrineToLocal)

	assert.Equal(t, 2, len(loaderi2b2.ListSensitiveConceptsShrine))
	assert.Equal(t, 3, len(loaderi2b2.ListSensitiveConceptsLocal))

	// Test #3
	loaderi2b2.ListSensitiveConceptsShrine = make(map[string]bool)
	loaderi2b2.ListSensitiveConceptsShrine["a"] = true

	loaderi2b2.ListSensitiveConceptsLocal = make(map[string][]string)

	localToShrine = make(map[string][]string)
	shrineToLocal = radix.New()

	shrineToLocal.Insert("a", []string{"1", "2", "3"})
	shrineToLocal.Insert("b", []string{"2", "4"})
	shrineToLocal.Insert("c", []string{"5", "4"})
	shrineToLocal.Insert("d", []string{"6"})

	localToShrine["1"] = []string{"a"}
	localToShrine["2"] = []string{"a", "b"}
	localToShrine["3"] = []string{"a"}
	localToShrine["4"] = []string{"c", "b"}
	localToShrine["5"] = []string{"c"}
	localToShrine["6"] = []string{"d"}

	loaderi2b2.StoreSensitiveLocalConcepts(localToShrine, shrineToLocal)

	assert.Equal(t, 3, len(loaderi2b2.ListSensitiveConceptsShrine))
	assert.Equal(t, 5, len(loaderi2b2.ListSensitiveConceptsLocal))
}

func TestConvertAdapterMappings(t *testing.T) {
	log.SetDebugVisible(2)

	loaderi2b2.ListSensitiveConceptsShrine = make(map[string]bool)
	loaderi2b2.ListSensitiveConceptsShrine[`\Admit Diagnosis\`] = true
	loaderi2b2.ListSensitiveConceptsShrine[`\Principal Diagnosis\`] = true
	loaderi2b2.ListSensitiveConceptsShrine[`\Secondary Diagnosis\`] = true
	//loaderI2B2.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\`] = true
	//loaderI2B2.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\`] = true
	//loaderI2B2.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\`] = true
	//loaderI2B2.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.8) Benign neoplasm of short bones of lower limb\`] = true
	loaderi2b2.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`] = true

	assert.Nil(t, loaderi2b2.ParseAdapterMappings())
	assert.Nil(t, loaderi2b2.ConvertAdapterMappings())
}

func TestDummyToPatient(t *testing.T) {
	log.SetDebugVisible(2)

	assert.Nil(t, loaderi2b2.ParseDummyToPatient())
}

func TestConvertPatientDimension(t *testing.T) {
	log.SetDebugVisible(2)
	setupEncryptEnv()

	loaderi2b2.ParseDummyToPatient()

	assert.Nil(t, loaderi2b2.ParsePatientDimension(publicKey))
	assert.Nil(t, loaderi2b2.ConvertPatientDimension(publicKey, true))

	local.CloseAll()
}

func TestConvertVisitDimension(t *testing.T) {
	log.SetDebugVisible(2)
	setupEncryptEnv()

	loaderi2b2.ParseDummyToPatient()

	loaderi2b2.ParsePatientDimension(publicKey)
	loaderi2b2.ConvertPatientDimension(publicKey, true)

	assert.Nil(t, loaderi2b2.ParseVisitDimension())
	assert.Nil(t, loaderi2b2.ConvertVisitDimension(true))

	local.CloseAll()
}

func TestUpdateChildrenEncryptIDs(t *testing.T) {
	loaderi2b2.TableShrineOntologyConceptEnc = make(map[string]*loaderi2b2.ShrineOntology)
	loaderi2b2.TableShrineOntologyModifierEnc = make(map[string][]*loaderi2b2.ShrineOntology)

	so0 := loaderi2b2.ShrineOntology{Fullname: "\\a\\", NodeEncryptID: 0}
	so1 := loaderi2b2.ShrineOntology{Fullname: "\\a\\b\\", NodeEncryptID: 1}
	so2 := loaderi2b2.ShrineOntology{Fullname: "\\a\\c\\", NodeEncryptID: 2}
	so3 := loaderi2b2.ShrineOntology{Fullname: "\\a\\c\\d", NodeEncryptID: 3}
	so4 := loaderi2b2.ShrineOntology{Fullname: "\\a\\c\\f", NodeEncryptID: 4}

	loaderi2b2.TableShrineOntologyConceptEnc["\\a\\"] = &so0
	loaderi2b2.TableShrineOntologyConceptEnc["\\a\\b\\"] = &so1
	loaderi2b2.TableShrineOntologyConceptEnc["\\a\\c\\"] = &so2
	loaderi2b2.TableShrineOntologyConceptEnc["\\a\\c\\d"] = &so3
	loaderi2b2.TableShrineOntologyConceptEnc["\\a\\c\\f"] = &so4

	soM0 := loaderi2b2.ShrineOntology{Fullname: "\\a\\", NodeEncryptID: 0}
	soM1 := loaderi2b2.ShrineOntology{Fullname: "\\a\\", NodeEncryptID: 0}
	soM2 := loaderi2b2.ShrineOntology{Fullname: "\\a\\b\\", NodeEncryptID: 1}
	soM3 := loaderi2b2.ShrineOntology{Fullname: "\\a\\b\\", NodeEncryptID: 1}

	loaderi2b2.TableShrineOntologyModifierEnc["\\a\\"] = []*loaderi2b2.ShrineOntology{&soM0, &soM1}
	loaderi2b2.TableShrineOntologyModifierEnc["\\a\\b\\"] = []*loaderi2b2.ShrineOntology{&soM2, &soM3}

	loaderi2b2.UpdateChildrenEncryptIDs()

	assert.Equal(t, 4, len(loaderi2b2.TableShrineOntologyConceptEnc["\\a\\"].ChildrenEncryptIDs))
	assert.Equal(t, 0, len(loaderi2b2.TableShrineOntologyConceptEnc["\\a\\b\\"].ChildrenEncryptIDs))
	assert.Equal(t, 2, len(loaderi2b2.TableShrineOntologyConceptEnc["\\a\\c\\"].ChildrenEncryptIDs))
	assert.Equal(t, 0, len(loaderi2b2.TableShrineOntologyConceptEnc["\\a\\c\\d"].ChildrenEncryptIDs))
	assert.Equal(t, 0, len(loaderi2b2.TableShrineOntologyConceptEnc["\\a\\c\\f"].ChildrenEncryptIDs))

	assert.Equal(t, []int64{1}, loaderi2b2.TableShrineOntologyModifierEnc["\\a\\"][0].ChildrenEncryptIDs)
	assert.Equal(t, []int64{1}, loaderi2b2.TableShrineOntologyModifierEnc["\\a\\"][1].ChildrenEncryptIDs)
}

func TestConvertShrineOntology(t *testing.T) {
	log.SetDebugVisible(2)

	loaderi2b2.ListSensitiveConceptsShrine = make(map[string]bool)
	loaderi2b2.ListSensitiveConceptsShrine[`\Admit Diagnosis\`] = true
	loaderi2b2.ListSensitiveConceptsShrine[`\Principal Diagnosis\`] = true
	loaderi2b2.ListSensitiveConceptsShrine[`\Secondary Diagnosis\`] = true
	//loaderI2B2.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\`] = true
	//loaderI2B2.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\`] = true
	//loaderI2B2.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\`] = true
	//loaderI2B2.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.8) Benign neoplasm of short bones of lower limb\`] = true
	loaderi2b2.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`] = true

	assert.Nil(t, loaderi2b2.ParseAdapterMappings())
	assert.Nil(t, loaderi2b2.ConvertAdapterMappings())

	assert.Nil(t, loaderi2b2.ParseShrineOntology())
	assert.Nil(t, loaderi2b2.ConvertShrineOntology())
}

func TestStripByLevel(t *testing.T) {

	test := `\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`
	result := loaderi2b2.StripByLevel(test, 1, true)
	assert.Equal(t, `\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`, result)

	test = `\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`
	result = loaderi2b2.StripByLevel(test, 2, true)
	assert.Equal(t, `\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`, result)

	test = `\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`
	result = loaderi2b2.StripByLevel(test, 3, true)
	assert.Equal(t, `\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`, result)

	test = `\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`
	result = loaderi2b2.StripByLevel(test, 1, false)
	assert.Equal(t, `\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\`, result)

	test = `\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`
	result = loaderi2b2.StripByLevel(test, 2, false)
	assert.Equal(t, `\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\`, result)

	test = `\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`
	result = loaderi2b2.StripByLevel(test, 10, true)
	assert.Equal(t, "", result)

	test = `\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`
	result = loaderi2b2.StripByLevel(test, 6, true)
	assert.Equal(t, "", result)
}

func TestConvertLocalOntology(t *testing.T) {
	log.SetDebugVisible(2)
	setupEncryptEnv()
	loaderi2b2.Testing = true

	loaderi2b2.ListSensitiveConceptsShrine = make(map[string]bool)
	loaderi2b2.ListSensitiveConceptsShrine[`\Admit Diagnosis\`] = true
	loaderi2b2.ListSensitiveConceptsShrine[`\Principal Diagnosis\`] = true
	loaderi2b2.ListSensitiveConceptsShrine[`\Secondary Diagnosis\`] = true
	//loaderI2B2.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\`] = true
	//loaderI2B2.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\`] = true
	//loaderI2B2.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\`] = true
	//loaderI2B2.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.8) Benign neoplasm of short bones of lower limb\`] = true
	loaderi2b2.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`] = true

	assert.Nil(t, loaderi2b2.ParseAdapterMappings())
	assert.Nil(t, loaderi2b2.ConvertAdapterMappings())

	assert.Nil(t, loaderi2b2.ParseShrineOntology())
	assert.Nil(t, loaderi2b2.ConvertShrineOntology())

	assert.Nil(t, loaderi2b2.ParseLocalOntology(el, 0))
	assert.Nil(t, loaderi2b2.ConvertLocalOntology())

	local.CloseAll()
}

func TestConvertConceptDimension(t *testing.T) {
	log.SetDebugVisible(2)
	setupEncryptEnv()
	loaderi2b2.Testing = true

	loaderi2b2.ListSensitiveConceptsShrine = make(map[string]bool)
	loaderi2b2.ListSensitiveConceptsShrine[`\Admit Diagnosis\`] = true
	loaderi2b2.ListSensitiveConceptsShrine[`\Principal Diagnosis\`] = true
	loaderi2b2.ListSensitiveConceptsShrine[`\Secondary Diagnosis\`] = true
	//loaderI2B2.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\`] = true
	//loaderI2B2.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\`] = true
	//loaderI2B2.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\`] = true
	//loaderI2B2.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.8) Benign neoplasm of short bones of lower limb\`] = true
	loaderi2b2.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`] = true

	assert.Nil(t, loaderi2b2.ParseAdapterMappings())
	assert.Nil(t, loaderi2b2.ConvertAdapterMappings())

	assert.Nil(t, loaderi2b2.ParseShrineOntology())
	assert.Nil(t, loaderi2b2.ConvertShrineOntology())

	assert.Nil(t, loaderi2b2.ParseLocalOntology(el, 0))
	assert.Nil(t, loaderi2b2.ConvertLocalOntology())

	assert.Nil(t, loaderi2b2.ParseConceptDimension())
	assert.Nil(t, loaderi2b2.ConvertConceptDimension())

	local.CloseAll()

}

func TestConvertModifierDimension(t *testing.T) {
	log.SetDebugVisible(2)
	setupEncryptEnv()
	loaderi2b2.Testing = true

	loaderi2b2.ListSensitiveConceptsShrine = make(map[string]bool)
	loaderi2b2.ListSensitiveConceptsShrine[`\Admit Diagnosis\`] = true
	loaderi2b2.ListSensitiveConceptsShrine[`\Principal Diagnosis\`] = true
	loaderi2b2.ListSensitiveConceptsShrine[`\Secondary Diagnosis\`] = true
	//loaderI2B2.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\`] = true
	//loaderI2B2.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\`] = true
	//loaderI2B2.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\`] = true
	//loaderI2B2.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.8) Benign neoplasm of short bones of lower limb\`] = true
	loaderi2b2.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`] = true

	assert.Nil(t, loaderi2b2.ParseAdapterMappings())
	assert.Nil(t, loaderi2b2.ConvertAdapterMappings())

	assert.Nil(t, loaderi2b2.ParseShrineOntology())
	assert.Nil(t, loaderi2b2.ConvertShrineOntology())

	assert.Nil(t, loaderi2b2.ParseLocalOntology(el, 0))
	assert.Nil(t, loaderi2b2.ConvertLocalOntology())

	assert.Nil(t, loaderi2b2.ParseModifierDimension())
	assert.Nil(t, loaderi2b2.ConvertModifierDimension())

	local.CloseAll()

}

func TestConvertObservationFact(t *testing.T) {
	log.SetDebugVisible(2)
	setupEncryptEnv()
	loaderi2b2.Testing = true

	loaderi2b2.ListSensitiveConceptsShrine = make(map[string]bool)
	loaderi2b2.ListSensitiveConceptsShrine[`\Admit Diagnosis\`] = true
	loaderi2b2.ListSensitiveConceptsShrine[`\Principal Diagnosis\`] = true
	loaderi2b2.ListSensitiveConceptsShrine[`\Secondary Diagnosis\`] = true
	//loaderI2B2.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\`] = true
	//loaderI2B2.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\`] = true
	//loaderI2B2.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\`] = true
	//loaderI2B2.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.8) Benign neoplasm of short bones of lower limb\`] = true
	loaderi2b2.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`] = true

	assert.Nil(t, loaderi2b2.ParseAdapterMappings())
	assert.Nil(t, loaderi2b2.ConvertAdapterMappings())

	log.LLvl1("--- Finished converting ADAPTER_MAPPINGS ---")

	assert.Nil(t, loaderi2b2.ParseShrineOntology())
	assert.Nil(t, loaderi2b2.ConvertShrineOntology())

	log.LLvl1("--- Finished converting SHRINE_ONTOLOGY ---")

	assert.Nil(t, loaderi2b2.ParseLocalOntology(el, 0))
	assert.Nil(t, loaderi2b2.ConvertLocalOntology())

	log.LLvl1("--- Finished converting LOCAL_ONTOLOGY ---")

	assert.Nil(t, loaderi2b2.ParseDummyToPatient())

	assert.Nil(t, loaderi2b2.ParsePatientDimension(publicKey))
	assert.Nil(t, loaderi2b2.ConvertPatientDimension(publicKey, true))

	log.LLvl1("--- Finished converting PATIENT_DIMENSION ---")

	assert.Nil(t, loaderi2b2.ParseVisitDimension())
	assert.Nil(t, loaderi2b2.ConvertVisitDimension(true))

	log.LLvl1("--- Finished converting VISIT_DIMENSION ---")

	assert.Nil(t, loaderi2b2.ParseConceptDimension())
	assert.Nil(t, loaderi2b2.ConvertConceptDimension())

	log.LLvl1("--- Finished converting CONCEPT_DIMENSION ---")

	assert.Nil(t, loaderi2b2.ParseModifierDimension())
	assert.Nil(t, loaderi2b2.ConvertModifierDimension())

	log.LLvl1("--- Finished converting MODIFIER_DIMENSION ---")

	assert.Nil(t, loaderi2b2.ParseObservationFact())
	assert.Nil(t, loaderi2b2.ConvertObservationFact())

	log.LLvl1("--- Finished converting OBSERVATION_FACT ---")

	local.CloseAll()
}
