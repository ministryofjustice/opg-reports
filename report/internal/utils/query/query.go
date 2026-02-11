package query

import (
	"fmt"
	"strings"
)

// base is the structure of a select
const base string = `
SELECT
{select}
FROM {from}
{joins}
{where}
{groupby}
;`

type Segment struct {
	Select  string
	Where   string
	GroupBy string
}

type Select struct {
	From     string
	Joins    string
	Segments map[string][]*Segment
}

func (self *Select) FromRequest(request map[string]interface{}) (query string) {
	query = fromRequest(self, self.Segments, request)
	return
}

func fromRequest(stmt *Select, segments map[string][]*Segment, request map[string]interface{}) (query string) {
	var (
		selectStr string = ""
		whereStr  string = ""
		groupStr  string = ""
	)

	query = strings.ReplaceAll(base, `{from}`, stmt.From)
	query = strings.ReplaceAll(query, `{joins}`, stmt.Joins)

	for key, value := range request {
		var useWhere = false
		var val = strings.ToLower(value.(string))
		if val != "true" {
			useWhere = true
		}

		for _, seg := range segments[key] {
			if seg.Select != "" {
				selectStr += fmt.Sprintf("\n  %s,", seg.Select)
			}

			if seg.GroupBy != "" {
				groupStr += fmt.Sprintf("\n  %s,", seg.GroupBy)
			}

			if useWhere && seg.Where != "" {
				whereStr += fmt.Sprintf("\n  %s AND", seg.Where)
			}
		}

	}
	// clean up select
	selectStr = strings.TrimSuffix(selectStr, ",")
	selectStr = strings.TrimPrefix(selectStr, "\n")
	query = strings.ReplaceAll(query, `{select}`, selectStr)

	// clean up where
	if whereStr != "" {
		whereStr = fmt.Sprintf("WHERE %s", whereStr)
	}
	whereStr = strings.TrimSuffix(whereStr, "AND")
	whereStr = strings.TrimPrefix(whereStr, "\n")
	query = strings.ReplaceAll(query, `{where}`, whereStr)

	// clean up group
	if groupStr != "" {
		groupStr = fmt.Sprintf("GROUP BY %s", groupStr)
	}
	groupStr = strings.TrimSuffix(groupStr, ",")
	groupStr = strings.TrimPrefix(groupStr, "\n")
	query = strings.ReplaceAll(query, `{groupby}`, groupStr)

	return
}
