package loader_test

import (
	"github.com/lca1/medco-loader/loader"
	"github.com/lca1/unlynx/lib"
	"github.com/stretchr/testify/assert"
	"gopkg.in/dedis/crypto.v0/abstract"
	"gopkg.in/dedis/onet.v1/log"
	"testing"
)

var publicKey abstract.Point
var secretKey abstract.Scalar

func setupEncryptEnv() {
	secretKey, publicKey = lib.GenKey()
}

func TestConvertAdapterMappings(t *testing.T) {
	log.SetDebugVisible(2)

	loader.ListSensitiveConcepts = make([]string, 0)
	loader.ListSensitiveConcepts = append(loader.ListSensitiveConcepts, `\Admit Diagnosis\`)
	loader.ListSensitiveConcepts = append(loader.ListSensitiveConcepts, `\Principal Diagnosis\`)
	loader.ListSensitiveConcepts = append(loader.ListSensitiveConcepts, `\Secondary Diagnosis\`)
	loader.ListSensitiveConcepts = append(loader.ListSensitiveConcepts, `\SHRINE\Diagnoses\Neoplasms (140-239.99)\`)
	loader.ListSensitiveConcepts = append(loader.ListSensitiveConcepts, `\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\`)
	loader.ListSensitiveConcepts = append(loader.ListSensitiveConcepts, `\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\`)
	loader.ListSensitiveConcepts = append(loader.ListSensitiveConcepts, `\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.8) Benign neoplasm of short bones of lower limb\`)
	loader.ListSensitiveConcepts = append(loader.ListSensitiveConcepts, `\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`)

	assert.Nil(t, loader.ConvertAdapterMappings())
}

func TestConvertPatientDimension(t *testing.T) {
	log.SetDebugVisible(2)
	setupEncryptEnv()

	assert.Nil(t, loader.ParsePatientDimension(publicKey))
	assert.Nil(t, loader.ConvertPatientDimension())
}

func TestConvertShrineOntology(t *testing.T) {
	log.SetDebugVisible(2)

	loader.ListSensitiveConcepts = make([]string, 0)

	// sensitive concepts
	loader.ListSensitiveConcepts = append(loader.ListSensitiveConcepts, `\SHRINE\Diagnoses\`)
	loader.ListSensitiveConcepts = append(loader.ListSensitiveConcepts, `\SHRINE\Diagnoses\Neoplasms (140-239.99)\`)
	loader.ListSensitiveConcepts = append(loader.ListSensitiveConcepts, `\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\`)
	loader.ListSensitiveConcepts = append(loader.ListSensitiveConcepts, `\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\`)
	loader.ListSensitiveConcepts = append(loader.ListSensitiveConcepts, `\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.8) Benign neoplasm of short bones of lower limb\`)
	loader.ListSensitiveConcepts = append(loader.ListSensitiveConcepts, `\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`)

	// sensitive concepts (modifiers)
	loader.ListSensitiveConcepts = append(loader.ListSensitiveConcepts, `\Admit Diagnosis\`)
	loader.ListSensitiveConcepts = append(loader.ListSensitiveConcepts, `\Principal Diagnosis\`)
	loader.ListSensitiveConcepts = append(loader.ListSensitiveConcepts, `\Secondary Diagnosis\`)

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

func TestFindLocalConceptAdapterMapping(t *testing.T) {
	loader.ConvertAdapterMappings()

	conceptToFind := `\i2b2\Diagnoses\Neoplasms (140-239)\`

	check, mapping := loader.FindLocalConceptAdapterMapping(conceptToFind)
	assert.Equal(t, true, check)
	assert.Equal(t, `\SHRINE\Diagnoses\Neoplasms (140-239.99)\`, mapping)

}

func TestConvertLocalOntology(t *testing.T) {
	log.SetDebugVisible(2)

	loader.ListSensitiveConcepts = make([]string, 0)

	// sensitive concepts
	loader.ListSensitiveConcepts = append(loader.ListSensitiveConcepts, `\SHRINE\Diagnoses\`)
	loader.ListSensitiveConcepts = append(loader.ListSensitiveConcepts, `\SHRINE\Diagnoses\Neoplasms (140-239.99)\`)
	loader.ListSensitiveConcepts = append(loader.ListSensitiveConcepts, `\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\`)
	loader.ListSensitiveConcepts = append(loader.ListSensitiveConcepts, `\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\`)
	loader.ListSensitiveConcepts = append(loader.ListSensitiveConcepts, `\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.8) Benign neoplasm of short bones of lower limb\`)
	loader.ListSensitiveConcepts = append(loader.ListSensitiveConcepts, `\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`)

	// sensitive concepts (modifiers)
	loader.ListSensitiveConcepts = append(loader.ListSensitiveConcepts, `\Admit Diagnosis\`)
	loader.ListSensitiveConcepts = append(loader.ListSensitiveConcepts, `\Principal Diagnosis\`)
	loader.ListSensitiveConcepts = append(loader.ListSensitiveConcepts, `\Secondary Diagnosis\`)

}
