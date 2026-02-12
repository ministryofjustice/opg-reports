package headers

import (
	"slices"
	"sort"
)

// Type used as enum constraint
type Type string

// types of headers we'd use in a table
const (
	KEY   Type = "key"
	DATA  Type = "data"
	EXTRA Type = "extra"
	END   Type = "end"
)

// Header represents a table heading and cotnains type and default value.
//
// Used for display and as part of table conversion from list to table row
// structure on complex data types like costs / uptime
type Header struct {
	Field   string      `json:"field"`   // name of the field / column this would be within map[string]interface{}
	Type    Type        `json:"type"`    // list of types of columns this aligns with
	Default interface{} `json:"default"` // default value to use for this field when creating a skeleton
}

// Headers is collection of all headers used within a table
//
// Used for display and as part of table conversion from list to table row
// structure on complex data types like costs / uptime
type Headers struct {
	Headers []*Header `json:"headers"`
}

func (self *Headers) AddKeyHeader(request map[string]string, exclude ...string) {
	for fieldName, value := range request {
		// if it has a value and its not an excluded header, add field
		if value != "" && !slices.Contains(exclude, fieldName) {
			self.Headers = append(self.Headers, &Header{
				Field:   fieldName,
				Type:    KEY,
				Default: "",
			})
		}
	}
}

// AddDataHeader allows add multiple data fields at once (used for adding months etc)
func (self *Headers) AddDataHeader(fields ...string) {
	for _, field := range fields {
		self.Headers = append(self.Headers, &Header{
			Field:   field,
			Type:    DATA,
			Default: 0.0,
		})
	}
}

// Get returns all headers that match the request names
func (self *Headers) Get(name string) (found *Header) {
	for _, h := range self.Headers {
		if h.Field == name {
			found = h
			return
		}
	}
	return
}

// ByType returns header field names grouped as a map by the type of field.
//
// Used in returning data from the api to avoid adding default fields
func (self *Headers) ByType() (list map[string][]string) {
	list = map[string][]string{
		"labels": fieldsByType(self.Headers, KEY),
		"data":   fieldsByType(self.Headers, DATA),
		"extra":  fieldsByType(self.Headers, EXTRA),
		"end":    fieldsByType(self.Headers, END),
	}
	return
}

// Keys returns the headers that have been set as a key, ordered by field name
func (self *Headers) Keys() (keys []*Header) {
	keys = byType(self.Headers, KEY)
	sort.Slice(keys, func(i, j int) bool {
		var a = keys[i].Field
		var b = keys[j].Field
		return (a < b)
	})
	return
}

// End returns the last column - usually row total / average
func (self *Headers) End() (header *Header) {
	var ends = byType(self.Headers, END)
	if len(ends) > 0 {
		header = ends[0]
	}
	return
}

// First returns the first header field
func (self *Headers) First() (header *Header) {
	var list = byType(self.Headers, KEY)
	if len(list) > 0 {
		header = list[0]
	}
	return
}

// Data returns all the data columns preset in the set of headers
func (self *Headers) Data() (data []*Header) {
	data = byType(self.Headers, DATA)
	sort.Slice(data, func(i, j int) bool {
		var a = data[i].Field
		var b = data[j].Field
		return (a < b)
	})
	return
}

// All returns set of all headers to use in a row, in order
func (self *Headers) All() (all []*Header) {
	all = []*Header{}
	// add labels
	all = append(all, byType(self.Headers, KEY)...)
	// add data
	all = append(all, byType(self.Headers, DATA)...)
	// add extras
	all = append(all, byType(self.Headers, EXTRA)...)
	// add endings
	all = append(all, byType(self.Headers, END)...)
	return
}

func KeyHeadersFromRequest(request map[string]interface{}) (keyHeaders []*Header) {
	keyHeaders = []*Header{}

	return
}

// byType filters
func byType(headers []*Header, t Type) (list []*Header) {
	list = []*Header{}
	for _, h := range headers {
		if h.Type == t {
			list = append(list, h)
		}
	}
	return
}

// fieldsByType filters just field names
func fieldsByType(headers []*Header, t Type) (list []string) {
	list = []string{}
	for _, h := range headers {
		if h.Type == t {
			list = append(list, h.Field)
		}
	}
	return
}
