package nav

import "net/http"

// Level checks only the current level (does not recurse) of the nav items
// passed to see if any "active" (url contains current url)
// If active is found then also run the activate check on the subitems
func Level(items []*Nav, r *http.Request) (nav map[string]*Nav, activeSection *Nav, activePage *Nav) {
	nav = map[string]*Nav{}
	activeSection = items[0]
	for _, n := range items {
		n.Active = n.InUrlPath(r.URL.Path)
		if n.Active {
			activeSection = n
		}
		nav[n.Uri] = n
	}
	activePage = activeSection
	if activeSection != nil {
		if ap := Activate(activeSection.Navigation, r); ap != nil {
			activePage = ap
		}
	}
	return
}

// Activate checks each nav item in the tree of items passed
// to see if the url exactly matches the current url path
// and sets the .Active flag when it does
//
// Used to help nav display
func Activate(items []*Nav, r *http.Request) (active *Nav) {
	flatMap := map[string]*Nav{}
	flat(items, flatMap)
	for _, n := range flatMap {
		n.Active = n.Matches(r.URL.Path)
		if n.Active {
			active = n
		}
	}
	return
}
