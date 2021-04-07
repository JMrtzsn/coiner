package coiner

import (
	"github.com/jmrtzsn/coiner/internal"
	"github.com/jmrtzsn/coiner/internal/binance"
)

// Exchange is the external interface coiner uses to structure exchanges
type Exchange interface {
	Init()
	OHLCV(symbol string, interval, start, end string) ([]internal.OHLCV, error)
}

// Ensure the handlers implement the required interfaces/types at compile time
var (
	_ Exchange = &binance.Binance{}
)

// TODO Exchange, Output [GCP, Local etc], time [day, hour, minute], Symbol, date_from, date_to
func main() {
	// Get command line args
	// "Binance, BTCUSD, 1 min, date_from, date_to
	// TODO: translate args into creation of an exchange interface object, init it and download
}
