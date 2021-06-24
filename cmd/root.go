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
	CfgFile string
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

// Run function is run when root is executed.
func Run(cmd *cobra.Command, args []string) {
	downloader := ToDownloader()
	downloader.Logger.Infof("Running on: %s ", downloader.String())
	downloader.Download()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringVarP(&CfgFile, "config", "c", "", "Name of the config file")
	rootCmd.Flags().StringP("exchange", "e", "", "Exchange")
	rootCmd.Flags().StringP("interval", "i", "", "Interval (optional) defaults to 1min")
	rootCmd.Flags().StringSlice("symbols", []string{}, "comma separated symbol list: --symbols=\"BTCUSDT,ETHUSDT\"")
	rootCmd.Flags().StringSlice("exports", []string{}, "comma separated output list: --exports=\"local,bucket\"")
	rootCmd.Flags().StringP("start", "s", "", "Start: 2019-01-01 (defaults to today)")
	rootCmd.Flags().StringP("end", "d", "", "End: 2019-01-02 (defaults to today)")

}

func initConfig() {
	if CfgFile != "" {
		LoadConfig(CfgFile)
	} else {
		LoadConfig("viper")
	}
}

func LoadConfig(name string) {
	viper.AddConfigPath(projectpath.Root)
	viper.SetConfigName(name)
	viper.SetConfigType("env")

	viper.BindPFlag("config", rootCmd.Flags().Lookup("config"))
	viper.BindPFlag("exchange", rootCmd.Flags().Lookup("exchange"))
	viper.BindPFlag("interval", rootCmd.Flags().Lookup("interval"))
	viper.BindPFlag("symbols", rootCmd.Flags().Lookup("symbols"))
	viper.BindPFlag("exports", rootCmd.Flags().Lookup("exports"))
	viper.BindPFlag("start", rootCmd.Flags().Lookup("start"))
	viper.BindPFlag("end", rootCmd.Flags().Lookup("end"))

	viper.SetDefault("Exchange", "binance")
	viper.SetDefault("Interval", "1min")
	viper.SetDefault("Start", time.Now().Format("2006-01-02"))
	viper.SetDefault("End", time.Now().Format("2006-01-02"))

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("unable to open config file, %v", err)
		} else {
			log.Fatalf("Unknown error, %v", err)
		}
	}
	fmt.Println("Using config file:", viper.ConfigFileUsed())
}
