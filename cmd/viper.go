package cmd

import (
	"context"
	"fmt"
	"github.com/jmrtzsn/coiner/internal"
	"github.com/jmrtzsn/coiner/internal/exchange"
	"github.com/jmrtzsn/coiner/internal/exchange/binance"
	"github.com/jmrtzsn/coiner/internal/export"
	"github.com/spf13/viper"
	"log"
	"time"
)

// Config contain program input (loaded from env hence the mapstructure all caps)
type Config struct {
	Exchange string   `mapstructure:"EXCHANGE"` // binance
	Interval string   `mapstructure:"INTERVAL"` // 1d, 1h, 15m, 1m
	Symbols  []string `mapstructure:"SYMBOLS"`  // BTCUSDT, ETHUSDT
	Exports  []string `mapstructure:"EXPORTS"`  // local, bucket
	From     string   `mapstructure:"FROM"`     // 2020-04-04
	To       string   `mapstructure:"TO"`       // 2020-04-05
	Key      string   `mapstructure:"KEY"`      // exchange key
	Secret   string   `mapstructure:"SECRET"`   // exchange secret
}

func unMarshalViper() *Config {
	config := &Config{}
	err := viper.Unmarshal(config)
	if err != nil {
		log.Fatalf("unable to decode config into struct, %v", err)
	}
	return config
}

// TODO validate input params
func ToDownloader() internal.Downloader {
	conf := *unMarshalViper()

	ctx := context.Background()

	inputExchange := setExchange(conf, ctx)
	inputExport := setExport(conf, ctx)
	From := FromTime(conf.From)
	To := ToTime(conf.To)

	downloader := internal.Downloader{
		Exchange: inputExchange,
		Exports:  inputExport,
		Interval: conf.Interval,
		Symbols:  conf.Symbols,
		From:     From,
		To:       To,
	}

	return downloader
}

func setExport(conf Config, ctx context.Context) []export.Export {
	var inputExport []export.Export
	for _, e := range conf.Exports {
		switch e {
		case "local":
			inputExport = append(inputExport, export.NewLocal(conf.Exchange))
		case "storage":
			bucket, err := export.NewBucket(ctx, conf.Exchange)
			if err != nil {
				panic(fmt.Sprintf("bucket export creation error: %s", e))
			}
			inputExport = append(inputExport, bucket)
		default:
			panic(fmt.Sprintf("export not found %s", e))
		}
	}
	return inputExport
}

func setExchange(conf Config, ctx context.Context) exchange.Exchange {
	var inputExchange exchange.Exchange
	switch conf.Exchange {
	case "binance":
		inputExchange = &binance.Binance{}
		inputExchange.Init(ctx, conf.Key, conf.Secret)
	default:
		panic(fmt.Sprintf("exchange not found %s", conf.Exchange))
	}
	return inputExchange
}

func FromTime(input string) time.Time {
	t, err := time.Parse(time.RFC3339, input+"T00:00:00.000Z")
	if err != nil {
		panic(fmt.Sprintf("failed to convert time %s", t))
	}
	return t
}

func ToTime(input string) time.Time {
	t, err := time.Parse(time.RFC3339, input+"T23:59:59.000Z")
	if err != nil {
		panic(fmt.Sprintf("failed to convert time %s", t))
	}
	return t
}
