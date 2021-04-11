package export

import (
	"encoding/csv"
	"fmt"
	"github.com/jmrtzsn/coiner/internal/projectpath"
	"log"
	"os"
)

func OpenCSV(symbol, name string) [][]string {
	csvfile, err := os.Open(createFilepath(symbol, name))
	defer csvfile.Close()

	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	r := csv.NewReader(csvfile)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatalln("Couldn't read the csv file", err)
	}
	return records
}

func SaveCSV(records [][]string, symbol, name string) error {
	path := createDirpath(symbol)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, os.ModePerm); err != nil{
			return err
		}
	}
	file, err := os.Create(createFilepath(symbol, name))
	if err != nil {
		return err
	}
	defer file.Close()
	// w implements io.Writer.
	w := csv.NewWriter(file)
	err = w.WriteAll(records)

	return nil
}
func createDirpath(symbol string) string {
	return fmt.Sprintf("%s/data/%s/", projectpath.Root, symbol)
}

func createFilepath(symbol, name string) string {
	return fmt.Sprintf("%s%s.csv", createDirpath(symbol), name)
}

