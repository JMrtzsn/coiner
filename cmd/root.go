package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var (
	// Binance
	Exchange string
	// day, hour, minute. Defaults to 1 min.
	Interval string
	// 2019-01-01
	From, To string
	Symbol []string
	Output []string
)

type config struct {
	Port int
	Name string
	PathMap string `mapstructure:"path_map"`
}

var (
	C config
	cfgFile     string
	rootCmd = &cobra.Command{
		Use:   "coiner",
		Short: "A common interface for popular crypto exchanges",
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().StringVarP(&Exchange, "exchange", "e", "", "Exchange (required)")
	rootCmd.MarkFlagRequired("exchange")
	rootCmd.Flags().StringVarP(&Interval, "interval", "r", "", "Interval (optional) defaults to 1min")

	rootCmd.Flags().StringSlice("Symbol", Symbol, "comma separated symbol list: BTCUSDT, ETHUSD")
	rootCmd.MarkFlagRequired("symbol")

	rootCmd.Flags().StringSlice("Output", Output,  "comma separated output list: local, storage")
	rootCmd.MarkFlagRequired("output")

	// Defaults to today
	rootCmd.Flags().StringVarP(&From, "from", "r", "", "From: 2019-01-01 (defaults to today)")
	rootCmd.Flags().StringVarP(&To, "to", "r", "", "To: 2019-01-01 (defaults to today)")

}
// TODO viper config -> load standard download stuff
// unmarshal json/yaml into config file use across project

func initConfig() {

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)



		err := viper.Unmarshal(&C)
		if err != nil {
			log.Fatal("unable to decode into struct, %v", err)
		}
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cobra")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

//
//func loadEnvs() {
//	err := godotenv.Load()
//	if err != nil {
//		log.Fatal("Error loading .env file")
//	}
//}
