package request

import (
	"net/http"

	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/request/get"
)

type Request struct {
	Parameters map[string]*get.GetParameter
	Order      []string
}

func (req *Request) Param(name string, r *http.Request) (value string) {
	if p, ok := req.Parameters[name]; ok {
		value = p.Value(r)
	}
	return
}

func New(params ...*get.GetParameter) (request *Request) {
	request = &Request{Parameters: map[string]*get.GetParameter{}, Order: []string{}}

	for _, param := range params {
		request.Parameters[param.Name] = param
		request.Order = append(request.Order, param.Name)
	}
	return
}
