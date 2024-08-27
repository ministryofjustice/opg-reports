package urls_test

import (
	"testing"

	"github.com/ministryofjustice/opg-reports/servers/shared/urls"
)

func TestSharedUrlsParse(t *testing.T) {

	if p := urls.Parse("http://", "localhost:80", "/test/a-b?path=1"); p.String() != "http://localhost:80/test/a-b/?path=1" {
		t.Errorf("faield to parse url [%s]", p.String())
	}

	if p := urls.Parse("http://", ":80", "/test/a-b?path=1"); p.String() != "http://localhost:80/test/a-b/?path=1" {
		t.Errorf("faield to parse url [%s]", p.String())
	}

	if p := urls.Parse("", "localhost:80", "/test/a-b?path=1"); p.String() != "http://localhost:80/test/a-b/?path=1" {
		t.Errorf("faield to parse url [%s]", p.String())
	}

	if p := urls.Parse("http://", "", "/test/a-b?path=1"); p.String() != "http://localhost/test/a-b/?path=1" {
		t.Errorf("faield to parse url [%s]", p.String())
	}

	if p := urls.Parse("http://", "", "localhost:80/test/a-b?path=1"); p.String() != "http://localhost:80/test/a-b/?path=1" {
		t.Errorf("faield to parse url [%s]", p.String())
	}

	if p := urls.Parse("", "", "http://localhost:80/test/a-b?path=1"); p.String() != "http://localhost:80/test/a-b/?path=1" {
		t.Errorf("faield to parse url [%s]", p.String())
	}
}
