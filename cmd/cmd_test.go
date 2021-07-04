package cmd

import (
	"context"
	"fmt"
	"github.com/jmrtzsn/coiner/internal/projectpath"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	log.Println("Setting up CMD testing suite!")
	if err := godotenv.Load(fmt.Sprintf("%s/test.env", projectpath.Root)); err != nil {
		log.Fatal(err)
	}
	exitVal := m.Run()
	log.Println("Completed CMD testing suite!")
	os.Exit(exitVal)
}

func TestLoadConfig(t *testing.T) {
	LoadConfig("test")
	conf := UnMarshal()
	got, err := conf.Downloader(context.Background())
	assert.Nil(t, err)
	want := "Exchange: Binance, Exports: [Local Bucket], Interval: 1m, Symbols: [BTCUSDT ETHUSDT]," +
		" From: 2019-01-01 00:00:00 +0000 UTC, To: 2019-01-02 23:59:59 +0000 UTC"
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
