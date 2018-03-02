package main

import (
	"gopkg.in/urfave/cli.v1"

	// Empty imports to have the init-functions called which should
	// register the protocol
	"github.com/dedis/onet/app"
	_ "github.com/lca1/medco/services"
	_ "github.com/lca1/unlynx/protocols"
)

func runServer(ctx *cli.Context) error {
	// first check the options
	config := ctx.String("config")

	app.RunServer(config)

	return nil
}
