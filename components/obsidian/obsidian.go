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

func RunObsidian(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 || exactArg(args, "--help") {
		fmt.Fprintln(stdout, "Usage: adlaire-ci-obsidian sync <plan|apply|rollback> [options]")
		return 0
	}
	cfg, obsErr := parseSyncArgs(args)
	if obsErr != nil {
		writeObsidianSyncError(stderr, obsErr)
		return obsErr.Exit
	}
	switch cfg.Action {
	case "plan":
		return executeSyncPlan(cfg, stdout, stderr)
	case "apply":
		return executeSyncApply(cfg, stdout, stderr)
	case "rollback":
		return executeSyncRollback(cfg, stdout, stderr)
	default:
		writeObsidianSyncError(stderr, newObsError("OBSIDIAN_SYNC_INVALID_OPTION", 2, "", 0, 0))
		return 2
	}
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

func writeObsidianSyncError(stderr io.Writer, obsErr *obsidianError) {
	fmt.Fprintf(stderr, "obsidian: %s\n", obsErr.Code)
}
