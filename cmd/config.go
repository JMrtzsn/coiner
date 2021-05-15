package cmd

import (
	"github.com/jmrtzsn/coiner/internal/projectpath"
	"github.com/spf13/viper"
	"log"
)

// Config contain program input (loaded from env or input vars)
type Config struct {
	// TODO create valid types (enums?) for vars (no checking?)
	// TODO: API_KEYS MAP? - Exchange object
	Exchanges []string
	Interval string
	Symbols []string // TODO create valid symbol (enum?(
	Output []string
	From string
	To string

}

// TODO read ENV vars and load into config
func InitConfig() *Config {
	viper.SetConfigName(".env")
	viper.AddConfigPath(projectpath.Root)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatal("unable to open config file, %v", err)
		} else {
			log.Fatal("Unknown error, %v", err)
		}
	}

	config := &Config{}
	err := viper.Unmarshal(config)
	if err != nil {
		log.Fatal("unable to decode config into struct, %v", err)
	}
	return config
}
