package htmlpage

import (
	"fmt"
	"html/template"
	"io"
	"opg-reports/report/config"
	"path/filepath"
	"strings"
)

type HtmlPage struct {
	files []string
	funcs template.FuncMap
}

func (self *HtmlPage) GetTemplate(name string) (tmpl *template.Template, err error) {
	tmpl, err = template.New(name).
		Funcs(self.funcs).
		ParseFiles(self.files...)
	return
}

func (self *HtmlPage) WriteToBuffer(buffer io.Writer, templateName string, data any) (err error) {
	var templ *template.Template

	templateName = pageName(templateName)
	templ, err = self.GetTemplate(templateName)
	if err != nil {
		return
	}
	err = templ.ExecuteTemplate(buffer, templateName, data)
	return
}

type HtmlPageContent struct {
	Name         string
	GovUKVersion string
	Signature    string
}

func New(templateFiles []string, templateFunctions template.FuncMap) *HtmlPage {
	return &HtmlPage{
		files: templateFiles,
		funcs: templateFunctions,
	}
}

func DefaultContent(conf *config.Config) HtmlPageContent {
	return HtmlPageContent{
		Name:         conf.Servers.Front.Name,
		GovUKVersion: strings.TrimPrefix(conf.GovUK.Front.ReleaseTag, "v"),
		Signature:    conf.Versions.Signature(),
	}
}

func GetTemplateFiles(directory string) (files []string) {
	pattern := filepath.Join(directory, "**/**")
	files, err := filepath.Glob(pattern)
	if err != nil {
		fmt.Println("err:" + err.Error())
	}
	return
}

func pageName(n string) (name string) {
	name = n
	name = filepath.Base(name)
	name = strings.Split(name, ".")[0]
	return
}
