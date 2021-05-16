package export

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

type Local struct {
	exchange string
	symbol   string
}

func newLocal(exchange, symbol string) *Local {
	return &Local{
		exchange: exchange,
		symbol: symbol,
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

func (l Local) Export(input *os.File, date string) error {
	output, err := func() (*os.File, error) {
		dir := dirPath(l.exchange, l.symbol)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err = os.MkdirAll(dir, os.ModePerm); err != nil {
				return nil, err
			}
		}
		output, err := os.Create(dir + fmt.Sprintf("%s.csv", date))
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
