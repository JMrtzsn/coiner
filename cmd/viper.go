package cmd

import (
	"context"
	"fmt"
	"github.com/jmrtzsn/coiner/internal/exchange"
	"github.com/jmrtzsn/coiner/internal/exchange/binance"
	"github.com/jmrtzsn/coiner/internal/export"
	"github.com/jmrtzsn/coiner/pkg"
	"github.com/spf13/viper"
	"github.com/xhit/go-str2duration/v2"
	"go.uber.org/zap"
	"log"
	"os"
	"time"
)

// Config contain program input (loaded from env hence the mapstructure all caps)
// TODO Load env vars ?
type Config struct {
	Exchange string   `mapstructure:"EXCHANGE"` // binance
	Interval string   `mapstructure:"INTERVAL"` // 1d, 1h, 15m, 1m
	Symbols  []string `mapstructure:"SYMBOLS"`  // BTCUSDT, ETHUSDT
	Exports  []string `mapstructure:"EXPORTS"`  // local, bucket
	Start    string   `mapstructure:"START"`    // 2020-04-04
	End      string   `mapstructure:"END"`      // 2020-04-05
	Key      string
	Secret   string
	Bucket   string
}

func UnMarshal() *Config {
	config := &Config{}
	err := viper.Unmarshal(config)
	if err != nil {
		log.Fatalf("unable to decode config into struct, %v", err)
	}

	config.Key = lookupStrict("KEY")
	config.Secret = lookupStrict("SECRET")
	config.Bucket = lookupStrict("BUCKET")
	_ = lookupStrict("GOOGLE_APPLICATION_CREDENTIALS")

	return config
}

func lookupStrict(key string) string {
	key, ok := os.LookupEnv(key)
	if !ok {
		log.Panic(ok)
	}
	return key
}

func (conf Config) Downloader(ctx context.Context) (*pkg.Downloader, error) {
	ex, err := conf.getExchange(ctx)
	if err != nil{
		return nil, err
	}
	exp, err := conf.getExports(ctx)
	if err != nil{
		return nil, err
	}
	s, err := start(conf.Start)
	if err != nil{
		return nil, err
	}
	e, err := end(conf.End)
	if err != nil{
		return nil, err
	}
	d, err := str2duration.ParseDuration(conf.Interval)
	if err != nil{
		return nil, err
	}
	l, err := newLogger()
	if err != nil{
		return nil, err
	}

	return &pkg.Downloader{
		Exchange: ex,
		Exports:  exp,
		Interval: conf.Interval,
		Duration: d,
		Symbols:  conf.Symbols,
		Start:    s,
		End:      e,
		Logger:   l,
	}, nil
}

func (conf Config) getExports(ctx context.Context) ([]export.Export, error) {
	var exports []export.Export
	for _, e := range conf.Exports {
		switch e {
		case "local":
			exports = append(exports, export.NewLocal(conf.Exchange))
		case "bucket":
			bucket, err := export.NewBucket(ctx, conf.Exchange, conf.Bucket)
			if err != nil {
				return nil, err
			}
			exports = append(exports, bucket)
		default:
			return nil, fmt.Errorf("export not found %s", e)
		}
	}
	return exports, nil
}

func (conf Config) getExchange(ctx context.Context) (exchange.Exchange, error) {
	switch conf.Exchange {
	case "binance":
		b := &binance.Binance{}
		b.Init(ctx, conf.Key, conf.Secret)
		return b, nil
	default:
		return nil, fmt.Errorf("exchange not found %s", conf.Exchange)
	}
}

func start(input string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, input+"T00:00:00.000Z")
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func end(input string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, input+"T23:59:59.000Z")
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func newLogger() (*zap.SugaredLogger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	defer logger.Sync()
	sugar := logger.Sugar()
	return sugar, nil
}
