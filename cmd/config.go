package cmd

import (
	"context"
	"fmt"
	"github.com/jmrtzsn/coiner/internal/exchanges"
	"github.com/jmrtzsn/coiner/internal/exchanges/binance"
	"github.com/jmrtzsn/coiner/internal/export"
	"github.com/jmrtzsn/coiner/internal/projectpath"
	"github.com/spf13/viper"
	"log"
)

// envConfig contain program input (loaded from env or input vars)
type envConfig struct {
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

type Config struct {
	Exchange exchanges.Exchange
	Interval string
	Symbols  []string
	Exports  map[string][]export.Command
	From     string
	To       string
}

// TODO read ENV vars and load into config
func LoadEnvConfig(name, typ string) *envConfig {
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

	config := &envConfig{}
	err := viper.Unmarshal(config)
	if err != nil {
		log.Fatalf("unable to decode config into struct, %v", err)
	}

	return config
}

func ToConfig(env envConfig) (*Config, error) {
	var err error
	config := &Config{}

	switch env.Exchange {
	case "binance":
		exchange := &binance.Binance{}
		exchange.Init(env.Key, env.Secret)
		config.Exchange = exchange
	default:
		return nil, fmt.Errorf("exchange not found %s", env.Exchange)
	}

	config.Exports, err = export.CreateCommands(context.Background(), env.Outputs, env.Symbols, env.Exchange)
	if err != nil {
		return nil, err
	}

	config.Interval = env.Interval
	config.Symbols = env.Symbols
	config.From = env.From
	config.To = env.To

	return config, nil
}
