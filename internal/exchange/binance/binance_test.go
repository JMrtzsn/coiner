package binance

import (
	"context"
	"github.com/jmrtzsn/coiner/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var res = model.Candle{
	DATE:   "2020-04-04T12:00:00Z",
	TS:     "1586001600",
	OPEN:   "6696.68000000",
	CLOSE:  "6717.68000000",
	HIGH:   "6717.68000000",
	LOW:    "6686.43000000",
	VOLUME: "155.99070000",
}

// Uses the public api key / secret doesnt matter
func TestOHLCV(t *testing.T) {
	var b = Binance{}
	b.Init(context.Background(), "TEST", "TEST")
	start := time.Date(2020, 4, 4,
		12, 0, 0, 0, time.UTC)
	end := time.Date(2020, 4, 4,
		12, 59, 0, 0, time.UTC)

	got, err := b.Candles("BTCUSDT", "1m", start, end)
	assert.Nil(t, err)
	assert.Equal(t, 60, len(got))
	assert.Equal(t, res, got[0])
}
