package main

import (
	"fmt"

	"github.com/ldsec/unlynx/lib"
	"github.com/urfave/cli"
	"go.dedis.ch/kyber/v3/util/encoding"
	"go.dedis.ch/kyber/v3/util/key"
	"go.dedis.ch/onet/v3/app"
	"go.dedis.ch/onet/v3/log"
	"go.dedis.ch/onet/v3/network"
)

// NonInteractiveSetup is used to setup the cothority node for unlynx in a non-interactive way (and without error checks)
func NonInteractiveSetup(c *cli.Context) error {

	// cli arguments
	serverBindingStr := c.String(optionServerBinding)
	description := c.String(optionDescription)
	privateTomlPath := c.String(optionPrivateTomlPath)
	publicTomlPath := c.String(optionPublicTomlPath)
	wsUrlStr := c.String(optionWsUrl)

	// provided keys (optional)
	providedPubKey := c.String(optionProvidedPubKey)
	providedPrivKey := c.String(optionProvidedPrivKey)

	if serverBindingStr == "" || description == "" || privateTomlPath == "" || publicTomlPath == "" {
		err := fmt.Errorf("arguments not OK")
		log.Error(err)
		return cli.NewExitError(err, 3)
	}

	var privStr, pubStr string
	var err error
	if providedPubKey != "" {
		privStr = providedPrivKey
		pubStr = providedPubKey
	} else {
		kp := key.NewKeyPair(libunlynx.SuiTe)

		privStr, err = encoding.ScalarToStringHex(libunlynx.SuiTe, kp.Private)
		if err != nil {
			log.Error("failed to convert scalar to hexadecimal")
			return cli.NewExitError(err, 3)
		}
		pubStr, err = encoding.PointToStringHex(libunlynx.SuiTe, kp.Public)
		if err != nil {
			log.Error("failed to convert point to hexadecimal")
			return cli.NewExitError(err, 3)
		}
	}

	public, err := encoding.StringHexToPoint(libunlynx.SuiTe, pubStr)
	if err != nil {
		log.Error("failed to convert hexadecimal to point")
		return cli.NewExitError(err, 3)
	}

	serverBinding := network.NewTLSAddress(serverBindingStr)
	services := app.GenerateServiceKeyPairs()

	conf := &app.CothorityConfig{
		Suite:       libunlynx.SuiTe.String(),
		Public:      pubStr,
		Private:     privStr,
		Address:     serverBinding,
		Services:    services,
		Description: description,
		URL:         wsUrlStr,
	}

	server := app.NewServerToml(libunlynx.SuiTe, public, serverBinding, conf.Description, services)
	server.URL = conf.URL

	group := app.NewGroupToml(server)

	if err := conf.Save(privateTomlPath); err != nil {
		err := fmt.Errorf("failed saving private.toml")
		log.Error(err)
		return cli.NewExitError(err, 3)
	}
	if err := group.Save(publicTomlPath); err != nil {
		err := fmt.Errorf("failed saving group.toml")
		log.Error(err)
		return cli.NewExitError(err, 3)
	}

	return nil
}
