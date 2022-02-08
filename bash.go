package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// profile is used to populate .bash_profile
// Its purpose is to source the bashrc file
const profile = `if [ -f ~/.bashrc ]; then
	source ~/.bashrc
fi
`

// rc is used to populate .bashrc
const rc = `if [ -f <<path>>/.bash_prompt ]; then
. <<path>>/.bash_prompt
fi

if [ -f <<path>>/.bash_aliases ]; then
. <<path>>/.bash_aliases
fi

if [ -f <<path>>/.bash_functions ]; then
. <<path>>/.bash_functions
fi

if [ -f <<path>>/.bash_paths ]; then
. <<path>>/.bash_paths
fi`

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

// paths is used to populate .bash_paths. The third line of PATH, removes any
// duplicates within the PATH.
const paths = `# Set paths
PYTHONPATH=$PYTHONPATH:<<codePath>>:<<workCodePath>>
export PYTHONPATH

PATH=$PATH:$HOME/bin
PATH=$PATH:$HOME/go/bin
PATH="$(perl -e 'print join(":", grep { not $seen{$_}++ } split(/:/, $ENV{PATH}))')"
export PATH`

// Bash sets up the terminal profile and bash settings using: .bashrc,
// .bash_profile, .bash_aliases, .bash_prompt, .bash_paths. Files are only
// created if they don't already exist otherwise text is appended to the end.
// The function of .bashrc and .bash_profile is to sit in the root of the OS
// and source other files which sit in the Abode path. This keeps all files
// clean and purposeful. If the file sourcing text already exists in .bashrc
// or .bash_profile, skip the append to keep the file clean.
func Bash(c Config) {
	SetTerminalSettings(c.TerminalProfile)

	abodeAlias := fmt.Sprint(`alias abode="cd `, c.Path, `"`)
	codePath := filepath.Join(c.Path, "code", "projects")
	workCodePath := filepath.Join(c.Path, "work", "code", "projects")

	replacer := strings.NewReplacer("<<codePath>>", codePath, "<<workCodePath>>", workCodePath)

	bf := map[string]string{
		".bashrc":       strings.Replace(rc, "<<path>>", c.Path, -1),
		".bash_profile": profile,
		".bash_prompt":  prompt,
		".bash_aliases": fmt.Sprint(alias, "\n\n", abodeAlias),
		".bash_paths":   replacer.Replace(paths),
	}

	// Loops through bf (bashFiles) to create the path which is
	// the map key and add create/append the file with the text
	for f, text := range bf {
		var path string

		if f == ".bashrc" || f == ".bash_profile" {
			path = filepath.Join(c.HomeDir, f)
			err := CheckStringInFile(path, text)

			if err != nil {
				continue
			}
		} else {
			path = filepath.Join(c.Path, f)
		}

		err := CheckPathAddText(path, text)

		if err != nil {
			panic(err)
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
