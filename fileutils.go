package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// Check if path exists, if not create it
func CheckDirAndMake(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)

		return err
	}
	return nil
}

// CheckPathExists returns a bool indicating if path exists
func CheckPathExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}

// CheckPathAddText tries to read a file path, if it can, it adds
// additional text to the file. Otherwise, it creates the file and
// writes the text to it.
func CheckPathAddText(filepath, text string) error {
	f, err := ioutil.ReadFile(filepath)

	if err != nil {
		newText := fmt.Sprint("#!/bin/bash\n\n", text)

		err = ioutil.WriteFile(filepath, []byte(newText), 0644)

		return err
	} else {
		newText := fmt.Sprint(string(f), "\n\n", text)

		err = ioutil.WriteFile(filepath, []byte(newText), 0644)

		return err
	}
}

// CheckStringInFile takes a path to a file and some text. It returns
// an error if that file contains the string, otherwise it returns nil.
func CheckStringInFile(filepath, text string) error {
	f, err := ioutil.ReadFile(filepath)

	if err != nil {
		log.Fatal(err)
	}

	ctx := string(f)

	if strings.Contains(ctx, text) {
		return fmt.Errorf("file already contains: %s", text)
	} else {
		return nil
	}
}
