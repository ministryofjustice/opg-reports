package tmpl

import (
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Title cleans up string for page title usage
func Title(s string) string {
	s = strings.ReplaceAll(s, "/", " ")
	s = strings.ReplaceAll(s, "_", " ")
	c := cases.Title(language.English)
	s = c.String(s)
	return s
}

func Functions() (funcs template.FuncMap) {
	funcs = map[string]interface{}{
		"Title":   Title,
		"ToLower": strings.ToLower,
		"ToUpper": strings.ToUpper,
	}
	return
}
