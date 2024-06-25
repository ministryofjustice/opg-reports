package server

import (
	"net/http"
	"opg-reports/services/front/cnf"
	"strings"
)

type Navigation struct {
	tree []*cnf.SiteSection
}

func (n *Navigation) Top(r *http.Request) (active *cnf.SiteSection, all []*cnf.SiteSection) {
	all = n.tree
	compare := r.RequestURI
	for _, nav := range all {
		if !nav.Exclude && strings.Contains(compare, nav.Href) {
			active = nav
		}
	}
	return
}

func (n *Navigation) Side(r *http.Request, top *cnf.SiteSection) (active *cnf.SiteSection, all []*cnf.SiteSection) {
	if top == nil {
		return
	}
	all = top.Sections
	compare := r.RequestURI
	for _, nav := range all {
		if !nav.Exclude && strings.Contains(compare, nav.Href) {
			active = nav
		}
	}
	return
}

func (n *Navigation) All() (flat map[string]*cnf.SiteSection) {
	flat = map[string]*cnf.SiteSection{}
	cnf.FlatSections(n.tree, flat)
	return
}

func (n *Navigation) Get(uri string) (active *cnf.SiteSection) {
	flat := n.All()
	for u, sect := range flat {
		if u == uri {
			active = sect
		}
	}
	return
}

func (n *Navigation) Active(r *http.Request) (active *cnf.SiteSection) {
	compare := r.RequestURI
	return n.Get(compare)
}

func (n *Navigation) Data(r *http.Request) map[string]interface{} {
	activeItem := n.Active(r)
	activeTop, allTop := n.Top(r)
	activeSide, allSide := n.Side(r, activeTop)

	return map[string]interface{}{
		"NavigationTop":        allTop,
		"NavigationTopActive":  activeTop,
		"NavigationSide":       allSide,
		"NavigationSideActive": activeSide,
		"NavigationActive":     activeItem,
	}
}
