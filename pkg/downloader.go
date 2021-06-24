package pkg

import (
	"fmt"
	"github.com/jmrtzsn/coiner/internal/exchange"
	"github.com/jmrtzsn/coiner/internal/export"
	"github.com/xhit/go-str2duration/v2"
	"go.uber.org/zap"
	"os"
	"time"
)

// TODO const per interval
const minutesInaDay = 1439
const batchsize = 500

type Downloader struct {
	Exchange exchange.Exchange
	Exports  []export.Export
	Interval string
	Symbols  []string
	Start    time.Time
	End      time.Time
	Logger   *zap.SugaredLogger
}

func (d Downloader) String() string {
	return fmt.Sprintf("Exchange: %s, Exports: %s, Interval: %v, Symbols: %v, From: %s, To: %s",
		d.Exchange, d.Exports, d.Interval, d.Symbols, d.Start, d.End)
}

// Downloads per day, will start with whole day and move to one hour if required.
func (d Downloader) Download() {
	defer d.Logger.Sync()

	// TODO handle batchsize

	// TODO switch intervsl to durstion formst

	for _, symbol := range d.Symbols {
		d.Logger.Infof("Downloading candles for Symbol: %s Start: %s End: %s using interval %s", symbol, d.Start, d.End, d.Interval)

		day, err := str2duration.ParseDuration("1d")
		if err != nil {
			d.Logger.Panic("Failed to parse duration %s ", err.Error())
		}
		minute, err := str2duration.ParseDuration("1m")
		if err != nil {
			d.Logger.Panic("Failed to parse duration %s ", err.Error())
		}

		// Day by Day TODO: Other
		for dayBegin := d.Start; dayBegin.Before(d.End); dayBegin = dayBegin.Add(day) {
			dayEnd := dayBegin.Add(minute * minutesInaDay)
			if dayEnd.After(d.End) {
				dayEnd = d.End
			}

			// TODO: move to object
			duration, err := str2duration.ParseDuration(d.Interval)
			if err != nil {
				d.Logger.Panic("Failed to parse duration %s ", err.Error())
			}

			records := d.Batch(symbol, dayBegin, dayEnd, duration)

			date := dayBegin.Format("2006-01-02")
			d.Export(symbol, date, records)
		}
	}
}

func (d Downloader) Batch(symbol string, dayBegin, dayEnd time.Time, duration time.Duration) [][]string {
	// Minute by Minute TODO: Other
	var records [][]string
	records = append(records, []string{"DATE", "TS", "OPEN", "CLOSE", "HIGH", "LOW", "VOLUME"})

	for begin := dayBegin; begin.Before(dayEnd); begin = begin.Add(duration * batchsize) {
		end := begin.Add(duration * batchsize)
		if end.After(dayEnd) {
			end = dayEnd
		}

		d.Logger.Infof("Candles for begin %s - end %s", begin, end)
		candles, err := d.Exchange.CandlesByPeriod(symbol, d.Interval, begin, end)
		if err != nil {
			// TODO: implement fallback
			d.Logger.Panicf("failed to CandlesByPeriod symbol: %s - err: %s", symbol, err.Error())
		}
		for _, candle := range candles {
			records = append(records, candle.Csv())
		}
	}
	return records
}

func (d Downloader) Export(symbol, date string, records [][]string ){
	d.Logger.Infof("Writing %s to temp file", symbol)
	temp, err := export.WriteToTempFile(records)
	if err != nil {
		// TODO: implement fallback
		d.Logger.Panicf("failed to create tempfile err: %s", err.Error())
	}

	for _, e := range d.Exports {
		err := e.Export(temp, date, symbol)
		if err != nil {
			d.Logger.Panicf("failed to export to %s - err: %s", e.String(), err.Error())
		}
	}

	err = temp.Close()
	if err != nil {
		d.Logger.Panicf("failed to close tempfile err: %s", err.Error())
	}
	err = os.Remove(temp.Name())
	if err != nil {
		d.Logger.Panicf("failed to remove tempfile err: %s", err.Error())
	}
}