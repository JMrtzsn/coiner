package export

import (
	"encoding/csv"
	"io/ioutil"
	"os"
)

// TODO: create export ... interface? or struct a common gateway to all exports
// TODO should contain common functions used..? by downstream packages? or have those packages further down?
type Export interface {
	Export()
	Read()
}

// CreateTempCSV create a temp CSV file from records
func CreateTempCSV(records [][]string) (*os.File, error) {
	file, err := ioutil.TempFile("", "file")
	if err != nil {
		return nil, err
	}
	w := csv.NewWriter(file)
	err = w.WriteAll(records)

	// Seek the pointer to the beginning
	file.Seek(0, 0)
	return file, nil
}
