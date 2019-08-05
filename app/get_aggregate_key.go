package main

import (
	"encoding/base64"
	"errors"
	"go.dedis.ch/onet/v3/app"
	"go.dedis.ch/onet/v3/log"
	"gopkg.in/urfave/cli.v1"
	"os"
	"path"
)

func getAggregateKey(c *cli.Context) error {
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
	b, err := el.Roster.Aggregate.MarshalBinary()

	// write aggregate key to file
	dir, _ := path.Split(groupTomlPath)
	pathToWrite := dir + "aggregate.txt"
	fWrite, err := os.Create(pathToWrite)
	if err != nil {
		return err
	}
	defer fWrite.Close()

	_, err = fWrite.WriteString(base64.URLEncoding.EncodeToString(b))
	if err != nil {
		return err
	}

	return nil
}
