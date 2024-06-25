package tmpl

import (
	"fmt"
	"opg-reports/shared/data"
	"opg-reports/shared/files"
	"strings"
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
		"heading": func(s string) string {
			s = strings.ReplaceAll(s, "_", " ")
			s = strings.ReplaceAll(s, "-", " ")
			return strings.Title(strings.ToLower(s))
		},
		"toIdx": func(k string, v string) string {
			return data.ToIdxKV(k, v)
		},
		"getIdx": func(i string, data map[string]any, def any) any {
			if v, ok := data[i]; ok {
				return v
			}
			return def
		},
		"dollars": func(s float64) string {
			return fmt.Sprintf("$%.2f", s)
		},
	}

}
