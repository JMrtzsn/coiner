package internal

import "github.com/jmrtzsn/coiner/internal/binance"

type Config interface {
	LoadEnv()
}

// Ensure the exchanges implement the required interfaces/types at compile time
var (
	_ Config = &binance.Config{}
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


