package main

import (
	"github.com/jmrtzsn/coiner/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
