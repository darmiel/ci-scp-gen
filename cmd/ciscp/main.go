package main

import (
	"flag"
	"fmt"
	"github.com/darmiel/ci-scp-gen/internal/scpgen"
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	IncludePattern *regexp.Regexp
)

func init() {
	var includePatternStr string
	flag.StringVar(&includePatternStr, "i", ".*\\.jar", "include pattern")
	flag.Parse()

	IncludePattern = regexp.MustCompile(includePatternStr)
}

func main() {
	var (
		err     error
		plugins = make(map[string]string) // plugin-name => plugin-path
	)

	// find all jar files
	fmt.Println("SCP :: Looking for files matching:", IncludePattern.String())
	if err = filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		var fname = d.Name()
		if d.IsDir() || !IncludePattern.MatchString(d.Name()) {
			return nil
		}
		var cname string
		if strings.Contains(fname, ".") {
			cname = fname[:strings.LastIndex(fname, ".")]
		} else {
			cname = fname
		}
		cname = strings.ToLower(cname) // always lower case plugin names
		// check if duplicate
		if _, ok := plugins[cname]; ok {
			return fmt.Errorf("ambiguous plugin '%s' found", cname)
		}
		plugins[cname] = path
		return nil
	}); err != nil {
		panic(err)
	}
	fmt.Println("SCP :: Found plugins:", plugins)

	n, err := scpgen.ReadRel("node05")
	if err != nil {
		panic(err)
	}
	fmt.Printf("node: %+v\n", n)
	fmt.Println("::", scpgen.Combine(n.SCPCommand("local.jar", "remote file.jar")))
}
