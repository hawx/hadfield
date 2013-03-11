package hadfield

import (
	"flag"
	"strings"
	"os"
)

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
	tmpl(os.Stdout, templates.Help, c)
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
