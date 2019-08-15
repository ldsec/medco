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

	// database settings
	optionDBhost      = "dbHost"
	optionDBhostShort = "dbH"

	optionDBport      = "dbPort"
	optionDBportShort = "dbP"

	optionDBname      = "dbName"
	optionDBnameShort = "dbN"

	optionDBuser      = "dbUser"
	optionDBuserShort = "dbU"

	optionDBpassword      = "dbPassword"
	optionDBpasswordShort = "dbPw"

	// #---- V0 ----#

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
			Name:  "debug, d",
			Value: 0,
			Usage: "debug-level: 1 for terse, 5 for maximal",
			EnvVar: "LOG_LEVEL",
		},
	}

	loaderFlagsCommon := []cli.Flag{
		cli.StringFlag{
			Name:  optionGroupFile + ", " + optionGroupFileShort,
			Usage: "UnLynx group definition file",
			EnvVar: "UNLYNX_GROUP_FILE_PATH",
		},
		cli.IntFlag{
			Name:  optionEntryPointIdx + ", " + optionEntryPointIdxShort,
			Usage: "Index (relative to the group definition file) of the collective authority server to load the data",
			EnvVar: "UNLYNX_GROUP_FILE_IDX",
		},
		cli.StringFlag{
			Name:  optionDBhost + ", " + optionDBhostShort,
			Usage: "Database hostname",
			EnvVar: "DB_HOST",
		},
		cli.IntFlag{
			Name:  optionDBport + ", " + optionDBportShort,
			Usage: "Database port",
			EnvVar: "DB_PORT",
		},
		cli.StringFlag{
			Name:  optionDBname + ", " + optionDBnameShort,
			Usage: "Database name",
			EnvVar: "DB_NAME",
		},
		cli.StringFlag{
			Name:  optionDBuser + ", " + optionDBuserShort,
			Usage: "Database user",
			EnvVar: "DB_USER",
		},
		cli.StringFlag{
			Name:  optionDBpassword + ", " + optionDBpasswordShort,
			Usage: "Database password",
			EnvVar: "DB_PASSWORD",
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
