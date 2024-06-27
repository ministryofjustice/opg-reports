package tmpl

import (
	"opg-reports/shared/dates"
	"opg-reports/shared/files"
	"opg-reports/shared/gh/comp"
	"opg-reports/shared/server/response"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func Files(fS *files.WriteFS, prefix string) []string {
	allFiles := files.All(fS, false)
	filtered := files.Filter(allFiles, `\.gotmpl$`)
	templateFiles := []string{}
	for _, f := range filtered {
		if prefix != "" {
			f.Path = prefix + f.Path
		}
		templateFiles = append(templateFiles, f.Path)
	}

	return templateFiles

}

func Funcs() map[string]interface{} {

	return map[string]interface{}{
		"append": func(slice []string, value string) []string {
			return append(slice, value)
		},
		"join": func(sep string, s ...string) string {
			return strings.Join(s, sep)
		},
		"dollars": func(s string) string {
			p := message.NewPrinter(language.English)
			f, _ := strconv.ParseFloat(s, 10)
			return p.Sprintf("$%.2f", f)
		},
		"month": func(d time.Time) string {
			return d.Format(dates.FormatYM)
		},
		"addi": func(a int, b int) any {
			return a + b
		},
		"complyI": func(row *response.Row[*response.Cell]) *comp.Compliance {
			return comp.FromRow(row)
		},
	}

}
