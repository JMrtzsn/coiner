package export

import (
	"io/ioutil"
	"os"
)

// TODO: create export ... interface? or struct a common gateway to all exports
// TODO should contain common functions used..? by downstream packages? or have those packages further down?
type Export interface {
	Export()
	Read()
}

func CreateTempFile() (*os.File, error) {
	file, err := ioutil.TempFile("", "file")
	if err != nil {
		return nil, err
	}
	return file, nil
}
