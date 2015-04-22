package hadfield_test

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	hadfield "."
	"github.com/stretchr/testify/assert"
)

var nilRun = func(c *hadfield.Command, args []string) {}

func captureStdout(f func()) string {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	out, _ := ioutil.ReadAll(r)
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
		Usage: `usage: test [command] [arguments]

  This is a test.

  Commands:{{range .}}
    {{.Name | printf "%-15s"}} # {{.Short}}{{end}}
`,
	}

	expectedOut := `usage: test [command] [arguments]

  This is a test.

  Commands:
    he              # he does stuff
    hey             # hey does other stuff
    bye             # bye goes away
`

	assert.Equal(t, expectedOut, captureStdout(func() {
		os.Args = []string{"me", "help"}
		hadfield.Exit = func(_ int) {}
		hadfield.Run(cmds, templates)
	}))
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
		Help: `usage: test {{.Usage}}
{{.Long}}
`,
	}

	expectedOut := `usage: test hey

  Hey is a command to say hey, HEY!

  Options:
    --later WHEN   # later is cool for now, BUT LATER (default: cool)
    --now REALLY   # now does stuff right NOW

`

	assert.Equal(t, expectedOut, captureStdout(func() {
		os.Args = []string{"me", "help", "hey"}
		hadfield.Exit = func(_ int) {}
		hadfield.Run(cmds, templates)
	}))
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
		Usage: `usage: test [command] [arguments]

  This is a test.

  Commands:{{range .}}{{if .Callable}}
    {{.Name | printf "%-15s"}} # {{.Short}}{{end}}{{end}}

  Additional help:{{range .}}{{if not .Callable}}
    {{.Name | printf "%-15s"}} # {{.Short}}{{end}}{{end}}
`,
		Help: `{{if .Callable}}usage: test {{.Usage}}
{{end}}{{.Long}}
`,
	}

	expectedUsage := `usage: test [command] [arguments]

  This is a test.

  Commands:
    hey             # runs hey

  Additional help:
    hey             # hey does other stuff
`

	expectedTopic := `
This is actually just documentation about the "hey" system.

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
