package main

import (
	"fmt"
	"os"

	"github.com/ldsec/medco"

	libunlynx "github.com/ldsec/unlynx/lib"
	"github.com/urfave/cli"
	"go.dedis.ch/onet/v3/app"
	"go.dedis.ch/onet/v3/log"
	"go.dedis.ch/onet/v3/network"
)

const (
	// BinaryName is the name of the binary
	BinaryName = "medco-unlynx"

	// DefaultGroupFile is the name of the default file to lookup for group definition
	DefaultGroupFile = "group.toml"

	optionConfig      = "config"
	optionConfigShort = "c"

	optionGroupFile      = "file"
	optionGroupFileShort = "f"

	optionDecryptKey      = "key"
	optionDecryptKeyShort = "k"

	// setup options
	optionServerBinding      = "serverBinding"
	optionServerBindingShort = "sb"

	optionDescription      = "description"
	optionDescriptionShort = "desc"

	optionPrivateTomlPath      = "privateTomlPath"
	optionPrivateTomlPathShort = "priv"

	optionPublicTomlPath      = "publicTomlPath"
	optionPublicTomlPathShort = "pub"

	optionProvidedPubKey = "pubKey"

	optionProvidedPrivKey = "privKey"

	optionProvidedSecrets      = "secrets"
	optionProvidedSecretsShort = "s"

	optionNodeIndex      = "nodeIndex"
	optionNodeIndexShort = "i"
)

