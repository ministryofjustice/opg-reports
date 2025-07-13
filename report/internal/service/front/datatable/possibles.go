package datatable

import (
	"fmt"
	"slices"
)

// uniqueValuesForEachIdentifier generates a slice of slices, with each slice representing
// each identifiers unique values that are within the data (where identifier is the map key)
func uniqueValuesForEachIdentifier(data []map[string]string, identifiers ...string) (combinations [][]string) {
	// sort and remove any duplicate keys
	slices.Sort(identifiers)
	identifiers = slices.Compact(identifiers)

	combinations = [][]string{}

	for _, key := range identifiers {
		var options = []string{}

		for _, item := range data {
			if v, ok := item[key]; ok {
				options = append(options, fmt.Sprintf("%s:%s^", key, v))
			}
		}
		// make unique
		slices.Sort(options)
		options = slices.Compact(options)
		combinations = append(combinations, options)
	}

	return
}

// PossibleCombinationsAsKeys used in converting an api response data set into a tabluar data
// structure. It find all the unique values for each identifier (think map key / column name)
// and generates a list of all possible combinations based on their values within `data`.
// These keys are then used to create table rows.
//
// The passed `data` map should be the value of the apiResponse.Data.
//
// Example:
//
//		Input
//			data = []map[string]string {
//				map[string]string{
//					"account": "A"
//					"region": "2024",
//					"cost": "100"
//				},
//				map[string]string{
//					"account": "B"
//					"region": "2025",
//					"cost": "100"
//				},
//				map[string]string{
//					"account": "A"
//					"region": "2024",
//					"cost": "100"
//				},
//			}
//			identifiers = "account", "region"
//	  Output
//			keys = []string{
//				"account:A^region:2024^",
//				"account:A^region:2025^",
//				"account:B^region:2024^",
//				"account:B^region:2025^"
//			}
//			uniques = [][]string{
//				[]string{"account:A^", "account:B^"}
//				[]string{"region:2024^", "region:2025^"}
//			}
//
// The `keys` returned can be used to reform the data grouped by the identifiers.
func PossibleCombinationsAsKeys(data []map[string]string, identifiers []string) (keys []string, uniques [][]string) {
	uniques = [][]string{}

	slices.Sort(identifiers)
	// get all unique values for each data key
	uniques = uniqueValuesForEachIdentifier(data, identifiers...)
	// generate flat list of possible values - these will be used for row keys
	keys = Permutations(uniques...)
	return

}
