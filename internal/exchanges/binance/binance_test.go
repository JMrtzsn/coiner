package binance

import (
	"github.com/jmrtzsn/coiner/internal"
	"log"
	"os"
	"testing"
)

var b = Binance{}

var res = internal.OHLCV{
	DATE:   "2020-04-04T12:00:00Z",
	TS:     "1586001600",
	OPEN:   "6696.68000000",
	CLOSE:  "6717.68000000",
	HIGH:   "6717.68000000",
	LOW:    "6686.43000000",
	VOLUME: "155.99070000",
}

func TestMain(m *testing.M) {
	log.Println("Setting up binance exchange testing suite!")
	b.Init()
	exitVal := m.Run()
	log.Println("Completed binance exchange testing suite!")
	os.Exit(exitVal)
}

func Test_isoToUnix(t *testing.T) {
	input := "2020-04-04T12:07:00Z"
	var want int64 = 1586002020000
	got, err := isoToUnix(input)
	internal.Check(t, err)
	internal.Compare(t, got, want)
}

func Test_unixToISO(t *testing.T) {
	var input int64 = 1586002020000
	var want = "2020-04-04T12:07:00Z"
	got := unixToISO(input)
	internal.Compare(t, got, want)
}

func TestBinance_OHLCV(t *testing.T) {
	start := "2020-04-04T12:00:00+00:00"
	end := "2020-04-04T12:59:00+00:00"
	got, err := b.OHLCV("BTCUSDT", "1m", start, end)
	internal.Check(t, err)
	internal.Compare(t,  len(got), 60)
	internal.Compare(t,  got[0], res)
}
