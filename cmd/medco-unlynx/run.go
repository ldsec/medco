package main

import (
	"fmt"
	"github.com/ldsec/medco/unlynx/services"
	"github.com/ldsec/unlynx/lib"
	"github.com/urfave/cli"
	"go.dedis.ch/kyber/v3/util/key"
	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/app"
	"go.dedis.ch/onet/v3/log"
	"os"
)

func runUnLynx(c *cli.Context) {
	el, err := openGroupToml("/medco-configuration/group.toml")
	log.ErrFatal(err, "Could not open group toml.")

	client := servicesmedco.NewMedCoClient(el.List[0], "client_test")

	keys := key.NewKeyPair(libunlynx.SuiTe)

	targetData := make(libunlynx.CipherVector, 0)
	targetData = append(targetData, *libunlynx.EncryptInt(el.Aggregate, int64(1)))


	_, res, _,err := client.SendSurveyKSRequest(el, "test", keys.Public, targetData, false)
	log.ErrFatal(err)
	log.Lvl1("PASSEIII", err, res)
}

func openGroupToml(tomlFileName string) (*onet.Roster, error) {
	f, err := os.Open(tomlFileName)
	if err != nil {
		return nil, err
	}
	el, err := app.ReadGroupDescToml(f)
	if err != nil {
		return nil, err
	}

	if len(el.Roster.List) <= 0 {
		return nil, fmt.Errorf("empty or invalid unlynx group file: %v", tomlFileName)
	}

	return el.Roster, nil
}