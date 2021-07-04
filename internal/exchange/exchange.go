package exchange

import (
	"context"
	"github.com/jmrtzsn/coiner/internal/exchange/binance"
	"github.com/jmrtzsn/coiner/internal/model"
	"time"
)

// Exchange common interface for exchange should they be implemented
type Exchange interface {
	Init(ctx context.Context, key, secret string)
	CandlesByPeriod(symbol, interval string, start, end time.Time) ([]model.Candle, error)
	String() string
}

var (
	_ Exchange = &binance.Binance{}
)
