package export

import (
	"encoding/csv"
	"fmt"
	"github.com/jmrtzsn/coiner/internal/projectpath"
	"io"
	"os"
	"path/filepath"
)

type Local struct {
	exchange string
}

// NewLoca
func NewLocal(exchange string) *Local {
	return &Local{
		exchange: exchange,
	}
}

func (l Local) String() string {
	return "Local"
}

// Read from string path
func (l Local) Read(path string) ([][]string, error) {
	csvfile, err := os.Open(path)
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

// Export copies file to symbol/date folder
func (l Local) Export(csv *os.File, date, symbol string) error {
	output, err := func() (*os.File, error) {
		dir := l.DirPath(symbol)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err = os.MkdirAll(dir, os.ModePerm); err != nil {
				return nil, err
			}
		}
		output, err := os.Create(dir + fmt.Sprintf("/%s.csv", date))
		if err != nil {
			return nil, err
		}
		return output, nil
	}()
	defer output.Close()

	csv.Seek(0, 0)
	_, err = io.Copy(output, csv)
	if err != nil {
		return err
	}
	return nil
}

// DirPath generates a exchange/symbol/date.csv path for local storage
func (l Local) DirPath(symbol string) string {
	return filepath.Join(projectpath.Root, l.exchange, symbol)
}
