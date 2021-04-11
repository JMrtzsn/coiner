package internal

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func Check(t *testing.T, err error) {
	if err != nil {
		t.Errorf("%v failed, got !nil error %v", t.Name(), err)
	}
}

func Compare(t *testing.T, got, want interface{}) {
	if !cmp.Equal(got, want) {
		t.Errorf("%v failed, Got %v want %v", t.Name(), got, want)
	}
}
