package hadfield

import (
	"flag"
	"strings"
	"os"
)

type CommandLike interface {
	// Runnable returns whether the command can be run or not.
	Runnable() bool

	// Run runs the command. It is passed the list of args that came after the
	// command name.
	Run(cmd *CommandLike, args []string)

	// Usage returns the one-line usage method. The first word on the line is
	// taken to be the command name.
	Usage()  string

	// Short is a single line description used in the listing.
	Short()  string

	// Long is the long, formatted message shown in the full help for the command.
	Long()   string
}

type Command struct {
	Run          func(cmd *Command, args []string)
	Usage        string
	Short        string
	Long         string
	Flag         flag.FlagSet
	CustomFlags  bool
}

type Commands []*Command

func (c *Command) Name() string {
	name := c.Usage
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

func (c *Command) PrintUsage(templates Templates) {
	templates.Help.Render(os.Stdout, c)
	os.Exit(0)
}

func (c *Command) Runnable() bool {
	return c.Run != nil
}

func (c *Command) Call(cmd *Command, templates Templates, args []string) {
	c.Flag.Usage = func() { cmd.PrintUsage(templates) }

	if c.CustomFlags {
		args = args[1:]
	} else {
		c.Flag.Parse(args[1:])
		args = c.Flag.Args()
	}

	c.Run(c, args)
	os.Exit(0)
}
