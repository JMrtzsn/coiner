// Contains common exchange structs, interfaces used as internal API for sub-packages
package model

// Candle mirrors the model used by pythons backtrader
// OPEN, HIGH, LOW, CLOSE and VOLUME
type Candle struct {
	DATE   string // dateTime 3933 format
	TS     string // unix timestamp milliseconds
	OPEN   string // float
	CLOSE  string // float
	HIGH   string // float
	LOW    string // float
	VOLUME string // float
}

func RecordsWithHeader() [][]string {
	return append([][]string{}, []string{"DATE", "TS", "OPEN", "CLOSE", "HIGH", "LOW", "VOLUME"})
}

func (d *Candle) Csv() []string {
	return []string{d.DATE, d.TS, d.OPEN, d.CLOSE, d.HIGH, d.LOW, d.VOLUME}
}
