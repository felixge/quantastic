// Command quantastic implements the quantastic backend.
package main

import (
	"github.com/felixge/quantastic/backend/app"
	"github.com/felixge/quantastic/backend/ui/text"
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
	textUI := text.NewTextUI(os.Stdout, os.Stderr, os.Args[1:])
	textUI.AddHandler(app.NewGetVersionHandler(buildVersion))
	if err := textUI.Wait(); err != nil {
		os.Exit(1)
	}
}
