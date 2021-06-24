package cmd

import (
	"context"
	"fmt"
	"github.com/jmrtzsn/coiner/internal/exchange"
	"github.com/jmrtzsn/coiner/internal/exchange/binance"
	"github.com/jmrtzsn/coiner/internal/export"
	"github.com/jmrtzsn/coiner/pkg"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"time"
)

// Config contain program input (loaded from env hence the mapstructure all caps)
type Config struct {
	Exchange string   `mapstructure:"EXCHANGE"` // binance
	Interval string   `mapstructure:"INTERVAL"` // 1d, 1h, 15m, 1m
	Symbols  []string `mapstructure:"SYMBOLS"`  // BTCUSDT, ETHUSDT
	Exports  []string `mapstructure:"EXPORTS"`  // local, bucket
	Start    string   `mapstructure:"START"`    // 2020-04-04
	End      string   `mapstructure:"END"`      // 2020-04-05
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
func ToDownloader() pkg.Downloader {
	conf := *unMarshalViper()
	ctx := context.Background()
	inputExchange := setExchange(conf, ctx)
	inputExport := setExport(conf, ctx)
	Start := StartTime(conf.Start)
	End := EndTime(conf.End)
	sugar := newLogger()

	downloader := pkg.Downloader{
		Exchange: inputExchange,
		Exports:  inputExport,
		Interval: conf.Interval,
		Symbols:  conf.Symbols,
		Start:    Start,
		End:      End,
		Logger:   sugar,
	}

	return downloader
}

func newLogger() *zap.SugaredLogger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()
	return sugar
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

func StartTime(input string) time.Time {
	t, err := time.Parse(time.RFC3339, input+"T00:00:00.000Z")
	if err != nil {
		panic(fmt.Sprintf("failed to convert time %s", t))
	}
	return t
}

func EndTime(input string) time.Time {
	t, err := time.Parse(time.RFC3339, input+"T23:59:59.000Z")
	if err != nil {
		panic(fmt.Sprintf("failed to convert time %s", t))
	}
	return t
}