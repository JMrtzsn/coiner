package cmd

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	log.Println("Setting up CMD testing suite!")

	exitVal := m.Run()
	log.Println("Completed CMD testing suite!")
	os.Exit(exitVal)
}

// TODO VALIDATE INPUT
func TestLoadEnvConfig(t *testing.T) {
	var env = Config{
		Exchange: "binance",
		Interval: "1m",
		Symbols:  []string{"BTCUSDT", "ETHUSDT"},
		Exports:  []string{"local"},
		Start:    "2019-01-01",
		End:      "2019-01-02",
		Key:      "test",
		Secret:   "test",
	}
	LoadConfig("test")
	got := unMarshalViper()
	assert.Equal(t, &env, got)
}

func TestEnvToConfig(t *testing.T) {
	LoadConfig("test")
	got := ToDownloader()
	want := "Exchange: Binance, Exports: [Local], Interval: 1m, Symbols: [BTCUSDT ETHUSDT]," +
		" From: 2019-01-01 00:00:00 +0000 UTC, To: 2019-01-02 23:59:59 +0000 UTC"
	assert.Equal(t, want, got.String())
}

func TestFromTime(t *testing.T) {
	got := StartTime("2019-01-01")
	want := int64(1546300800)
	assert.Equal(t, want, got.Unix())
}

func TestToTime(t *testing.T) {
	got := EndTime("2019-01-01")
	want := int64(1546387199)
	assert.Equal(t, want, got.Unix())
}