package v1

import (
	"fmt"
	"strings"

	"github.com/ministryofjustice/opg-reports/models"
)

func AccountFromUnit(unit *models.Unit, ts string) (account *models.AwsAccount) {
	account = &models.AwsAccount{
		Environment: "production",
		Ts:          ts,
		Name:        strings.ToLower(fmt.Sprintf("%s production", unit.Name)),
		Label:       strings.ToLower(unit.Name),
	}

	switch unit.Name {
	case "digideps":
		account.Number = "515688267891"
	case "make":
		account.Number = "980242665824"
	case "modernise":
		account.Number = "313879017102"
	case "serve":
		account.Number = "933639921819"
	case "sirius":
		account.Number = "649098267436"
	case "use":
		account.Number = "690083044361"
	}

	return
}
