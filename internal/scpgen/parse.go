package scpgen

import (
	"encoding/json"
	"os"
	"path"
	"strconv"
	"strings"
)

const KeyDir = "_nodes"

var nodeCache = make(map[string]*Node)

func ReadAbs(path string) (n *Node, err error) {
	var data []byte
	if data, err = os.ReadFile(path); err != nil {
		return
	}
	// update raw param
	if err = json.Unmarshal(data, &n); err == nil {
		var raw string
		// remove leading path
		// _nodes/node05.json -> node05.json
		if strings.ContainsRune(path, os.PathSeparator) {
			raw = path[strings.LastIndex(path, string(os.PathSeparator))+1:]
		}
		// remove .json extension
		// node05.json -> node05
		if strings.HasSuffix(raw, ".json") {
			raw = raw[:len(raw)-5]
		}
		n.Raw = raw
	}
	return
}

func ReadRel(server string) (*Node, error) {
	// in "cache"?
	if n, ok := nodeCache[server]; ok {
		return n, nil
	}
	// has extension .json?
	if !strings.HasSuffix(server, ".json") {
		server += ".json"
	}
	return ReadAbs(path.Join("_nodes", server))
}

func writeAll(b *strings.Builder, val ...string) {
	for _, v := range val {
		(*b).WriteString(v)
	}
}

func Combine(val []string) string {
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
