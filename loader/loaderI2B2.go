package loader

import (
	"encoding/xml"
	"os"
	"gopkg.in/dedis/onet.v1/log"
	"io/ioutil"
	"io"
)

// ListConceptsPaths list all the sensitive concepts (paths)
var ListConceptsPaths []string

// The different paths and handlers for all the file both for input and/or output
var (
	InputFilePaths = map[string]string{
		"ADAPTER_MAPPINGS": "../data/original/AdapterMappings.xml",
	}

	OutputFilePaths = map[string]string{
		"ADAPTER_MAPPINGS": "../data/converted/AdapterMappings.xml",
	}
)

const (
	// A generic XML header suitable for use with the output of Marshal.
	// This is not automatically added to any output of this package,
	// it is provided as a convenience.
	Header = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n"
)

// ConvertAdapterMappings converts the old AdapterMappings.xml file. This file maps a shrine concept code to an i2b2 concept code
func ConvertAdapterMappings() error{
	xmlInputFile, err := os.Open(InputFilePaths["ADAPTER_MAPPINGS"])
	if err != nil {
		log.Fatal("Error opening [AdapterMappings].xml xmlInputFile")
		return err
	}
	defer xmlInputFile.Close()

	b, _ := ioutil.ReadAll(xmlInputFile)

	var am AdapterMappings

	xml.Unmarshal(b, &am)

	// filter out sensitive entries
	numElementsDel := FilterSensitiveEntries(&am)
	log.Lvl2(numElementsDel,"entries deleted")

	xmlOutputFile, err := os.Create(OutputFilePaths["ADAPTER_MAPPINGS"])
	if err != nil {
		log.Fatal("Error creating converted [AdapterMappings].xml xmlOutputFile")
		return err
	}
	xmlOutputFile.Write([]byte(Header))

	xmlWriter := io.Writer(xmlOutputFile)

	enc := xml.NewEncoder(xmlWriter)
	enc.Indent("", "\t")
	err = enc.Encode(am)
	if err != nil {
		log.Fatal("Error writing converted [AdapterMappings].xml xmlOutputFile")
	}
	return nil
}


// FilterSensitiveEntries filters out (removes) the <key>, <values> pair(s) that belong to sensitive concepts
func FilterSensitiveEntries(am *AdapterMappings) int{
	m := am.ListEntries

	deleted := 0
	for i := range m {
		j := i - deleted
		if containsArrayString(ListConceptsPaths, m[j].Key){
			m = m[:j+copy(m[j:], m[j+1:])]
			deleted++
		}
	}

	return deleted
}

func containsArrayString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
