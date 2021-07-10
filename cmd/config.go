package cmd

import (
	"context"
	"fmt"
	"github.com/jmrtzsn/coiner/pkg/downloader"
	"github.com/jmrtzsn/coiner/pkg/exchange"
	"github.com/jmrtzsn/coiner/pkg/exchange/binance"
	"github.com/jmrtzsn/coiner/pkg/export"
	"github.com/jmrtzsn/coiner/pkg/notification"
	"github.com/jmrtzsn/coiner/pkg/projectpath"
	"github.com/spf13/viper"
	"github.com/xhit/go-str2duration/v2"
	"go.uber.org/zap"
	"log"
	"os"
	"time"
)

// Config contain program input (loaded from env hence the mapstructure all caps)
// Uses: https://github.com/mitchellh/mapstructure
type Config struct {
	Exchange        string   `mapstructure:"EXCHANGE,omitempty"`                       // binance
	Interval        string   `mapstructure:"INTERVAL,omitempty"`                       // 1d, 1h, 15m, 1m
	Symbols         []string `mapstructure:"SYMBOLS,omitempty"`                        // BTCUSDT, ETHUSDT
	Exports         []string `mapstructure:"EXPORTS,omitempty"`                        // local, bucket
	Start           string   `mapstructure:"START,omitempty"`                          // 2020-04-04
	End             string   `mapstructure:"END,omitempty"`                            // 2020-04-05
	Key             string   `mapstructure:"KEY,omitempty"`                            // Exchange key
	Secret          string   `mapstructure:"SECRET,omitempty"`                         // Exchange secret
	Bucket          string   `mapstructure:"BUCKET,omitempty"`                         // GCP Bucket
	Credentials     string   `mapstructure:"GOOGLE_APPLICATION_CREDENTIALS,omitempty"` // GCP credentials file
	TelegramBot     string   `mapstructure:"TELEGRAM_BOT,omitempty"`                   // Bot string
	TelegramChannel string   `mapstructure:"TELEGRAM_CHANNEL,omitempty"`               // Channel id, check bot/getUpdates
}

func UnMarshal() *Config {
	config := &Config{}
	err := viper.Unmarshal(config)
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", projectpath.Root+config.Credentials)
	if err != nil {
		log.Fatalf("unable to decode config into struct, %v", err)
	}
	return config
}

// TODO: Downloader opts pattern
// NewDownloader move to Downloader implement default values etc
func (conf Config) NewDownloader(ctx context.Context) (*downloader.Downloader, error) {
	exchang, err := conf.setupExchange(ctx)
	if err != nil {
		return nil, err
	}
	exports, err := conf.setupExports(ctx)
	if err != nil {
		return nil, err
	}
	s, err := start(conf.Start)
	if err != nil {
		return nil, err
	}
	e, err := end(conf.End)
	if err != nil {
		return nil, err
	}
	dur, err := str2duration.ParseDuration(conf.Interval)
	if err != nil {
		return nil, err
	}
	l, err := newLogger()
	if err != nil {
		return nil, err
	}

	notifier := notification.NewTelegram(
		conf.TelegramBot,
		conf.TelegramChannel,
	)

	return &downloader.Downloader{
		Exchange: exchang,
		Exports:  exports,
		Interval: conf.Interval,
		Duration: dur,
		Symbols:  conf.Symbols,
		Start:    s,
		End:      e,
		Logger:   l,
		Telegram: notifier,
	}, nil
}

func (conf Config) setupExports(ctx context.Context) ([]export.Export, error) {
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

func (conf Config) setupExchange(ctx context.Context) (exchange.Exchange, error) {
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
