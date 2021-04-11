package internal

import (
	"log"
	"os"
	"testing"
)

var data = []OHLCV{
	{
		DATE:   "2020-04-04T12:00:00Z",
		TS:     "1586001600",
		OPEN:   "6696.68000000",
		CLOSE:  "6717.68000000",
		HIGH:   "6717.68000000",
		LOW:    "6686.43000000",
		VOLUME: "155.99070000",
	},
	{
		DATE:   "2020-04-04T12:01:00Z",
		TS:     "1586001660",
		OPEN:   "6716.03000000",
		CLOSE:  "6720.17000000",
		HIGH:   "6722.40000000",
		LOW:    "6707.45000000",
		VOLUME: "173.49916200",
	},
	{
		DATE:   "2020-04-04T12:02:00Z",
		TS:     "1586001720",
		OPEN:   "6720.18000000",
		CLOSE:  "6715.26000000",
		HIGH:   "6722.05000000",
		LOW:    "6711.59000000",
		VOLUME: "56.91619600",
	},
}

func TestMain(m *testing.M) {
	log.Println("Setting up internal testing suite!")
	exitVal := m.Run()
	log.Println("Completed internal testing suite!")
	os.Exit(exitVal)
}

func Test_OHLCV_ToCSV(t *testing.T) {
	want := []string{"2020-04-04T12:00:00Z", "1586001600", "6696.68000000", "6717.68000000", "6717.68000000", "6686.43000000", "155.99070000"}
	got := data[0].csv()
	Compare(t, got, want)
}

func Test_toCSV(t *testing.T) {
	got := toCSV(data)
	Compare(t, got[0], []string{"DATE", "TS", "OPEN", "CLOSE", "HIGH", "LOW", "VOLUME"})
	Compare(t, got[1], []string{"2020-04-04T12:00:00Z", "1586001600", "6696.68000000", "6717.68000000", "6717.68000000", "6686.43000000", "155.99070000"})
}

func Test_saveCSVLocal(t *testing.T) {
	records := toCSV(data)
	test := "test"
	err := SaveCSV(records, test)
	Check(t, err)

	got := OpenCSV("test")

	Compare(t, got[0], []string{"DATE", "TS", "OPEN", "CLOSE", "HIGH", "LOW", "VOLUME"})
	Compare(t, got[1], []string{"2020-04-04T12:00:00Z", "1586001600", "6696.68000000", "6717.68000000", "6717.68000000", "6686.43000000", "155.99070000"})

	err = os.Remove(GetFilepath(test))
	Check(t, err)
}