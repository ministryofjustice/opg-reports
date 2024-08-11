package navigation

import (
	"strings"

	"github.com/ministryofjustice/opg-reports/servers/front/config/page"
)

type NavigationItem struct {
	Name        string            `json:"name"`
	Uri         string            `json:"uri"`
	Template    string            `json:"template"`
	IsHeader    bool              `json:"is_header"`
	DataSources page.Data         `json:"data_sources"`
	Navigation  []*NavigationItem `json:"navigation"`
	Active      bool              `json:"-"`
	Registered  bool              `json:"-"`
}

func (n *NavigationItem) ClassName() string {
	str := "sect-"
	str = str + strings.ToLower(n.Name)
	str = strings.ReplaceAll(str, " ", "-")
	return str
}

func (n *NavigationItem) InUrlPath(url string) bool {
	return strings.HasPrefix(url, n.Uri)
}

func (n *NavigationItem) Matches(url string) bool {
	return n.Uri == url
}
