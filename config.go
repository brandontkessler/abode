package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config represents the configuration details. Path gets populated
// via the GetPath method and returns the real path for setting up
// dependent on the mode being test or prod.
type Config struct {
	Path          string
	SetupPath     string   `json:"setup_path"`
	TestSetupPath string   `json:"test_setup_path"`
	Structure     []string `json:"structure"`
}

// GetRealPath creates the correct path based on the mode passed in
// the config (prod/test). It also replaces the ~ alias path with
// the HOME environment variable.
func (c *Config) GetPath(mode string) {
	root := os.Getenv("HOME")

	switch mode {
	case "test":
		path := strings.Replace(c.TestSetupPath, "~", root, -1)
		c.Path = filepath.Join(path, "Abode")
	case "prod":
		path := strings.Replace(c.SetupPath, "~", root, -1)
		c.Path = filepath.Join(path, "Abode")
	default:
		path := strings.Replace(c.TestSetupPath, "~", root, -1)
		c.Path = filepath.Join(path, "Abode")
	}

	fmt.Sprintf("Created RealPath at: %s", c.Path)
}

// GetConfig reads the configuration file. It has one parameter called
// file that is the path and name of file containing config. It returns
// a Config struct.
func GetConfig(file string) Config {
	configFile, err := os.Open(file)

	if err != nil {
		panic(err)
	}

	defer configFile.Close()

	var c Config

	err = json.NewDecoder(configFile).Decode(&c)

	if err != nil {
		panic(err)
	}

	return c
}
