package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const mode = "test"

func main() {
	// Get config settings
	config := GetConfig("config.json")
	config.GetPath(mode)

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

	// if mode == "test" {
	// 	Teardown(config)
	// }

	// git config --global init.defaultBranch <name>
	// add this as part of git setup
}
