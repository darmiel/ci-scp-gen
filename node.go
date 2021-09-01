package main

import (
	"path"
	"strconv"
	"strings"
)

type Node struct {
	// Raw: name of file
	Raw string `json:"-"`

	Friendly string `json:"friendly"`
	SSH      struct {
		Host string `json:"host"`
		Port uint8  `json:"port"`
		User string `json:"user"`
		Auth struct {
			Key     bool   `json:"key"`
			KeyPath string `json:"key_path"`
			Pass    string `json:"pass"`
		} `json:"auth"`
	} `json:"ssh"`
}

// generates:
// scp -i <key> -P <port> <lfile> user@host:<rfile>
// or:
// sshpass -p <password> scp -P <port> <lfile> user@host:<rfile>
func (n *Node) SCPCommand(lfile, rfile string) (args []string) {
	// do we need sshpass?
	if !n.SSH.Auth.Key {
		// no, we'll only use password for authentication
		// -> sshpass
		args = append(args, "sshpass", "-p", n.SSH.Auth.Pass)
	}

	// base scp command
	args = append(args, "scp", "-P", strconv.Itoa(int(n.SSH.Port)))

	if n.SSH.Auth.Key {
		// key path
		var p string
		if p = n.SSH.Auth.KeyPath; p == "" {
			// auto path
			p = path.Join(KeyDir, strings.Join([]string{
				n.Raw,
				".pri",
			}, ""))
		}

		args = append(args, "-i", p)
	}

	// append lfile
	args = append(args, lfile)

	// create user@host:file
	var b strings.Builder
	writeAll(&b, n.SSH.User, "@", n.SSH.Host, ":", rfile)
	args = append(args, b.String())

	return
}
