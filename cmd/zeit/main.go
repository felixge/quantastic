package main

import (
	"fmt"
	"github.com/felixge/zeit"
	"os"
)

func main() {
	z, err := zeit.NewZeit("config.xml")
	if err != nil {
		fmt.Printf("init error: %s", err)
		os.Exit(1)
	}
	if err := z.Run(); err != nil {
		fmt.Printf("runtime error: %s", err)
		os.Exit(1)
	}
}
