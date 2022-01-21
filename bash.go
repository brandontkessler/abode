package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
)

// profile is used to populate .bash_profile
// Its purpose is to source the bashrc file
const profile = `if [ -f ~/.bashrc ]; then
	source ~/.bashrc
fi
`

// rc is used to populate .bashrc
const rc = `if [ -f ~/.bash_prompt ]; then
. ~/.bash_prompt
fi

if [ -f ~/.bash_aliases ]; then
. ~/.bash_aliases
fi

if [ -f ~/.bash_functions ]; then
. ~/.bash_functions
fi

PYTHONPATH=$PYTHONPATH:<workspace/code>
export PYTHONPATH

PATH=$PATH:$HOME/bin
export PATH`

// prompt is used to populate .bash_prompt
const prompt = `blue=$(tput setaf 33);
yellow=$(tput setaf 11);
red=$(tput setaf 124);
green=$(tput setaf 64);
white=$(tput setaf 7);
bold=$(tput bold);
reset=$(tput sgr0);

PS1="\[${blue}\]\u";
PS1+="\[${white}\]@";
PS1+="\[${yellow}\]\h";
PS1+="\[${white}\]: ";
PS1+="\[${green}\]\w";
PS1+="\n";
PS1+="\[${white}\]$ \[${reset}\]";

export PS1;`

// alias is used to populate .bash_aliases
const alias = `# Set Python
alias python=python3

alias la="ls -la"`

// Bash sets up the terminal profile and bash settings using .bashrc,
// .bash_profile, .bash_aliases, .bash_prompt
func Bash(c Config) {
	if c.GenerateBash {
		SetTerminalSettings(c.TerminalProfile)

		abodeAlias := fmt.Sprint(`alias abode="cd `, c.Path, `"`)

		bf := map[string]string{
			".bash_profile": profile,
			".bashrc":       rc,
			".bash_prompt":  prompt,
			".bash_aliases": fmt.Sprint(alias, "\n\n", abodeAlias),
		}

		// Loops through bf (bashFiles) to create the path which is
		// the map key and add create/append the file with the text
		for f, text := range bf {
			path := filepath.Join(c.HomeDir, f)

			err := CheckPathAddText(path, text)

			if err != nil {
				panic(err)
			}
		}
	}
}

// SetTerminalSettings sets the default terminal profile and the startup window
// settings as what is provided in Config.TerminalProfile
func SetTerminalSettings(profile string) {
	baseCmd := "defaults write com.apple.Terminal"
	defSetting := fmt.Sprint("\"Default Window Settings\"", " ", profile)
	startupSetting := fmt.Sprint("\"Startup Window Settings\"", " ", profile)

	defCmd := exec.Command("bash", "-c", fmt.Sprint(baseCmd, " ", defSetting))
	defErr := defCmd.Run()

	if defErr != nil {
		panic(defErr)
	}

	strtCmd := exec.Command("bash", "-c", fmt.Sprint(baseCmd, " ", startupSetting))
	strtErr := strtCmd.Run()

	if strtErr != nil {
		panic(strtErr)
	}
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
