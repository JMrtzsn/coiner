package export

import (
	"encoding/csv"
	"io/ioutil"
	"os"
)

// WriteToTempFile create a temp CSV file from records
func WriteToTempFile(records [][]string) (*os.File, error) {
	file, err := ioutil.TempFile("", "file")
	if err != nil {
		return nil, err
	}
	w := csv.NewWriter(file)
	err = w.WriteAll(records)

	// Seek the pointer to the beginning
	// TODO return buffer? function that resets file when reading
	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// TODO Verify integrity of data function(s)
// LOCAL, BUCKET
// Should loop through a symbol and log all data points where data is complete
