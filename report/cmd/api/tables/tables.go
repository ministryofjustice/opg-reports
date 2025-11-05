package tables

import (
	"bytes"
	"fmt"
	"log/slog"
	"opg-reports/report/internal/service/api"
	"opg-reports/report/internal/utils"
	"slices"
	"sort"
	"strings"
)

type ListToTableConfig struct {
	TextColumns     []string // these are the text / non-data fields are the start of the table, used for row headers etc
	DataColumns     []string // the columns that will be added to each row as data items, generally values for each date interval
	ValueField      string   // the field from each list record to use for the sorce of data
	DataSourceField string   // the field used to compare to the DataColumn items - generallt `date`
	DefaultValue    string   // default value used for a blank cell
}

func ListToTable[T api.Model](
	log *slog.Logger,
	cfg *ListToTableConfig,
	dbRecords []T,

) (table map[string]map[string]string, err error) {
	var (
		records []map[string]string = []map[string]string{} // converted version of records into a generic slice map
		rowKeys []string            = []string{}            // contains all the possible row keys based on the values of the group columns
	)
	// init
	log = log.With("operation", "ListToTable", "cfg", cfg)
	table = map[string]map[string]string{}
	// sort dates
	sort.Strings(cfg.DataColumns)

	// convert to slice map from the records
	log.Debug("converting []T to slice of maps")
	err = utils.Convert(dbRecords, &records)
	if err != nil {
		log.Error("error converting T[] to slice of maps")
		return
	}
	// generate row keys
	if len(cfg.TextColumns) <= 0 {
		err = fmt.Errorf("require at least 1 column to generate a table in ListToTable")
		return
	}
	log.Debug("generating row keys from data and columns")
	rowKeys, _ = PossibleCombinationsAsKeys(records, cfg.TextColumns)
	// now create a skeleton table from the rowKeys & date values
	log.Debug("generating skeleton table structure")
	table = Skeleton(rowKeys, cfg.DataColumns, cfg.DefaultValue)

	// now populate the table
	log.Debug("populating table skeleton with real data")
	table = Populate(&PopulateConfig{
		Skeleton:        table,
		TextColumns:     cfg.TextColumns,
		ValueField:      cfg.ValueField,
		DataSourceField: cfg.DataSourceField,
		DefaultValue:    cfg.DefaultValue,
	}, records)

	return
}

type PopulateConfig struct {
	Skeleton        map[string]map[string]string
	TextColumns     []string // these are the text / non-data fields are the start of the table, used for row headers etc
	ValueField      string   // the field from each list record to use for the sorce of data
	DataSourceField string   // the field used to compare to the DataColumn items
	DefaultValue    string   // default value used for a blank cell
}

// Populate loops over each data item and inject its value (from valueColumn)
// into the transformed tables colum (transformColumn) based on the items unique
// key (from identifiers)
//
// This generates a full table, with values in every column merged from the raw
// dataset
func Populate(
	cfg *PopulateConfig,
	data []map[string]string,
) (table map[string]map[string]string) {
	// assign to self
	table = cfg.Skeleton

	for _, item := range data {
		key := combinationKey(item, cfg.TextColumns...)
		column := item[cfg.DataSourceField]

		if row, ok := table[key]; ok {
			row[column] = item[cfg.ValueField]
		}
	}

	// remove empty rows...
	for id, row := range table {
		var empty = true
		for k, v := range row {
			if !slices.Contains(cfg.TextColumns, k) && v != cfg.DefaultValue {
				empty = false
			}
		}
		if empty {
			delete(table, id)
		}
	}

	return table
}

// Skeleton reates a series of skeleton rows from the known
// column values and date range to form a map of table rows.
//
// `keys` is the output from `PossibleCombinationsAsKeys`
// `cells` is normally the date columns
//
// Output:
//
//	map[string]map[string]string{}{
//		"environment:dev^service:ec2^" : map[string]string{}{
//			"service": "",
//			"environment": "",
//			"2024-01": 0.0,
//			"2024-02": 0.0,
//		},
//		"environment:dev^service:ecs^" : map[string]string{}{
//			"service": "",
//			"environment": "",
//			"2024-01": 0.0,
//			"2024-02": 0.0,
//		},
//		"environment:prod^service:ec2^" : map[string]string{}{
//			"service": "",
//			"environment": "",
//			"2024-01": 0.0,
//			"2024-02": 0.0,
//		},
//		"environment:prod^service:ecs^" : map[string]string{}{
//			"service": "",
//			"environment": "",
//			"2024-01": 0.0,
//			"2024-02": 0.0,
//		},
//	}
func Skeleton(keys []string, cells []string, emptyCell string) (table map[string]map[string]string) {
	table = map[string]map[string]string{}

	for _, key := range keys {
		row := map[string]string{}
		// recreate the column name and value from the formatted key
		// 	- "environment:backup^account:A" => {"environment":"backup", "account":"A"}
		key = strings.TrimSuffix(key, "^")
		for _, columnAndValue := range strings.Split(key, "^") {
			sp := strings.Split(columnAndValue, ":")
			col, val := sp[0], sp[1]
			row[col] = val
		}
		// append the extra cells on to the table
		for _, cell := range cells {
			row[cell] = emptyCell
		}
		table[key+"^"] = row
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
	keys = permutations(uniques...)
	return

}

// AddColumnToEachRow injects column in every row in the table as empty
func AddColumnToEachRow(table map[string]map[string]string, col string) {
	for _, row := range table {
		row[col] = ""
	}
}

// permutations merges the values of parts together to find all the
// possible combinations
//
// Input:
//
//	[][]string {
//		[]string{"A", "B", "C"}
//		[]string{"1", "2"}
//	}
//
// Output:
//
//	[]string {"A1", "A2", "B1", "B2", "C1", "C2"}
//
// Resulting length is length of each part passed multiplied by each
// other. So in the example above its 3 x 2 = 6
//
// Is the basis for generating complete table rows from api data.
func permutations(parts ...[]string) (ret []string) {
	{
		var n = 1
		for _, ar := range parts {
			n *= len(ar)
		}
		ret = make([]string, 0, n)
	}
	var at = make([]int, len(parts))
	var buf bytes.Buffer
loop:
	for {
		// increment position counters
		for i := len(parts) - 1; i >= 0; i-- {
			if at[i] > 0 && at[i] >= len(parts[i]) {
				if i == 0 || (i == 1 && at[i-1] == len(parts[0])-1) {
					break loop
				}
				at[i] = 0
				at[i-1]++
			}
		}
		// construct permutated string
		buf.Reset()
		for i, ar := range parts {
			var p = at[i]
			if p >= 0 && p < len(ar) {
				buf.WriteString(ar[p])
			}
		}
		ret = append(ret, buf.String())
		at[len(parts)-1]++
	}
	return ret
}

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

// combinationKey generates a key from the item and config
func combinationKey(item map[string]string, identifiers ...string) (key string) {
	slices.Sort(identifiers)
	identifiers = slices.Compact(identifiers)

	key = ""
	for _, k := range identifiers {
		key += fmt.Sprintf("%s:%s^", k, item[k])
	}
	return
}
