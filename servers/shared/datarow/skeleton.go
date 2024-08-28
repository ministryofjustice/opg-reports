package datarow

import "github.com/ministryofjustice/opg-reports/shared/convert"

// Skeleton uses the set of columns and intervals to generate a row for each possible type
func Skeleton(columnsAndPossible map[string][]string, intervals map[string][]string) (skel map[string]map[string]interface{}) {
	skel = map[string]map[string]interface{}{}

	keys := convert.PermuteStrings(flatKeys(columnsAndPossible)...)
	for _, key := range keys {
		if _, ok := skel[key]; !ok {
			skel[key] = map[string]interface{}{}
		}
		for intName, values := range intervals {
			i := map[string]interface{}{}
			for _, val := range values {
				i[val] = emptyVal
			}
			skel[key][intName] = i
		}

		skel[key]["columns"] = splitKey(key)
	}
	return
}
