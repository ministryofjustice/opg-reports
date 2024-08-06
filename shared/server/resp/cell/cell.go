// Package Cell handles a single peice of information within the table
// structure. It uses Name & Value to represent content and
// some bool flags (IsHeader, IsSupplementary) to give insight
// on what type of cells this may be (header, footer etc)

package cell

type Cell struct {
	Name            string      `json:"name"`
	Value           interface{} `json:"value"`
	IsHeader        bool        `json:"is_header"`
	IsSupplementary bool        `json:"is_supplementary"`
}

func New(name string, value interface{}, isHeader bool, isExtra bool) *Cell {
	return &Cell{Name: name, Value: value, IsHeader: isHeader, IsSupplementary: isExtra}
}

func Default(name string, value interface{}) *Cell {
	return New(name, value, false, false)
}
