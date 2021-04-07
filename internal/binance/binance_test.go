package binance

import (
	"github.com/google/go-cmp/cmp"
	"github.com/jmrtzsn/coiner/internal"
	"log"
	"os"
	"testing"
)

var b = Binance{}

func TestMain(m *testing.M) {
	log.Println("Setting up binance exchange testing suite!")
	b.Init()
	exitVal := m.Run()
	log.Println("Completed binance exchange testing suite!")
	os.Exit(exitVal)
}

func Test_isoToUnix(t *testing.T) {
	input := "2020-04-04T12:07:00+00:00"
	var want int64 = 1586002020000
	got, err := isoToUnix(input)
	checkErr(t, err)
	if !cmp.Equal(got, want) {
		t.Errorf("Got %v want %v", got, want)
	}
}

func Test_unixToUnix(t *testing.T) {
	var input int64 = 1586002020000
	var want = "2020-04-04T12:07:00+00:00"
	got := unixToISO(input)
	if !cmp.Equal(got, want) {
		t.Errorf("Got %v want %v", got, want)
	}
}

func TestBinance_OHLCV(t *testing.T) {
	start := "2020-04-04T12:00:00+00:00"
	end := "2020-04-04T12:59:00+00:00"
	want := 60
	want2 := internal.OHLCV{
		"2020-04-04T14:00:00+02:00",
		"1586001600000",
		"6696.68000000",
		"6717.68000000",
		"6717.68000000",
		"6686.43000000",
		"155.99070000",
	}
	got, err := b.OHLCV("BTCUSDT", "1m", start, end)
	checkErr(t, err)
	if !cmp.Equal(len(got), want) {
		t.Errorf("Got %v want %v", len(got), want)
	}
}

func checkErr(t *testing.T, err error) {
	if err != nil {
		t.Errorf("got !nil error: %s", err)
	}
}
