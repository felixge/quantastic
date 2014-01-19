// Command quantastic implements the quantastic backend.
package main

import (
	"github.com/felixge/quantastic/backend/ui/cli"
	"github.com/felixge/quantastic/backend/handlers"
	"github.com/felixge/quantastic/backend/version"
	"os"
)

// populated by the build system
var (
	buildRelease string
	buildCommit  string
)

func main() {
	buildVersion := version.NewVersion(buildRelease, buildCommit)
	cliUI := cli.NewCLI(os.Stdout, os.Stderr, os.Args[1:])
	cliUI.AddHandler(handlers.NewGetVersion(buildVersion))
	if err := cliUI.Wait(); err != nil {
		os.Exit(1)
	}
}
