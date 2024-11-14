package adaptors

import (
	"testing"

	"github.com/ministryofjustice/opg-reports/internal/dbs"
)

var _ dbs.Moder = &ReadOnly{}
var _ dbs.Moder = &ReadWrite{}

func TestAdaptorsModes(t *testing.T) {
	var r = &ReadOnly{}
	var rw = &ReadWrite{}

	if !r.Read() {
		t.Errorf("mode should be false when not set")
	}
	if r.Write() {
		t.Errorf("write should be false")
	}

	if !rw.Read() {
		t.Errorf("mode should be false when not set")
	}
	if !rw.Write() {
		t.Errorf("write should be true")
	}
}
