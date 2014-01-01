// Command quantastic provides the quantastic command line application.
package main

// @TODO -h/--help flag to print usage
// @TODO -v/--version flag to print version

import (
	"flag"
	"fmt"
	pkgconfig "github.com/felixge/quantastic/config"
	pkgserver "github.com/felixge/quantastic/server"
	"os"
)

func main() {
	var (
		configPath = flag.Arg(0)
		config     pkgconfig.Config
	)
	if configPath == "" {
		configPath = "config.yml"
	}
	if err := pkgconfig.Load(configPath, &config); err != nil {
		fmt.Printf("Could not load config. path=%s err=%s", configPath, err)
		os.Exit(1)
	}
	s, err := pkgserver.NewServer(config)
	if err != nil {
		fmt.Printf("Could not create server. err=%s", err)
		os.Exit(1)
	}
	if err := s.Run(); err != nil {
		fmt.Printf("System failure. err=%s", err)
		os.Exit(1)
	}
}
