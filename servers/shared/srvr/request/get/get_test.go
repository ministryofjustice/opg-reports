package get_test

import (
	"net/http"
	"testing"

	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/request/get"
	"github.com/ministryofjustice/opg-reports/shared/testhelpers"
)

func TestSharedSrvrRequestGet(t *testing.T) {
	var (
		route string
		r     *http.Request
		g     *get.GetParameter
	)

	// -- test some simple ones
	route = "/home/2024-02/test/?archive=foo&team=true"
	_, r = testhelpers.WRGet(route)

	g = get.New("team", "bar")
	if g.Value(r) != "true" {
		t.Errorf("get param failed")
	}

	g = get.New("foo", "bar")
	if g.Value(r) != g.Default {
		t.Errorf("should have used default [%s]", g.Value(r))
	}

	// -- test multiple values for the same param, so first is used
	route = "/home/2024-02/test/?archive=foo&team=1&team=2"
	_, r = testhelpers.WRGet(route)
	g = get.New("team", "bar")
	if g.Value(r) != "1" {
		t.Errorf("get param failed")
	}

	// -- test param with limited options
	route = "/home/2024-02/test/?archive=foo&archive=false&archive=true"
	_, r = testhelpers.WRGet(route)
	g = get.WithChoices("archive", []string{"true", "false"})
	if g.Value(r) != "false" {
		t.Errorf("failed to get correct param: [%s]", g.Value(r))
	}

}
