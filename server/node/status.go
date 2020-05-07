package node

import (
	"github.com/ldsec/medco-connector/server/genomicannotations"
	"github.com/ldsec/medco-connector/wrappers/i2b2"
	"github.com/ldsec/medco-connector/wrappers/unlynx"
)

// CheckStatus checks the status of the MedCo node
func CheckStatus() (message string, status bool) {

	okI2b2 := checkI2b2()
	if !okI2b2 {
		message += "{i2b2 error} "
	}

	okUnlynx := checkUnlynx()
	if !okUnlynx {
		message += "{Unlynx error} "
	}

	okGenomicAnnotations := checkGenomicAnnotations()
	if !okGenomicAnnotations {
		message += "{Genomic annotations error}"
	}

	return message, okI2b2 && okUnlynx && okGenomicAnnotations

}

func checkI2b2() bool {
	return i2b2.TestI2b2()
}

func checkUnlynx() bool {
	return unlynx.TestUnlynx()
}

func checkGenomicAnnotations() bool {
	return genomicannotations.TestGenomicAnnotations()
}
