package main

import (
	"errors"
	"github.com/ldsec/unlynx/lib"
	"github.com/urfave/cli"
	"go.dedis.ch/onet/v3/log"
	"io"
	"os"
)

func keyGenerationFromApp(c *cli.Context) error {
	if c.NArg() != 0 {
		err := errors.New("wrong number of arguments (none allowed, except for the flags)")
		log.Error(err)
		return cli.NewExitError(err, 3)
	}

	secKey, pubKey := libunlynx.GenKey()
	secKeySer, err1 := libunlynx.SerializeScalar(secKey)
	pubKeySer, err2 := libunlynx.SerializePoint(pubKey)

	if err1 != nil {
		log.Error("Error while serializing.", err1)
		return cli.NewExitError(err1, 4)
	}
	if err2 != nil {
		log.Error("Error while serializing.", err2)
		return cli.NewExitError(err2, 4)
	}

	// output in xml format on stdout
	resultString := "<key_pair><public>" + pubKeySer + "</public><private>" + secKeySer + "</private></key_pair>\n"
	_, err := io.WriteString(os.Stdout, resultString)
	if err != nil {
		log.Error("Error while writing result.", err)
		return cli.NewExitError(err, 4)
	}

	return nil
}
