package internal

import "github.com/jmrtzsn/coiner/internal/binance"

type Config interface {
	LoadEnv()
}

// Ensure the exchanges implement the required interfaces/types at compile time
var (
	_ Config = &binance.Config{}
)

