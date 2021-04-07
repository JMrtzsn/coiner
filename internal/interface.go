// Contains common exchange structs, interfaces
package internal

type Config interface {
	LoadEnv()
}

// OPEN, HIGH, LOW, CLOSE and VOLUME
type OHLCV struct {
	DATE   string // dateTime
	TS     string // Timestamp
	OPEN   string
	CLOSE  string
	HIGH   string
	LOW    string
	VOLUME string
}
