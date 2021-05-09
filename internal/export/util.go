package export

import (
	"encoding/csv"
	"fmt"
	"github.com/jmrtzsn/coiner/internal/projectpath"
	"io/ioutil"
	"os"
)

// CreateTempCSV create a temp CSV file from records
func CreateTempCSV(records [][]string) (*os.File, error) {
	file, err := ioutil.TempFile("", "file")
	if err != nil {
		return nil, err
	}
	w := csv.NewWriter(file)
	err = w.WriteAll(records)

	// Seek the pointer to the beginning
	// TODO return buffer? function that resets file when reading
	file.Seek(0, 0)
	return file, nil
}

func dirPath(symbol string) string {
	// TODO filepath.Join()
	return fmt.Sprintf("%s/data/%s/", projectpath.Root, symbol)
}

// storagePath generates a folder/file.csv
func storagePath(symbol, name string) string {
	return fmt.Sprintf("%s/%s.csv", symbol, name)
}
