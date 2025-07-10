package restr

import "net/http"

type Model interface{}

// RepositoryRestGetter represents the Get (as in http method get) option to
// call a URI and return its content converted into T
type RepositoryRestGetter interface {
	Get(client http.Client, uri string, result interface{}) (statuscode int, err error)
}
