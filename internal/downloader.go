package internal

import (
	"fmt"
	"github.com/jmrtzsn/coiner/internal/exchange"
	"github.com/jmrtzsn/coiner/internal/export"
	"time"
)

const batchSize = 500

type Downloader struct {
	Exchange exchange.Exchange
	Exports  []export.Export
	Interval string
	Symbols  []string
	From     time.Time
	To       time.Time
}

func (d Downloader) String() string{
	return fmt.Sprintf("Exchange: %s, Exports: %s, Interval: %v, Symbols: %v, From: %s, To: %s",
		d.Exchange, d.Exports, d.Interval, d.Symbols, d.From, d.To)
}
