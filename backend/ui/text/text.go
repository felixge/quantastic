package text

import (
	"fmt"
	"github.com/felixge/quantastic/backend/app"
	"io"
)

func NewTextUI(stdout, stderr io.Writer, args []string) *TextUI {
	u := &TextUI{
		stdout: stdout,
		stderr: stderr,
		args:   args,
		quit:   make(chan error),
	}
	go u.loop()
	return u
}

type TextUI struct {
	stdout, stderr io.Writer
	quit           chan error
	args           []string
	handlers       []interface{}
}

// @TODO really make this a loop / support interactive mode
func (u *TextUI) loop() {
	err := u.dispatch(getVersionRequest{})
	if err != nil {
		u.printError(err)
	}
	u.quit <- err
}

func (u *TextUI) dispatch(request interface{}) error {
	for _, handler := range u.handlers {
		switch r := request.(type) {
		case getVersionRequest:
			if h, ok := handler.(app.GetVersionHandler); ok {
				response := h.GetVersion(r)
				fmt.Fprintf(u.stdout, "%s\n", response.Version())
				return nil
			}
		}
	}
	return fmt.Errorf("Could not find handler for: %#v", request)
}

func (u *TextUI) printError(err error) {
	fmt.Fprintf(u.stderr, "Error: %s\n", err)
}

// @TODO Panic if an existing handler implements the same interface already.
func (u *TextUI) AddHandler(handler interface{}) {
	u.handlers = append(u.handlers, handler)
}

func (u *TextUI) Wait() error {
	return <-u.quit
}

type getVersionRequest struct{}
