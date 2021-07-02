package export

import (
	"os"
)

type Export interface {
	// TODO: pass buffer?
	Export(csv *os.File, symbol string, date string) error
	String() string
}

// Ensure the handlers implement the required interfaces/types at compile time
var (
	_ Export = &Local{}
	_ Export = &Bucket{}
)
