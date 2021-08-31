package main

import (
	"encoding/json"
	"fmt"
	"os"
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
			p = path.Join("_keys", strings.Join([]string{
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

func writeAll(b *strings.Builder, val ...string) {
	for _, v := range val {
		(*b).WriteString(v)
	}
}

func combine(val []string) string {
	var b strings.Builder
	for i, v := range val {
		if i != 0 {
			b.WriteRune(' ')
		}
		// escape quotes
		v = strings.Replace(v, `"`, `\"`, -1)
		if strings.Contains(v, " ") {
			b.WriteString(strconv.Quote(v))
		} else {
			b.WriteString(v)
		}
	}
	return b.String()
}

func readAbs(server string) (*Node,error) {
	data, err := os.ReadFile(path.Join("_nodes"))
}

func readRel(server string) (*Node,error) {
	return readAbs(path.Join("_nodes", ))
}

func main() {
	data, err := os.ReadFile(path.Join("_nodes", "node05.json"))
	if err != nil {
		panic(err)
	}
	n := new(Node)
	if err = json.Unmarshal(data, n); err != nil {
		panic(err)
	}
	gen := n.SCPCommand("~/local.file", "~/remote.file")
	fmt.Println("Generated:", gen)
	fmt.Println("Combined:", combine(gen))
}
