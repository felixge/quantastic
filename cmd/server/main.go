// Command quantastic implements the quantastic server.

package main

// @TODO -h/--help flag to print usage
// @TODO -v/--version flag to print version

import (
	"flag"
	"fmt"
	pkgconfig "github.com/felixge/quantastic/config"
	"os"
)

func main() {
	var (
		configPath = flag.Arg(0)
		config     pkgconfig.Server
	)
	if configPath == "" {
		configPath = "config.yml"
	}
	if err := pkgconfig.Load(configPath, &config); err != nil {
		fmt.Printf("Could not load config. path=%s err=%s", configPath, err)
		os.Exit(1)
	}
	s, err := pkgconfig.NewServer(config)
	if err != nil {
		fmt.Printf("Could not create server. err=%s", err)
		os.Exit(1)
	}
	if err := s.Run(); err != nil {
		fmt.Printf("System failure. err=%s", err)
		os.Exit(1)
	}
}
