package utilclient

import (
	"encoding/csv"
	"fmt"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

const comma rune = ','
const useCarriageReturnThenLineFeed bool = false
const filePermission os.FileMode = 777

// ResultCSV holds a file descriptor and abstract create/write/close function
type ResultCSV struct {
	path             string
	lock             *sync.Mutex
	isStandardOutput bool
	file             *os.File
	csvWriter        *csv.Writer
}

// NewCSV creates ResultCSV instance. It creates a file at provided path if the file does not exists yet, otherwise truncates it.
// Empty path makes the ResultCSV instance printing directly to standard output.
func NewCSV(filePath string) (*ResultCSV, error) {

	var file *os.File
	var isStdout bool
	var err error
	if filePath == "" {
		logrus.Info("No output file provided, dumping directly to stdout")
		isStdout = true
		file = os.Stdout
	} else {
		logrus.Infof("Opening output file %s", filePath)

		file, err = os.Create(filePath)
		if err != nil {
			logrus.Error("error opening file: ", err)
			return nil, err
		}
		err = os.Chmod(filePath, filePermission)
		if err != nil {
			logrus.Error("error setting permissions on file: ", err)
			return nil, err
		}
		isStdout = false
	}

	csvWriter := csv.NewWriter(file)
	csvWriter.Comma = comma
	csvWriter.UseCRLF = useCarriageReturnThenLineFeed
	return &ResultCSV{
		path:             filePath,
		lock:             &sync.Mutex{},
		file:             file,
		csvWriter:        csvWriter,
		isStandardOutput: isStdout,
	}, nil

}

// Write writes a single CSV record. A record is a slice of strings.
func (csv *ResultCSV) Write(record []string) error {
	csv.lock.Lock()
	defer csv.lock.Unlock()
	err := csv.csvWriter.Write(record)
	if err != nil {
		return fmt.Errorf("while writing record %v: %s", record, err.Error())
	}
	return nil
}

// WriteAll writes multiple CSV records.
func (csv *ResultCSV) WriteAll(records [][]string) error {
	csv.lock.Lock()
	defer csv.lock.Unlock()
	err := csv.csvWriter.WriteAll(records)
	if err != nil {
		return fmt.Errorf("while writing records %v: %s", records, err.Error())
	}
	return nil
}

// Flush writes buffered data to the output file
func (csv *ResultCSV) Flush() error {
	csv.lock.Lock()
	defer csv.lock.Unlock()
	csv.csvWriter.Flush()
	err := csv.csvWriter.Error()
	if err != nil {
		return fmt.Errorf("while flushing profiling CSV writer: %s", err.Error())
	}
	return nil
}

// Close closes output file if it is not standard output, otherwise does nothing
func (csv *ResultCSV) Close() error {
	if csv.isStandardOutput {
		return nil
	}
	csv.lock.Lock()
	defer csv.lock.Unlock()
	err := csv.file.Close()
	if err != nil {
		return fmt.Errorf("while closing profiling csv file: %s", err.Error())
	}
	return nil
}
