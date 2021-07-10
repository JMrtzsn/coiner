package downloader_test

import (
	"context"
	"fmt"
	"github.com/jmrtzsn/coiner/cmd"
	"github.com/jmrtzsn/coiner/pkg/export"
	"github.com/jmrtzsn/coiner/pkg/projectpath"
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
	eos := "EOSUSDT"

	t.Run("Local", func(t *testing.T) {
		// Read and assert everything is correct
		local, _ := downloader.Exports[0].(*export.Local)

		btc := local.DirPath(eos) + fmt.Sprintf("/%s.csv", date1)
		got, err := local.Read(btc)
		assert.Nil(t, err)
		if len(got) > 1 {
			assert.Equal(t, []string{"2019-01-01T00:00:00Z", "1546300800", "2.53450000", "2.53860000", "2.54010000", "2.53310000", "2838.73000000"}, got[1])
		} else {
			t.Errorf("Read 0 records from file")
		}
		// Cleanup
		err = os.RemoveAll(local.DirPath(eos))
		assert.Nil(t, err)
	})

	t.Run("Bucket", func(t *testing.T) {

		// Read and assert everything is correct
		bucket, _ := downloader.Exports[1].(*export.Bucket)
		got, err := bucket.Read(date1, "EOSUSDT")
		assert.Nil(t, err)
		if len(got) > 1 {
			assert.Equal(t, []string{"2019-01-01T00:00:00Z", "1546300800", "2.53450000", "2.53860000", "2.54010000", "2.53310000", "2838.73000000"}, got[1])
		} else {
			t.Errorf("Read 0 records from gcp storage")
		}
	})

}
