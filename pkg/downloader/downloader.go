package downloader

import (
	"encoding/csv"
	"fmt"
	"github.com/jmrtzsn/coiner/pkg/exchange"
	"github.com/jmrtzsn/coiner/pkg/export"
	"github.com/jmrtzsn/coiner/pkg/model"
	"github.com/jmrtzsn/coiner/pkg/notification"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"time"
)

const (
	minutesInADay = 1439
	batchSize     = 500
	YMD           = "2006-01-02"
	day           = time.Hour * 24
)

type Downloader struct {
	Exchange exchange.Exchange
	Exports  []export.Export
	Interval string
	Duration time.Duration
	Symbols  []string
	Start    time.Time
	End      time.Time
	Logger   *zap.SugaredLogger
	Telegram notification.Telegram
}

func (d Downloader) String() string {
	return fmt.Sprintf("Exchange: %s, Exports: %s, Interval: %v, Symbols: %v, From: %s, To: %s",
		d.Exchange, d.Exports, d.Interval, d.Symbols, d.Start.Format(YMD), d.End.Format(YMD))
}

// Download main function of coiner, downloads files daily
// TODO make "size" an input param?
func (d Downloader) Download() {
	defer d.Logger.Sync()
	for _, symbol := range d.Symbols {
		d.Logger.Infof("Downloading candles for Symbol: "+
			"%s Start: %s End: %s using interval %s",
			symbol, d.Start, d.End, d.Interval)

		for begin := d.Start; begin.Before(d.End); begin = begin.Add(day) {
			date := begin.Format(YMD)
			end := begin.Add(time.Minute * minutesInADay)
			if end.After(d.End) {
				end = d.End
			}
			d.Logger.Infof("Downloading Candles %s for date: %s", symbol, begin)
			records, err := d.batch(symbol, begin, end, d.Duration)
			if err != nil {
				err := d.Telegram.OnError(err)
				if err != nil {
					d.Logger.Errorf(err.Error())
				}
				d.Logger.Panicf(err.Error())
			}

			if len(records) > 2 {
				if err = d.Export(symbol, date, records); err != nil {
					d.Logger.Errorf(err.Error())
					err := d.Telegram.OnError(err)
					if err != nil {
						d.Logger.Errorf(err.Error())
					}
				}
			} else {
				d.Logger.Infof("Recieved empty response from exchange for symbol: %s"+
					" skipping export. date :%s - %s", symbol, begin, end)
			}
		}
	}
	// TODO pretty
	if d.Telegram.Bot != "" {
		err := d.Telegram.Notify(fmt.Sprintf("✅ Download Completed: \n%s", d.String()))
		if err != nil {
			d.Logger.Errorf(err.Error())
		}
	}
}

// Minute by Minute TODO fix
// batch creates a set of records containing data between from and to splitting by the duration, returning a day of records
func (d Downloader) batch(symbol string, from, to time.Time, duration time.Duration) ([][]string, error) {
	records := model.RecordsWithHeader()
	for begin := from; begin.Before(to); begin = begin.Add(duration * batchSize) {
		end := begin.Add(duration * batchSize)
		if end.After(to) {
			end = to
		}
		candles, err := d.Exchange.Candles(symbol, d.Interval, begin, end)
		if err != nil {
			// TODO: switch on error type and try again if necessary
			return nil, fmt.Errorf("failed to execute CandlesByPeriod symbol: %s - err: %s", symbol, err.Error())
		}
		for _, candle := range candles {
			records = append(records, candle.CSV())
		}
	}
	return records, nil
}

func (d Downloader) Export(symbol, date string, records [][]string) error {
	temp, err := writeToTempFile(records)
	if err != nil {
		return fmt.Errorf("failed to create tempfile err: %s", err.Error())
	}

	for _, e := range d.Exports {
		err := e.Export(temp, date, symbol)
		if err != nil {
			d.Logger.Errorf("failed to export to %s - err: %s", e.String(), err.Error())
		}
	}

	err = temp.Close()
	if err != nil {
		return fmt.Errorf("failed to close tempfile err: %s", err.Error())
	}
	err = os.Remove(temp.Name())
	if err != nil {
		return fmt.Errorf("failed to remove tempfile err: %s", err.Error())
	}
	return nil
}

// writeToTempFile create a temp CSV file from records
func writeToTempFile(records [][]string) (*os.File, error) {
	file, err := ioutil.TempFile("", "file")
	if err != nil {
		return nil, err
	}
	w := csv.NewWriter(file)
	err = w.WriteAll(records)

	// Seek the pointer to the beginning
	// TODO return buffer? function that resets file when reading
	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, err
	}
	return file, nil
}
