package utils

import (
	"html/template"
	"strings"
)

func TemplateFunctions() (funcs template.FuncMap) {
	funcs = map[string]interface{}{
		// simple strings
		"ToLower": strings.ToLower,
		"ToUpper": strings.ToUpper,
	}

	return
}
