//go:build unit_test
// +build unit_test

package loadergenomic_test

import (
	"encoding/base64"
	"os"
	"testing"

	"github.com/ldsec/medco/loader"
	loadergenomic "github.com/ldsec/medco/loader/genomic"
	libunlynx "github.com/ldsec/unlynx/lib"
	"github.com/stretchr/testify/assert"
	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/app"
	"go.dedis.ch/onet/v3/log"
)

func init() {
	dpath := os.Getenv("DEFAULT_DATA_PATH")
	if dpath == "" {
		DefaultDataPath = "../../test/data/"
	} else {
		DefaultDataPath = dpath
	}
}

// DefaultDataPath is the default path for the data folder
var DefaultDataPath string

var (
	clinicalOntology = "genomic/tcga_cbio/8_clinical_data.csv"
	genomicOntology  = "genomic/tcga_cbio/8_mutation_data.csv"
	clinicalFile     = "genomic/tcga_cbio/8_clinical_data.csv"
	genomicFile      = "genomic/tcga_cbio/8_mutation_data.csv"

	//clinicalOntology = "genomic/tcga_cbio/clinical_data.csv"
	//genomicOntology  = "genomic/tcga_cbio/mutation_data.csv"
	//clinicalFile     = "genomic/tcga_cbio/clinical_data.csv"
	//genomicFile      = "genomic/tcga_cbio/mutation_data.csv"
)

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

func generateFiles(t *testing.T, el *onet.Roster, entryPointIdx int) {
	log.SetDebugVisible(1)

	clinicalOntology = DefaultDataPath + clinicalOntology
	genomicOntology = DefaultDataPath + genomicOntology

	fOntologyClinical, err := os.Open(clinicalOntology)
	assert.True(t, err == nil, err)
	fOntologyGenomic, err := os.Open(genomicOntology)
	assert.True(t, err == nil, err)

	clinicalFile = DefaultDataPath + clinicalFile
	genomicFile = DefaultDataPath + genomicFile

	fClinical, err := os.Open(clinicalFile)
	assert.True(t, err == nil, err)
	fGenomic, err := os.Open(genomicFile)
	assert.True(t, err == nil, err)

	// init global variables
	loadergenomic.FileHandlers = make([]*os.File, 0)
	loadergenomic.Testing = true
	loadergenomic.OntValues = make(map[loadergenomic.ConceptPath]loadergenomic.ConceptID)
	loadergenomic.TextSearchIndex = int64(1)
	loadergenomic.OutputFilePath = DefaultDataPath + "genomic/"

	for i := range loadergenomic.FilePathsOntology {
		loadergenomic.FilePathsOntology[i] = loadergenomic.OutputFilePath + loadergenomic.FilePathsOntology[i]
		fp, err := os.Create(loadergenomic.FilePathsOntology[i])
		assert.True(t, err == nil, err)
		loadergenomic.FileHandlers = append(loadergenomic.FileHandlers, fp)
	}

	for i := range loadergenomic.FilePathsData {
		loadergenomic.FilePathsData[i] = loadergenomic.OutputFilePath + loadergenomic.FilePathsData[i]
		fp, err := os.Create(loadergenomic.FilePathsData[i])
		assert.True(t, err == nil, err)
		loadergenomic.FileHandlers = append(loadergenomic.FileHandlers, fp)
	}

	mapSensitive := make(map[string]struct{}, 2) // DO NOT FORGET!! to modify the '11' value depending on the number of sensitive attributes
	/*mapSensitive["AJCC_PATHOLOGIC_TUMOR_STAGE"] = struct{}{}
	mapSensitive["CANCER_TYPE"] = struct{}{}
	mapSensitive["CANCER_TYPE_DETAILED"] = struct{}{}
	mapSensitive["HISTOLOGICAL_DIAGNOSIS"] = struct{}{}
	mapSensitive["ICD_O_3_HISTOLOGY"] = struct{}{}
	mapSensitive["ICD_O_3_SITE"] = struct{}{}
	mapSensitive["SAMPLE_TYPE"] = struct{}{}
	mapSensitive["TISSUE_SOURCE_SITE"] = struct{}{}
	mapSensitive["TUMOR_TISSUE_SITE"] = struct{}{}
	mapSensitive["VITAL_STATUS"] = struct{}{}
	mapSensitive["CLIN_M_STAGE"] = struct{}{}*/

	loadergenomic.AllSensitive = true

	err = loadergenomic.GenerateOntologyFiles(el, entryPointIdx, fOntologyClinical, fOntologyGenomic, mapSensitive)
	assert.True(t, err == nil, err)

	err = loadergenomic.GenerateDataFiles(el, fClinical, fGenomic)
	assert.True(t, err == nil, err)

	for _, f := range loadergenomic.FileHandlers {
		f.Close()
	}

	fClinical.Close()
	fGenomic.Close()

	fOntologyClinical.Close()
	fOntologyGenomic.Close()
}

func TestSanitizeHeader(t *testing.T) {
	ex := "AJCC_PATHOLOGIC_TUMOR_STAGE"
	res := loadergenomic.SanitizeHeader(ex)
	assert.Equal(t, "Ajcc Pathologic Tumor Stage", res)

	ex = "CANCER_TYPE"
	res = loadergenomic.SanitizeHeader(ex)
	assert.Equal(t, "Cancer Type", res)

	ex = "CANCER_TYPE_DETAILED"
	res = loadergenomic.SanitizeHeader(ex)
	assert.Equal(t, "Cancer Type Detailed", res)
}

func TestGenerateFilesLocalTest(t *testing.T) {
	el, local, err := getRoster("")
	assert.True(t, err == nil, err)
	generateFiles(t, el, 0)
	local.CloseAll()
}

func TestGeneratePubKey(t *testing.T) {
	t.Skip() //TODO let's skip also this test for the moment
	el, _, err := getRoster(DefaultDataPath + "genomic/group.toml")
	assert.True(t, err == nil, err)

	b, err := el.Aggregate.MarshalBinary()
	assert.True(t, err == nil, err)

	log.LLvl1("Aggregate Key:", base64.StdEncoding.EncodeToString(b))
}

func TestGenerateFilesGroupFile(t *testing.T) {
	t.Skip()
	// increase maximum in onet.tcp.go to allow for big packets (for now is the max value for uint32)
	el, _, err := getRoster(DefaultDataPath + "genomic/group.toml")

	assert.True(t, err == nil, err)
	generateFiles(t, el, 0)
}

func TestReplayDataset(t *testing.T) {
	t.Skip()
	err := loadergenomic.ReplayDataset(genomicFile, 2)
	assert.True(t, err == nil)
}

func TestGenerateLoadingScript(t *testing.T) {
	dbSettings := loader.DBSettings{DBhost: "localhost", DBport: 5434, DBname: "medcodeployment", DBuser: "postgres", DBpassword: "prigen2017"}
	err := loadergenomic.GenerateLoadingOntologyScript(dbSettings, dbSettings)
	assert.True(t, err == nil)
	err = loadergenomic.GenerateLoadingDataScript(dbSettings)
	assert.True(t, err == nil)
}

func TestLoadDataFiles(t *testing.T) {
	t.Skip()
	err := loadergenomic.LoadDataFiles()
	assert.True(t, err == nil)
}
