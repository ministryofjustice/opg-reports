package template_functions

import (
	"fmt"
	"strings"
	"time"

	"github.com/ministryofjustice/opg-reports/shared/dates"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

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
		"month": func(d time.Time) string {
			return d.Format(dates.FormatYM)
		},
		"day": func(d time.Time) string {
			return d.Format(dates.FormatYMD)
		},
	}

}
