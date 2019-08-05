package main

import (
	"errors"
	"github.com/lca1/medco-unlynx/services"
	libunlynx "github.com/lca1/unlynx/lib"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/onet/v3/app"
	"go.dedis.ch/onet/v3/log"
	"gopkg.in/urfave/cli.v1"
	"os"
	"path"
	"strconv"
	"strings"
)

func generateTaggingSecrets(c *cli.Context) error {
	// cli arguments
	groupTomlPath := c.String("file")
	providedSecretsString := c.String("secrets")
	nodeIndex := c.Int("nodeIndex")

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

	// parse provided secrets
	var providedSecrets []kyber.Scalar
	if providedSecretsString != "" {

		providedSecretsStringSplit := strings.Split(providedSecretsString, ",")
		if len(providedSecretsStringSplit) != len(el.Roster.List) {
			err := errors.New("provided secrets list does not match the length of the roster list")
			log.Error(err, len(providedSecretsStringSplit), " != ", len(el.Roster.List))
			return err
		}

		for _, s := range providedSecretsStringSplit {
			providedSecretScalar, err := libunlynx.DeserializeScalar(s)
			if err != nil {
				log.Error(err)
				return err
			}

			providedSecrets = append(providedSecrets, providedSecretScalar)
		}
	}

	// setup secrets
	dir, _ := path.Split(groupTomlPath)

	for i, dest := range el.Roster.List {
		var err error
		if len(providedSecrets) > 0 {
			_, err = servicesmedco.CheckDDTSecrets(
				dir+"srv"+strconv.FormatInt(int64(nodeIndex), 10)+"-ddtsecrets.toml",
				dest.Address,
				providedSecrets[i])
		} else {
			_, err = servicesmedco.CheckDDTSecrets(
				dir+"srv"+strconv.FormatInt(int64(nodeIndex), 10)+"-ddtsecrets.toml",
				dest.Address,
				nil)
		}

		if err != nil {
			return err
		}
	}

	return nil
}
