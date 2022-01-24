package main

import (
	"encoding/json"
	"log"
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
	Name             string            `json:"name"`
	SetupPath        string            `json:"setup_path"`
	Structure        []string          `json:"structure"`
	VsCodeSettings   map[string]string `json:"vscode_settings"`
	VsCodeExtensions []string          `json:"vscode_extensions"`
	TerminalProfile  string            `json:"terminal_profile"`
	ProjectGitRepos  []string          `json:"project_git_repos"`
	NotesGitRepos    []string          `json:"notes_git_repos"`
}

// GetPath creates the correct path based on the mode passed in
// the config (prod/test). It also replaces the ~ alias path with
// the HOME environment variable.
func (c *Config) GetPath(mode string) {
	root := os.Getenv("HOME")

	wd, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	path := strings.Replace(wd, "~", root, -1)
	c.HomeDir = root
	c.Path = filepath.Join(path, c.Name)
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
