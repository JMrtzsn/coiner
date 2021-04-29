package local

import (
	"encoding/csv"
	"fmt"
	"github.com/jmrtzsn/coiner/internal/projectpath"
	"io"
	"os"
)

var projectID = "YOUR_PROJECT_ID"

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
func Write(input *os.File, symbol, name string) error {
	// Todo move out path / folder creation?
	path := CreateDirpath(symbol)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}
	output, err := os.Create(CreateFilepath(symbol, name))
	if err != nil {
		return err
	}
	defer output.Close()

	_, err = io.Copy(output, input)
	if err != nil {
		return err
	}
	return nil
}

//TODO: private
func CreateDirpath(symbol string) string {
	// TODO filepath.Join()
	return fmt.Sprintf("%s/data/%s/", projectpath.Root, symbol)
}

//TODO: private
func CreateFilepath(symbol, name string) string {
	// TODO filepath.Join()
	return fmt.Sprintf("%s%s.csv", CreateDirpath(symbol), name)
}
