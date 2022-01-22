package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const mode = "test"

func main() {
	start := time.Now()
	// Get config settings
	config := GetConfig("config.json")
	config.GetPath(mode)

	// start with teardown
	if mode == "test" {
		Teardown(config)
	}

	// Build the folder structure
	Structure(config)

	// Loops through each base path that was built during the Structure
	// function and adds a unique .code-workspace file that can be used
	// to open a different vscode workspace.
	for _, path := range GetBasePaths(config) {
		codeWorkspace := MakeVsWorkspace(config, path)
		fileName := fmt.Sprintf("%s.code-workspace", path)
		workspaceFile, err := os.Create(filepath.Join(config.Path, path, fileName))

		if err != nil {
			fmt.Println("Error creating code-workspace file:", err)
		}

		encoder := json.NewEncoder(workspaceFile)
		encoder.SetIndent("", "    ")
		err = encoder.Encode(codeWorkspace)

		if err != nil {
			fmt.Println("Error encoding json to code-workspace file:", err)
		}

		err = workspaceFile.Close()

		if err != nil {
			fmt.Println("Error closing code-workspace file:", err)
		}
	}

	// Sets up bash and terminal settings
	Bash(config)

	// Source .bashrc if in production
	if mode == "prod" {
		rcPath := filepath.Join(config.HomeDir, ".bashrc")
		cmd := exec.Command("bash", "-c", fmt.Sprint("source ", rcPath))
		err := cmd.Run()

		if err != nil {
			panic(err)
		}
	}

	sshPath := filepath.Join(config.HomeDir, ".ssh")
	keyGenPassword := ""
	err := GenKeys(sshPath, "idrsa", keyGenPassword)

	if err != nil {
		panic(err)
	}

	err = SshConfig(sshPath)

	if err != nil {
		panic(err)
	}

	// Set up Git projects and Notes
	projPath := getPath(config.Structure, "code", "projects")

	if projPath == "" {
		panic(fmt.Errorf("no code -> project path for git repos"))
	}

	projLoc := filepath.Join(config.Path, projPath)
	RepoSetup(projLoc, config.ProjectGitRepos)

	notePath := getPath(config.Structure, "notes")

	if notePath == "" {
		panic(fmt.Errorf("no notes path for git repo"))
	}

	noteLoc := filepath.Join(config.Path, notePath)
	RepoSetup(noteLoc, config.NotesGitRepos)

	// git config --global init.defaultBranch main
	// add this as part of git setup
	fmt.Println(time.Since(start))
}

// getPath finds the first path of a list of paths that contains strings
// provided as args.
// Example, if Paths = ["foo/bar", "hello", "hello/world"],
// getPath(Paths, "hello", "world") returns "hello/world"
func getPath(pathList []string, contains ...string) string {
	for _, path := range pathList {
		matches := 0

		for _, arg := range contains {
			if strings.Contains(path, arg) {
				matches++
			}

			if matches == len(contains) {
				return path
			}
		}
	}

	return ""
}
