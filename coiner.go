package coiner

import (
	"github.com/jmrtzsn/coiner/internal"
	"github.com/jmrtzsn/coiner/internal/exchanges/binance"
)

// Exchange is the external interface coiner uses to structure exchanges
type Exchange interface {
	// Init ENV stuff
	Init()
	// OHLCV downloads and converts whatever the API exposes to a uniform format
	OHLCV(symbol string, interval, start, end string) ([]internal.OHLCV, error)
	// Download breaks up a longer intervals into manageable chunks
	// Download(symbol string, interval, start, end string)
}

// Ensure the handlers implement the required interfaces/types at compile time
var (
	_ Exchange = &binance.Binance{}
)

// TODO Exchange, Output [GCP, Local etc], time [day, hour, minute], Symbol, date_from, date_to
func main() {
	// "Binance, BTCUSD, 1 min, date_from, date_to
	// TODO: translate args into creation of an exchange interface object, init it and download
}
