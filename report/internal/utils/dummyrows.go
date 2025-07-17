package utils

// DummyRows generates fake rows - this is normally added before the
// PossibleCombinationsAsKeys calls to make sure tihngs like all dates
// within a given range are present, when there is a chance that
// they might not be a cost for a particular month.
func DummyRows(extras []string, key string) (dummys []map[string]string) {
	dummys = []map[string]string{}
	for _, d := range extras {
		dummys = append(dummys, map[string]string{key: d})
	}
	return
}
