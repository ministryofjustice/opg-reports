package awsaccount

import "github.com/ministryofjustice/opg-reports/report/internal/interfaces"

// check interface applies
var _ interfaces.ImporterExistingCommand = Existing
