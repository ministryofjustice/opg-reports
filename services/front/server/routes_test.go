package server

import (
	"opg-reports/services/front/cnf"
	"testing"
)

func TestFrontServerRegister(t *testing.T) {

	mux := testMux()
	conf, _ := cnf.Load([]byte(testCfg))
	s := New(conf, nil, "", "")
	s.Register(mux)

	home := s.Nav.Get("/")
	if home == nil || home.Registered != true {
		t.Errorf("home page not registered")
	}
}
