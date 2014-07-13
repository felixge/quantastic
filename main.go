package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/felixge/quantastic/db"
)

var commands = []*command{
	cmdTimeLog,
	cmdTimeStart,
	cmdTimeEnd,
	cmdTimeEdit,
	cmdTimeRm,
	cmdTimeInvoice,
}

func usage() string {
	b := &bytes.Buffer{}
	fmt.Fprintf(b, "usage: quantastic <command>\n\n")
	ml := 0 // max command name length
	for _, cmd := range commands {
		if l := len(cmd.name); l > ml {
			ml = l
		}
	}
	for _, cmd := range commands {
		name := cmd.name + strings.Repeat(" ", ml-len(cmd.name)+1)
		fmt.Fprintf(b, "  %s %s\n", name, cmd.description)
	}
	return b.String()
}

func main() {
	args := os.Args
	if len(args) < 2 {
		fatal("%s", usage())
	}
	args = args[1:]

outer:
	for _, cmd := range commands {
		cmdParts := strings.Split(cmd.name, " ")
		if len(cmdParts) > len(args) {
			continue
		}
		for i, cmdPart := range cmdParts {
			if cmdPart != args[i] {
				continue outer
			}
		}
		args = args[len(cmdParts):]
		db, err := db.OpenDb(filepath.Join(os.Getenv("HOME"), ".quantastic"))
		if err != nil {
			fatal("Could not open db: %s", err)
		}
		cmd.fn(&Context{Db: db, Cmd: cmd, Args: args})
		return
	}
	cmdStr := strings.Join(args, " ")
	fatal("Unknown command: %s\n\n%s", cmdStr, usage())
}

type command struct {
	name        string
	description string
	usage       string
	fn          func(c *Context)
}

func (c *command) Usage() string {
	return fmt.Sprintf("usage: quantastic %s %s", c.name, c.usage)
}
