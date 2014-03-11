package main

import (
	"fmt"
	"github.com/felixge/quantastic/backend/server/config"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: quantastic-server <config-file>")
		os.Exit(1)
	}

	var conf config.Config
	if err := config.ReadFile(os.Args[1], &conf); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Could not load config: %s", err)
		os.Exit(1)
	}
}
