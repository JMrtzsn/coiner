package pkg

import (
	"fmt"
	"github.com/jmrtzsn/coiner/internal/exchange"
	"github.com/jmrtzsn/coiner/internal/export"
	"go.uber.org/zap"
	"os"
	"time"
)

// TODO const per interval
const minutesInADay = 1439
const hoursInADay = 23
const batchsize = 500 // Binance max batchsize

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
			end := begin.Add(time.Minute * minutesInADay)
			if end.After(d.End) {
				end = d.End
			}

			records := d.batch(symbol, begin, end, d.Duration)

			date := begin.Format("2006-01-02")
			d.Export(symbol, date, records)
		}
	}
}

// Minute by Minute TODO: Opts -> daily weekly etc
// TODO
// batch creates a set of records containing data between from and to splitting by the duration
func (d Downloader) batch(symbol string, from, to time.Time, duration time.Duration) [][]string {

	var records [][]string
	records = append(records, []string{"DATE", "TS", "OPEN", "CLOSE", "HIGH", "LOW", "VOLUME"})
	for begin := from; begin.Before(to); begin = begin.Add(duration * batchsize) {
		end := begin.Add(duration * batchsize)
		if end.After(to) {
			end = to
		}

		d.Logger.Infof("Candles for begin %s - end %s", begin, end)
		candles, err := d.Exchange.CandlesByPeriod(symbol, d.Interval, begin, end)
		if err != nil {
			// TODO: implement fallback
			d.Logger.Panicf("failed to CandlesByPeriod symbol: %s - err: %s", symbol, err.Error())
		}
		if len(candles) > 1 {
			for _, candle := range candles {
				records = append(records, candle.Csv())
			}
		} else {
			d.Logger.Warnf("Recieved empty response from exchange for symbol: %s skipping date :%s - %s", symbol, begin, end)
		}

	}
	return records
}

func (d Downloader) Export(symbol, date string, records [][]string) {
	d.Logger.Infof("Writing %s to temp file", symbol)
	temp, err := export.WriteToTempFile(records)
	if err != nil {
		// TODO: implement fallback
		d.Logger.Panicf("failed to create tempfile err: %s", err.Error())
	}

	for _, e := range d.Exports {
		d.Logger.Infof("Exporting to %s", e)
		err := e.Export(temp, date, symbol)
		if err != nil {
			// TODO: implement fallback
			d.Logger.Panicf("failed to export to %s - err: %s", e.String(), err.Error())
		}
	}

	// TODO: necessary?
	d.Logger.Infof("Cleaning up temp file")
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
