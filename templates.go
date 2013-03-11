package hadfield

import (
	"text/template"
	"strings"
	"unicode"
	"unicode/utf8"
	"io"
)

type Templates struct {
	Usage  Template
	Help   Template
}

type Template string

func (text *Template) Render(w io.Writer, data interface{}) {
	t := template.New("top")
	t.Funcs(template.FuncMap{
		"trim": strings.TrimSpace,
		"capitalize": capitalize,
	})

	template.Must(t.Parse(string(*text)))
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
