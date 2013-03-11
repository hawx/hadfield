package hadfield

import (
	"io"
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
