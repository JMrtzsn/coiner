package coiner

import (
	"internal/binance"
	"net/http"
)

// Exchange is the external interface coiner uses to structure exchanges
type Exchange interface {
	Init()
	Klines() []Kline
}

// Ensure the handlers implement the required interfaces/types at compile time
var (
	_ Exchange = Exchange(&Binance)
	_ http.Handler = http.HandlerFunc((&controller{}).health)
	_ middleware   = (&controller{}).logging
	_ middleware   = (&controller{}).tracing
)

// Data model for KLINE data
type Kline struct  {
	DATE string
	TS string
	OPEN float64
	CLOSE float64
	HIGH float64
	LOW float64
	VOLUME float64
}

// TODO Exchange, Output [GCP, Local etc], time [day, hour, minute], Symbol, date_from, date_to
func main(){
  // Get command line args
	// "Binance, BTCUSD, 1 min, date_from, date_to
	// TODO: translate args into creation of an exchange interface object, init it and download
}
