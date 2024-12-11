package v1_test

import (
	"testing"

	v1 "github.com/ministryofjustice/opg-reports/convertor/v1"
	"github.com/ministryofjustice/opg-reports/internal/structs"
	"github.com/ministryofjustice/opg-reports/models"
)

// TestV1AwsUptimeMarshal small test to make sure the
// cost is mapped
func TestV1AwsUptimeMarshal(t *testing.T) {
	og := v1.AwsUptime{
		Unit:    "sirius",
		Average: 99.5,
		Date:    "2024-01-01",
	}

	b, _ := og.MarshalJSON()
	model := &models.AwsUptime{}
	structs.Unmarshal(b, model)

	if og.Date != model.Date {
		t.Errorf("error mapping fields - expected [%s] actual [%s]", og.Date, model.Date)
	}

	if og.Unit != model.Unit.Name {
		t.Errorf("error mapping fields - expected [%s] actual [%s]", og.Unit, model.Unit.Name)
	}
	if og.Unit != model.AwsAccount.Label {
		t.Errorf("error mapping fields - expected [%s] actual [%s]", og.Unit, model.AwsAccount.Label)
	}

}
