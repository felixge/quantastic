package cli

import (
	"fmt"
	"github.com/felixge/quantastic/backend/app"
	"io"
)

func NewCLI(stdout, stderr io.Writer, args []string) *CLI {
	c := &CLI{
		stdout: stdout,
		stderr: stderr,
		args:   args,
		quit:   make(chan error),
	}
	go c.loop()
	return c
}

type CLI struct {
	stdout, stderr io.Writer
	quit           chan error
	args           []string
	app       []interface{}
}

// @TODO really make this a loop / support interactive mode
func (c *CLI) loop() {
	err := c.dispatch(getVersionRequest{})
	if err != nil {
		c.printError(err)
	}
	c.quit <- err
}

func (c *CLI) dispatch(request interface{}) error {
	for _, handler := range c.app {
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
	c.app = append(c.app, handler)
}

func (c *CLI) Wait() error {
	return <-c.quit
}

type getVersionRequest struct{}
