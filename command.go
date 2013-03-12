package hadfield

import (
	"flag"
	"strings"
	"os"
)

type Interface interface {
	Name() string
	Data() interface{}
	Runnable() bool
	Call(cmd Interface, templates Templates, args []string)
}

type Command struct {
	// Run runs the command. It is passed the list of args that came after the
	// command name.
	Run          func(cmd *Command, args []string)

	// Usage returns the one-line usage string. The first word on the line is
	// taken to be the command name.
	Usage        string

	// Short is a single line description used in the help listing.
	Short        string

	// Long is the detailed and formatted message shown in the full help for the
	// command.
	Long         string

	Flag         flag.FlagSet
	CustomFlags  bool
}

type Commands []Interface

func (c *Command) Name() string {
	name := c.Usage
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

func printUsage(c Interface, templates Templates) {
	templates.Help.Render(os.Stdout, c.Data())
	os.Exit(0)
}

func (c *Command) Runnable() bool {
	return c.Run != nil
}

func (c *Command) Call(cmd Interface, templates Templates, args []string) {
	c.Flag.Usage = func() { printUsage(cmd, templates) }

	if c.CustomFlags {
		args = args[1:]
	} else {
		c.Flag.Parse(args[1:])
		args = c.Flag.Args()
	}

	c.Run(c, args)
	os.Exit(0)
}

func (c *Command) Data() interface{} {
	return c
}
