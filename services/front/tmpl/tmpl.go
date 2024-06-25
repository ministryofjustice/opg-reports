package tmpl

import "opg-reports/shared/files"

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
