package cmd

import (
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
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

var (
	rootCmd = &cobra.Command{
		Use:   "coiner",
		Short: "A common interface for popular crypto exchanges",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(loadEnvs)
	rootCmd.Flags().StringVarP(&Exchange, "exchange", "r", "", "Exchange (required)")
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

func loadEnvs() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
