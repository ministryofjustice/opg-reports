package v1_test

import (
	"testing"

	"github.com/ministryofjustice/opg-reports/convertor/lib"
	v1 "github.com/ministryofjustice/opg-reports/convertor/v1"
	"github.com/ministryofjustice/opg-reports/internal/structs"
	"github.com/ministryofjustice/opg-reports/models"
)

var _ lib.V1er = &v1.AwsCost{}

// TestV1AwsCostMarshal small test to make sure the
// cost is mapped
func TestV1AwsCostMarshal(t *testing.T) {
	og := v1.AwsCost{
		Organisation: "opg",
		AccountID:    "01",
		AccountName:  "test",
		Unit:         "sirius",
		Label:        "test-prod",
		Environment:  "production",
		Service:      "ECS",
		Region:       "eu-west-1",
		Date:         "2024-01-01",
		Cost:         "0.15",
	}

	b, _ := og.MarshalJSON()
	model := &models.AwsCost{}
	structs.Unmarshal(b, model)

	if og.Service != model.Service {
		t.Errorf("error mapping fields - expected [%s] actual [%s]", og.Service, model.Service)
	}
	if og.Date != model.Date {
		t.Errorf("error mapping fields - expected [%s] actual [%s]", og.Date, model.Date)
	}
	if og.Region != model.Region {
		t.Errorf("error mapping fields - expected [%s] actual [%s]", og.Region, model.Region)
	}
	if og.Cost != model.Cost {
		t.Errorf("error mapping fields - expected [%s] actual [%s]", og.Cost, model.Cost)
	}
	if og.AccountID != model.AwsAccount.Number {
		t.Errorf("error mapping fields - expected [%s] actual [%s]", og.AccountID, model.AwsAccount.Number)
	}
	if og.AccountName != model.AwsAccount.Name {
		t.Errorf("error mapping fields - expected [%s] actual [%s]", og.AccountName, model.AwsAccount.Name)
	}
	if og.Label != model.AwsAccount.Label {
		t.Errorf("error mapping fields - expected [%s] actual [%s]", og.Label, model.AwsAccount.Label)
	}
	if og.Environment != model.AwsAccount.Environment {
		t.Errorf("error mapping fields - expected [%s] actual [%s]", og.Environment, model.AwsAccount.Environment)
	}
	if og.Unit != model.Unit.Name {
		t.Errorf("error mapping fields - expected [%s] actual [%s]", og.Unit, model.Unit.Name)
	}

}
