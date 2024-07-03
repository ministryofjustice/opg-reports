package tmpl

import (
	"fmt"
	"opg-reports/services/front/cnf"
	"opg-reports/shared/data"
	"opg-reports/shared/dates"
	"opg-reports/shared/files"
	"opg-reports/shared/github/std"
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
		"percent": func(got int, total int) string {
			x := float64(got)
			y := float64(total)
			p := x / (y / 100)
			return fmt.Sprintf("%.2f", p)
		},
		"title": func(s string) string {
			s = strings.ReplaceAll(s, "_", " ")
			s = strings.ReplaceAll(s, "-", " ")
			c := cases.Title(language.English)
			s = c.String(s)
			return s
		},
		"lower": func(s string) string {
			return strings.ToLower(s)
		},
		"dict": func(values ...any) (dict map[string]any) {
			dict = map[string]any{}
			if len(values)%2 != 0 {
				return
			}
			// if the key isnt a string, this will crash!
			for i := 0; i < len(values); i += 2 {
				var key string = values[i].(string)
				var v any = values[i+1]
				dict[key] = v
			}
			return
		},
		// "append": func(slice []string, value string) []string {
		// 	return append(slice, value)
		// },
		// "join": func(sep string, s ...string) string {
		// 	return strings.Join(s, sep)
		// },
		"month": func(d time.Time) string {
			return d.Format(dates.FormatYM)
		},
		"day": func(d time.Time) string {
			return d.Format(dates.FormatYMD)
		},
		// Costs
		"dollar": func(f float64) string {
			p := message.NewPrinter(language.English)
			return p.Sprintf("$%.2f", f)
		},
		"dollars": func(s string) string {
			p := message.NewPrinter(language.English)
			f, _ := strconv.ParseFloat(s, 10)
			return p.Sprintf("$%.2f", f)
		},
		"addi": func(a int, b int) any {
			return a + b
		},
		// Compliance
		"getComplianceItem": func(row *response.Row[*response.Cell]) *std.Repository {
			return data.FromRow[*std.Repository](row)
		},
		"repoSetStandards": func(c *std.Repository, standards *cnf.RepoStandards) error {
			c.SetStandards(standards)
			return nil
		},
		"repoStandardPassed": func(c *std.Repository, standards []string) (pass bool) {
			pass, _, _ = c.Compliant(standards)
			return
		},
		"repoStandardDetail": func(c *std.Repository, standards []string) (detail map[string]bool) {
			_, detail, _ = c.Compliant(standards)
			return
		},
		"repoStandardValues": func(c *std.Repository, fields []string) (detail map[string]interface{}) {
			detail = map[string]interface{}{}
			m, _ := data.ToMap(c)
			for _, k := range fields {
				detail[k] = m[k]
			}
			return
		},
		"totalCountPassed": func(rows []*response.Row[*response.Cell], standards []string) (count int) {
			count = 0
			for _, row := range rows {
				c := data.FromRow[*std.Repository](row)
				if pass, _, _ := c.Compliant(standards); pass {
					count += 1
				}
			}
			return
		},
	}

}
