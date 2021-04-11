package internal

import (
	"encoding/csv"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/jmrtzsn/coiner/internal/projectpath"
	"log"
	"os"
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


func OpenCSV(name string) [][]string {
	csvfile, err := os.Open(GetFilepath(name))
	defer csvfile.Close()

	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	r := csv.NewReader(csvfile)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatalln("Couldn't read the csv file", err)
	}
	return records
}

func SaveCSV(records [][]string, name string) error {
	file, err := os.Create(GetFilepath(name))
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	err = w.WriteAll(records) // calls Flush internally

	return nil
}

func GetFilepath(name string) string {
	return fmt.Sprintf("%s/data/%s.csv", projectpath.Root, name)
}
