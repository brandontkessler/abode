package main

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

func RepoSetup(path string, repos []string) {
	var venvPaths []string

	for _, repo := range repos {
		var flags []string
		sp := strings.Split(repo, " ")
		url := sp[0]

		repoPath := strings.Split(url, "/")
		repoName := repoPath[len(repoPath)-1]
		proj := strings.Replace(repoName, ".git", "", -1)

		fp := filepath.Join(path, proj)

		if len(sp) > 1 {
			flags = sp[1:]
		}

		// clone repos
		clone := fmt.Sprintf("git clone %s %s", url, fp)
		cmd := exec.Command("bash", "-c", clone)
		err := cmd.Run()

		if err != nil {
			panic(err)
		}

		// Setup venv if -v flag exists and break
		for _, v := range flags {
			if v == "-v" {
				venvPaths = append(venvPaths, fp)
				break
			}
		}
	}

	// Use goroutine to install packages into venv concurrently
	var wg sync.WaitGroup
	wg.Add(len(venvPaths))

	for _, env := range venvPaths {
		go func(env string) {
			defer wg.Done()
			generateVenvs(env)
		}(env)
	}

	wg.Wait()
}

func generateVenvs(path string) {
	log.Printf("started: setting up venv at %s\n", path)
	Venv(path)
	log.Printf("finished: setting up venv at %s\n", path)
}

func Venv(path string) {
	mk := fmt.Sprintf("python3 -m venv %s/venv", path)
	act := fmt.Sprintf("source %s/venv/bin/activate", path)
	upg := "pip install --upgrade pip"
	reqs := fmt.Sprintf("pip install -r %s/requirements.txt", path)
	script := fmt.Sprintf("%s && %s && %s && %s", mk, act, upg, reqs)

	cmd := exec.Command("bash", "-c", script)
	err := cmd.Run()

	if err != nil {
		panic(err)
	}
}
