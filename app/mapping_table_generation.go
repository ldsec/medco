package main

import (
	"errors"
	"github.com/urfave/cli"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/suites"
	"go.dedis.ch/onet/v3/log"
	"os"
	"strconv"
)

func mappingTableGenFromApp(c *cli.Context) error {

	var PointToInt = make(map[string]int64, 0)
	var suite = suites.MustFind("Ed25519")

	// cli arguments
	outputFile := c.String("outputFile")     // mandatory
	outputFormat := c.String("outputFormat") // typescript
	nbMappings := c.Int64("nbMappings")      // optional default to 1000
	checkNeg := c.Bool("checkNeg")           // optional default to false

	var Bi kyber.Point
	B := suite.Point().Base()
	var m int64

	// generate mapping in memory
	for Bi, m = suite.Point().Null(), 0; m < nbMappings; Bi, m = Bi.Add(Bi, B), m+1 {
		PointToInt[Bi.String()] = m

		if checkNeg {
			neg := suite.Point().Mul(suite.Scalar().SetInt64(int64(-m)), B)
			PointToInt[neg.String()] = -m
		}
	}

	// open file
	file, err := os.Create(outputFile)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	defer file.Close()

	// write mapping to disk
	switch outputFormat {
	case "typescript":
		err = writeMapToTSFile(file, PointToInt)
	case "go":
		err = writeMapToGoFile(file, PointToInt)
	default:
		err = errors.New("format selected is incorrect: " + outputFormat)
	}

	if err != nil {
		return cli.NewExitError(err, 2)
	}

	log.Info("Successfully generated mapping file with ", nbMappings, " mappings to ", outputFile)
	return nil
}

func writeMapToTSFile(file *os.File, pointToInt map[string]int64) (err error) {
	//export let PointToInt: Record<string, number> = {
	//	"edc876d6831fd2105d0b4389ca2e283166469289146e2ce06faefe98b22548df": 5,
	//	"f47e49f9d07ad2c1606b4d94067c41f9777d4ffda709b71da1d88628fce34d85": 6,
	//}

	_, err = file.WriteString("export let PointToInt: Record<string, number> = {\n")
	if err != nil {
		return
	}
	for k, v := range pointToInt {
		_, err = file.WriteString("\t" + `"` + k + `": ` + strconv.FormatInt(v, 10) + ",\n")
		if err != nil {
			return
		}
	}
	_, err = file.WriteString("};")
	if err != nil {
		return
	}
	return
}

func writeMapToGoFile(file *os.File, pointToInt map[string]int64) (err error) {
	//package main
	//var PointToInt = map[string]int64{
	//	"00022ddff3737fda59ef096dae2ea2876a5893510442fde25cb37486ed8b97c3": 7414,
	//}

	_, err = file.WriteString("package main \nvar PointToInt = map[string]int64{\n")
	if err != nil {
		return
	}

	for k, v := range pointToInt {
		_, err = file.WriteString("\t" + `"` + k + `": ` + strconv.FormatInt(v, 10) + ",\n")
		if err != nil {
			return
		}
	}

	_, err = file.WriteString("}")
	if err != nil {
		return
	}
	return
}
