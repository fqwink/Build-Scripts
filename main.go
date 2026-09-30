package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/fqwink/build-scripts/components"
)

func main() {
	os.Exit(dispatchMain(filepath.Base(os.Args[0]), os.Args[1:], os.Stdout, os.Stderr))
}

func dispatchMain(name string, args []string, stdout io.Writer, stderr io.Writer) int {
	switch name {
	case "adlaire-ci-build":
		return components.RunBuild(args, stdout, stderr)
	case "adlaire-ci-runner":
		return components.RunRunner(args, stdout, stderr)
	case "adlaire-ci-api":
		return components.RunAPI(args, os.Stdin, stdout, stderr)
	case "adlaire-ci-admin":
		return components.RunAdmin(args, stdout, stderr)
	default:
		fmt.Fprintf(stderr, "unknown command: %s\n", name)
		return 2
	}
}
