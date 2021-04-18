// Contains common exchange structs, interfaces used as internal API for sub-packages
package internal

// OHLCV is the main communications struct for downloading historical data
// OPEN, HIGH, LOW, CLOSE and VOLUME
type OHLCV struct {
	DATE   string // dateTime 3933 format
	TS     string // unix timestamp milliseconds
	OPEN   string // float
	CLOSE  string // float
	HIGH   string // float
	LOW    string // float
	VOLUME string // float
}

// ToCSV converts an OHLCV slice to CSV records
func ToCSV(data []OHLCV) [][]string {
	var csvs [][]string
	csvs = append(csvs, []string{"DATE", "TS", "OPEN", "CLOSE", "HIGH", "LOW", "VOLUME"})
	for _, row := range data {
		csvs = append(csvs, row.csv())
	}
	return csvs
}

func (d *OHLCV) csv() []string {
	return []string{d.DATE, d.TS, d.OPEN, d.CLOSE, d.HIGH, d.LOW, d.VOLUME}
}


