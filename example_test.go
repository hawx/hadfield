package hadfield_test

import (
	"fmt"
	"os"

	"hawx.me/code/hadfield"
)

var cmdGreet = &hadfield.Command{
	Usage: "greet [options]",
	Short: "displays a greeting",
	Long: `
  Greet displays a formatted greeting to a person in the language specified.

    --person <name>     # Name of person to greet
    --lang <en|fr>      # Language to use, English or French
`,
}

func runGreet(cmd *hadfield.Command, args []string) {
	switch greetLang {
	case "en":
		fmt.Println("Hello", greetPerson)
	case "fr":
		fmt.Println("Bonjour", greetPerson)
	default:
		os.Exit(2)
	}
}

var greetPerson, greetLang string

func init() {
	cmdGreet.Run = runGreet

	cmdGreet.Flag.StringVar(&greetPerson, "person", "someone?", "")
	cmdGreet.Flag.StringVar(&greetLang, "lang", "en", "")
}

var templates = hadfield.Templates{
	Usage: `usage: example [command] [arguments]

  This is an example.

  Commands: {{range .}}
    {{.Name | printf "%-15s"}} # {{.Short}}{{end}}
`,
	Help: `usage: example {{.Usage}}
{{.Long}}
`,
}

var commands = hadfield.Commands{
	cmdGreet,
}

func Example() {
	hadfield.Run(commands, templates)
}
