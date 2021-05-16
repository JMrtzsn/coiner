package cmd

import (
	"github.com/jmrtzsn/coiner/internal/exchanges/binance"
	"github.com/jmrtzsn/coiner/internal/export"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

var env = envConfig{
	Exchange: "binance",
	Interval: "1min",
	Symbols:  []string{"BTCUSDT", "ETHUSDT"},
	Outputs:  []string{"local", "storage"},
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
	got := LoadEnvConfig("test", "env")
	assert.Equal(t, &env, got )
}

func TestEnvToConfig(t *testing.T){
	got, err := ToConfig(env)
	assert.Nil(t, err)

	ex := &binance.Binance{}
	ex.Init(env.Key, env.Secret)

	exports := make(map[string][]export.Command)
	exports["BTCUSDT"] = []export.Command{&export.Local{}}
	exports["ETHUSDT"] = []export.Command{&export.Local{}}

	want := &Config{
		Exchange: ex,
		Interval: "1min",
		Symbols:  []string{"BTCUSDT", "ETHUSDT"},
		Exports:  exports,
		From:     "2019-01-01",
		To:       "2019-01-02",
	}
	assert.Equal(t, want, got )
}
