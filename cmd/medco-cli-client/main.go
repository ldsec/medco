package main

import (
	"os"
	"strings"

	"github.com/ldsec/medco"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	exploreclient "github.com/ldsec/medco/connector/client/explore"
	querytoolsclient "github.com/ldsec/medco/connector/client/querytools"
	survivalclient "github.com/ldsec/medco/connector/client/survivalanalysis"
)

func main() {

	cliApp := cli.NewApp()
	cliApp.Name = "medco-cli-client"
	cliApp.Usage = "Command-line query tool for MedCo."
	cliApp.Version = medco.Version

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
		cli.StringFlag{
			Name:  "outputFile, o",
			Usage: "Output file for the result. Printed to stdout if omitted.",
			Value: "",
		},
	}

	//--- query command flags
	queryFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "timing, t",
			Usage: "Query timing: any|samevisit|sameinstancenum",
			Value: "any",
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
		cli.StringFlag{
			Name:  "parameterFile, p",
			Usage: "YAML parameter file URL",
			Value: "",
		},
		cli.StringFlag{
			Name:  "dumpFile, d",
			Usage: "Output file for the timers CSV. Printed to stdout if omitted.",
			Value: "",
		},
		cli.IntFlag{
			Name:  "limit, l",
			Usage: "Max limit of survival analysis. Unit depends on chosen granularity, default: day",
		},
		cli.StringFlag{
			Name:  "granularity, g",
			Usage: "Time resolution, one of [day, week, month, year]",
			Value: "day",
		},
		// this is supposed to be a required argument, but we need -1 for testing, and -1 is not possible to pass as an argument here
		cli.StringFlag{
			Name:  "cohortName, c",
			Usage: "Cohort identifier",
		},
		cli.StringFlag{
			Name:  "startConcept, s",
			Usage: "Survival start concept",
		},
		cli.StringFlag{
			Name:  "endConcept, e",
			Usage: "Survival end concept",
		},
	}

	//--- query tools command flags
	getCohortFlag := []cli.Flag{
		cli.IntFlag{
			Name:  "limit, l",
			Usage: "Limits the number of retrieved cohorts. 0 means no limit.",
			Value: 10,
		},
	}

	//--- query tools command flags
	postCohortFlag := []cli.Flag{
		// cli.IntSlice produces wrong results
		cli.StringFlag{
			Name:     "patientSetIDs, p",
			Usage:    "List of patient set IDs, there must be one per node",
			Required: true,
		},
		cli.StringFlag{
			Name:     "cohortName, c",
			Usage:    "Name of the new cohort",
			Required: true,
		},
	}

	//--- query tools command flags
	putCohortFlag := []cli.Flag{
		// cli.IntSlice produces wrong results
		cli.StringFlag{
			Name:     "patientSetIDs, p",
			Usage:    "List of patient set IDs, there must be one per node",
			Required: true,
		},
		cli.StringFlag{
			Name:     "cohortName, c",
			Usage:    "Name of the existing cohort",
			Required: true,
		},
	}

	//--- query tools command flags
	removeCohortFlag := []cli.Flag{
		// cli.IntSlice produces wrong results
		cli.StringFlag{
			Name:     "cohortName, c",
			Usage:    "Name of the new cohort",
			Required: true,
		},
	}

	//--- query tools command flags
	cohortsPatientListFlag := []cli.Flag{
		cli.StringFlag{
			Name:     "cohortName, c",
			Usage:    "Name of the new cohort",
			Required: true,
		},
		cli.StringFlag{
			Name:     "dumpFile, d",
			Usage:    "File for dumping ",
			Required: false,
		},
	}

	// --- app commands
	cliApp.Commands = []cli.Command{
		{
			Name:      "concept-children",
			Aliases:   []string{"con-c"},
			Usage:     "Get the concept children (both concepts and modifiers)",
			ArgsUsage: "conceptPath",
			Action: func(c *cli.Context) error {
				return exploreclient.ExecuteClientSearchConceptChildren(
					c.GlobalString("token"),
					c.GlobalString("user"),
					c.GlobalString("password"),
					c.Args().Get(0),
					c.GlobalString(("outputFile")),
					c.GlobalBool("disableTLSCheck"))
			},
		},

		{
			Name:      "modifier-children",
			Aliases:   []string{"mod-c"},
			Usage:     "Get the modifier children",
			ArgsUsage: "modifierPath appliedPath appliedConcept",
			Action: func(c *cli.Context) error {
				return exploreclient.ExecuteClientSearchModifierChildren(
					c.GlobalString("token"),
					c.GlobalString("user"),
					c.GlobalString("password"),
					c.Args().Get(0),
					c.Args().Get(1),
					c.Args().Get(2),
					c.GlobalString("outputFile"),
					c.GlobalBool("disableTLSCheck"))
			},
		},

		{
			Name:      "concept-info",
			Aliases:   []string{"con-i"},
			Usage:     "Get the concept info",
			ArgsUsage: "conceptPath",
			Action: func(c *cli.Context) error {
				return exploreclient.ExecuteClientSearchConceptInfo(
					c.GlobalString("token"),
					c.GlobalString("user"),
					c.GlobalString("password"),
					c.Args().Get(0),
					c.GlobalString(("outputFile")),
					c.GlobalBool("disableTLSCheck"))
			},
		},

		{
			Name:      "modifier-info",
			Aliases:   []string{"mod-i"},
			Usage:     "Get the modifier info",
			ArgsUsage: "modifierPath appliedPath",
			Action: func(c *cli.Context) error {
				return exploreclient.ExecuteClientSearchModifierInfo(
					c.GlobalString("token"),
					c.GlobalString("user"),
					c.GlobalString("password"),
					c.Args().Get(0),
					c.Args().Get(1),
					c.GlobalString("outputFile"),
					c.GlobalBool("disableTLSCheck"))
			},
		},

		{
			Name:      "query",
			Aliases:   []string{"q"},
			Usage:     "Query the MedCo network",
			Flags:     queryFlags,
			ArgsUsage: "[-t timing] query_string",
			Action: func(c *cli.Context) error {
				return exploreclient.ExecuteClientQuery(
					c.GlobalString("token"),
					c.GlobalString("user"),
					c.GlobalString("password"),
					strings.Join(c.Args(), " "),
					c.String("timing"),
					c.GlobalString("outputFile"),
					c.GlobalBool("disableTLSCheck"),
				)
			},
		},

		{
			Name:      "ga-get-values",
			Aliases:   []string{"ga-val"},
			Usage:     "Get the values of the genomic annotations of type *annotation* whose values contain *value*",
			Flags:     genomicAnnotationsGetValuesFlag,
			ArgsUsage: "[-l limit] annotation value",
			Action: func(c *cli.Context) error {
				return exploreclient.ExecuteClientGenomicAnnotationsGetValues(
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
			Name:      "ga-get-variant",
			Aliases:   []string{"ga-var"},
			Usage:     "Get the variant ID of the genomic annotation of type *annotation* and value *value*",
			Flags:     genomicAnnotationsGetVariantsFlag,
			ArgsUsage: "[-z zygosity] [-e] annotation value",
			Description: "zygosity can be either heterozygous, homozygous, unknown or a combination of the three separated by |\n" +
				"If omitted, the command will execute as if zygosity was equal to \"heterozygous|homozygous|unknown|\"",
			Action: func(c *cli.Context) error {
				return exploreclient.ExecuteClientGenomicAnnotationsGetVariants(
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
			Name:      "survival-analysis",
			Aliases:   []string{"srva"},
			Usage:     "Run a survival analysis",
			Flags:     survivalAnalysisFlag,
			ArgsUsage: "[-p parameterFile |  [-g granularity] -c cohortName -s startConcept [-x startModifier] -e endConcept [-y endModifier]]",
			Description: "Returns the points of the survival curve with the provided parameters." +
				"Instead of using command line arguments, paramters can also be written in parameter file." +
				"If both parameter file URL and command line argument set are used," +
				"definitions are overridden by the parameter file.",
			Action: func(c *cli.Context) error {
				return survivalclient.ExecuteClientSurvival(
					c.GlobalString("token"),
					c.String("parameterFile"),
					c.GlobalString("user"),
					c.GlobalString("password"),
					c.GlobalBool("disableTLSCheck"),
					c.GlobalString("outputFile"),
					c.String("dumpFile"),
					c.Int("limit"),
					c.String("cohortName"),
					c.String("granularity"),
					c.String("startConcept"),
					c.String("endConcept"),
				)

			},
		},

		{
			Name:        "get-saved-cohorts",
			Aliases:     []string{"getsc"},
			Usage:       "get cohorts",
			Flags:       getCohortFlag,
			ArgsUsage:   "[-l limit]",
			Description: "Gets the list of cohorts.",
			Action: func(c *cli.Context) error {
				return querytoolsclient.ExecuteGetCohorts(
					c.GlobalString("token"),
					c.GlobalString("user"),
					c.GlobalString("password"),
					c.GlobalBool("disableTLSCheck"),
					c.GlobalString("outputFile"),
					c.Int("limit"),
				)
			},
		},

		{
			Name:        "add-saved-cohorts",
			Aliases:     []string{"addsc"},
			Usage:       "Create a new cohort.",
			Flags:       postCohortFlag,
			ArgsUsage:   "-c cohortName -p patientSetIDs",
			Description: "Creates a new cohort with given name. The patient set IDs correspond to explore query result IDs.",
			Action: func(c *cli.Context) error {
				return querytoolsclient.ExecutePostCohorts(
					c.GlobalString("token"),
					c.GlobalString("user"),
					c.GlobalString("password"),
					c.String("cohortName"),
					c.String("patientSetIDs"),
					c.GlobalBool("disableTLSCheck"),
				)
			},
		},

		{
			Name:        "update-saved-cohorts",
			Aliases:     []string{"upsc"},
			Usage:       "Updates an existing cohort.",
			Flags:       putCohortFlag,
			ArgsUsage:   "-c cohortName -p patientSetIDs",
			Description: "Updates a new cohort with given name. The patient set IDs correspond to explore query result IDs.",
			Action: func(c *cli.Context) error {
				return querytoolsclient.ExecutePutCohorts(
					c.GlobalString("token"),
					c.GlobalString("user"),
					c.GlobalString("password"),
					c.String("cohortName"),
					c.String("patientSetIDs"),
					c.GlobalBool("disableTLSCheck"),
				)
			},
		},

		{
			Name:        "remove-saved-cohorts",
			Aliases:     []string{"rmsc"},
			Usage:       "Remove a cohort.",
			Flags:       removeCohortFlag,
			ArgsUsage:   "-c cohortName",
			Description: "Removes a cohort for a given name. If the user does not have a cohort with this name in DB, an error is sent.",
			Action: func(c *cli.Context) error {
				return querytoolsclient.ExecuteRemoveCohorts(
					c.GlobalString("token"),
					c.GlobalString("user"),
					c.GlobalString("password"),
					c.String("cohortName"),
					c.GlobalBool("disableTLSCheck"),
				)
			},
		},
		{
			Name:        "cohorts-patient-list",
			Aliases:     []string{"cpl"},
			Usage:       "Request the patients' numbers.",
			Flags:       cohortsPatientListFlag,
			ArgsUsage:   "-c cohortName [-d timer dump file]",
			Description: "Retrieves the numbers of the patients associated to a cohort.",
			Action: func(c *cli.Context) error {
				return querytoolsclient.ExecuteCohortsPatientList(
					c.GlobalString("token"),
					c.GlobalString("user"),
					c.GlobalString("password"),
					c.String("cohortName"),
					c.GlobalString("outputFile"),
					c.String("dumpFile"),
					c.GlobalBool("disableTLSCheck"),
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
