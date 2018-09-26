package loader_test

import (
	"github.com/armon/go-radix"
	"github.com/dedis/onet/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRadixTreeTest(t *testing.T) {
	// Create a tree
	r := radix.New()
	r.Insert(`\Admit Diagnosis\`, 1)
	r.Insert(`\Principal Diagnosis\`, 1)
	r.Insert(`\Secondary Diagnosis\`, 1)
	r.Insert(`\SHRINE\Diagnoses\`, 1)
	r.Insert(`\SHRINE\Diagnoses\Neoplasms (140-239.99)\`, 1)
	r.Insert(`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\`, 1)
	r.Insert(`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\`, 1)
	r.Insert(`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.8) Benign neoplasm of short bones of lower limb\`, 1)
	r.Insert(`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\Benign neoplasm of bone and articular cartilage (213)\(213.9) Benign neoplasm of bone and articular cartilage, site unspecified\`, 1)

	assert.Equal(t, 9, r.Len())
	r.WalkPrefix(`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\`, func(s string, v interface{}) bool {
		log.LLvl1(s)
		return false
	})
	assert.Equal(t, 4, r.DeletePrefix(`\SHRINE\Diagnoses\Neoplasms (140-239.99)\Benign neoplasms (210-229.99)\`))
	assert.Equal(t, 5, r.Len())
}
