package main

import (
	"fmt"
	"github.com/ldsec/unlynx/lib"
	"github.com/urfave/cli"
	"go.dedis.ch/onet/v3/log"
	"io"
	"os"
	"strconv"
)

func decryptIntFromApp(c *cli.Context) error {

	// cli arguments
	secKeySerialized := c.String("key")
	secKey, err := libunlynx.DeserializeScalar(secKeySerialized)
	if err != nil {
		log.Error(err)
		return cli.NewExitError(err, 4)
	}

	if c.NArg() != 1 {
		err := fmt.Errorf("wrong number of arguments (only 1 allowed, except for the flags)")
		log.Error(err)
		return cli.NewExitError(err, 3)
	}

	// value to decrypt
	toDecryptSerialized := c.Args().Get(0)
	toDecrypt, err := libunlynx.NewCipherTextFromBase64(toDecryptSerialized)
	if err != nil {
		return err
	}

	// decryption
	decVal := libunlynx.DecryptInt(secKey, *toDecrypt)

	// output in xml format on stdout
	resultString := "<decrypted>" + strconv.FormatInt(decVal, 10) + "</decrypted>\n"
	_, err = io.WriteString(os.Stdout, resultString)
	if err != nil {
		log.Error("Error while writing result.", err)
		return cli.NewExitError(err, 4)
	}

	return nil
}
