package internal

import (
	"context"
	"encoding/csv"
	"github.com/jmrtzsn/coiner/internal/exchange"
	"github.com/jmrtzsn/coiner/internal/export"
	"os"
	"time"
)

const batchSize = 500

type Downloader struct {
	Exchange exchange.Exchange
	Exports  []export.Export
	Interval string
	Symbols  []string
	From     string
	To       string
}

// Todo adds storage
func NewDownloader(exchange exchange.Exchange, exports []export.Export) Downloader {
	return Downloader{
		exchange: exchange,
		exports: exports,
	}
}

type Parameters struct {
	Start time.Time
	End   time.Time
}

type Option func(*Parameters)

func WithInterval(start, end time.Time) Option {
	return func(parameters *Parameters) {
		parameters.Start = start
		parameters.End = end
	}
}

// TODO: input as a config file
func (d Downloader) Download(ctx context.Context, symbol, timeframe string, output string, options ...Option) error {
	recordFile, err := os.Create(output)
	file, err := CreateTempCSV(records)
	if err != nil {
		return err
	}

	parameters := &Parameters{
		Start: now.AddDate(0, -1, 0),
		End:   now,
	}

	for _, option := range options {
		option(parameters)
	}

	parameters.Start = time.Date(parameters.Start.Year(), parameters.Start.Month(), parameters.Start.Day(),
		0, 0, 0, 0, time.UTC)
	parameters.End = time.Date(parameters.End.Year(), parameters.End.Month(), parameters.End.Day(),
		0, 0, 0, 0, time.UTC)

	writer := csv.NewWriter(recordFile)
	for begin := parameters.Start; begin.Before(parameters.End); begin = begin.Add(interval * batchSize) {
		end := begin.Add(interval * batchSize)
		if end.After(parameters.End) {
			end = parameters.End
		}

		candles, err := d.exchange.CandlesByPeriod(ctx, symbol, timeframe, begin, end)
		if err != nil {
			return err
		}

		for _, candle := range candles {
			err := writer.Write(candle.ToSlice())
			if err != nil {
				return err
			}
		}
	}
	writer.Flush()
	log.Info("Done!")
	return writer.Error()
}
