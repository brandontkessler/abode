package main

import (
	"os"
	"path/filepath"
)

// Structure builds the folder structure within the base directory
// using the structure provided in the config file.
func Structure(c Config) {
	for _, v := range c.Structure {
		newPath := filepath.Join(c.Path, v)
		err := os.MkdirAll(newPath, os.ModePerm)

		if err != nil {
			panic(err)
		}
	}
}
