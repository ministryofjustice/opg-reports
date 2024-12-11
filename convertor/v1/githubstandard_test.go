package v1_test

import (
	"testing"

	v1 "github.com/ministryofjustice/opg-reports/convertor/v1"
	"github.com/ministryofjustice/opg-reports/internal/structs"
	"github.com/ministryofjustice/opg-reports/models"
)

// TestV1GithubStandardMarshal small test to make sure the
// cost is mapped
func TestV1GithubStandardMarshal(t *testing.T) {
	og := v1.GithubStandard{
		Name:       "test",
		IsArchived: 1,
		HasLicense: 1,
		Teams:      "opg-webops#OPG#sirius#organisation-security-auditor#",
	}

	b, _ := og.MarshalJSON()
	model := &models.GitHubRepositoryStandard{}
	structs.Unmarshal(b, model)
	// pretty.Print(model)
}
