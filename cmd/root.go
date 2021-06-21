package cmd

import (
	"fmt"
	"github.com/jmrtzsn/coiner/internal/projectpath"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"time"
)

var (
	Exchange string
	Interval string
	To       string
	From     string
	Symbols  []string
	Exports  []string
)

var (
	rootCmd = &cobra.Command{
		Use:   "coiner",
		Short: "coiner: A common interface downloading data from popular crypto exchanges",
		Run:   Run,
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

func initConfig() {
	LoadConfig("viper", "env")
	viper.AutomaticEnv()
}

func Run(cmd *cobra.Command, args []string) {
	fmt.Println(Exchange)
	fmt.Println(Interval)
	fmt.Printf("%v \n", len(Symbols))
	fmt.Printf("%v \n", Exports)
	fmt.Println(From)
	fmt.Println(To)
}

// TODO run with args -> run with env
func setupFlags(rootCmd *cobra.Command) {
	rootCmd.Flags().StringVarP(&Exchange, "exchange", "e", "", "Exchange")
	rootCmd.Flags().StringVarP(&Interval, "interval", "i", "", "Interval (optional) defaults to 1min")
	rootCmd.Flags().StringArrayP("symbols", "y", Symbols, "comma separated symbol list: --symbols=\"BTCUSDT,ETHUSDT\"")
	rootCmd.Flags().StringArrayP("exports", "x", Exports, "comma separated output list: --symbols=\"local,bucket\"")

	// Defaults to today
	rootCmd.Flags().StringVarP(&From, "from", "f", "", "From: 2019-01-01 (defaults to today)")
	rootCmd.Flags().StringVarP(&To, "to", "t", "", "To: 2019-01-02 (defaults to today)")
}

func LoadConfig(name, typ string) {
	viper.AddConfigPath(projectpath.Root)
	viper.SetConfigName(name)
	viper.SetConfigType(typ)

	viper.SetDefault("Exchange", "binance")
	viper.SetDefault("Interval", "1min")
	viper.SetDefault("From", time.Now().Format("2006-01-02"))
	viper.SetDefault("To", time.Now().Format("2006-01-02"))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("unable to open config file, %v", err)
		} else {
			log.Fatalf("Unknown error, %v", err)
		}
	}
	fmt.Println("Using config file:", viper.ConfigFileUsed())
}
