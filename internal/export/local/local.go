package local

import (
	"encoding/csv"
	"fmt"
	"github.com/jmrtzsn/coiner/internal/projectpath"
	"os"
)

func Read(file string) ([][]string, error) {
	// TODO read from byte?
	csvfile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer csvfile.Close()
	records, err := csv.NewReader(csvfile).ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}

// TODO either use to create a "temp" local file, or change structure to create a temp then copy
func Write(records [][]string, symbol, name string) error {
	path := CreateDirpath(symbol)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}

	// TODO: move this logic outside (file is created first then exported)
	file, err := os.Create(CreateFilepath(symbol, name))
	if err != nil {
		return err
	}
	defer file.Close()

	// Todo pass writer?
	w := csv.NewWriter(file)
	err = w.WriteAll(records)

	return nil
}

//TODO: private
func CreateDirpath(symbol string) string {
	return fmt.Sprintf("%s/data/%s/", projectpath.Root, symbol)
}

//TODO: private
func CreateFilepath(symbol, name string) string {
	return fmt.Sprintf("%s%s.csv", CreateDirpath(symbol), name)
}
