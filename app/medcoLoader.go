package main

import (
	"os"

	"go.dedis.ch/onet/v3/log"
	"go.dedis.ch/onet/v3/network"
	"gopkg.in/urfave/cli.v1"
)

const (
	// #---- COMMON ----#

	optionGroupFile      = "group"
	optionGroupFileShort = "g"

	optionSensitiveFile      = "sensitive"
	optionSensitiveFileShort = "sen"

	optionEntryPointIdx      = "entryPointIdx"
	optionEntryPointIdxShort = "entry"

	// i2b2 database settings
	optionI2b2DBhost      = "i2b2DbHost"
	optionI2b2DBhostShort = "i2b2H"

	optionI2b2DBport      = "i2b2DbPort"
	optionI2b2DBportShort = "i2b2P"

	optionI2b2DBname      = "i2b2DbName"
	optionI2b2DBnameShort = "i2b2N"

	optionI2b2DBuser      = "i2b2DbUser"
	optionI2b2DBuserShort = "i2b2U"

	optionI2b2DBpassword      = "i2b2DbPassword"
	optionI2b2DBpasswordShort = "i2b2Pw"

	// #---- V0 ----#

	// genomic annotations database settings
	optionGaDBhost      = "gaDbHost"
	optionGaDBhostShort = "gaH"

	optionGaDBport      = "gaDbPort"
	optionGaDBportShort = "gaP"

	optionGaDBname      = "gaDbName"
	optionGaDBnameShort = "gaN"

	optionGaDBuser      = "gaDbUser"
	optionGaDBuserShort = "gaU"

	optionGaDBpassword      = "gaDbPassword"
	optionGaDBpasswordShort = "gaPw"

	// DefaultOntologyClinical is the name of the default clinical file (dataset)
	DefaultOntologyClinical = "../../data/genomic/tcga_cbio/clinical_data.csv"
	// DefaultOntologyGenomic is the name of the default clinical file (dataset)
	DefaultOntologyGenomic = "../../data/genomic/tcga_cbio/mutation_data.csv"
	// DefaultClinicalFile is the name of the default clinical file (dataset)
	DefaultClinicalFile = "../../data/genomic/tcga_cbio/clinical_data.csv"
	// DefaultGenomicFile is the name of the default genomic file (dataset)
	DefaultGenomicFile = "../../data/genomic/tcga_cbio/mutation_data.csv"
	// DefaultOutputPath is the output path for the generated .csv files
	DefaultOutputPath = "../data/genomic/"

	// dataset settings (for now we have no incremental loading, and so we require both ontology and dataset files)
	optionOntologyClinical      = "ont_clinical"
	optionOntologyClinicalShort = "oc"

	optionOntologyGenomic      = "ont_genomic"
	optionOntologyGenomicShort = "og"

	optionClinicalFile      = "clinical"
	optionClinicalFileShort = "cl"

	optionGenomicFile      = "genomic"
	optionGenomicFileShort = "gen"

	optionOutputPath     = "output"
	optionOuputPathShort = "o"

	// #---- V1 ----#

	// DefaultDataFiles is the name of the default toml file with the file paths
	DefaultDataFiles = "files.toml"

	optionDataFiles      = "files"
	optionDataFilesShort = "f"

	optionEmpty      = "empty"
	optionEmptyShort = "e"
)

