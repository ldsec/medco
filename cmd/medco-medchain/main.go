package main

import (
	"os"
	"time"

	"github.com/ldsec/medco"
	"github.com/urfave/cli"
	"go.dedis.ch/cothority/v3/byzcoin"
	"go.dedis.ch/cothority/v3/byzcoin/bcadmin/lib"
	"go.dedis.ch/cothority/v3/darc"
	"go.dedis.ch/onet/v3/log"
	"golang.org/x/xerrors"
)

const (
	// BinaryName is the name of the binary
	BinaryName = "medco-medchain"

	optionRosterTomlPath      = "rosterTomlPath"
	optionRosterTomlPathShort = "roster"

	optionBlockInterval      = "interval"
	optionBlockIntervalShort = "i"

	optionConfigDirectoryPath      = "configDirPath"
	optionConfigDirectoryPathShort = "c"
)

func main() {
	cliApp := cli.NewApp()
	cliApp.Name = "medco-medchain"
	cliApp.Usage = "Interact with Medchain"
	cliApp.Version = medco.Version

	nonInteractiveSetupFlags := []cli.Flag{
		cli.StringFlag{
			Name:  optionRosterTomlPath + ", " + optionRosterTomlPathShort,
			Usage: "Roster toml file path",
			Value: "/medco-configuration/group.toml",
		},
		cli.StringFlag{
			Name:  optionConfigDirectoryPath + ", " + optionConfigDirectoryPathShort,
			Usage: "directory where the configurations will be stored",
			Value: "/medco-configuration",
		},
		cli.DurationFlag{
			Name:  optionBlockInterval + ", " + optionBlockIntervalShort,
			Usage: "block interval for the medchain ledger",
			Value: 5 * time.Second,
		},
	}

	cliApp.Commands = []cli.Command{
		{
			Name:   "create",
			Usage:  "Creates a new byzcoin configuration for setting up Medchain",
			Action: nonInteractiveSetup,
			Flags:  nonInteractiveSetupFlags,
		},
	}

	err := cliApp.Run(os.Args)
	log.ErrFatal(err)
}

func nonInteractiveSetup(c *cli.Context) error {
	// cli args
	lib.ConfigPath = c.String(optionConfigDirectoryPath)

	rosterPath := c.String(optionRosterTomlPath)
	if rosterPath == "" {
		return xerrors.New("--rosterTomlPath is required")
	}

	r, err := lib.ReadRoster(rosterPath)
	if err != nil {
		return xerrors.Errorf("failed to read roster: %v", err)
	}

	interval := c.Duration(optionBlockInterval)

	owner := darc.NewSignerEd25519(nil, nil)

	req, err := byzcoin.DefaultGenesisMsg(byzcoin.CurrentVersion, r,
		[]string{"spawn:longTermSecret"}, owner.Identity())
	if err != nil {
		log.Error(err)
		return xerrors.Errorf("failed to spawn genesis msg: %v", err)
	}

	req.BlockInterval = interval

	cl, resp, err := byzcoin.NewLedger(req, false)
	if err != nil {
		return xerrors.Errorf("failed to create new byzcoin ledger: %v", err)
	}

	cfg := lib.Config{
		ByzCoinID:     resp.Skipblock.SkipChainID(),
		Roster:        *r,
		AdminDarc:     req.GenesisDarc,
		AdminIdentity: owner.Identity(),
	}

	configPath, err := lib.SaveConfig(cfg)
	if err != nil {
		return xerrors.Errorf("failed to save byzcoin config: %v", err)
	}

	err = lib.SaveKey(owner)
	if err != nil {
		return err
	}

	log.Infof("Created ByzCoin with ID %x.\nConfiguration saved at %s\n", cfg.ByzCoinID, configPath)

	return lib.WaitPropagation(c, cl)
}
