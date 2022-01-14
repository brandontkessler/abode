package main

import (
	"os"
	"path/filepath"
)

func Structure(c Config) {
	for _, v := range c.Structure {
		newPath := filepath.Join(c.Path, v)
		err := os.MkdirAll(newPath, os.ModePerm)

		if err != nil {
			panic(err)
		}
	}
}
