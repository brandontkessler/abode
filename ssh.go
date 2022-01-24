package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
)

// sshConfig populates the config file in .ssh folder allowing easy ssh
// functionality. It is populated with examples that can be uncommented
// and used as necessary.
const sshConfigText = `### Basic Host ###
# Host dev
#     HostName dev.example.com
#     User john
#     Port 2322

### Multi-Jump Host (Jump first to mux1 then to mta1) ###
# Host jumphost
# 	HostName mux1.accuenbi.com
# 	User brandon_kessler
# 	Port 22
# 	IdentityFile ~/.ssh/id_rsa
# 	RequestTTY force
# 	RemoteCommand ssh mta1

### Proxy Jump ###
# Host *
# 	ForwardAgent yes

# Host bastion
# 	Hostname public.domain.com
# 	User alex
# 	Port 50482
# 	IdentityFile ~/.ssh/id_ed25519

# Host lanserver
# 	Hostname 192.168.1.1
# 	User alex
# 	ProxyJump bastion`

// GenKeys takes a path and keyname to generate a pub/priv keypair
func GenKeys(path, keyname, passphrase string) error {
	err := CheckDirAndMake(path)

	if err != nil {
		return fmt.Errorf("cannot check and make directory")
	}

	fp := filepath.Join(path, keyname)

	if pe := CheckPathExists(fp); pe {
		fmt.Printf("key already exists: %s\n", fp)
		return nil
	}

	key := fmt.Sprint("ssh-keygen -t rsa -b 4096 -f ", fp, " -q -N \"", passphrase, "\"")

	cmd := exec.Command("bash", "-c", key)
	err = cmd.Run()

	return err
}

// SshConfig adds a config file to .ssh folder if not already
// exists. It uses the const sshConfigText to populate the file. If the
// additional text already exists in the file, do nothing.
func SshConfig(path string) error {
	fp := filepath.Join(path, "config")
	txt := sshConfigText

	// if err != nil, file does not exist
	if f, err := ioutil.ReadFile(fp); err != nil {
		err = ioutil.WriteFile(fp, []byte(txt), 0644)

		return err
	} else {
		err := CheckStringInFile(fp, txt)

		if err != nil {
			return nil
		}

		newTxt := fmt.Sprint(string(f), "\n\n", txt)

		err = ioutil.WriteFile(fp, []byte(newTxt), 0644)

		return err
	}
}
