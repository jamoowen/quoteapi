package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/jamoowen/quoteapi/internal/csv"
)

func main() {
	logger := log.Default()
	logger.Printf("Starting csv migration")

	args := os.Args
	if len(args) != 3 {
		logger.Fatal("ERROR: Need to provide path to unordered csv && desired name of ordered csv")
	}

	// create target path to write csv to
	targetFileName := args[2]
	if !strings.Contains(targetFileName, ".csv") {
		targetFileName = targetFileName + ".csv"
	}
	wd, err := os.Getwd()
	targetPath := path.Join(wd, "/data", targetFileName)

	unorderedCSVPath := args[1]
	records, err := csv.ReadCsv(unorderedCSVPath)
	if err != nil {
		logger.Fatal(err.Error())
	}

	// remove any records with blank quote or author
	cleanedRecords := make([][]string, 0, len(records))
	for i, record := range records {
		if i == 0 {
			continue
		}
		if len(record) != 2 {
			fmt.Printf("malformed record: %v", record)
			continue
		}
		cleanedRecords = append(cleanedRecords, record)
	}

	slices.SortFunc(cleanedRecords, func(a, b []string) int {
		if n := strings.Compare(strings.ToLower(a[0]), strings.ToLower(b[0])); n != 0 {
			return n
		}
		return strings.Compare(strings.ToLower(a[1]), strings.ToLower(b[1]))
	})

	// insert header again
	cleanedRecords = slices.Insert(cleanedRecords, 0, []string{"Author", "Message"})
	err = csv.OverwriteCsv(cleanedRecords, targetPath)
	if err != nil {
		logger.Fatal(err.Error())
	}

	fmt.Printf("ordered csv successfully written to path:%v length: %v\nRecords head: %v %v\n", targetPath, len(cleanedRecords), cleanedRecords[0][0], cleanedRecords[0][1])
	// isolate the header initially
}