/*
Return system error codes signification
0: success
1: failed to init client
*/
func main() {
	// increase maximum in onet.tcp.go to allow for big packets (for now is the max value for uint32)
	network.MaxPacketSize = network.Size(^uint32(0))

	cliApp := cli.NewApp()
	cliApp.Name = "MedCo Loader"
	cliApp.Usage = "Software tool to manipulate i2b2/medco data"

	binaryFlags := []cli.Flag{
		cli.IntFlag{
			Name:   "debug, d",
			Value:  0,
			Usage:  "debug-level: 1 for terse, 5 for maximal",
			EnvVar: "LOG_LEVEL",
		},
	}

	loaderFlagsCommon := []cli.Flag{
		cli.StringFlag{
			Name:   optionGroupFile + ", " + optionGroupFileShort,
			Usage:  "UnLynx group definition file",
			EnvVar: "UNLYNX_GROUP_FILE_PATH",
		},
		cli.IntFlag{
			Name:   optionEntryPointIdx + ", " + optionEntryPointIdxShort,
			Usage:  "Index (relative to the group definition file) of the collective authority server to load the data",
			EnvVar: "UNLYNX_GROUP_FILE_IDX",
		},
		cli.StringFlag{
			Name:   optionI2b2DBhost + ", " + optionI2b2DBhostShort,
			Usage:  "I2B2 database hostname",
			EnvVar: "I2B2_DB_HOST",
		},
		cli.IntFlag{
			Name:   optionI2b2DBport + ", " + optionI2b2DBportShort,
			Usage:  "I2B2 database port",
			EnvVar: "I2B2_DB_PORT",
		},
		cli.StringFlag{
			Name:   optionI2b2DBname + ", " + optionI2b2DBnameShort,
			Usage:  "I2B2 database name",
			EnvVar: "I2B2_DB_NAME",
		},
		cli.StringFlag{
			Name:   optionI2b2DBuser + ", " + optionI2b2DBuserShort,
			Usage:  "I2B2 database user",
			EnvVar: "I2B2_DB_USER",
		},
		cli.StringFlag{
			Name:   optionI2b2DBpassword + ", " + optionI2b2DBpasswordShort,
			Usage:  "I2B2 database password",
			EnvVar: "I2B2_DB_PASSWORD",
		},
	}

	loaderFlagsv0 := []cli.Flag{
		cli.StringFlag{
			Name:  optionOntologyClinical + ", " + optionOntologyClinicalShort,
			Value: DefaultOntologyClinical,
			Usage: "Clinical ontology to load",
		},
		cli.StringFlag{
			Name:  optionSensitiveFile + ", " + optionSensitiveFileShort,
			Usage: "File with the list of clinical sensitive attributes (e.g., CANCER_TYPE_DETAILED). The entry 'all' means all attributes are considered sensitive)",
		},
		cli.StringFlag{
			Name:  optionOntologyGenomic + ", " + optionOntologyGenomicShort,
			Value: DefaultOntologyGenomic,
			Usage: "Genomic ontology to load",
		},
		cli.StringFlag{
			Name:  optionClinicalFile + ", " + optionClinicalFileShort,
			Value: DefaultClinicalFile,
			Usage: "Clinical file to load",
		},
		cli.StringFlag{
			Name:  optionGenomicFile + ", " + optionGenomicFileShort,
			Value: DefaultGenomicFile,
			Usage: "Genomic file to load",
		},
		cli.StringFlag{
			Name:  optionOutputPath + ", " + optionOuputPathShort,
			Value: DefaultOutputPath,
			Usage: "Output path for the .csv files",
		},
		cli.StringFlag{
			Name:   optionGaDBhost + ", " + optionGaDBhostShort,
			Usage:  "Genomic annotations database hostname",
			EnvVar: "GA_DB_HOST",
		},
		cli.IntFlag{
			Name:   optionGaDBport + ", " + optionGaDBportShort,
			Usage:  "Genomic annotations database port",
			EnvVar: "GA_DB_PORT",
		},
		cli.StringFlag{
			Name:   optionGaDBname + ", " + optionGaDBnameShort,
			Usage:  "Genomic annotations database name",
			EnvVar: "GA_DB_NAME",
		},
		cli.StringFlag{
			Name:   optionGaDBuser + ", " + optionGaDBuserShort,
			Usage:  "Genomic annotations database user",
			EnvVar: "GA_DB_USER",
		},
		cli.StringFlag{
			Name:   optionGaDBpassword + ", " + optionGaDBpasswordShort,
			Usage:  "Genomic annotations database password",
			EnvVar: "GA_DB_PASSWORD",
		},
	}
	loaderFlagsv0 = append(loaderFlagsCommon, loaderFlagsv0...)

	loaderFlagsv1 := []cli.Flag{
		cli.StringFlag{
			Name:  optionDataFiles + ", " + optionDataFilesShort,
			Value: DefaultDataFiles,
			Usage: "Configuration toml with the path of the all the necessary i2b2 files",
		},
		cli.StringFlag{
			Name:  optionSensitiveFile + ", " + optionSensitiveFileShort,
			Usage: `File with the list of sensitive concepts (e.g., \i2b2\Diagnoses\Neoplasms (140-239)\Benign neoplasms (210-229)\(216) Benign neoplasm of skin\). The entry 'all' means all concepts are considered sensitive)`,
		},
		cli.BoolFlag{
			Name:  optionEmpty + ", " + optionEmptyShort,
			Usage: "Empty patient and visit dimension tables",
		},
	}
	loaderFlagsv1 = append(loaderFlagsCommon, loaderFlagsv1...)

	cliApp.Commands = []cli.Command{
		// BEGIN CLIENT: DATA LOADER ----------
		{
			Name:    "v0",
			Aliases: []string{"v0"},
			Usage:   "Load genomic data (e.g. tcga_bio dataset)",
			Flags:   loaderFlagsv0,
			Action:  loadV0,
		},
		{
			Name:    "v1",
			Aliases: []string{"v1"},
			Usage:   "Convert existing i2b2 data model",
			Flags:   loaderFlagsv1,
			Action:  loadV1,
		},
		// CLIENT END: DATA LOADER ------------
	}

	cliApp.Flags = binaryFlags
	cliApp.Before = func(c *cli.Context) error {
		log.SetDebugVisible(c.GlobalInt("debug"))
		return nil
	}
	err := cliApp.Run(os.Args)
	log.ErrFatal(err)
}
