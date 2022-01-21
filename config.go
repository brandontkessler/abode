package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// Config represents the configuration details. Path gets populated
// via the GetPath method and returns the real path for setting up
// dependent on the mode being test or prod. HomeDir represents the
// HOME directory of the system. If in test mode, HomeDir will be
// the same as Path.
type Config struct {
	Path             string
	HomeDir          string
	SetupPath        string            `json:"setup_path"`
	TestSetupPath    string            `json:"test_setup_path"`
	Structure        []string          `json:"structure"`
	VsCodeSettings   map[string]string `json:"vscode_settings"`
	VsCodeExtensions []string          `json:"vscode_extensions"`
	GenerateBash     bool              `json:"generate_bash_files"`
	TerminalProfile  string            `json:"terminal_profile"`
}

// GetPath creates the correct path based on the mode passed in
// the config (prod/test). It also replaces the ~ alias path with
// the HOME environment variable.
func (c *Config) GetPath(mode string) {
	root := os.Getenv("HOME")

	switch mode {
	case "test":
		path := strings.Replace(c.TestSetupPath, "~", root, -1)
		c.HomeDir = filepath.Join(path, "Abode")
		c.Path = filepath.Join(path, "Abode")
	case "prod":
		path := strings.Replace(c.SetupPath, "~", root, -1)
		c.HomeDir = root
		c.Path = filepath.Join(path, "Abode")
	default:
		path := strings.Replace(c.TestSetupPath, "~", root, -1)
		c.HomeDir = filepath.Join(path, "Abode")
		c.Path = filepath.Join(path, "Abode")
	}
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

func GetBasePaths(c Config) []string {
	base_paths := []string{}
	for _, v := range c.Structure {
		base := strings.Split(v, "/")[0]
		base_paths = append(base_paths, base)
	}

	return base_paths
}
