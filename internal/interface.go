// Contains common exchange structs, interfaces
package internal

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
func toCSV(data []OHLCV) [][]string {
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


