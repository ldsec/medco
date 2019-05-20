package main

import (
	"errors"
	"github.com/lca1/medco-unlynx/services"
	"go.dedis.ch/onet/v3/app"
	"go.dedis.ch/onet/v3/log"
	"gopkg.in/urfave/cli.v1"
	"os"
	"path"
	"strconv"
)

func generateTaggingSecrets(c *cli.Context) error {
	// cli arguments
	groupTomlPath := c.String("file")

	if groupTomlPath == "" {
		err := errors.New("arguments not OK")
		log.Error(err)
		return cli.NewExitError(err, 3)
	}

	fRead, err := os.Open(groupTomlPath)
	if err != nil {
		log.Error("Error while opening group file", err)
		return err
	}
	defer fRead.Close()

	el, err := app.ReadGroupDescToml(fRead)
	if err != nil {
		log.Error("Error while reading group file", err)
		return err
	}
	if len(el.Roster.List) <= 0 {
		log.Error("Empty or invalid group file", err)
		return err
	}

	for i := range el.Roster.List {
		dir, _ := path.Split(groupTomlPath)

		for _, dest := range el.Roster.List {
			_, err := servicesmedco.CheckDDTSecrets(dir+"srv"+strconv.FormatInt(int64(i), 10)+"-ddtsecrets.toml", dest.Address)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
