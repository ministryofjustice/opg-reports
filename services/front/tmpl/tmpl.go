package tmpl

import (
	"opg-reports/services/front/cnf"
	"opg-reports/shared/data"
	"opg-reports/shared/dates"
	"opg-reports/shared/files"
	"opg-reports/shared/gh/comp"
	"opg-reports/shared/server/response"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/cases"
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
		"title": func(s string) string {
			s = strings.ReplaceAll(s, "_", " ")
			s = strings.ReplaceAll(s, "-", " ")
			c := cases.Title(language.English)
			s = c.String(s)
			return s
		},
		"append": func(slice []string, value string) []string {
			return append(slice, value)
		},
		"join": func(sep string, s ...string) string {
			return strings.Join(s, sep)
		},
		"month": func(d time.Time) string {
			return d.Format(dates.FormatYM)
		},
		// Costs
		"dollars": func(s string) string {
			p := message.NewPrinter(language.English)
			f, _ := strconv.ParseFloat(s, 10)
			return p.Sprintf("$%.2f", f)
		},
		"addi": func(a int, b int) any {
			return a + b
		},
		// Compliance
		"getComplianceItem": func(row *response.Row[*response.Cell]) *comp.Compliance {
			return comp.FromRow(row)
		},
		"repoSetStandards": func(c *comp.Compliance, standards *cnf.RepoStandards) error {
			c.SetStandards(standards)
			return nil
		},
		"repoStandardPassed": func(c *comp.Compliance, standards []string) (pass bool) {
			pass, _, _ = c.Compliant(standards)
			return
		},
		"repoStandardDetail": func(c *comp.Compliance, standards []string) (detail map[string]bool) {
			_, detail, _ = c.Compliant(standards)
			return
		},
		"repoStandardValues": func(c *comp.Compliance, fields []string) (detail map[string]interface{}) {
			detail = map[string]interface{}{}
			m, _ := data.ToMap(c)
			for _, k := range fields {
				detail[k] = m[k]
			}
			return
		},
	}

}
