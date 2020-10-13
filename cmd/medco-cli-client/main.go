package main

import (
	medcoclient "github.com/ldsec/medco/connector/client"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
	"strings"
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
		cli.StringFlag{
			Name:  "outputFile, o",
			Usage: "Output file for the result. Printed to stdout if omitted.",
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

	// --- app commands
	cliApp.Commands = []cli.Command{
		{
			Name:      "concept-children",
			Aliases:   []string{"conc"},
			Usage:     "Get the children (concepts and modifiers) of a concept",
			ArgsUsage: "conceptPath",
			Action: func(c *cli.Context) error {
				return medcoclient.ExecuteClientSearchConcept(
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
			Aliases:   []string{"modc"},
			Usage:     "Get the children of a modifier",
			ArgsUsage: "modifierPath appliedPath appliedConcept",
			Action: func(c *cli.Context) error {
				return medcoclient.ExecuteClientSearchModifier(
					c.GlobalString("token"),
					c.GlobalString("user"),
					c.GlobalString("password"),
					c.Args().Get(0),
					c.Args().Get(1),
					c.Args().Get(2),
					c.GlobalString(("outputFile")),
					c.GlobalBool("disableTLSCheck"))
			},
		},

		{
			Name:    "query",
			Aliases: []string{"q"},
			Usage:   "Query the MedCo network",
			ArgsUsage: "[patient_list|count_per_site|count_per_site_obfuscated|count_per_site_shuffled|" +
				"count_per_site_shuffled_obfuscated|count_global|count_global_obfuscated] [query string]",
			Action: func(c *cli.Context) error {
				return medcoclient.ExecuteClientQuery(
					c.GlobalString("token"),
					c.GlobalString("user"),
					c.GlobalString("password"),
					c.Args().First(),
					strings.Join(c.Args().Tail(), " "),
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
			Name:      "ga-get-variant",
			Aliases:   []string{"ga-var"},
			Usage:     "Get the variant ID of the genomic annotation of type *annotation* and value *value*",
			Flags:     genomicAnnotationsGetVariantsFlag,
			ArgsUsage: "[-z zygosity] [-e] annotation value",
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
			Name:    "node-status",
			Aliases: []string{"nod-stat"},
			Usage:   "Get node status",
			Action: func(c *cli.Context) error {
				return medcoclient.ExecuteClientGetNodeStatus(
					c.GlobalString("token"),
					c.GlobalString("user"),
					c.GlobalString("password"),
					c.GlobalString("outputFile"),
					c.GlobalBool("disableTLSCheck"),
				)
			},
		},

		{
			Name:    "network-info",
			Aliases: []string{"net-info"},
			Usage:   "Get network info metadata",
			Action: func(c *cli.Context) error {
				return medcoclient.ExecuteClientGetNetwork(
					c.GlobalString("token"),
					c.GlobalString("user"),
					c.GlobalString("password"),
					c.GlobalString("outputFile"),
					c.GlobalBool("disableTLSCheck"),
				)
			},
		},

		{
			Name:    "network-status",
			Aliases: []string{"net-stat"},
			Usage:   "Get network status",
			Action: func(c *cli.Context) error {
				return medcoclient.ExecuteClientGetNetworkStatus(
					c.GlobalString("token"),
					c.GlobalString("user"),
					c.GlobalString("password"),
					c.GlobalString("outputFile"),
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
