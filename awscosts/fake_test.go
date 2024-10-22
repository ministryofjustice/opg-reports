package awscosts_test

import (
	"testing"

	"github.com/ministryofjustice/opg-reports/awscosts"
)

func TestAwsCostsFake(t *testing.T) {

	start := &awscosts.Cost{Organisation: "foobar"}

	faked := awscosts.Fake(start)
	if faked.Organisation != start.Organisation {
		t.Errorf("organisation was not copied over")
	}

	fakes := awscosts.Fakes(10, start)
	randFields := map[string]int{}
	for _, f := range fakes {
		if f.Organisation != start.Organisation {
			t.Errorf("organisation was not copied over")
		}
		if _, ok := randFields[f.Unit]; !ok {
			randFields[f.Unit] = 1
		}
		if _, ok := randFields[f.AccountID]; !ok {
			randFields[f.AccountID] = 1
		}
	}

	for _, count := range randFields {
		if count > 1 {
			t.Errorf("field was not very random")
		}
	}

}
