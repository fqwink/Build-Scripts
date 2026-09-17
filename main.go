package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/fqwink/build-scripts/components"
)

func main() {
	if strings.Contains(filepath.Base(os.Args[0]), "adlaire-ci-runner") {
		os.Exit(components.RunRunner(os.Args[1:], os.Stdout, os.Stderr))
	}
	os.Exit(components.RunBuild(os.Args[1:], os.Stdout, os.Stderr))
}
