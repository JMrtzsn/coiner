package pkg

import (
	"fmt"
	"github.com/jmrtzsn/coiner/internal/exchange"
	"github.com/jmrtzsn/coiner/internal/export"
	"github.com/jmrtzsn/coiner/internal/model"
	"go.uber.org/zap"
	"os"
	"time"
)

// TODO const per interval
const minutesInADay = 1439
const batchSize = 500 // Binance max batchSize
const dateFmt = "2006-01-02"

// TODO input param
const day = time.Hour * 24

type Downloader struct {
	Exchange exchange.Exchange
	Exports  []export.Export
	Interval string        // xm, xh, xd
	Duration time.Duration // minute, hour, day
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
// Download main function of coiner
func (d Downloader) Download() {
	defer d.Logger.Sync()
	for _, symbol := range d.Symbols {
		d.Logger.Infof("Downloading candles for Symbol: %s Start: %s End: %s using interval %s",
			symbol, d.Start, d.End, d.Interval)
		for begin := d.Start; begin.Before(d.End); begin = begin.Add(day) {
			date := begin.Format(dateFmt)
			end := begin.Add(time.Minute * minutesInADay)
			if end.After(d.End) {
				end = d.End
			}

			d.Logger.Infof("Downloading Candles %s for date: %s", symbol, begin)
			records := d.batch(symbol, begin, end, d.Duration)
			if len(records) > 2 {
				d.Export(symbol, date, records)
			} else {
				d.Logger.Warnf("Recieved empty response from exchange for symbol: %s skipping date :%s - %s", symbol, begin, end)
			}
		}
	}
}

// Minute by Minute TODO: Opts -> daily weekly etc
// TODO
// batch creates a set of records containing data between from and to splitting by the duration, returning a day of records
func (d Downloader) batch(symbol string, from, to time.Time, duration time.Duration) [][]string {

	records := model.RecordsWithHeader()
	for begin := from; begin.Before(to); begin = begin.Add(duration * batchSize) {
		end := begin.Add(duration * batchSize)
		if end.After(to) {
			end = to
		}

		candles, err := d.Exchange.CandlesByPeriod(symbol, d.Interval, begin, end)
		if err != nil {
			// TODO: implement fallback
			d.Logger.Panicf("failed to execute CandlesByPeriod symbol: %s - err: %s", symbol, err.Error())
		}

		// TODO - How to handle empty dates
		for _, candle := range candles {
			records = append(records, candle.Csv())
		}

	}
	return records
}

func (d Downloader) Export(symbol, date string, records [][]string) {
	temp, err := export.WriteToTempFile(records)
	if err != nil {
		// TODO: implement fallback
		d.Logger.Panicf("failed to create tempfile err: %s", err.Error())
	}

	for _, e := range d.Exports {
		err := e.Export(temp, date, symbol)
		if err != nil {
			// TODO: implement fallback
			d.Logger.Panicf("failed to export to %s - err: %s", e.String(), err.Error())
		}
	}

	// TODO: necessary?
	err = temp.Close()
	if err != nil {
		// TODO: implement fallback
		d.Logger.Panicf("failed to close tempfile err: %s", err.Error())
	}
	err = os.Remove(temp.Name())
	if err != nil {
		// TODO: implement fallback
		d.Logger.Panicf("failed to remove tempfile err: %s", err.Error())
	}
}
