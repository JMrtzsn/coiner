// Contains common exchange structs, interfaces
package internal

import (
	"encoding/csv"
	"fmt"
	"github.com/jmrtzsn/coiner/internal/projectpath"
	"os"
)

type Config interface {
	LoadEnv()
}

// OPEN, HIGH, LOW, CLOSE and VOLUME
type OHLCV struct {
	DATE   string // dateTime
	TS     string // Timestamp
	OPEN   string // float
	CLOSE  string // float
	HIGH   string // float
	LOW    string // float
	VOLUME string // float
}

// ToCSV converts to CSV format
func toCSV(data []OHLCV) []string {
	var csvs []string
	csvs = append(csvs, "DATE,TS,OPEN,CLOSE,HIGH,LOW,VOLUME")
	for _, row := range data {
		csvs = append(csvs, row.csv())
	}
	return csvs
}

func (d *OHLCV) csv() string {
	return fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s",
		d.DATE, d.TS, d.OPEN, d.CLOSE, d.HIGH, d.LOW, d.VOLUME)
}

func LocalOutput(records [][]string, name string) error {
	file, err := os.Create(getFilepath(name))
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	err = w.WriteAll(records) // calls Flush internally

	return nil
}

func getFilepath(name string) string {
	return fmt.Sprintf("%s/data/%s.csv", projectpath.Root, name)
}
