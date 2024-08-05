package endpoint

import "net/http"

type Queryable struct {
	Allowed []string
	Found   map[string][]string
}

func (q *Queryable) Parse(r *http.Request) map[string][]string {
	q.Found = map[string][]string{}

	queryStr := r.URL.Query()
	for _, field := range q.Allowed {
		if val, ok := queryStr[field]; ok {
			if _, set := q.Found[field]; !set {
				q.Found[field] = []string{}
			}
			q.Found[field] = val
		} else if r.PathValue(field) != "" {
			q.Found[field] = []string{r.PathValue(field)}
		}
	}
	return q.Found
}

func NewQueryable(allow []string) *Queryable {
	return &Queryable{
		Allowed: allow,
	}
}
