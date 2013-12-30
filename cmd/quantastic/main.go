package main

import (
	"fmt"
	"github.com/felixge/quantastic"
	"os"
)

func main() {
	s, err := quantastic.NewServer("config.yml")
	if err != nil {
		fmt.Printf("init error: %s", err)
		os.Exit(1)
	}
	if err := s.Run(); err != nil {
		fmt.Printf("runtime error: %s", err)
		os.Exit(1)
	}
}
