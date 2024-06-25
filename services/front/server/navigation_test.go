package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"opg-reports/services/front/cnf"
	"testing"
)

func TestFrontServerNavigationTop(t *testing.T) {
	conf, _ := cnf.Load([]byte(testCfg))
	s := New(conf, nil, "", "")

	r := httptest.NewRequest(http.MethodGet, "/s1/2/", nil)
	active, all := s.Nav.Top(r)

	if len(all) != 3 {
		t.Errorf("failed to get top nav items")
		fmt.Println(all)
	}

	if active == nil {
		t.Errorf("failed to find active top nav")
	} else if active.Href != "/s1/" {
		t.Errorf("found incorrect top nav")
	}

	r = httptest.NewRequest(http.MethodGet, "/random/123/", nil)
	active, all = s.Nav.Top(r)
	if len(all) != 3 {
		t.Errorf("failed to get top nav items")
	}
	if active != nil {
		t.Errorf("top should not have been found")
	}

}

func TestFrontServerNavigationSide(t *testing.T) {
	conf, _ := cnf.Load([]byte(testCfg))
	s := New(conf, nil, "", "")

	r := httptest.NewRequest(http.MethodGet, "/s2/2/1/", nil)
	activeT, _ := s.Nav.Top(r)
	active, side := s.Nav.Side(r, activeT)

	if len(side) != 2 {
		t.Errorf("failed to get side nav items")
	}

	if active == nil {
		t.Errorf("should have found active side item")
	} else if active.Href != "/s2/2/" {
		t.Errorf("found incorrect side item")
	}

}

func TestFrontServerNavigationActive(t *testing.T) {
	conf, _ := cnf.Load([]byte(testCfg))
	s := New(conf, nil, "", "")

	r := httptest.NewRequest(http.MethodGet, "/s2/2/1/", nil)
	active := s.Nav.Active(r)

	if active == nil {
		t.Errorf("failed to get active item")
	} else if active.Href != "/s2/2/1/" || active.Name != "S2.2.1" {
		t.Errorf("failed to get correct active item")
	}

	r = httptest.NewRequest(http.MethodGet, "/s5/2/1/", nil)
	active = s.Nav.Active(r)
	if active != nil {
		t.Errorf("found an item when should not have")
	}
}
