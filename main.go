package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fqwink/build-scripts/components/admin"
	"github.com/fqwink/build-scripts/components/api"
	"github.com/fqwink/build-scripts/components/builder"
	"github.com/fqwink/build-scripts/components/mcp"
	"github.com/fqwink/build-scripts/components/release"
	"github.com/fqwink/build-scripts/components/runner"
	"github.com/fqwink/build-scripts/components/setup"
)

var binaryVersion = "V.0.0-dev"

func main() {
	os.Exit(dispatchMain(filepath.Base(os.Args[0]), os.Args[1:], os.Stdout, os.Stderr))
}

func dispatchMain(name string, args []string, stdout io.Writer, stderr io.Writer) int {
	name = canonicalBinaryName(name)
	builder.SetBinaryVersion(binaryVersion)
	runner.SetBinaryVersion(binaryVersion)
	api.SetBinaryVersion(binaryVersion)
	admin.SetBinaryVersion(binaryVersion)
	setup.SetBinaryVersion(binaryVersion)
	release.SetBinaryVersion(binaryVersion)
	mcp.SetBinaryVersion(binaryVersion)
	if hasMainExactArg(args, "--version") && isStandardBinaryName(name) {
		fmt.Fprintf(stdout, "%s %s go=%s\n", name, binaryVersion, runtime.Version())
		return 0
	}
	switch name {
	case "adlaire-ci-build":
		return builder.RunBuild(args, stdout, stderr)
	case "adlaire-ci-runner":
		return runner.RunRunner(args, stdout, stderr)
	case "adlaire-ci-api":
		return api.RunAPI(args, os.Stdin, stdout, stderr)
	case "adlaire-ci-admin":
		return admin.RunAdmin(args, stdout, stderr)
	case "adlaire-ci-setup":
		return setup.RunSetup(args, stdout, stderr)
	case "adlaire-ci-release":
		return release.RunRelease(args, stdout, stderr)
	case "adlaire-ci-mcp":
		return mcp.RunMCP(args, stdout, stderr)
	default:
		fmt.Fprintf(stderr, "unknown command: %s\n", name)
		return 2
	}
}

func canonicalBinaryName(name string) string {
	if strings.HasSuffix(name, "-linux-amd64") {
		trimmed := strings.TrimSuffix(name, "-linux-amd64")
		if isStandardBinaryName(trimmed) {
			return trimmed
		}
	}
	return name
}

func isStandardBinaryName(name string) bool {
	switch name {
	case "adlaire-ci-build", "adlaire-ci-runner", "adlaire-ci-api", "adlaire-ci-admin", "adlaire-ci-setup", "adlaire-ci-release", "adlaire-ci-mcp":
		return true
	default:
		return false
	}
}

func hasMainExactArg(args []string, want string) bool {
	for _, arg := range args {
		if arg == want {
			return true
		}
	}
	return false
}
