package exchange

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsoToUnix(t *testing.T) {
	input := "2020-04-04T12:07:00Z"
	var want int64 = 1586002020000
	got, err := IsoToUnix(input)
	assert.Nil(t, err)
	assert.Equal(t, want,got)
}

func TestUnixToISO(t *testing.T) {
	var input int64 = 1586002020000
	var want = "2020-04-04T12:07:00Z"
	got := UnixToISO(input)
	assert.Equal(t, want,got)
}
