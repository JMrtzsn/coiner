package coiner

import (
	"github.com/joho/godotenv"
	"log"
)

// TODO Exchange, Output [export, Local etc], time [day, hour, minute], Symbol, date_from, date_to
func main() {
	// "Binance, BTCUSD, 1 min, date_from, date_to
	// TODO: translate args into creation of an exchange interface object, init it and download
}

// TODO: setup
func loadEnvs(){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}