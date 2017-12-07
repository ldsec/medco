package loader_test

import (
	"testing"
	"github.com/lca1/medco-loader/loader"
	"gopkg.in/dedis/onet.v1/log"
	"github.com/lca1/unlynx/lib"
	"gopkg.in/dedis/crypto.v0/abstract"
	"github.com/stretchr/testify/assert"
)

var publicKey abstract.Point
var secretKey abstract.Scalar

func setupEncryptEnv() {
	secretKey, publicKey = lib.GenKey()
}

func TestConvertAdapterMappings(t *testing.T) {
	log.SetDebugVisible(2)

	loader.ListSensitiveConcepts = make([]string,0)
	loader.ListSensitiveConcepts = append(loader.ListSensitiveConcepts, `\SHRINE\Demographics\Age\0-9 years old\`)
	loader.ListSensitiveConcepts = append(loader.ListSensitiveConcepts, `\SHRINE\Demographics\Age\0-9 years old\0 years old\`)

	assert.Nil(t,loader.ConvertAdapterMappings())
}

func TestConvertPatientDimension(t *testing.T) {
	log.SetDebugVisible(2)
	setupEncryptEnv()

	assert.Nil(t,loader.ParsePatientDimension(publicKey))
	assert.Nil(t,loader.ConvertPatientDimension())
}

func TestConvertShrineOntology(t *testing.T) {
	log.SetDebugVisible(2)

	loader.ListSensitiveConcepts = make([]string,0)

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