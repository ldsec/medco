package main

import (
	"github.com/BurntSushi/toml"
	"github.com/dedis/onet/log"
	"gopkg.in/urfave/cli.v1"
)

// Files is the object structure behind the files.toml
type Files struct {
	AdapterMappings   string
	I2B2              string
	SHRINE            string
	DummyToPatient    string
	PatientDimension  string
	VisitDimension    string
	ConceptDimension  string
	ModifierDimension string
	ObservationFact   string
}

func loadData(c *cli.Context) error {
	return nil
}

func convertI2B2DataModel(c *cli.Context) error {
	// data set file paths
	groupFilePath := c.String("group")
	dataFilesPath := c.String("files")
	//entryPointIdx := c.Int("entryPointIdx")

	// db settings
	/*dbHost := c.String("dbHost")
	dbPort := c.Int("dbPort")
	dbName := c.String("dbName")
	dbUser := c.String("dbUser")
	dbPassword := c.String("dbPassword")*/

	var files Files
	if _, err := toml.Decode(dataFilesPath, &files); err != nil {
		log.Fatal("Error while reading [files].toml:", err)
	}

	log.LLvl1(groupFilePath)

	return nil
}
