package cmd

import (
	"context"
	"fmt"
	"github.com/jmrtzsn/coiner/internal"
	"github.com/jmrtzsn/coiner/internal/exchange"
	"github.com/jmrtzsn/coiner/internal/exchange/binance"
	"github.com/jmrtzsn/coiner/internal/export"
	"github.com/jmrtzsn/coiner/internal/projectpath"
	"github.com/spf13/viper"
	"log"
	"time"
)

// Viper contain program input (loaded from env or input vars)
type Viper struct {
	// TODO create valid types (enums?) for vars (no checking?)
	// TODO: API_KEYS MAP? - Exchange object
	Exchange string   `mapstructure:"EXCHANGE"` // binance .. TODO slice
	Interval string   `mapstructure:"INTERVAL"` // 1d, 1h, 15m, 1m
	Symbols  []string `mapstructure:"SYMBOLS"`  // BTCUSDT, etc
	Exports  []string `mapstructure:"EXPORT"`   // local, bucket
	From     string   `mapstructure:"FROM"`     // 2020-04-04
	To       string   `mapstructure:"TO"`       // 2020-04-05
	Key      string   `mapstructure:"KEY"`      // exchange key
	Secret   string   `mapstructure:"SECRET"`   // exchange secret
}

// TODO read ENV vars and load into config
func LoadConfig(name, typ string) *Viper {
	viper.AddConfigPath(projectpath.Root)
	viper.SetConfigName(name)
	viper.SetConfigType(typ)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("unable to open config file, %v", err)
		} else {
			log.Fatalf("Unknown error, %v", err)
		}
	}
	config := &Viper{}
	err := viper.Unmarshal(config)
	if err != nil {
		log.Fatalf("unable to decode config into struct, %v", err)
	}

	return config
}

// TODO validate input params
func ToDownloader(conf Viper) internal.Downloader {
	var inputExchange exchange.Exchange
	ctx := context.Background()
	switch conf.Exchange {
	case "binance":
		inputExchange = &binance.Binance{}
		inputExchange.Init(ctx, conf.Key, conf.Secret)
	default:
		panic(fmt.Sprintf("exchange not found %s", conf.Exchange))
	}

	var inputExport []export.Export
	for _, e := range conf.Exports {
		switch e {
		case "local":
			for _, symbol := range conf.Symbols {
				inputExport = append(inputExport, export.NewLocal(conf.Exchange, symbol))
			}
		case "storage":
			for _, symbol := range conf.Symbols {
				bucket, err := export.NewBucket(ctx, conf.Exchange, symbol)
				if err != nil {
					panic(fmt.Sprintf("bucket export creation error: %s", e))
				}
				inputExport = append(inputExport, bucket)
			}
		default:
			panic(fmt.Sprintf("export not found %s", e))
		}
	}

	downloader := internal.Downloader{
		Exchange: inputExchange,
		Exports:  inputExport,
		Interval: conf.Interval,
		Symbols:  conf.Symbols,
		From:     ToTime(conf.From),
		To:       ToEndTime(conf.To),
	}

	return downloader
}

func ToTime(input string) time.Time {
	t, err := time.Parse(time.RFC3339, input+"T00:00:00.000Z")
	if err != nil {
		panic(fmt.Sprintf("failed to convert time %s", t))
	}
	return t
}

func ToEndTime(input string) time.Time {
	t, err := time.Parse(time.RFC3339, input+"T23:59:59.000Z")
	if err != nil {
		panic(fmt.Sprintf("failed to convert time %s", t))
	}
	return t
}
