package loader_test

import (
	"github.com/lca1/medco-loader/loader"
	"github.com/lca1/unlynx/lib"
	"github.com/stretchr/testify/assert"
	"gopkg.in/dedis/crypto.v0/abstract"
	"gopkg.in/dedis/onet.v1/log"
	"gopkg.in/dedis/onet.v1"
	"gopkg.in/dedis/onet.v1/app"
	"testing"
	"os"
)

var publicKey abstract.Point
var secretKey abstract.Scalar
var el 		  *onet.Roster
var local     *onet.LocalTest

func getRoster(groupFilePath string) (*onet.Roster, *onet.LocalTest, error) {

	// empty string: make localtest
	if len(groupFilePath) == 0 {
		log.Info("Creating local test roster")

		local := onet.NewLocalTest()
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
		el, err := app.ReadGroupToml(f)
		if err != nil {
			log.Error("Error while reading group file", err)
			return nil, nil, err
		}
		if len(el.List) <= 0 {
			log.Error("Empty or invalid group file", err)
			return nil, nil, err
		}

		return el, nil, nil
	}
}

func setupEncryptEnv() {
	elAux, localAux, err := getRoster("")
	if err != nil {
		log.Fatal("Something went wrong when creating a testing environment!")
	}
	el = elAux
	local = localAux

	secretKey, publicKey = lib.GenKey()
}

func TestConvertAdapterMappings(t *testing.T) {
	log.SetDebugVisible(2)

	loader.ListSensitiveConceptsShrine = make(map[string]bool)
	loader.ListSensitiveConceptsShrine[`\Admit Diagnosis\`] = true
	loader.ListSensitiveConceptsShrine[`\Principal Diagnosis\`] = true
	loader.ListSensitiveConceptsShrine[`\Secondary Diagnosis\`] = true
	loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\`] = true
	loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\`] = true
	loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\`] = true
	loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\`] = true
	loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.8) Benign neoplasm of short bones of lower limb\`] = true
	loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`] = true

	assert.Nil(t, loader.ConvertAdapterMappings())
}

func TestConvertPatientDimension(t *testing.T) {
	log.SetDebugVisible(2)
	setupEncryptEnv()

	assert.Nil(t, loader.ParsePatientDimension(publicKey))
	assert.Nil(t, loader.ConvertPatientDimension())

	local.CloseAll()
}

func TestConvertShrineOntology(t *testing.T) {
	log.SetDebugVisible(2)

	loader.ListSensitiveConceptsShrine = make(map[string]bool)
	loader.ListSensitiveConceptsShrine[`\Admit Diagnosis\`] = true
	loader.ListSensitiveConceptsShrine[`\Principal Diagnosis\`] = true
	loader.ListSensitiveConceptsShrine[`\Secondary Diagnosis\`] = true
	loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\`] = true
	loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\`] = true
	loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\`] = true
	loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\`] = true
	loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.8) Benign neoplasm of short bones of lower limb\`] = true
	loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`] = true

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
	loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\`] = true
	loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\`] = true
	loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\`] = true
	loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\`] = true
	loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.8) Benign neoplasm of short bones of lower limb\`] = true
	loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`] = true

	assert.Nil(t, loader.ConvertAdapterMappings())

	assert.Nil(t, loader.ParseShrineOntology())
	assert.Nil(t, loader.ConvertShrineOntology())

	assert.Nil(t, loader.ParseLocalOntology(el,0))
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
	loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\`] = true
	loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\`] = true
	loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\`] = true
	loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\`] = true
	loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.8) Benign neoplasm of short bones of lower limb\`] = true
	loader.ListSensitiveConceptsShrine[`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`] = true

	assert.Nil(t, loader.ConvertAdapterMappings())

	assert.Nil(t, loader.ParseShrineOntology())
	assert.Nil(t, loader.ConvertShrineOntology())

	assert.Nil(t, loader.ParseLocalOntology(el,0))
	assert.Nil(t, loader.ConvertLocalOntology())

	assert.Nil(t, loader.ParseConceptDimension())
	assert.Nil(t, loader.ConvertConceptDimension())

	local.CloseAll()

}
