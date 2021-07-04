package cmd

import (
	"context"
	"fmt"
	"github.com/jmrtzsn/coiner/internal/projectpath"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"time"
)

const dateFmt = "2006-01-02"

var (
	cfg     string
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
	conf := UnMarshal()
	downloader, err := conf.Downloader(context.Background())
	if err != nil {
		panic(err)
	}
	downloader.Logger.Infof("Running on: %s ", downloader.String())
	downloader.Download()
}

// Init on cobra start
func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringVarP(&cfg, "config", "c", "", "Name of the config file")
	rootCmd.Flags().StringP("exchange", "e", "", "Exchange")
	rootCmd.Flags().StringP("interval", "i", "", "Interval (optional) defaults to 1m")
	rootCmd.Flags().StringSlice("symbols", []string{}, "comma separated symbol list: --symbols=\"BTCUSDT,ETHUSDT\"")
	rootCmd.Flags().StringSlice("exports", []string{}, "comma separated output list: --exports=\"local,bucket\"")
	rootCmd.Flags().StringP("start", "s", "", "Start: 2019-01-01 (defaults to today)")
	rootCmd.Flags().StringP("end", "d", "", "End: 2019-01-02 (defaults to today)")

}

// initConfig loads env vars
func initConfig() {
	LoadConfig(cfg)
	if err := godotenv.Load(fmt.Sprintf("%s/%s.env", projectpath.Root, cfg)); err != nil {
		log.Fatal(err)
	}
}

func LoadConfig(name string) {
	// Currently only env files are supported
	viper.AddConfigPath(projectpath.Root)
	viper.SetConfigName(name)
	viper.SetConfigType("env")

	// Bind Cobra flags to Viper
	viper.BindPFlag("config", rootCmd.Flags().Lookup("config"))
	viper.BindPFlag("exchange", rootCmd.Flags().Lookup("exchange"))
	viper.BindPFlag("interval", rootCmd.Flags().Lookup("interval"))
	viper.BindPFlag("symbols", rootCmd.Flags().Lookup("symbols"))
	viper.BindPFlag("exports", rootCmd.Flags().Lookup("exports"))
	viper.BindPFlag("start", rootCmd.Flags().Lookup("start"))
	viper.BindPFlag("end", rootCmd.Flags().Lookup("end"))

	// Set default values
	viper.SetDefault("Exchange", "binance")
	viper.SetDefault("Interval", "1m")
	viper.SetDefault("Start", time.Now().Format(dateFmt))
	viper.SetDefault("End", time.Now().Format(dateFmt))

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