/*
Return system error codes signification
0: success
1: failed to init client
2: error in the XML query parsing or during query
*/
func main() {
	// increase maximum in onet.tcp.go to allow for big packets (for now is the max value for uint32)
	network.MaxPacketSize = network.Size(^uint32(0))

	cliApp := cli.NewApp()
	cliApp.Name = "medco-unlynx"
	cliApp.Usage = "Query medical information securely and privately"
	cliApp.Version = medco.Version

	binaryFlags := []cli.Flag{
		cli.IntFlag{
			Name:  "debug, d",
			Value: 0,
			Usage: "debug-level: 1 for terse, 5 for maximal",
		},
	}

	encryptFlags := []cli.Flag{
		cli.StringFlag{
			Name:  optionGroupFile + ", " + optionGroupFileShort,
			Value: DefaultGroupFile,
			Usage: "Unlynx group definition file",
		},
	}

	decryptFlags := []cli.Flag{
		cli.StringFlag{
			Name:  optionDecryptKey + ", " + optionDecryptKeyShort,
			Usage: "Base64-encoded key to decrypt a value",
		},
	}

	mappingTableGenFlags := []cli.Flag{
		cli.StringFlag{
			Name:     "outputFile",
			Usage:    "Path to the file that will be generated.",
			Required: true,
		},
		cli.StringFlag{
			Name:     "outputFormat",
			Usage:    "Format of the output file. Value: go|typescript. Default: typescript.",
			Required: false,
			Value:    "typescript",
		},
		cli.Int64Flag{
			Name:     "nbMappings",
			Usage:    "Number of mappings to generate. Default: 1000.",
			Required: false,
			Value:    1000,
		},
		cli.BoolFlag{
			Name:     "checkNeg",
			Usage:    "Whether to check for negative values. Default: false.",
			Required: false,
		},
	}

	serverFlags := []cli.Flag{
		cli.StringFlag{
			Name:  optionConfig + ", " + optionConfigShort,
			Usage: "Configuration file of the server",
		},
	}

	nonInteractiveSetupFlags := []cli.Flag{
		cli.StringFlag{
			Name:  optionServerBinding + ", " + optionServerBindingShort,
			Usage: "Server binding address in the form of address:port",
		},
		cli.StringFlag{
			Name:  optionDescription + ", " + optionDescriptionShort,
			Usage: "Description of the node for the toml files",
		},
		cli.StringFlag{
			Name:  optionPrivateTomlPath + ", " + optionPrivateTomlPathShort,
			Usage: "Private toml file path",
		},
		cli.StringFlag{
			Name:  optionPublicTomlPath + ", " + optionPublicTomlPathShort,
			Usage: "Public toml file path",
		},
		cli.StringFlag{
			Name:  optionProvidedPubKey,
			Usage: "Provided public key (optional)",
		},
		cli.StringFlag{
			Name:  optionProvidedPrivKey,
			Usage: "Provided private key (optional)",
		},
	}

	getAggregateKeyFlags := []cli.Flag{
		cli.StringFlag{
			Name:  optionGroupFile + ", " + optionGroupFileShort,
			Value: DefaultGroupFile,
			Usage: "Unlynx group definition file",
		},
		cli.StringFlag{
			Name:  optionProvidedSecrets + ", " + optionProvidedSecretsShort,
			Usage: "Provided secrets (optional). Repeat the option to put several.",
		},
		cli.IntFlag{
			Name:  optionNodeIndex + ", " + optionNodeIndexShort,
			Usage: "Node index of the server for which the secrets are generated",
		},
	}

	cliApp.Commands = []cli.Command{
		// BEGIN CLIENT: DATA ENCRYPTION ----------
		{
			Name:    "encrypt",
			Aliases: []string{"e"},
			Usage:   "Encrypt an integer with the public key of the collective authority",
			Action:  encryptIntFromApp,
			Flags:   encryptFlags,
		},
		// CLIENT END: DATA ENCRYPTION ------------

		// BEGIN CLIENT: DATA DECRYPTION ----------
		{
			Name:    "decrypt",
			Aliases: []string{"d"},
			Usage:   "Decrypt an integer with the provided private key",
			Action:  decryptIntFromApp,
			Flags:   decryptFlags,
		},
		// CLIENT END: DATA DECRYPTION ------------

		// BEGIN CLIENT: KEY GENERATION ----------
		{
			Name:    "keygen",
			Aliases: []string{"k"},
			Usage:   "Generate a pair of public/private keys.",
			Action:  keyGenerationFromApp,
		},
		// CLIENT END: KEY GENERATION ------------

		// BEGIN CLIENT: MAPPING TABLE GENERATION ----------
		{
			Name:    "mappingtablegen",
			Aliases: []string{"m"},
			Usage:   "Generate a point-integer mapping table.",
			Action:  mappingTableGenFromApp,
			Flags:   mappingTableGenFlags,
		},
		// CLIENT END: MAPPING TABLE GENERATION ------------

		// BEGIN SERVER --------
		{
			Name:  "server",
			Usage: "Start UnLynx MedCo server",
			Action: func(c *cli.Context) error {
				if err := runServer(c); err != nil {
					return err
				}
				return nil
			},
			Flags: serverFlags,
			Subcommands: []cli.Command{
				{
					Name:    "setup",
					Aliases: []string{"s"},
					Usage:   "Setup server configuration (interactive)",
					Action: func(c *cli.Context) error {
						if c.String(optionConfig) != "" {
							return fmt.Errorf("[-] Configuration file option cannot be used for the 'setup' command")
						}
						if c.GlobalIsSet("debug") {
							return fmt.Errorf("[-] Debug option cannot be used for the 'setup' command")
						}
						app.InteractiveConfig(libunlynx.SuiTe, BinaryName)
						return nil
					},
				},
				{
					Name:    "setupNonInteractive",
					Aliases: []string{"sni"},
					Usage:   "Setup server configuration (non-interactive)",
					Action:  NonInteractiveSetup,
					Flags:   nonInteractiveSetupFlags,
				},
				{
					Name:    "getAggregateKey",
					Aliases: []string{"gak"},
					Usage:   "Get AggregateTarget Key from group.toml",
					Action:  getAggregateKey,
					Flags:   getAggregateKeyFlags,
				},
				{
					Name:    "generateTaggingSecrets",
					Aliases: []string{"gs"},
					Usage:   "Generate DDT Secrets for the participating nodes",
					Action:  generateTaggingSecrets,
					Flags:   getAggregateKeyFlags,
				},
			},
		},
		// SERVER END ----------
	}

	cliApp.Flags = binaryFlags
	cliApp.Before = func(c *cli.Context) error {
		log.SetDebugVisible(c.GlobalInt("debug"))
		return nil
	}
	err := cliApp.Run(os.Args)
	log.ErrFatal(err)
}
