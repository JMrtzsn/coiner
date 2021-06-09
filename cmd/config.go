package cmd

import (
	"context"
	"fmt"
	"github.com/jmrtzsn/coiner/internal"
	"github.com/jmrtzsn/coiner/internal/exchange"
	"github.com/jmrtzsn/coiner/internal/exchange/binance"
	"github.com/jmrtzsn/coiner/internal/projectpath"
	"github.com/spf13/viper"
	"log"
)

// Viper contain program input (loaded from env or input vars)
type Viper struct {
	// TODO create valid types (enums?) for vars (no checking?)
	// TODO: API_KEYS MAP? - Exchange object
	Exchange string   `mapstructure:"EXCHANGE"`
	Interval string   `mapstructure:"INTERVAL"`
	Symbols  []string `mapstructure:"SYMBOLS"`
	Outputs  []string `mapstructure:"OUTPUT"`
	From     string   `mapstructure:"FROM"`
	To       string   `mapstructure:"TO"`
	Key      string   `mapstructure:"KEY"`
	Secret   string   `mapstructure:"SECRET"`
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

func ToDownloader(conf Viper) (*internal.Downloader, error) {

	downloader := &internal.Downloader{}

	var exchange exchange.Exchange
	ctx := context.Background()
	switch conf.Exchange {
	case "binance":
		exchange = &binance.Binance{}
		exchange.Init(ctx, conf.Key, conf.Secret)
	default:
		return nil, fmt.Errorf("exchange not found %s", conf.Exchange)
	}



	return downloader, nil
}
