package export

import (
	"fmt"
	"github.com/jmrtzsn/coiner/internal"
	"github.com/jmrtzsn/coiner/internal/export/local"
	"github.com/jmrtzsn/coiner/internal/export/storage"
	"github.com/jmrtzsn/coiner/internal/projectpath"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
)

var data = []internal.OHLCV{
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
var (
	records [][]string
	test    = "test"
	bucket  = os.Getenv("CLOUD_BUCKET")
)

func TestMain(m *testing.M) {
	log.Println("Setting up export testing suite!")
	records = internal.ToCSV(data)
	CreateTempFile()
	if err := godotenv.Load(fmt.Sprintf("%s/.env",projectpath.Root)); err != nil {
		log.Fatal(err)
	}
	bucket = os.Getenv("CLOUD_BUCKET")
	exitVal := m.Run()
	log.Println("Completed export testing suite!")
	os.Exit(exitVal)
}

func TestExportLocalCSV(t *testing.T) {
	err := local.Write(records, test, test)
	internal.Check(t, err)

	// Read and assert everything is correct
	file := local.CreateFilepath(test, test)
	got, err := local.Read(file)
	internal.Check(t, err)
	internal.Compare(t, got[0], []string{"DATE", "TS", "OPEN", "CLOSE", "HIGH", "LOW", "VOLUME"})
	internal.Compare(t, got[1], []string{"2020-04-04T12:00:00Z", "1586001600", "6696.68000000", "6717.68000000", "6717.68000000", "6686.43000000", "155.99070000"})
	// Cleanup
	err = os.Remove(file)
	internal.Check(t, err)
	err = os.Remove(local.CreateDirpath(test))
	internal.Check(t, err)
}

func TestExportStorageCSV(t *testing.T) {
	var s storage.Storage
	err := s.Init(bucket)
	internal.Check(t, err)
	path := storage.Path(test, test)
	err = s.Write(path, file, records)
	internal.Check(t, err)

	// Read and assert everything is correct
	path := storage.Path(test, test)
	got, err := s.Read(path)
	internal.Check(t, err)
	internal.Compare(t, got[0], []string{"DATE", "TS", "OPEN", "CLOSE", "HIGH", "LOW", "VOLUME"})
	internal.Compare(t, got[1], []string{"2020-04-04T12:00:00Z", "1586001600", "6696.68000000", "6717.68000000", "6717.68000000", "6686.43000000", "155.99070000"})

	// Cleanup
	// TODO remove test file from GCP
	err = s.Delete(path)
	internal.Check(t, err)
}
