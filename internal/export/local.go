package export

import (
	"encoding/csv"
	"io"
	"os"
)

type Local struct {
	symbol string
	date   string
}

func newLocal(symbol, date string) *Local {
	return &Local{
		symbol: symbol,
		date:   date,
	}
}

func (l Local) Read(file string) ([][]string, error) {
	csvfile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer csvfile.Close()
	// TODO read safely
	records, err := csv.NewReader(csvfile).ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (l Local) Export(input *os.File) error {
	output, err := func() (*os.File, error) {
		dir := dirpath(l.symbol)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err = os.MkdirAll(dir, os.ModePerm); err != nil {
				return nil, err
			}
		}
		output, err := os.Create(filepath(l.symbol, l.date))
		if err != nil {
			return nil, err
		}
		return output, nil
	}()
	defer output.Close()

	input.Seek(0, 0)
	_, err = io.Copy(output, input)
	if err != nil {
		return err
	}
	return nil
}
