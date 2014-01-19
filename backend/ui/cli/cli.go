package cli

import (
	"fmt"
	"github.com/felixge/quantastic/backend/app"
	"github.com/felixge/quantastic/backend/ui/cli/commands"
	"io"
)

func NewCLI(stdout, stderr io.Writer, args []string) *CLI {
	c := &CLI{
		stdout:   stdout,
		stderr:   stderr,
		args:     args,
		quit:     make(chan error),
		commands: commands.Commands,
	}
	go c.loop()
	return c
}

type CLI struct {
	stdout, stderr io.Writer
	quit           chan error
	args           []string
	handlers       []interface{}
	commands       []commands.Command
}

// @TODO really make this a loop / support interactive mode
func (c *CLI) loop() {
	r, err := c.parseRequest(c.args)
	if err == nil {
		if r == nil {
			fmt.Fprintf(c.stdout, "@TODO: Show usage\n")
		} else {
			err = c.dispatch(r)
		}
	}
	if err != nil {
		c.printError(err)
	}
	c.quit <- err
}

func (c *CLI) parseRequest(args []string) (interface{}, error) {
	if len(args) == 0 {
		return nil, nil
	}

	switch cmd := args[0]; cmd {
	case "-v", "--version", "version":
		return getVersionRequest{}, nil
	case "time":
		return getVersionRequest{}, nil
	default:
		return nil, fmt.Errorf("Unknown command: %s", cmd)
	}
}

func (c *CLI) dispatch(request interface{}) error {
	for _, handler := range c.handlers {
		switch r := request.(type) {
		case getVersionRequest:
			if h, ok := handler.(app.GetVersionHandler); ok {
				response := h.GetVersion(r)
				fmt.Fprintf(c.stdout, "%s\n", response.Version())
				return nil
			}
		}
	}
	return fmt.Errorf("Could not find handler for: %#v", request)
}

func (c *CLI) printError(err error) {
	fmt.Fprintf(c.stderr, "Error: %s\n", err)
}

// @TODO Panic if an existing handler implements the same interface already.
func (c *CLI) AddHandler(handler interface{}) {
	c.handlers = append(c.handlers, handler)
}

func (c *CLI) Wait() error {
	return <-c.quit
}

type getVersionRequest struct{}
