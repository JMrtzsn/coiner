package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCSV(t *testing.T) {
	var data = []Candle{
		{
			DATE:   "2020-04-04T12:00:00Z",
			TS:     "1586001600",
			OPEN:   "6696.68000000",
			CLOSE:  "6717.68000000",
			HIGH:   "6717.68000000",
			LOW:    "6686.43000000",
			VOLUME: "155.99070000",
		},
	}
	got := RecordsWithHeader()
	for _, row := range data {
		got = append(got, row.CSV())
	}
	assert.Equal(t, got[0], []string{"DATE", "TS", "OPEN", "CLOSE", "HIGH", "LOW", "VOLUME"})
	assert.Equal(t, got[1], []string{"2020-04-04T12:00:00Z", "1586001600", "6696.68000000", "6717.68000000", "6717.68000000", "6686.43000000", "155.99070000"})
}
