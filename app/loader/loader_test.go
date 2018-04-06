package loader_test

import (
	"github.com/dedis/onet"
	"github.com/dedis/onet/app"
	"github.com/dedis/onet/log"
	"github.com/lca1/medco/app/loader"
	"github.com/lca1/unlynx/lib"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const (
	clinicalOntology = "files/tcga_cbio/clinical_data.csv"
	genomicOntology  = "files/tcga_cbio/mutation_data.csv"
	clinicalFile     = "files/tcga_cbio/clinical_data.csv"
	genomicFile      = "files/tcga_cbio/mutation_data.csv"
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

	fOntologyClinical, err := os.Open(clinicalOntology)
	assert.True(t, err == nil, err)
	fOntologyGenomic, err := os.Open(genomicOntology)
	assert.True(t, err == nil, err)

	fClinical, err := os.Open(clinicalFile)
	assert.True(t, err == nil, err)
	fGenomic, err := os.Open(genomicFile)
	assert.True(t, err == nil, err)

	// init global variables
	loader.FileHandlers = make([]*os.File, 0)
	loader.Testing = true
	loader.OntValues = make(map[loader.ConceptPath]loader.ConceptID)
	loader.TextSearchIndex = int64(1)

	for _, f := range loader.FilePaths {
		fp, err := os.Create(f)
		assert.True(t, err == nil, err)
		loader.FileHandlers = append(loader.FileHandlers, fp)
	}

	listSensitive := make([]string, 0)
	listSensitive = append(listSensitive, "AJCC_PATHOLOGIC_TUMOR_STAGE")
	listSensitive = append(listSensitive, "CANCER_TYPE")
	listSensitive = append(listSensitive, "CANCER_TYPE_DETAILED")
	listSensitive = append(listSensitive, "CLIN_M_STAGE")
	listSensitive = append(listSensitive, "HISTOLOGICAL_DIAGNOSIS")
	listSensitive = append(listSensitive, "ICD_O_3_HISTOLOGY")
	listSensitive = append(listSensitive, "ICD_O_3_SITE")
	listSensitive = append(listSensitive, "SAMPLE_TYPE")
	listSensitive = append(listSensitive, "TISSUE_SOURCE_SITE")
	listSensitive = append(listSensitive, "TUMOR_TISSUE_SITE")
	listSensitive = append(listSensitive, "VITAL_STATUS")
	listSensitive = append(listSensitive, "ICD_O_3_SITE")


	err = loader.GenerateOntologyFiles(el, entryPointIdx, fOntologyClinical, fOntologyGenomic, listSensitive)
	assert.True(t, err == nil, err)

	err = loader.GenerateDataFiles(el, fClinical, fGenomic)
	assert.True(t, err == nil, err)

	for _, f := range loader.FileHandlers {
		f.Close()
	}

	fClinical.Close()
	fGenomic.Close()

	fOntologyClinical.Close()
	fOntologyGenomic.Close()
}

func TestGenerateFilesLocalTest(t *testing.T) {
	el, local, err := getRoster("")
	assert.True(t, err == nil, err)
	generateFiles(t, el, 0)
	local.CloseAll()
}

func TestGenerateFilesGroupFile(t *testing.T) {
	t.Skip()
	// todo: fix hardcoded path
	// increase maximum in onet.tcp.go to allow for big packets (for now is the max value for uint32)
	el, _, err := getRoster("/Users/jagomes/Documents/EPFL/MedCo/medco-deployment/configuration-profiles/dev-3nodes-samehost/group.toml")
	assert.True(t, err == nil, err)
	generateFiles(t, el, 0)
}

func TestReplayDataset(t *testing.T) {
	t.Skip()
	err := loader.ReplayDataset(genomicFile, 2)
	assert.True(t, err == nil)
}

func TestGenerateLoadingScript(t *testing.T) {
	t.Skip()
	err := loader.GenerateLoadingScript(loader.DBSettings{DBhost: "localhost", DBport: 5434, DBname: "medcodeployment", DBuser: "postgres", DBpassword: "prigen2017"})
	assert.True(t, err == nil)
}

func TestLoadDataFiles(t *testing.T) {
	t.Skip()
	err := loader.LoadDataFiles()
	assert.True(t, err == nil)
}
