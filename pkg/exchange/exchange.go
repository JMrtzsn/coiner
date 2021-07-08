package exchange

import (
	"context"
	binance2 "github.com/jmrtzsn/coiner/pkg/exchange/binance"
	model2 "github.com/jmrtzsn/coiner/pkg/model"
	"time"
)

// Exchange common interface for exchange should they be implemented
type Exchange interface {
	Init(ctx context.Context, key, secret string)
	Candles(symbol, interval string, start, end time.Time) ([]model2.Candle, error)
	String() string
}

var (
	_ Exchange = &binance2.Binance{}
)
