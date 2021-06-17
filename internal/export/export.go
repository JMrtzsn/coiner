package export

import (
	"os"
)

type Export interface {
	Export(*os.File, string) error
	String() string
}

// Ensure the handlers implement the required interfaces/types at compile time
var (
	_ Export = &Local{}
	_ Export = &Bucket{}
)
