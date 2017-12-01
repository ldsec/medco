package loader_test

import (
	"testing"
	"github.com/lca1/medco-loader/loader"
	"gopkg.in/dedis/onet.v1/log"
)

func TestConvertAdapterMappings(t *testing.T) {
	log.SetDebugVisible(2)

	loader.ListConceptsPaths = make([]string,0)
	loader.ListConceptsPaths = append(loader.ListConceptsPaths, `\\SHRINE\SHRINE\Demographics\Age\0-9 years old\`)
	loader.ListConceptsPaths = append(loader.ListConceptsPaths, `\\SHRINE\SHRINE\Demographics\Age\0-9 years old\0 years old\`)

	loader.ConvertAdapterMappings()
}