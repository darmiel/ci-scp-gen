package main

import (
	"fmt"
	"github.com/darmiel/ci-scp-gen/internal/scpgen"
)

func main() {
	n, err := scpgen.ReadRel("node05")
	if err != nil {
		panic(err)
	}
	fmt.Printf("node: %+v\n", n)
	fmt.Println("::", scpgen.Combine(n.SCPCommand("local.jar", "remote file.jar")))
}
