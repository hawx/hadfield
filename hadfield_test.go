package hadfield_test

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	hadfield "hawx.me/code/hadfield"
)

var nilRun = func(c *hadfield.Command, args []string) {}

func captureStderr(f func()) string {
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f()

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = old

	return string(out)
}

func captureStdout(f func()) string {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = oldStdout

	return string(out)
}

func TestHadfield(t *testing.T) {
	receivedArgs := make(chan []string, 1)

	cmds := hadfield.Commands{
		&hadfield.Command{Usage: "he", Run: nilRun},
		&hadfield.Command{
			Usage: "hey",
			Run: func(c *hadfield.Command, args []string) {
				receivedArgs <- args
			},
		},
		&hadfield.Command{Usage: "bye", Run: nilRun},
	}

	os.Args = []string{"me", "hey", "you"}
	hadfield.Exit = func(_ int) {}
	hadfield.Run(cmds, hadfield.Templates{})

	select {
	case args := <-receivedArgs:
		assert.Equal(t, []string{"you"}, args)
	case <-time.After(time.Second):
		t.Fatal("timeout")
	}
}

func TestHadfieldUnknownSubcommand(t *testing.T) {
	exitCode := -1

	assert.Equal(t, "unknown subcommand \"hey\"\n", captureStderr(func() {
		os.Args = []string{"me", "hey", "you"}
		hadfield.Exit = func(i int) { exitCode = i }
		hadfield.Run(hadfield.Commands{}, hadfield.Templates{})
	}))

	assert.Equal(t, 2, exitCode)
}

func TestHadfieldWithFlags(t *testing.T) {
	flagString := ""
	receivedArgs := make(chan []string, 1)
	receivedFlag := make(chan string, 1)

	heyCmd := &hadfield.Command{
		Usage: "hey",
		Run: func(c *hadfield.Command, args []string) {
			receivedArgs <- args
			receivedFlag <- flagString
		},
	}

	heyCmd.Flag.StringVar(&flagString, "flag", "", "")

	cmds := hadfield.Commands{
		&hadfield.Command{Usage: "he", Run: nilRun},
		heyCmd,
		&hadfield.Command{Usage: "bye", Run: nilRun},
	}

	os.Args = []string{"me", "hey", "--flag", "value"}
	hadfield.Exit = func(_ int) {}
	hadfield.Run(cmds, hadfield.Templates{})

	select {
	case args := <-receivedArgs:
		assert.Equal(t, []string{}, args)
	case <-time.After(time.Second):
		t.Fatal("timeout")
	}

	select {
	case args := <-receivedFlag:
		assert.Equal(t, "value", args)
	case <-time.After(time.Second):
		t.Fatal("timeout")
	}
}

func TestHadfieldHelp(t *testing.T) {
	cmds := hadfield.Commands{
		&hadfield.Command{
			Usage: "he",
			Short: "he does stuff",
			Run:   nilRun,
		},
		&hadfield.Command{
			Usage: "hey",
			Short: "hey does other stuff",
			Run:   nilRun,
		},
		&hadfield.Command{
			Usage: "bye",
			Short: "bye goes away",
			Run:   nilRun,
		},
	}

	var templates = hadfield.Templates{
		Help: `usage: test [command] [arguments]

  This is a test.

  Commands:{{range .}}
    {{.Name | printf "%-15s"}} # {{.Short | capitalize}}{{end}}
`,
	}

	expectedOut := `usage: test [command] [arguments]

  This is a test.

  Commands:
    he              # He does stuff
    hey             # Hey does other stuff
    bye             # Bye goes away
`

	exitCode := -1

	assert.Equal(t, expectedOut, captureStdout(func() {
		os.Args = []string{"me", "help"}
		hadfield.Exit = func(i int) { exitCode = i }
		hadfield.Run(cmds, templates)
	}))
	assert.Equal(t, 0, exitCode)

	assert.Equal(t, expectedOut, captureStderr(func() {
		os.Args = []string{"me"}
		hadfield.Exit = func(i int) { exitCode = i; panic("") }

		defer func() { recover() }()
		hadfield.Run(cmds, templates)
	}))
	assert.Equal(t, 1, exitCode)
}

func TestHadfieldHelpCommand(t *testing.T) {
	cmd := &hadfield.Command{
		Usage: "hey",
		Short: "hey does other stuff",
		Long: `
  Hey is a command to say hey, HEY!

  Options:
    --later WHEN   # later is cool for now, BUT LATER (default: cool)
    --now REALLY   # now does stuff right NOW
`,
		Run: nilRun,
	}

	cmd.Flag.String("now", "", "")
	cmd.Flag.String("later", "cool", "")

	cmds := hadfield.Commands{cmd}

	var templates = hadfield.Templates{
		Command: `usage: test {{.Usage}}

  {{.Long | trim}}
`,
	}

	expectedOut := `usage: test hey

  Hey is a command to say hey, HEY!

  Options:
    --later WHEN   # later is cool for now, BUT LATER (default: cool)
    --now REALLY   # now does stuff right NOW
`

	exitCode := -1

	assert.Equal(t, expectedOut, captureStdout(func() {
		os.Args = []string{"me", "help", "hey"}
		hadfield.Exit = func(i int) { exitCode = i }
		hadfield.Run(cmds, templates)
	}))
	assert.Equal(t, 0, exitCode)

	assert.Equal(t, "unknown help topic \"what\"\n", captureStderr(func() {
		os.Args = []string{"me", "help", "what"}
		hadfield.Exit = func(i int) { exitCode = i }
		hadfield.Run(cmds, templates)
	}))
	assert.Equal(t, 1, exitCode)

	exitCode = -1
	assert.Equal(t, "help given too many arguments\n", captureStderr(func() {
		os.Args = []string{"me", "help", "what", "and"}
		hadfield.Exit = func(i int) { exitCode = i; panic("") }

		defer func() { recover() }()
		hadfield.Run(cmds, templates)
	}))
	assert.Equal(t, 1, exitCode)
}

func TestHadfieldHelpNonCallable(t *testing.T) {
	cmds := hadfield.Commands{
		&hadfield.Command{
			Usage: "hey",
			Short: "hey does other stuff",
			Long: `
This is actually just documentation about the "hey" system.
`,
		},
		&hadfield.Command{
			Usage: "hey",
			Short: "runs hey",
			Run:   nilRun,
		},
	}

	var templates = hadfield.Templates{
		Help: `usage: test [command] [arguments]

  This is a test.

  Commands:{{range .}}{{if eq .Category "Command"}}
    {{.Name | printf "%-15s"}} # {{.Short}}{{end}}{{end}}

  Additional help:{{range .}}{{if not .Callable}}
    {{.Name | printf "%-15s"}} # {{.Short}}{{end}}{{end}}
`,
		Command: `{{if .Callable}}usage: test {{.Usage}}
{{end}}{{.Long | trim}}
`,
	}

	expectedUsage := `usage: test [command] [arguments]

  This is a test.

  Commands:
    hey             # runs hey

  Additional help:
    hey             # hey does other stuff
`

	expectedTopic := `This is actually just documentation about the "hey" system.
`

	assert.Equal(t, expectedUsage, captureStdout(func() {
		os.Args = []string{"me", "help"}
		hadfield.Exit = func(_ int) {}
		hadfield.Run(cmds, templates)
	}))

	assert.Equal(t, expectedTopic, captureStdout(func() {
		os.Args = []string{"me", "help", "hey"}
		hadfield.Exit = func(_ int) {}
		hadfield.Run(cmds, templates)
	}))
}
