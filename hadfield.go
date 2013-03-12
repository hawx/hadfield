// Package hadfield implements a basic subcommand system to complement the
// standard library's 'flag' package. To use, in your 'main' function you must
// call the Run function with a list of Commands and Templates.
package hadfield

import (
	"os"
	"fmt"
	"flag"
	"log"
)

// Run executes the correct subcommand. It delegates the subcommand 'help' to
// the help function below. If it can not find a subcommand with the name given
// it exits with an error.
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
		if cmd.Name() == args[0] && cmd.Callable() {
			cmd.Call(cmd, templates, args)
			os.Exit(0)
		}
	}

	fmt.Fprintf(os.Stderr, "unknown subcommand %q\n", args[0])
	os.Exit(2)
}

// usage prints a usage message and then exits.
func usage(templates Templates, cmds Commands) {
	templates.Usage.Render(os.Stderr, cmds.Data())
	os.Exit(0)
}


// help controls the "help" pseudo-command. It will print the usage message if
// given an empty list of arguments. It prints the associated help text if given
// a signle argument. And otherwise exists with an error.
func help(templates Templates, cmds Commands, args []string) {
	if len(args) == 0 {
		templates.Usage.Render(os.Stdout, cmds.Data())
		return
	}
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "help given too many arguments\n")
		os.Exit(2)
	}

	arg := args[0]

	for _, cmd := range cmds {
		if cmd.Name() == arg {
			templates.Help.Render(os.Stdout, cmd.Data())
			return
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown help topic %#q\n", arg)
	os.Exit(2)
}
