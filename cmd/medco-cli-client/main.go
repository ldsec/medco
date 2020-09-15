package main

import (
	"os"
	"strings"

	survivalclient "github.com/ldsec/medco-connector/survival/client"

	medcoclient "github.com/ldsec/medco-connector/client"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {

	cliApp := cli.NewApp()
	cliApp.Name = "medco-cli-client"
	cliApp.Usage = "Command-line query tool for MedCo."
	cliApp.Version = "1.0.0" // todo: dynamically get version from build process

	// from env / config: whatever is in the config of GB : debug,
	// cli: whatever is user input

	// --- global app flags
	cliApp.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "user, u",
			Usage: "OIDC user login",
		},
		cli.StringFlag{
			Name:  "password, p",
			Usage: "OIDC password login",
		},
		cli.StringFlag{
			Name:  "token, t",
			Usage: "OIDC token",
		},
		cli.BoolFlag{
			Name:  "disableTLSCheck",
			Usage: "Disable check of TLS certificates",
		},
	}

	// --- search command flags
	//searchCommandFlags := []cli.Flag{
	//	cli.StringFlag{
	//		Name:  "path",
	//		Usage: "File containing the query definition",
	//
	//	},
	//}

	//--- query command flags
	queryCommandFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "resultFile, r",
			Usage: "Output file for the result CSV. Printed to stdout if omitted.",
			Value: "",
		},
	}

	//--- genomic annotations get values command flags
	genomicAnnotationsGetValuesFlag := []cli.Flag{
		cli.Int64Flag{
			Name:  "limit, l",
			Usage: "Maximum number of returned values",
			Value: 0,
		},
	}

	//--- genomic annotations get variants command flags
	genomicAnnotationsGetVariantsFlag := []cli.Flag{
		cli.StringFlag{
			Name:  "zygosity, z",
			Usage: "Variant zygosysty",
			Value: "",
		},
		cli.BoolFlag{
			Name:  "encrypted, e",
			Usage: "Return encrypted variant id",
		},
	}

	//--- survival analysis command flags
	survivalAnalysisFlag := []cli.Flag{
		// TODO as CLI is ran in a docker container, some adjustments are needed for I/O files to work
		//cli.StringFlag{
		//	Name:  "parameterFile, p",
		//	Usage: "YAML parameter file URL",
		//	Value: "",
		//},
		//cli.StringFlag{
		//	Name:  "resultFile, r",
		//	Usage: "Output file for the result CSV. Printed to stdout if omitted.",
		//	Value: "",
		//},
		//cli.StringFlag{
		//	Name:  "dumpFile, d",
		//	Usage: "Output file for the timers CSV. Printed to stdout if omitted.",
		//	Value: "",
		//},
		cli.IntFlag{
			Name:     "limit, l",
			Usage:    "Max limit of survival analysis.",
			Required: true,
		},
		cli.StringFlag{
			Name:  "granularity, g",
			Usage: "Time resolution, one of [day, week, month, year]",
			Value: "day",
		},
		// this is supposed to be a required argument, but we need -1 for testing, and -1 is not possible to pass as an argument here
		cli.IntFlag{
			Name:  "cohortID, c",
			Usage: "Cohort identifier",
			Value: -1,
		},
		cli.StringFlag{
			Name:     "startConcept, s",
			Usage:    "Survival start concept",
			Required: true},
		cli.StringFlag{
			Name:  "startModifier, x",
			Usage: "Survival start modifier",
			Value: "@",
		},
		cli.StringFlag{
			Name:     "endConcept, e",
			Usage:    "Survival end concept",
			Required: true,
		},
		cli.StringFlag{
			Name:  "endModifier, y",
			Usage: "Survival end modifier",
			Value: "@",
		},
	}

	// --- app commands
	cliApp.Commands = []cli.Command{
		//{
		//	Name:    "search",
		//	Aliases: []string{"s"},
		//	Usage:   "Browse the MedCo tree ontology",
		//	Action:  encryptIntFromApp,
		//	Flags:   searchCommandFlags,
		//	ArgsUsage: "",
		//
		//},
		{
			Name:    "query",
			Aliases: []string{"q"},
			Usage:   "Query the MedCo network",
			Flags:   queryCommandFlags,
			ArgsUsage: "[patient_list|count_per_site|count_per_site_obfuscated|count_per_site_shuffled|" +
				"count_per_site_shuffled_obfuscated|count_global|count_global_obfuscated] [query string]",
			Action: func(c *cli.Context) error {
				return medcoclient.ExecuteClientQuery(
					c.GlobalString("token"),
					c.GlobalString("user"),
					c.GlobalString("password"),
					c.Args().First(),
					strings.Join(c.Args().Tail(), " "),
					c.String("resultFile"),
					c.GlobalBool("disableTLSCheck"),
				)
			},
		},

		{
			Name:      "genomic-annotations-get-values",
			Aliases:   []string{"gval"},
			Usage:     "Get genomic annotations values",
			Flags:     genomicAnnotationsGetValuesFlag,
			ArgsUsage: "[-l limit] [annotation] [value]",
			Action: func(c *cli.Context) error {
				return medcoclient.ExecuteClientGenomicAnnotationsGetValues(
					c.GlobalString("token"),
					c.GlobalString("user"),
					c.GlobalString("password"),
					c.Args().Get(0),
					c.Args().Get(1),
					c.Int64("limit"),
					c.GlobalBool("disableTLSCheck"),
				)
			},
		},

		{
			Name:      "genomic-annotations-get-variants",
			Aliases:   []string{"gvar"},
			Usage:     "Get genomic annotations variants",
			Flags:     genomicAnnotationsGetVariantsFlag,
			ArgsUsage: "[-z zygosity] [-e] [annotation] [value]",
			Description: "zygosity can be either heterozygous, homozygous, unknown or a combination of the three separated by |\n" +
				"If omitted, the command will execute as if zygosity was equal to \"heterozygous|homozygous|unknown|\"",
			Action: func(c *cli.Context) error {
				return medcoclient.ExecuteClientGenomicAnnotationsGetVariants(
					c.GlobalString("token"),
					c.GlobalString("user"),
					c.GlobalString("password"),
					c.Args().Get(0),
					c.Args().Get(1),
					c.String("zygosity"),
					c.Bool("encrypted"),
					c.GlobalBool("disableTLSCheck"),
				)
			},
		},

		{
			Name:        "survival-analysis",
			Aliases:     []string{"srva"},
			Usage:       "Run a survival analysis",
			Flags:       survivalAnalysisFlag,
			ArgsUsage:   "-l limit [-g granularity] -c cohortID -s startConcept [-x startModifier] -e endConcept -y endModifier",
			Description: "Returns the points of the survival curve",
			Action: func(c *cli.Context) error {
				return survivalclient.ExecuteClientSurvival(
					c.GlobalString("token"),
					"",
					c.GlobalString("user"),
					c.GlobalString("password"),
					c.GlobalBool("disableTLSCheck"),
					"",
					"",
					c.Int("limit"),
					c.Int("cohortID"),
					c.String("granularity"),
					c.String("startConcept"),
					c.String("startModifier"),
					c.String("endConcept"),
					c.String("endModifier"),
				)

			},
		},
	}

	//cliApp.Before = func(c *cli.Context) error {
	//	log.SetDebugVisible(c.GlobalInt("debug"))
	//	return nil
	//}
	err := cliApp.Run(os.Args)
	if err != nil {
		logrus.Error(err)
	}
}
