package hadfield

import (
	"flag"
	"strings"
	"os"
	"io"
	"text/template"
	"unicode/utf8"
	"unicode"
	"fmt"
	"log"
)

type Templates struct {
	Usage  string
	Help   string
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
	tmpl(os.Stdout, templates.Help, c)
	os.Exit(0)
}

func (c *Command) Runnable() bool {
	return c.Run != nil
}


func Run(cmds Commands, templates Templates) {
	flag.Usage = func() { usage(templates, cmds) }
	flag.Parse()
	log.SetFlags(0)

	args := flag.Args()

	if len(args) < 1 {
		usage(templates, cmds)
	}

	if args[0] == "help" {
		help(templates, cmds, args[1:])
		return
	}

	for _, cmd := range cmds {
		if cmd.Name() == args[0] && cmd.Run != nil {
			cmd.Flag.Usage = func() { cmd.PrintUsage(templates) }
			if cmd.CustomFlags {
				args = args[1:]
			} else {
				cmd.Flag.Parse(args[1:])
				args = cmd.Flag.Args()
			}
			cmd.Run(cmd, args)
			os.Exit(0)
		}
	}

	fmt.Fprintf(os.Stderr, "unknown subcommand %q\n", args[0])
	os.Exit(2)
}

func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	t.Funcs(template.FuncMap{"trim": strings.TrimSpace, "capitalize": capitalize})
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	r,n := utf8.DecodeRuneInString(s)
	return string(unicode.ToTitle(r)) + s[n:]
}

func printUsage(w io.Writer, tmpls Templates, cmds Commands) {
	tmpl(w, tmpls.Usage, cmds)
}

func usage(tmpls Templates, cmds Commands) {
	printUsage(os.Stderr, tmpls, cmds)
	os.Exit(0)
}

func help(tmpls Templates, cmds Commands, args []string) {
	if len(args) == 0 {
		printUsage(os.Stdout, tmpls, cmds)
		return
	}
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "help given too many arguments\n")
	}

	arg := args[0]

	for _, cmd := range cmds {
		if cmd.Name() == arg {
			tmpl(os.Stdout, tmpls.Help, cmd)
			return
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown help topic %#q\n", arg)
	os.Exit(2)
}
