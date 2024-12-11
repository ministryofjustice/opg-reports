package v1

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/ministryofjustice/opg-reports/internal/dateformats"
	"github.com/ministryofjustice/opg-reports/internal/dateutils"
	"github.com/ministryofjustice/opg-reports/models"
)

type AwsUptime struct {
	ID      int     `json:"id"`
	Ts      string  `json:"ts"`
	Unit    string  `json:"unit"`
	Average float64 `json:"average"`
	Date    string  `json:"date"`
}

// MarshalJSON converts from current version to model version so
// when writing this struct to a json file it will take the form
// of a new model
func (self *AwsUptime) MarshalJSON() (bytes []byte, err error) {
	var (
		account *models.AwsAccount
		uptime  *models.AwsUptime
		unit    *models.Unit
		ts      string
		now     string = time.Now().UTC().Format(dateformats.Full)
	)
	if self.Ts == "" {
		self.Ts = now
	}
	ts = dateutils.Reformat(self.Ts, dateformats.Full)

	unit = &models.Unit{
		Ts:   ts,
		Name: strings.ToLower(self.Unit),
	}
	account = AccountFromUnit(unit, self.Ts)
	uptime = &models.AwsUptime{
		Ts:         ts,
		Date:       self.Date,
		Average:    self.Average,
		Unit:       (*models.UnitForeignKey)(unit),
		AwsAccount: (*models.AwsAccountForeignKey)(account),
	}
	bytes, err = json.MarshalIndent(uptime, "", "  ")

	return
}
