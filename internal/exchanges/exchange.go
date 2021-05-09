package exchanges

import (
	"github.com/jmrtzsn/coiner/internal/exchanges/binance"
	"github.com/jmrtzsn/coiner/internal/model"
)

// Exchange common interface for exchanges should they be implemented
type Exchange interface {
	// Init ENV stuff
	Init(apiKey, apiSecret string)
	// OHLCV downloads and converts whatever the API exposes to a uniform format
	OHLCV(symbol string, interval, start, end string) ([]model.OHLCV, error)
	// Download breaks up a longer intervals into manageable chunks
	// Download(symbol string, interval, start, end string)
}

// Ensure the handlers implement the required interfaces/types at compile time
var (
	_ Exchange = &binance.Binance{}
)