package cmd

import (
	"github.com/jmrtzsn/coiner/internal/projectpath"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Port    int
	Name    string
	PathMap string `mapstructure:"path_map"`
}

// TODO read ENV vars and load into config
func InitConfig()*Config {
	viper.SetConfigName(".env")
	viper.AddConfigPath(projectpath.Root)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatal("unable to open config file, %v", err)
		} else {
			log.Fatal("Unknown error, %v", err)
		}
	}

	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatal("unable to decode config into struct, %v", err)
	}

}
