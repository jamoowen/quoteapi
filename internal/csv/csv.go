package csv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
)

//
// read csv?
// write new csv?
// overwrite csv?

func ReadCsv(path string) ([][]string, error) {
	// need file first
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("ERROR: failed to open file: %v", err))
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("ERRORe failed to read csv: %v", err))
	}
	return records, nil
}

func OverwriteCsv(data [][]string, path string) error {
	// create will truncate the file in that path
	fileToWriteTo, err := os.Create(path)
	if err != nil {
		return errors.New(fmt.Sprintf("ERROR: failed to create file: %v", err))
	}
	writer := csv.NewWriter(fileToWriteTo)
	err = writer.WriteAll(data)

	if err != nil {
		return errors.New(fmt.Sprintf("ERROR: failed to write new csv: %v", err))
	}
	return nil
}
