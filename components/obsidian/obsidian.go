package obsidian

import (
	"fmt"
	"io"

	"github.com/fqwink/build-scripts/components/builder"
)

func RunBuild(args []string, stdout, stderr io.Writer) int {
	cfg, handled, obsErr := parseBuildArgs(args)
	if obsErr != nil {
		writeObsidianError(stderr, obsErr)
		return obsErr.Exit
	}
	if !handled {
		return builder.RunBuild(args, stdout, stderr)
	}
	return executeVaultBuild(cfg, stdout, stderr)
}

func writeObsidianError(stderr io.Writer, obsErr *obsidianError) {
	fmt.Fprintln(stderr, obsErr.Code)
	if obsErr.Target != "" {
		if obsErr.Line > 0 && obsErr.Column > 0 {
			fmt.Fprintf(stderr, "%s:%d:%d\n", obsErr.Target, obsErr.Line, obsErr.Column)
			return
		}
		fmt.Fprintln(stderr, obsErr.Target)
	}
}
