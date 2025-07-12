package datatable

import "slices"

// PopulateTable loops over each data item and inject its value (from valueColumn)
// into the transformed tables colum (transformColumn) based on the items unique
// key (from identifiers)
//
// This generates a full table, with values in every column merged from the raw
// dataset
func PopulateTable(
	data []map[string]string, table map[string]map[string]string,
	identifiers []string,
	transformColumn string, valueColumn string,
) map[string]map[string]string {

	for _, item := range data {
		key := CombinationKey(item, identifiers...)
		column := item[transformColumn]

		if row, ok := table[key]; ok {
			row[column] = item[valueColumn]
		}
	}

	// remove empty rows...
	for id, row := range table {
		var empty = true
		for k, v := range row {
			if !slices.Contains(identifiers, k) && v != emptyCell {
				empty = false
			}
		}
		if empty {
			delete(table, id)
		}
	}

	return table
}
