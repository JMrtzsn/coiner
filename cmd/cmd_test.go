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

func TestInitConfig(t *testing.T){
	got := InitConfig()
	want := Config{
		Exchanges: []string{"Binance"},
		Interval:  "1min",
		Symbols:   []string{"BTCUSDT"},
		Output:    []string{"csv"},
		From:      "2019-01-01",
		To:        "2019-01-05",
	}
	assert.Equal(t, got,want )
}
