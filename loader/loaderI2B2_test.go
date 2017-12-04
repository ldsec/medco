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

	loader.ListConceptsPaths = make([]string,0)
	loader.ListConceptsPaths = append(loader.ListConceptsPaths, `\\SHRINE\SHRINE\Demographics\Age\0-9 years old\`)
	loader.ListConceptsPaths = append(loader.ListConceptsPaths, `\\SHRINE\SHRINE\Demographics\Age\0-9 years old\0 years old\`)

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

	assert.Nil(t, loader.ParseShrineOntology())
	assert.Nil(t, loader.ConvertShrineOntology())
}