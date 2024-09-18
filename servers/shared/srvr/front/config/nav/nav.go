package nav

import (
	"strings"

	"github.com/ministryofjustice/opg-reports/shared/convert"
)

type Nav struct {
	Name        string            `json:"name"`
	Uri         string            `json:"uri"`
	Template    string            `json:"template"`
	IsHeader    bool              `json:"is_header"`
	DataSources map[string]string `json:"data_sources"`
	Navigation  []*Nav            `json:"navigation"`
	Active      bool              `json:"-"`
	Registered  bool              `json:"-"`
}

func (n *Nav) ClassName() string {
	str := n.Uri
	str = strings.TrimPrefix(str, "/")
	str = strings.TrimSuffix(str, "/")
	str = strings.ReplaceAll(str, "/", "-")

	return str
}

func (n *Nav) InUrlPath(url string) bool {
	return strings.HasPrefix(url, n.Uri)
}

func (n *Nav) Matches(url string) bool {
	return n.Uri == url
}

func New(content []byte) (navList []*Nav) {
	if set, err := convert.Unmarshals[*Nav](content); err == nil {
		navList = set
	}
	return
}
