package template_helpers

import (
	"strings"
	"time"

	"github.com/ministryofjustice/opg-reports/datastore/github_standards/ghs"
	"github.com/ministryofjustice/opg-reports/servers/api/aws_costs"
	"github.com/ministryofjustice/opg-reports/shared/convert"
	"github.com/ministryofjustice/opg-reports/shared/dates"
)

func Funcs() map[string]interface{} {

	return map[string]interface{}{
		"currency": func(s interface{}, symbol string) string {
			return convert.Curr(s, symbol)
		},
		"add": func(a float64, b float64) float64 {
			return a + b
		},
		"addInt": func(a int, b int) int {
			return a + b
		},
		"percent": func(got int, total int) string {
			return convert.Percent(got, total)
		},
		"stripI": func(s string) string {
			return convert.StripIntPrefix(s)
		},
		"title": func(s string) string {
			return convert.Title(s)
		},
		"titles": func(strs ...string) (str string) {
			str = ""
			for _, s := range strs {
				str += s + " - "
			}
			str = strings.TrimSuffix(str, " - ")
			return
		},
		"lower": func(s string) string {
			return strings.ToLower(s)
		},
		"col": func(i string, mapped map[string]string) string {
			return mapped[i]
		},
		"dict": func(values ...any) map[string]any {
			return convert.Dict(values...)
		},
		"month": func(d time.Time) string {
			return d.Format(dates.FormatYM)
		},
		"day": func(t string) string {
			d := dates.Time(t)
			return d.Format(dates.FormatYMD)
		},
		"dayBefore": func(t string) string {
			d := dates.Time(t).AddDate(0, 0, -1)
			return d.Format(dates.FormatYMD)
		},
		"costIdx": func(set []*aws_costs.CommonResult, i int) any {
			return set[i]
		},
		// -- casting
		"modelGHS": func(m map[string]interface{}) (g *ghs.GithubStandard) {
			g, _ = convert.Unmap[*ghs.GithubStandard](m)
			return
		},
	}

}
