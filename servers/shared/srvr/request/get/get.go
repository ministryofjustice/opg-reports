// Package get handles parasing and fetching get parameters from query string
package get

import (
	"net/http"
	"slices"
)

type GetParameter struct {
	value   string
	Default string
	Allowed []string
	Name    string
}

func (g *GetParameter) Value(r *http.Request) string {
	var (
		value       = g.value
		values      []string
		allowed     []string
		queryString = r.URL.Query()
	)
	//  update the value
	if found, ok := queryString[g.Name]; ok {
		values = found
	}
	// if theres no allow list then use all values directly
	// otherwise filter the values found by the ones allowed
	if len(g.Allowed) == 0 {
		allowed = values
	} else {
		for _, val := range values {
			if slices.Contains(g.Allowed, val) {
				allowed = append(allowed, val)
			}
		}
	}

	if len(allowed) > 0 {
		value = allowed[0]
	}

	// if default is set, but we have no current value, use it
	if value == "" && g.Default != "" {
		value = g.Default
	}

	g.value = value

	return g.value
}

func New(name string, defaultV string) *GetParameter {
	return &GetParameter{Default: defaultV, Name: name, value: ""}
}

// WithChoices limits the allowed values of the get parameter
// Note: sets the first item from allowed as the default value
func WithChoices(name string, allowed []string) (param *GetParameter) {
	param = New(name, allowed[0])
	param.Allowed = allowed
	return
}
