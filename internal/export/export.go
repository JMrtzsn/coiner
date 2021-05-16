package export

import (
	"context"
	"fmt"
	"github.com/jmrtzsn/coiner/internal/model"
	"os"
)

// Command Command pattern
type Command interface {
	Export(*os.File, string) error
}

// TODO Implement command pattern
// Ensure the handlers implement the required interfaces/types at compile time
var (
	_ Command = &Local{}
	_ Command = &Storage{}
)

// Export to outputs [local, storage] Usage: export symbol data using
func Export(commands []Command, data []model.OHLCV, date string) error {
	records := model.ToCSV(data)
	file, err := CreateTempCSV(records)
	if err != nil {
		return err
	}
	defer file.Close()
	defer os.Remove(file.Name())

	for _, c := range commands {
		if err := c.Export(file, date); err != nil {
			// TODO log why export failed, try again on scenarios + redundancy
			fmt.Printf("Failed to export: %s", err)
		}
	}
	return nil
}

func CreateCommands(ctx context.Context, outputs, symbols []string, exchange string) (map[string][]Command, error) {
	exports := make(map[string][]Command)
	for _, symbol := range symbols {
		for _, output := range outputs {
			var commands []Command
			switch output {
			case "local":
				local := newLocal(exchange, symbol)
				commands = append(commands, local)
			case "storage":
				storage, err := newStorage(ctx, exchange, symbol)
				if err != nil {
					return nil, err
				}
				commands = append(commands, storage)
			default:
				return nil, fmt.Errorf("recevied invalid output %s", output)
			}
			exports[symbol] = commands
		}
	}
	return exports, nil
}
