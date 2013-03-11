package hadfield

import (
	"os"
	"fmt"
	"flag"
	"log"
)

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
		if cmd.Name() == args[0] && cmd.Runnable() == true {
			cmd.Call(cmd, templates, args)
		}
	}

	fmt.Fprintf(os.Stderr, "unknown subcommand %q\n", args[0])
	os.Exit(2)
}

func usage(templates Templates, cmds Commands) {
	templates.Usage.Render(os.Stderr, cmds)
	os.Exit(0)
}

func help(templates Templates, cmds Commands, args []string) {
	if len(args) == 0 {
		templates.Usage.Render(os.Stdout, cmds)
		return
	}
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "help given too many arguments\n")
	}

	arg := args[0]

	for _, cmd := range cmds {
		if cmd.Name() == arg {
			templates.Help.Render(os.Stdout, cmd)
			return
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown help topic %#q\n", arg)
	os.Exit(2)
}
