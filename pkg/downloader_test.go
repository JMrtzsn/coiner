package pkg_test

import (
	"fmt"
	"github.com/jmrtzsn/coiner/cmd"
	"github.com/jmrtzsn/coiner/internal/export"
	"github.com/jmrtzsn/coiner/pkg"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

var downloader pkg.Downloader

func TestMain(m *testing.M) {
	log.Println("Setting up download testing suite!")
	cmd.LoadConfig("test")
	downloader = cmd.ToDownloader()
	exitVal := m.Run()
	log.Println("Completed download testing suite!")
	os.Exit(exitVal)
}

func TestDownloader_Download(t *testing.T) {
	// Will download 2 days
	downloader.Download()

	// Read and assert everything is correct
	local, _ := downloader.Exports[0].(*export.Local)

	csvPath := local.DirPath("BTCUSDT") + fmt.Sprintf("/%s.csv", "2019-01-01")
	got, err := local.Read(csvPath)
	assert.Nil(t, err)

	if len(got) > 1 {
		assert.Equal(t, []string{"DATE", "TS", "OPEN", "CLOSE", "HIGH", "LOW", "VOLUME"}, got[0])
		assert.Equal(t, []string{"2019-01-01T00:00:00Z", "1546300800", "3701.23000000", "3702.46000000", "3703.72000000", "3701.09000000", "17.10011000"}, got[1])
	} else {
		t.Errorf("Read 0 records from file")
	}

	// Cleanup

	err = os.RemoveAll(local.DirPath("BTCUSDT"))
	assert.Nil(t, err)
	err = os.RemoveAll(local.DirPath("ETHUSDT"))
	assert.Nil(t, err)
}
