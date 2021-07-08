package pkg_test

import (
	"context"
	"fmt"
	"github.com/jmrtzsn/coiner/cmd"
	"github.com/jmrtzsn/coiner/internal/export"
	"github.com/jmrtzsn/coiner/internal/projectpath"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	log.Println("Setting up download testing suite!")
	if err := godotenv.Load(fmt.Sprintf("%s/test.env", projectpath.Root)); err != nil {
		log.Fatal(err)
	}
	exitVal := m.Run()
	log.Println("Completed download testing suite!")
	os.Exit(exitVal)
}

func TestDownload(t *testing.T) {
	// Will download 2 days
	cmd.LoadConfig("test")
	conf := cmd.UnMarshal()
	downloader, err := conf.NewDownloader(context.Background())
	assert.Nil(t, err)
	downloader.Download()

	date1 := "2019-01-01"
	date2 := "2019-01-02"
	btcusdt := "BTCUSDT"
	ethusdt := "ETHUSDT"

	t.Run("Local", func(t *testing.T) {
		// Read and assert everything is correct
		local, _ := downloader.Exports[0].(*export.Local)

		btc := local.DirPath(btcusdt) + fmt.Sprintf("/%s.csv", date1)
		got, err := local.Read(btc)
		assert.Nil(t, err)
		if len(got) > 1 {
			assert.Equal(t, []string{"2019-01-01T00:00:00Z", "1546300800", "3701.23000000", "3702.46000000", "3703.72000000", "3701.09000000", "17.10011000"}, got[1])
		} else {
			t.Errorf("Read 0 records from file")
		}

		eth := local.DirPath(ethusdt) + fmt.Sprintf("/%s.csv", date2)
		got, err = local.Read(eth)
		assert.Nil(t, err)
		if len(got) > 1 {
			assert.Equal(t, []string{"2019-01-02T00:00:00Z", "1546387200", "139.10000000", "139.43000000", "140.00000000", "139.06000000", "3180.53061000"}, got[1])
		} else {
			t.Errorf("Read 0 records from file")
		}

		// Cleanup
		err = os.RemoveAll(local.DirPath(btcusdt))
		assert.Nil(t, err)
		err = os.RemoveAll(local.DirPath(ethusdt))
		assert.Nil(t, err)
	})

	t.Run("Bucket", func(t *testing.T) {

		// Read and assert everything is correct
		bucket, _ := downloader.Exports[1].(*export.Bucket)
		got, err := bucket.Read(date1, btcusdt)
		assert.Nil(t, err)
		if len(got) > 1 {
			assert.Equal(t, []string{"2019-01-01T00:00:00Z", "1546300800", "3701.23000000", "3702.46000000", "3703.72000000", "3701.09000000", "17.10011000"}, got[1])
		} else {
			t.Errorf("Read 0 records from gcp storage")
		}

		got, err = bucket.Read(date2, ethusdt)
		assert.Nil(t, err)
		if len(got) > 1 {
			assert.Equal(t, []string{"2019-01-02T00:00:00Z", "1546387200", "139.10000000", "139.43000000", "140.00000000", "139.06000000", "3180.53061000"}, got[1])
		} else {
			t.Errorf("Read 0 records from gcp storage")
		}
	})

}
