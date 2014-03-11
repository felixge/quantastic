package commands

//func NewFlag(longName, shortName string) {
//}

//NewFlag("version", "v")

type Flag struct{}

// What is a command?
// * a name, along with *aliases
// * a high level description of what it does
// * a list of flags, along of descriptions what they do
// * can have sub-commands
// I'd like a declaritive approach

type Command interface {
	Names() []string
	ParseRequest(args []string) (interface{}, error)
}

var Commands = []Command{
	GetVersionCommand{},
}

type GetVersionCommand struct{}

func (c GetVersionCommand) Names() []string {
	return []string{"version", "-v", "--version"}
}

func (c GetVersionCommand) ParseRequest(args []string) (interface{}, error) {
	return nil, nil
}
