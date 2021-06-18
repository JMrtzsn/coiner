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
	Symbols  []string `mapstructure:"SYMBOLS"`  // BTCUSDT, ETHUSDT
	Exports  []string `mapstructure:"EXPORTS"`  // local, bucket
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

	viper.SetDefault("Exchange", "binance")
	viper.SetDefault("Interval", "1min")
	viper.SetDefault("From", time.Now().Format("2006-01-02"))
	viper.SetDefault("To", time.Now().Format("2006-01-02"))

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

func setExport(conf Viper, ctx context.Context) []export.Export {
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
	return inputExport
}

func setExchange(conf Viper, ctx context.Context) exchange.Exchange {
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
