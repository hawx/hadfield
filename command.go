package hadfield

import (
	"flag"
	"os"
	"strings"
)

// Interface defines the common behaviour of subcommands. The default
// implementation is Command, but it is possible to mix these with other types
// of subcommand such as those dynamically discovered when used.
type Interface interface {
	// Name is the word used to call the subcommand.
	Name() string

	// Data is used when rendering help and usage templates. Values that should be
	// expected are:
	//
	// Callable:
	//   i.e. Callable()
	// Category:
	//   i.e. Category().
	// Usage:
	//   The string starting with the commands name.
	// Short:
	//   A short, one-line, description.
	// Long:
	//   A long description.
	// Name:
	//   i.e. Name().
	Data() interface{}

	// Category is the type of subcommand, it can be anything and can be used to
	// group subcommands together.
	Category() string

	// Callable is true if the Call() method actually does something.
	Callable() bool

	// Call runs the subcommand.
	Call(cmd Interface, templates Templates, args []string)
}

// CommandUsage displays a help message for the subcommand to Stdout, then exits.
func CommandUsage(c Interface, templates Templates) {
	templates.Help.Render(os.Stdout, c.Data())
	Exit(0)
}

type Commands []Interface

func (cs Commands) Data() []interface{} {
	is := make([]interface{}, len(cs))

	for i, c := range cs {
		is[i] = c.Data()
	}

	return is
}

type Command struct {
	// Run runs the command. It is passed the list of args that came after the
	// command name.
	Run func(cmd *Command, args []string)

	// Usage returns the one-line usage string. The first word on the line is
	// taken to be the command name.
	Usage string

	// Short is a single line description used in the help listing.
	Short string

	// Long is the detailed and formatted message shown in the full help for the
	// command.
	Long string

	Flag        flag.FlagSet
	CustomFlags bool
}

// Name returns the first word in Usage.
func (c *Command) Name() string {
	name := c.Usage
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

// Category returns "Command" if Run is defined, and otherwise "Documentation".
func (c *Command) Category() string {
	if c.Run != nil {
		return "Command"
	}
	return "Documentation"
}

// Callable returns true if Run is defined, and false otherwise.
func (c *Command) Callable() bool {
	return c.Run != nil
}

// Call parses the flags if CustomFlags is not set, then calls the function
// defined by Run, and finally exits.
func (c *Command) Call(cmd Interface, templates Templates, args []string) {
	c.Flag.Usage = func() { CommandUsage(cmd, templates) }

	if c.CustomFlags {
		args = args[1:]
	} else {
		c.Flag.Parse(args[1:])
		args = c.Flag.Args()
	}

	c.Run(c, args)
	Exit(0)
}

type cmdData struct {
	Callable bool
	Category string
	Usage    string
	Short    string
	Long     string
	Name     string
}

// Data returns the information required for rendering the help templates.
func (c *Command) Data() interface{} {
	return cmdData{
		Callable: c.Callable(),
		Category: c.Category(),
		Usage:    c.Usage,
		Short:    c.Short,
		Long:     c.Long,
		Name:     c.Name(),
	}
}
