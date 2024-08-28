package nav

// ForTemplate creates a list of all nav items from the list passed
// whose template matches the templateName
func ForTemplate(templateName string, navs []*Nav) (found []*Nav) {
	flatMap := map[string]*Nav{}
	found = []*Nav{}
	flat(navs, flatMap)
	for _, nav := range flatMap {
		if nav.Template == templateName {
			found = append(found, nav)
		}
	}
	return
}
