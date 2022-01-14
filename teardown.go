package main

import (
	"os"
)

func Teardown(c Config) {
	err := os.RemoveAll(c.Path)

	if err != nil {
		panic(err)
	}
}
