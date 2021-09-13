package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	path2 "path"
	"path/filepath"
	"regexp"
	"strings"
)

var ErrRenamedFirst = errors.New("file renamed (fail first)")
var ErrRenamedAny = errors.New("file/s renamed (fail any)")

var (
	Included           *regexp.Regexp
	IncludeDirectories bool
	DryRun             bool
	TargetDirectory    string
	FailFirst          bool
	FailAny            bool
)

func init() {
	// tmp vars for processing after flag parsing
	var (
		include string
	)

	// flags
	flag.StringVar(&include, "i", ".*\\.jar", "RegEx Include Pattern")
	flag.StringVar(&TargetDirectory, "t", ".", "Target directory")
	flag.BoolVar(&FailFirst, "fy", false, "Exit != 0 on first rename")
	flag.BoolVar(&FailAny, "fl", true, "Exit != 0 if any file was renamed after checking through every file")
	flag.BoolVar(&IncludeDirectories, "d", false, "Include Directories")
	flag.BoolVar(&DryRun, "n", false, "Dry Run (default)")

	flag.Parse()

	// compile regex for excluded files
	Included = regexp.MustCompile(include)
}

func main() {
	// walk through files
	var renamed bool
	var err error
	if err = filepath.WalkDir(TargetDirectory, func(path string, d fs.DirEntry, err error) error {
		// check if >file< matches expression
		if (!IncludeDirectories && d.IsDir()) || !Included.MatchString(d.Name()) {
			return nil
		}
		dir, name := filepath.Dir(path), d.Name()
		fmt.Print("[*] Found file: ", name, " ")
		// check if contains hyphen
		if !strings.Contains(name, "-") {
			fmt.Println("... No Hyphens! :)")
			return nil
		}

		// new name
		newName := strings.Split(name, "-")[0] + ".jar"
		newPath := path2.Join(dir, newName)

		// rename
		fmt.Print("-> ", newName)

		if !DryRun {
			if err = os.Rename(path, newPath); err != nil {
				return err
			}
			fmt.Println("... Ok! :)")
		} else {
			fmt.Println("... Dry Run :)")
		}

		// fail?
		renamed = true
		// Fail if first first was activated
		if FailFirst {
			return ErrRenamedFirst
		}
		return nil
	}); err != nil {
		fmt.Println()
		panic(err)
	}
	if FailAny && renamed {
		panic(ErrRenamedAny)
	}
}
