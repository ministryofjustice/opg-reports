package nav

import "net/http"

// flat recurses down the navItems passed adding each the flatmap
// based on the uri - which should be unique
func flat(navItems []*Nav, flatmap map[string]*Nav) {
	for _, item := range navItems {
		flatmap[item.Uri] = item
		if item.Navigation != nil && len(item.Navigation) > 0 {
			flat(item.Navigation, flatmap)
		}
	}
}

// Flat creates a flatMap from the original items passed.
// It marks any nav item whose url is contained within
// the current request url path as being active
func Flat(items []*Nav, r *http.Request) (flatMap map[string]*Nav) {
	flatMap = map[string]*Nav{}
	flat(items, flatMap)
	for _, n := range flatMap {
		n.Active = n.InUrlPath(r.URL.Path)
	}
	return
}
