package main

import (
	// Empty imports to have the init-functions called which should
	// register the protocol
	_ "github.com/ldsec/medco/unlynx/services"
	_ "github.com/ldsec/unlynx/protocols"
	"github.com/urfave/cli"
	_ "go.dedis.ch/cothority/v3/byzcoin"
	_ "go.dedis.ch/cothority/v3/byzcoin/contracts"
	_ "go.dedis.ch/cothority/v3/skipchain"
	"go.dedis.ch/onet/v3/app"
)

func runServer(ctx *cli.Context) error {
	// first check the options
	config := ctx.String("config")
	app.RunServer(config)
	return nil
}
