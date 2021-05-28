package export

import (
	"context"
	"fmt"
	"github.com/jmrtzsn/coiner/internal/model"
	"github.com/jmrtzsn/coiner/internal/projectpath"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

var data = []model.OHLCV{
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
)

func TestMain(m *testing.M) {
	log.Println("Setting up export testing suite!")
	records = model.ToCSV(data)
	if err := godotenv.Load(fmt.Sprintf("%s/.env", projectpath.Root)); err != nil {
		log.Fatal(err)
	}
	exitVal := m.Run()
	log.Println("Completed export testing suite!")
	os.Exit(exitVal)
}

func TestExportLocalCSV(t *testing.T) {
	file, err := CreateTempCSV(records)
	l := newLocal(test, test)
	assert.Nil(t, err)

	defer file.Close()
	defer os.Remove(file.Name())

	err = l.Export(file, test)
	assert.Nil(t, err)

	// Read and assert everything is correct
	path := dirPath(test, test) + fmt.Sprintf("%s.csv", test)
	got, err := l.Read(path)
	assert.Nil(t, err)

	if len(got) > 1 {
		assert.Equal(t, got[0], []string{"DATE", "TS", "OPEN", "CLOSE", "HIGH", "LOW", "VOLUME"})
		assert.Equal(t, got[1], []string{"2020-04-04T12:00:00Z", "1586001600", "6696.68000000", "6717.68000000", "6717.68000000", "6686.43000000", "155.99070000"})
	} else {
		t.Errorf("Read 0 records from file")
	}

	// Cleanup
	err = os.Remove(path)
	assert.Nil(t, err)
	err = os.Remove(dirPath(test, test))
	assert.Nil(t, err)
}

func TestExportBucketCSV(t *testing.T) {
	file, err := CreateTempCSV(records)
	assert.Nil(t, err)
	defer file.Close()
	defer os.Remove(file.Name())

	ctx := context.Background()
	bucket, err := New(ctx, test, test)
	assert.Nil(t, err)

	err = bucket.Export(file, test)
	assert.Nil(t, err)

	// Read and assert everything is correct
	got, err := bucket.Read(test)
	assert.Nil(t, err)
	if len(got) > 1 {
		assert.Equal(t, got[0], []string{"DATE", "TS", "OPEN", "CLOSE", "HIGH", "LOW", "VOLUME"})
		assert.Equal(t, got[1], []string{"2020-04-04T12:00:00Z", "1586001600", "6696.68000000", "6717.68000000", "6717.68000000", "6686.43000000", "155.99070000"})
	} else {
		t.Errorf("Read 0 records from gcp storage")
	}

	// Cleanup
	// TODO remove test file from GCP
	err = bucket.Delete(bucket.Path(test))
	assert.Nil(t, err)
}
