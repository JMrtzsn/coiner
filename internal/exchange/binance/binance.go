package binance

import (
	"context"
	"github.com/adshao/go-binance/v2"
	"github.com/jmrtzsn/coiner/internal/model"
	"strconv"
	"time"
)

type Binance struct {
	client binance.Client
	ctx        context.Context
}

// Init loads env variables and creates a binance rest client
func (e *Binance) Init(ctx context.Context, key, secret string) {
	e.client = *binance.NewClient(key, secret)
	e.ctx = ctx
}

// CandlesByPeriod validats and parses inputs.
// Returns standard trade data, according to the Open, High, Low, Volume Format
// symbol: BUYSELL = BTCUSD
// interval: 1d, 1h, 15m, 1m
// start: datetime ISO RFC3339 - "2020-04-04 T12:07:00"
// end: UNIX datetime - 1499040000000
// limit: rows returned - 10
func (e *Binance) 	CandlesByPeriod(symbol, interval string, start, end time.Time) ([]model.OHLCV, error) {
	candles, err := e.client.
		NewKlinesService().
		Symbol(symbol).
		Interval(interval).
		StartTime(start.Unix()).
		EndTime(end.Unix()).
		Do(e.ctx)
	if err != nil {
		return nil, err
	}

	var ohlcvs []model.OHLCV
	for _, c := range candles {
		ohlcvs = append(ohlcvs, OHLCV(c))
	}
	return ohlcvs, nil
}

func OHLCV(b *binance.Kline) model.OHLCV {
	dt := unixToISO(b.OpenTime)
	ts := strconv.Itoa(int(b.OpenTime) / 1000)
	o := model.OHLCV{
		DATE:   dt,
		TS:     ts,
		OPEN:   b.Open,
		CLOSE:  b.Close,
		HIGH:   b.High,
		LOW:    b.Low,
		VOLUME: b.Volume,
	}
	return o
}

func unixToISO(date int64) string {
	dt := time.Unix(date/1000, 0)
	return dt.UTC().Format(time.RFC3339)
}
