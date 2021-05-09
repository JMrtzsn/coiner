package export

import (
	"context"
	"fmt"
	"github.com/jmrtzsn/coiner/internal/model"
	"os"
)

// Command Command pattern
type Command interface {
	Export(*os.File) error
}

// TODO Implement command pattern
// Ensure the handlers implement the required interfaces/types at compile time
var (
	_ Command = &Local{}
	_ Command = &Storage{}
)

// Export to outputs [local, storage]
// TODO add format options
func Export(commands []Command, data []model.OHLCV) error {
	records := model.ToCSV(data)
	file, err := CreateTempCSV(records)
	if err != nil {
		return err
	}
	defer file.Close()
	defer os.Remove(file.Name())

	for _, c := range commands {
		if err := c.Export(file); err != nil {
			// TODO log why export failed, try again on scenarios + redundancy
			fmt.Printf("Failed to export: %s", err)
		}
	}
	return nil
}

func Commands(ctx context.Context, outputs []string, symbol string, date string) ([]Command, error) {
	var exports []Command
	for _, output := range outputs {
		switch output {
		case "local":
			local := newLocal(symbol, date)
			exports = append(exports, local)
		case "storage":
			storage, err := newStorage(ctx, storagePath(symbol, date))
			if err != nil {
				return nil, err
			}
			exports = append(exports, storage)
		default:
			// TODO
			fmt.Printf("Recieved invalid command option %s", output)
		}
	}
	return exports, nil
}
