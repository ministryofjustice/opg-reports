package front

import (
	"fmt"
	"net/http"
	"opg-reports/report/internal/repository/restr"
)

// GetFromAPI fetches the data from the endpoint and converts it to T, and the result of this is then passed into each
// postProcessor in turn.
//
// The postProcessor allows reformatting / data changes to be handled; so the raw API result can be re-strucutured
// into a suitable format for the front end display (list -> tablular etc)
//
//   - endpoint should have any substituions replaced and query strings added, the api host is added in this function via config values
//   - postProcessers should handle nil checks on the response result passed in
func (self *Service[T]) GetFromAPI(client restr.RepositoryRestGetter, endpoint string, postProcessors ...func(result T) (err error)) (result T, err error) {
	var (
		code  int
		uri   = fmt.Sprintf("%s/%s", self.conf.Servers.Api.Addr, endpoint)
		httpc = http.Client{Timeout: self.conf.Servers.Front.Timeout}
		log   = self.log.With("operation", "GetFromAPI", "uri", uri)
	)

	log.Debug("calling api ... ")
	// fetch the result from the api, and convert to T
	code, err = client.Get(httpc, uri, &result)
	if err != nil {
		log.Error("error calling endpoint", "err", err.Error())
		return
	}
	if code != http.StatusOK {
		err = fmt.Errorf("status code error, expected [%v] actual [%v]", http.StatusOK, code)
		return
	}
	// run all post processors against the result
	for _, pf := range postProcessors {
		log.Debug("calling post processor ...")
		err = pf(result)
	}

	return
}
