package export

import (
	"os"
)

// Command Command pattern
type Export interface {
	// Export file to export, date string
	Export(*os.File, string) error
}

// TODO Implement command pattern
// Ensure the handlers implement the required interfaces/types at compile time
var (
	_ Export = &Local{}
	_ Export = &Bucket{}
)
