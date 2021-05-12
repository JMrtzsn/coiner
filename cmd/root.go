package cmd

import (
	"github.com/spf13/cobra"
)

var (
	// Binance
	Exchange string
	// day, hour, minute. Defaults to 1 min.
	Interval string
	// 2019-01-01
	From, To string
	Symbol   []string
	Output   []string
)



var (
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
	setupFlags(rootCmd)
}

// TODO viper config -> load standard download stuff
// unmarshal json/yaml into config file use across project
func setupFlags(rootCmd *cobra.Command) {
	rootCmd.Flags().StringVarP(&Exchange, "exchange", "e", "", "Exchange (required)")
	rootCmd.MarkFlagRequired("exchange")
	rootCmd.Flags().StringVarP(&Interval, "interval", "r", "", "Interval (optional) defaults to 1min")

	rootCmd.Flags().StringSlice("Symbol", Symbol, "comma separated symbol list: BTCUSDT, ETHUSD")
	rootCmd.MarkFlagRequired("symbol")

	rootCmd.Flags().StringSlice("Output", Output, "comma separated output list: local, storage")
	rootCmd.MarkFlagRequired("output")

	// Defaults to today
	rootCmd.Flags().StringVarP(&From, "from", "r", "", "From: 2019-01-01 (defaults to today)")
	rootCmd.Flags().StringVarP(&To, "to", "r", "", "To: 2019-01-01 (defaults to today)")
}

//
//func loadEnvs() {
//	err := godotenv.Load()
//	if err != nil {
//		log.Fatal("Error loading .env file")
//	}
//}
