package front_templates

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ministryofjustice/opg-reports/datastore/aws_costs/awsc"
	"github.com/ministryofjustice/opg-reports/datastore/github_standards/ghs"
	"github.com/ministryofjustice/opg-reports/shared/convert"
	"github.com/ministryofjustice/opg-reports/shared/dates"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func Funcs() map[string]interface{} {

	return map[string]interface{}{
		"currency": func(s interface{}, symbol string) string {
			p := message.NewPrinter(language.English)
			switch s.(type) {
			case string:
				f, _ := strconv.ParseFloat(s.(string), 10)
				return symbol + p.Sprintf("%.2f", symbol, f)
			case float64:
				return symbol + p.Sprintf("%.2f", s.(float64))
			}
			return symbol + "0.0"
		},
		"add": func(a float64, b float64) float64 {
			return a + b
		},
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
		"col": func(i string, mapped map[string]string) string {
			return mapped[i]
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
		"month": func(d time.Time) string {
			return d.Format(dates.FormatYM)
		},
		"day": func(t string) string {
			d := dates.Time(t)
			return d.Format(dates.FormatYMD)
		},
		// -- casting
		"modelGHS": func(m map[string]interface{}) (g *ghs.GithubStandard) {
			g, _ = convert.Unmap[*ghs.GithubStandard](m)
			return
		},
		"modelAWSCTax": func(m map[string]interface{}) (i *awsc.MonthlyTotalsTaxSplitRow) {
			i, _ = convert.Unmap[*awsc.MonthlyTotalsTaxSplitRow](m)
			return
		},
		"modelAWSMonthlyPerUnit": func(m map[string]interface{}) (i *awsc.MonthlyCostsPerUnitRow) {
			i, _ = convert.Unmap[*awsc.MonthlyCostsPerUnitRow](m)
			return
		},
	}

}
