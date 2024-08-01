package server

import (
	th "opg-reports/internal/testhelpers"
	"opg-reports/services/front/cnf"
	"opg-reports/shared/logger"
	"testing"
)

func TestFrontServerRegister(t *testing.T) {
	logger.LogSetup()
	mux := th.Mux()
	conf, _ := cnf.Load([]byte(testServerCfg))
	s := New(conf, nil, "", "")
	s.Register(mux)

	home := s.Nav.Get("/")
	if home == nil || home.Registered != true {
		t.Errorf("home page not registered")
	}
}
