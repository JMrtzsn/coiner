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
	got := LoadConfig("test", "env")
	assert.Equal(t, &env, got )
}

func TestEnvToConfig(t *testing.T){
	got, err := ToDownloader(env)
	assert.Nil(t, err)

	assert.Equal(t, "want", got )
}
