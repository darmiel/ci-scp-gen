package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

const (
	NodeDir = "_nodes"
	KeyDir  = "_keys"
)

var nodes = make(map[string]*Node)

func init() {
	if x := filepath.WalkDir(NodeDir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		n, err := readRel(entry.Name())
		if err != nil {
			return err
		}
		nodes[n.Raw] = n
		fmt.Printf("~ Parsed node: %+v\n", n)
		return nil
	}); x != nil {
		panic(x)
	}
}

func main() {
	_, err := readRel("node05")
	if err != nil {
		panic(err)
	}
}
