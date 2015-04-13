# Hadfield [![docs](http://godoc.org/github.com/hawx/hadfield?status.svg)](http://godoc.org/github.com/hawx/hadfield)

A basic subcommands to complement the [flag package][flag].

## Example

First you must import the hadfield package (along with the others we need).

``` go
import (
  "github.com/hawx/hadfield"
  "fmt"
  "os"
)
```

Then you need to create instances of `hadfield.Command` for each subcommand. I
am prefixing their variable names with `cmd` to make life easier. The first word
in `Usage` must be the subcommand's name.

``` go
var cmdGreet = &hadfield.Command{
  Usage: "greet [options]",
  Short: "displays a greeting",
  Long: `
  Greet displays a formatted greeting to a person in the language specified.

    --person <name>        # Name of person to greet
    --lang <en|fr>         # Language to use, English or French
`,
}
```

Now we need to define a function that will be called when the subcommand is run.

``` go
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
```

We used some values in the run function that need to be taken from command line
flags. Let's define the flags needed in an `init` function.

``` go
var greetPerson, greetLang string

func init() {
  // Setup the command to use our previously defined function
  cmdGreet.Run = runGreet

  // Define the flags needed on command.Flag
  cmdGreet.Flag.StringVar(&greetPerson, "person", "", "")
  cmdGreet.Flag.StringVar(&greetLang,   "lang", "en", "")
}
```

We've nearly finished, we just need to create some templates for the help
documentation. Read the documentation for [text/template][] for more information
on what can be done here.

``` go
var templates = hadfield.Templates{
Usage: `usage: test [command] [arguments]

  This is a test.

  Commands: {{range .}}
    {{.Name | printf "%-15s"}} # {{.Short}}{{end}}
`,
Help: `usage: test {{.Usage}}
{{.Long}}
`,
}
```

Finally we call `hadfield.Run` in `main`. This executes the correct subcommand,
parsing flags as it goes.

``` go
var commands = hadfield.Commands{
  cmdGreet,
}

func main() {
  hadfield.Run(commands, templates)
}
```


[flag]: http://golang.org/pkg/flag/
[text/template]: http://golang.org/pkg/text/template/
