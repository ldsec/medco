package main

import (
	servicesmedco "github.com/ldsec/medco-unlynx/services"
	"time"

	// Empty imports to have the init-functions called which should
	// register the protocol
	_ "github.com/ldsec/medco-unlynx/services"
	_ "github.com/ldsec/unlynx/protocols"
	"github.com/urfave/cli"
	"go.dedis.ch/onet/v3/app"
)

func runServer(ctx *cli.Context) error {
	// first check the options
	config := ctx.String("config")
	timeout := ctx.Int64("timeout")

	app.RunServer(config)

	servicesmedco.TimeoutService = time.Duration(timeout) * time.Minute
	return nil
}
