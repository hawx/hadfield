package hadfield

import (
	"io"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"
)

// Templates defines how the help screens are displayed given in the
// text/template format.
type Templates struct {
	// Help is the template rendered to display the help for the executable, that
	// is the one shown when "$0 help", "$0 -h" or "$0 --help" are called.
	Help Template

	// Command is the template rendered to display help for a particular
	// subcommand, it is shown when "$0 help [subcommand]", "$0 [subcommand] -h"
	// or "$0 [subcommand] --help" are called.
	Command Template
}

var templateFuncs = template.FuncMap{
	"trim":       strings.TrimSpace,
	"capitalize": capitalize,
}

type Template string

func (text *Template) Render(w io.Writer, data interface{}) {
	t := template.New("top")
	t.Funcs(templateFuncs)

	template.Must(t.Parse(string(*text)))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToTitle(r)) + s[n:]
}
