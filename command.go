package hadfield

import (
	"flag"
	"os"
	"strings"
)

type Commands []*Command

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

	// Flag defines a set of flags to be parsed when the command is activated.
	Flag flag.FlagSet

	// CustomFlags, if set, will prevent the parsing of any flags when this
	// command is activated.
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

// Call parses the flags if CustomFlags is not set, then calls the function
// defined by Run, and finally exits.
func (c *Command) Call(cmd *Command, templates Templates, args []string) {
	c.Flag.Usage = func() { templates.Command.Render(os.Stdout, cmd) }

	if c.CustomFlags {
		args = args[1:]
	} else {
		c.Flag.Parse(args[1:])
		args = c.Flag.Args()
	}

	c.Run(c, args)
}
