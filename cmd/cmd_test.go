package cmd

import (
	"context"
	"fmt"
	projectpath2 "github.com/jmrtzsn/coiner/pkg/projectpath"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	log.Println("Setting up CMD testing suite!")
	if err := godotenv.Load(fmt.Sprintf("%s/test.env", projectpath2.Root)); err != nil {
		log.Fatal(err)
	}
	exitVal := m.Run()
	log.Println("Completed CMD testing suite!")
	os.Exit(exitVal)
}

func TestLoadConfig(t *testing.T) {
	LoadConfig("test")
	conf := UnMarshal()
	got, err := conf.NewDownloader(context.Background())
	assert.Nil(t, err)
	want := "\nExchange: Binance, \nExports: [Local Bucket], \nInterval: 1m, \nSymbols: [BTCUSDT ETHUSDT], \nFrom: 2019-01-01, \nTo: 2019-01-02"
	assert.Equal(t, want, got.String())
}

func TestFromTime(t *testing.T) {
	got, _ := start("2019-01-01")
	want := int64(1546300800)
	assert.Equal(t, want, got.Unix())
}

func TestToTime(t *testing.T) {
	got, _ := end("2019-01-01")
	want := int64(1546387199)
	assert.Equal(t, want, got.Unix())
}
