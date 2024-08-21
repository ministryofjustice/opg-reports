package query

import (
	"net/http"
)

type Source string

const (
	GET  Source = "GET"
	PATH Source = "PATH"
)

type Query struct {
	Name string
	From Source
}

func (q *Query) Values(r *http.Request) (vals []string) {

	vals = []string{}
	switch q.From {
	case PATH:
		vals = fromQueryPath(q.Name, r)
	case GET:
		vals = fromGetParameters(q.Name, r)
	}
	return
}

func fromQueryPath(name string, r *http.Request) (v []string) {
	v = []string{}
	if val := r.PathValue(name); val != "" {
		v = append(v, val)
	}
	return
}

func fromGetParameters(name string, r *http.Request) (v []string) {
	queryStr := r.URL.Query()
	if val, ok := queryStr[name]; ok {
		v = val
	}
	return
}

func First(strs []string) string {
	if len(strs) > 0 {
		return strs[0]
	}
	return ""
}

func AllD(strs []string, defaultValue string) (v []string) {
	v = []string{}
	if len(strs) > 0 {
		v = strs
	} else {
		v = append(v, defaultValue)
	}
	return
}

func FirstD(strs []string, defaultValue string) string {
	if len(strs) > 0 && strs[0] != "" {
		return strs[0]
	}
	return defaultValue
}

func Get(name string) *Query {
	return &Query{
		Name: name,
		From: GET,
	}
}

func Path(name string) *Query {
	return &Query{
		Name: name,
		From: PATH,
	}
}
