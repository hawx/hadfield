package hadfield

import (
	"flag"
	"strings"
	"os"
)

type Interface interface {
	Name()      string
	Data()      interface{}
	Category()  string
	Callable()  bool
	Call(cmd Interface, templates Templates, args []string)
}

func PrintUsage(c Interface, templates Templates) {
	templates.Help.Render(os.Stdout, c.Data())
	os.Exit(0)
}

type Commands []Interface

func (cs Commands) PrintUsage(templates Templates) {
	templates.Usage.Render(os.Stderr, cs)
	os.Exit(0)
}

func (cs Commands) Data() []interface{} {
	is := make([]interface{}, len(cs))

	for i,c := range cs {
		is[i] = c.Data()
	}

	return is
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

func (c *Command) Name() string {
	name := c.Usage
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

func (c *Command) Category() string {
	if c.Run != nil {
		return "Runnable"
	}
	return ""
}

func (c *Command) Callable() bool {
	return c.Category() == "Runnable"
}

func (c *Command) Call(cmd Interface, templates Templates, args []string) {
	c.Flag.Usage = func() { PrintUsage(cmd, templates) }

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
	return map[string]interface{}{
		"Callable": c.Callable(),
		"Category": c.Category(),
		"Usage":    c.Usage,
		"Short":    c.Short,
		"Long":     c.Long,
		"Name":     c.Name(),
	}
}
