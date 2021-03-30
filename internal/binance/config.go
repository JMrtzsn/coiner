package binance

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	apiKey string
	apiSecret string
}

func (c *Config) loadEnv(){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	c.apiKey = os.Getenv("BINANCE_KEY")
	c.apiSecret = os.Getenv("BINANCE_SECRET")
}