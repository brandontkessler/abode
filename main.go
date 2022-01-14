package main

import (
	"fmt"
)

const mode = "test"

func main() {
	config := GetConfig("config.json")
	config.GetPath(mode)

	fmt.Println(config)

	// Structure(config)

	// if mode == "test" {
	// 	Teardown(config)
	// }

	// git config --global init.defaultBranch <name>
	// add this as part of git setup
}
