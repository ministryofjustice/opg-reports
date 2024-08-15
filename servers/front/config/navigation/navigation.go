package navigation

import (
	"net/http"

	"github.com/ministryofjustice/opg-reports/shared/convert"
	"github.com/ministryofjustice/opg-reports/shared/must"
)

func flat(ni []*NavigationItem, f map[string]*NavigationItem) {
	for _, i := range ni {
		f[i.Uri] = i
		if i.Navigation != nil && len(i.Navigation) > 0 {
			flat(i.Navigation, f)
		}
	}
}

func Level(items []*NavigationItem, r *http.Request) (nav map[string]*NavigationItem, active *NavigationItem) {
	nav = map[string]*NavigationItem{}
	for _, n := range items {
		n.Active = n.InUrlPath(r.URL.Path)
		if n.Active {
			active = n
		}
		nav[n.Uri] = n
	}
	return
}

func Flat(items []*NavigationItem, r *http.Request) (f map[string]*NavigationItem) {
	f = map[string]*NavigationItem{}
	flat(items, f)
	for _, n := range f {
		n.Active = n.InUrlPath(r.URL.Path)
	}
	return
}

func ForTemplate(templateName string, navs []*NavigationItem) (found *NavigationItem) {
	fl := map[string]*NavigationItem{}
	flat(navs, fl)
	for _, n := range fl {
		if n.Template == templateName {
			found = n
		}
	}
	return
}

func NewNav(content []byte) ([]*NavigationItem, error) {
	// n := []*NavigationItem{}
	return convert.Unmarshals[*NavigationItem](content)
}

func New(content []byte) []*NavigationItem {
	return must.Must[[]*NavigationItem](NewNav(content))
}
