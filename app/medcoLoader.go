package main

import (
	"os"

	"github.com/dedis/onet/log"
	"github.com/dedis/onet/network"
	"gopkg.in/urfave/cli.v1"
)

const (
	// #---- COMMON ----#

	// DefaultGroupFile is the name of the default file to lookup for group definition
	DefaultGroupFile = "group.toml"
	// DefaultSensitiveFile is the name of the file that lists all sensitive attributes
	DefaultSensitiveFile = "sensitive.txt"
	// DefaultDBhost is the name of the default database hostname
	DefaultDBhost = "localhost"
	// DefaultDBport is the value of the default database access port
	DefaultDBport = 5434
	// DefaultDBname is the name of the default database name
	DefaultDBname = "medcodeployment"
	// DefaultDBuser is the name of the default user
	DefaultDBuser = "postgres"
	// DefaultDBpassword is the name of the default password
	DefaultDBpassword = "prigen2017"

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
	DefaultOntologyClinical = "data_clinical_skcm_broad.csv"
	// DefaultOntologyGenomic is the name of the default clinical file (dataset)
	DefaultOntologyGenomic = "data_mutations_extended_skcm_broad.csv"
	// DefaultClinicalFile is the name of the default clinical file (dataset)
	DefaultClinicalFile = "data_clinical_skcm_broad.csv"
	// DefaultGenomicFile is the name of the default genomic file (dataset)
	DefaultGenomicFile = "data_mutations_extended_skcm_broad.csv"
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

	optionOutputPath 		= "output"
	optionOuputPathShort 	= "o"

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
		},
	}

	loaderFlagsCommon := []cli.Flag{
		cli.StringFlag{
			Name:  optionGroupFile + ", " + optionGroupFileShort,
			Value: DefaultGroupFile,
			Usage: "Unlynx group definition file",
		},
		cli.IntFlag{
			Name:  optionEntryPointIdx + ", " + optionEntryPointIdxShort,
			Usage: "Index (relative to the group definition file) of the collective authority server to load the data",
		},
		cli.StringFlag{
			Name:  optionSensitiveFile + ", " + optionSensitiveFileShort,
			Value: DefaultSensitiveFile,
			Usage: "File containing a list of sensitive concepts",
		},
		cli.StringFlag{
			Name:  optionDBhost + ", " + optionDBhostShort,
			Value: DefaultDBhost,
			Usage: "Database hostname",
		},
		cli.IntFlag{
			Name:  optionDBport + ", " + optionDBportShort,
			Value: DefaultDBport,
			Usage: "Database port",
		},
		cli.StringFlag{
			Name:  optionDBname + ", " + optionDBnameShort,
			Value: DefaultDBname,
			Usage: "Database name",
		},
		cli.StringFlag{
			Name:  optionDBuser + ", " + optionDBuserShort,
			Value: DefaultDBuser,
			Usage: "Database user",
		},
		cli.StringFlag{
			Name:  optionDBpassword + ", " + optionDBpasswordShort,
			Value: DefaultDBpassword,
			Usage: "Database password",
		},
	}

	loaderFlagsv0 := []cli.Flag{
		cli.StringFlag{
			Name:  optionOntologyClinical + ", " + optionOntologyClinicalShort,
			Value: DefaultOntologyClinical,
			Usage: "Clinical ontology to load",
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
		cli.BoolFlag{
			Name:  optionEmpty + ", " + optionEmptyShort,
			Usage: "Empty patient and visit dimension tables (y/n)",
		},
	}
	loaderFlagsv1 = append(loaderFlagsCommon, loaderFlagsv1...)

	cliApp.Commands = []cli.Command{
		// BEGIN CLIENT: DATA LOADER ----------
		{
			Name:    "version0",
			Aliases: []string{"v0"},
			Usage:   "Load genomic data (e.g. tcga_bio and skcm_broad datasets)",
			Flags:   loaderFlagsv0,
			Action:  loadGenomicData,
		},
		{
			Name:    "version1",
			Aliases: []string{"v1"},
			Usage:   "Convert existing i2b2 data model",
			Flags:   loaderFlagsv1,
			Action:  loadi2b2Data,
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
