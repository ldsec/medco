package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/ldsec/medco-loader/loader"
	"github.com/ldsec/medco-loader/loader/genomic"
	"github.com/ldsec/medco-loader/loader/i2b2"
	_ "github.com/lib/pq"
	"go.dedis.ch/onet/v3/app"
	"go.dedis.ch/onet/v3/log"
	"gopkg.in/urfave/cli.v1"
	"os"
	"path/filepath"
)

// Loader functions
//______________________________________________________________________________________________________________________

//----------------------------------------------------------------------------------------------------------------------
//#----------------------------------------------- LOAD DATA -----------------------------------------------------------
//----------------------------------------------------------------------------------------------------------------------

func loadV0(c *cli.Context) error {

	// data set file paths
	clinicalOntologyPath := c.String("ont_clinical")
	genomicOntologyPath := c.String("ont_genomic")
	clinicalFilePath := c.String("clinical")
	genomicFilePath := c.String("genomic")
	groupFilePath := c.String("group")
	entryPointIdx := c.Int("entryPointIdx")
	sensitiveFilePath := c.String("sensitive")
	replaySize := c.Int("replay")
	outputPath := c.String("output")

	// i2b2 db settings
	i2b2DbHost := c.String("i2b2DbHost")
	i2b2DbPort := c.Int("i2b2DbPort")
	i2b2DbName := c.String("i2b2DbName")
	i2b2DbUser := c.String("i2b2DbUser")
	i2b2DbPassword := c.String("i2b2DbPassword")

	// genomic annotations db settings
	gaDbHost := c.String("gaDbHost")
	gaDbPort := c.Int("gaDbPort")
	gaDbName := c.String("gaDbName")
	gaDbUser := c.String("gaDbUser")
	gaDbPassword := c.String("gaDbPassword")

	i2b2DB := loader.DBSettings{DBhost: i2b2DbHost, DBport: i2b2DbPort, DBname: i2b2DbName, DBuser: i2b2DbUser, DBpassword: i2b2DbPassword}
	gaDB := loader.DBSettings{DBhost: gaDbHost, DBport: gaDbPort, DBname: gaDbName, DBuser: gaDbUser, DBpassword: gaDbPassword}

	// check if db connection works
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", i2b2DbHost, i2b2DbPort, i2b2DbUser, i2b2DbPassword, i2b2DbName)
	db, err := sql.Open("postgres", psqlInfo)
	err = db.Ping()
	if err != nil {
		log.Error("Error while connecting to i2b2 database", err)
		return cli.NewExitError(err, 1)
	}
	db.Close()

	psqlInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", gaDbHost, gaDbPort, gaDbUser, gaDbPassword, gaDbName)
	db, err = sql.Open("postgres", psqlInfo)
	err = db.Ping()
	if err != nil {
		log.Error("Error while connecting to genomic annotations database", err)
		return cli.NewExitError(err, 1)
	}
	db.Close()

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
	allSensitive := false
	for scanner.Scan() {
		line := scanner.Text()
		if line == "all" {
			allSensitive = true
			break
		}
		mapSensitive[line] = struct{}{}
	}

	if replaySize < 0 {
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

	err = loadergenomic.LoadGenomicData(el.Roster, entryPointIdx, fOntClinical, fOntGenomic, fClinical, fGenomic, outputPath, allSensitive, mapSensitive, i2b2DB, gaDB, false)
	if err != nil {
		log.Fatal("Error while loading client data:", err)
	}

	return nil
}

func loadV1(c *cli.Context) error {
	// data set file paths
	groupFilePath := c.String("group")
	dataFilesPath := c.String("files")
	sensitiveFilePath := c.String("sensitive")
	entryPointIdx := c.Int("entryPointIdx")
	empty := c.Bool("empty")

	// db settings
	i2b2DbHost := c.String("i2b2DbHost")
	i2b2DbPort := c.Int("i2b2DbPort")
	i2b2DbName := c.String("i2b2DbName")
	i2b2DbUser := c.String("i2b2DbUser")
	i2b2DbPassword := c.String("i2b2DbPassword")

	i2b2DB := loader.DBSettings{DBhost: i2b2DbHost, DBport: i2b2DbPort, DBname: i2b2DbName, DBuser: i2b2DbUser, DBpassword: i2b2DbPassword}

	// check if db connection works
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", i2b2DbHost, i2b2DbPort, i2b2DbUser, i2b2DbPassword, i2b2DbName)
	db, err := sql.Open("postgres", psqlInfo)
	err = db.Ping()
	if err != nil {
		log.Error("Error while connecting to i2b2 database", err)
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
	directory := filepath.Dir(dataFilesPath)

	// get the list of sensitiveConcepts
	f, err = os.Open(sensitiveFilePath)
	if err != nil {
		log.Error("Error while reading [sensitive].txt:", err)
		return cli.NewExitError(err, 1)
	}

	// place all sensitive attributes in map set to allow for faster search
	mapSensitive := make(map[string]struct{}, 0)
	scanner := bufio.NewScanner(f)
	allSensitive := false
	for scanner.Scan() {
		line := scanner.Text()
		if line == "all" {
			allSensitive = true
			break
		}
		mapSensitive[line] = struct{}{}
	}

	err = loaderi2b2.LoadI2B2Data(el.Roster, entryPointIdx, directory, files, allSensitive, mapSensitive, i2b2DB, empty)
	if err != nil {
		log.Error("Error while converting I2B2 data:", err)
		return cli.NewExitError(err, 1)
	}

	return nil
}
