// Package hadfield implements a basic subcommand system to complement the
// flag package.
package hadfield

import (
	"flag"
	"fmt"
	"os"
)

const (
	unknownSubcommand = "unknown subcommand %q\n"
	unknownHelpTopic  = "unknown help topic %q\n"
	helpTooManyArgs   = "help given too many arguments\n"
)

// Customisable Exit function. This is used for exiting in various places
// throughout and can be overriden for testing purposes or to perform other
// tasks.
var Exit = os.Exit

// Run executes the correct subcommand.
//
// The special subcommand 'help' is defined and displays either the usage
// message, or if called with an argument the help message for a particular
// subcommand.
//
// If the subcommand cannot be found a message is displayed and it exits with
// status code 2.
func Run(cmds Commands, tpcs Topics, templates Templates) {
	flag.Usage = func() { Usage(cmds, tpcs, templates) }
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		data := struct {
			Commands Commands
			Topics   Topics
		}{
			Commands: cmds,
			Topics:   tpcs,
		}

		templates.Help.Render(os.Stderr, data)
		Exit(1)
	}

	if args[0] == "help" {
		help(templates, cmds, tpcs, args[1:])
		return
	}

	for _, cmd := range cmds {
		if cmd.Name() == args[0] {
			cmd.Call(cmd, templates, args)
			return
		}
	}

	fmt.Fprintf(os.Stderr, unknownSubcommand, args[0])
	Exit(1)
}

// Usage writes a help message to Stdout.
func Usage(cmds Commands, tpcs Topics, templates Templates) {
	data := struct {
		Commands Commands
		Topics   Topics
	}{
		Commands: cmds,
		Topics:   tpcs,
	}

	templates.Help.Render(os.Stdout, data)
}

// help controls the "help" pseudo-command. It will print the usage message if
// given an empty list of arguments. It prints the associated help text if given
// a signle argument. And otherwise exists with an error.
func help(templates Templates, cmds Commands, tpcs Topics, args []string) {
	if len(args) == 0 {
		Usage(cmds, tpcs, templates)
		return
	}
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, helpTooManyArgs)
		Exit(1)
	}

	arg := args[0]

	for _, cmd := range cmds {
		if cmd.Name() == arg {
			templates.Command.Render(os.Stdout, cmd)
			return
		}
	}
	for _, tpc := range tpcs {
		if tpc.Name == arg {
			templates.Topic.Render(os.Stdout, tpc)
			return
		}
	}

	fmt.Fprintf(os.Stderr, unknownHelpTopic, arg)
	Exit(1)
}
