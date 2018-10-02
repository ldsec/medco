package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/dedis/onet/app"
	"github.com/dedis/onet/log"
	"github.com/lca1/medco-loader/loader"
	"github.com/lca1/medco-loader/loader/genomic"
	"github.com/lca1/medco-loader/loader/i2b2"
	"gopkg.in/urfave/cli.v1"
	"os"
)

// Loader functions
//______________________________________________________________________________________________________________________

//----------------------------------------------------------------------------------------------------------------------
//#----------------------------------------------- LOAD DATA -----------------------------------------------------------
//----------------------------------------------------------------------------------------------------------------------

func loadData(c *cli.Context) error {

	// data set file paths
	clinicalOntologyPath := c.String("ont_clinical")
	genomicOntologyPath := c.String("ont_genomic")
	clinicalFilePath := c.String("clinical")
	genomicFilePath := c.String("genomic")
	groupFilePath := c.String("file")
	entryPointIdx := c.Int("entryPointIdx")
	sensitiveFilePath := c.String("sensitive")
	replaySize := c.Int("replay")

	// db settings
	dbHost := c.String("dbHost")
	dbPort := c.Int("dbPort")
	dbName := c.String("dbName")
	dbUser := c.String("dbUser")
	dbPassword := c.String("dbPassword")

	databaseS := loader.DBSettings{DBhost: dbHost, DBport: dbPort, DBname: dbName, DBuser: dbUser, DBpassword: dbPassword}

	// check if db connection works
	/*psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Error("Error while opening database", err)
		return cli.NewExitError(err, 1)
	}
	db.Close()*/

	// generate el with group file
	f, err := os.Open(groupFilePath)
	if err != nil {
		log.Error("Error while opening group file", err)
		return cli.NewExitError(err, 1)
	}
	el, err := app.ReadGroupDescToml(f)
	if err != nil {
		log.Error("Error while reading group file", err)
		return cli.NewExitError(err, 1)
	}
	if len(el.Roster.List) <= 0 {
		log.Error("Empty or invalid group file", err)
		return cli.NewExitError(err, 1)
	}

	fOntClinical, err := os.Open(clinicalOntologyPath)
	if err != nil {
		log.Error("Error while opening the clinical ontology file", err)
		return cli.NewExitError(err, 1)
	}

	fOntGenomic, err := os.Open(genomicOntologyPath)
	if err != nil {
		log.Error("Error while opening the genomic ontology file", err)
		return cli.NewExitError(err, 1)
	}

	fClinical, err := os.Open(clinicalFilePath)
	if err != nil {
		log.Error("Error while opening the clinical file", err)
		return cli.NewExitError(err, 1)
	}

	fGenomic, err := os.Open(genomicFilePath)
	if err != nil {
		log.Error("Error while opening the genomic file", err)
		return cli.NewExitError(err, 1)
	}

	// get the list of sensitiveConcepts
	f, err = os.Open(sensitiveFilePath)
	if err != nil {
		log.Error("Error while reading [sensitive].txt:", err)
		return cli.NewExitError(err, 1)
	}

	// place all sensitive attributes in map set to allow for faster search
	mapSensitive := make(map[string]struct{}, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		mapSensitive[line] = struct{}{}
	}

	if replaySize < 1 {
		log.Error("Wrong file size value (1>)", err)
		return cli.NewExitError(err, 1)
	} else if replaySize > 1 {
		fGenomic.Close()
		loadergenomic.ReplayDataset(genomicFilePath, replaySize)

		fGenomic, err = os.Open(genomicFilePath)
		if err != nil {
			log.Error("Error while opening the new genomic file", err)
			return cli.NewExitError(err, 1)
		}
	}

	err = loadergenomic.LoadClient(el.Roster, entryPointIdx, fOntClinical, fOntGenomic, fClinical, fGenomic, mapSensitive, databaseS, false)
	if err != nil {
		log.Fatal("Error while loading client data:", err)
	}

	return nil
}

func convertI2B2DataModel(c *cli.Context) error {
	// data set file paths
	groupFilePath := c.String("group")
	dataFilesPath := c.String("files")
	sensitiveFilePath := c.String("sensitive")
	entryPointIdx := c.Int("entryPointIdx")
	empty := c.Bool("empty")

	// db settings
	dbHost := c.String("dbHost")
	dbPort := c.Int("dbPort")
	dbName := c.String("dbName")
	dbUser := c.String("dbUser")
	dbPassword := c.String("dbPassword")

	databaseS := loader.DBSettings{DBhost: dbHost, DBport: dbPort, DBname: dbName, DBuser: dbUser, DBpassword: dbPassword}

	// check if db connection works
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Error("Error while opening database", err)
		return cli.NewExitError(err, 1)
	}
	db.Close()

	// generate el with group file
	f, err := os.Open(groupFilePath)
	if err != nil {
		log.Error("Error while opening group file:", err)
		return cli.NewExitError(err, 1)
	}
	el, err := app.ReadGroupDescToml(f)
	if err != nil {
		log.Error("Error while reading group file:", err)
		return cli.NewExitError(err, 1)
	}
	if len(el.Roster.List) <= 0 {
		log.Error("Empty or invalid group file:", err)
		return cli.NewExitError(err, 1)
	}
	f.Close()

	// get all files to convert
	var files loaderi2b2.Files
	if _, err := toml.DecodeFile(dataFilesPath, &files); err != nil {
		log.Error("Error while reading [files].toml:", err)
		return cli.NewExitError(err, 1)
	}

	// get the list of sensitiveConcepts
	f, err = os.Open(sensitiveFilePath)
	if err != nil {
		log.Error("Error while reading [sensitive].txt:", err)
		return cli.NewExitError(err, 1)
	}

	// place all sensitive attributes in map set to allow for faster search
	mapSensitive := make(map[string]bool, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		mapSensitive[line] = true
	}

	loaderi2b2.ConvertI2B2(el.Roster, entryPointIdx, files, mapSensitive, databaseS, empty)
	if err != nil {
		log.Error("Error while converting I2B2 data:", err)
		return cli.NewExitError(err, 1)
	}

	return nil
}
