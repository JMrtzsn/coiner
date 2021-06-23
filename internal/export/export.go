package export

import (
	"os"
)

type Export interface {
	// USAGE EXPORT(file, symbol, date)
	Export(*os.File, string, string) error
	String() string
}

// Ensure the handlers implement the required interfaces/types at compile time
var (
	_ Export = &Local{}
	_ Export = &Bucket{}
)
