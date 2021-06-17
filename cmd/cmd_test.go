package cmd

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

var env = Viper{
	Exchange: "binance",
	Interval: "1min",
	Symbols:  []string{"BTCUSDT", "ETHUSDT"},
	Exports:  []string{"local"},
	From:     "2019-01-01",
	To:       "2019-01-02",
	Key:      "test",
	Secret:   "test",
}

func TestMain(m *testing.M) {
	log.Println("Setting up CMD testing suite!")

	exitVal := m.Run()
	log.Println("Completed CMD testing suite!")
	os.Exit(exitVal)
}

func TestLoadEnvConfig(t *testing.T){
	got := LoadConfig("viper", "env")
	assert.Equal(t, &env, got )
}

func TestEnvToConfig(t *testing.T){
	env := *LoadConfig("viper", "env")
	got := ToDownloader(env)
	want := "Exchange: Binance, Exports: [Local Local], Interval: 1min, Symbols: [BTCUSDT ETHUSDT]," +
		" From: 2019-01-01 00:00:00 +0000 UTC, To: 2019-01-02 23:59:59 +0000 UTC"
	assert.Equal(t, want, got.String() )
}

func TestToTime(t *testing.T){
	got := FromTime("2019-01-01")
	want := int64(1546300800)
	assert.Equal(t, want, got.Unix())
}
