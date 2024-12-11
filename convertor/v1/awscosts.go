package v1

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/ministryofjustice/opg-reports/internal/dateformats"
	"github.com/ministryofjustice/opg-reports/internal/dateutils"
	"github.com/ministryofjustice/opg-reports/models"
)

type AwsCost struct {
	ID           int    `json:"id"`
	Ts           string `json:"ts"`
	Organisation string `json:"organisation"`
	AccountID    string `json:"account_id"`
	AccountName  string `json:"account_name"`
	Unit         string `json:"unit"`
	Label        string `json:"label"`
	Environment  string `json:"environment"`
	Service      string `json:"service"`
	Region       string `json:"region"`
	Date         string `json:"date"`
	Cost         string `json:"cost"`
}

// MarshalJSON converts from current version to model version so
// when writing this struct to a json file it will take the form
// of a new model
func (self *AwsCost) MarshalJSON() (bytes []byte, err error) {
	var (
		unit    *models.Unit
		account *models.AwsAccount
		cost    *models.AwsCost
		ts      string
		now     string = time.Now().UTC().Format(dateformats.Full)
	)
	// costs have a different time format, so convert between
	if self.Ts == "" {
		self.Ts = now
	} else {
		self.Ts = dateutils.Convert(self.Ts, dateformats.Old, dateformats.Full)
	}
	ts = self.Ts
	unit = &models.Unit{
		Ts:   ts,
		Name: strings.ToLower(self.Unit),
	}
	account = &models.AwsAccount{
		Ts:          ts,
		Number:      self.AccountID,
		Name:        strings.ToLower(self.AccountName),
		Label:       strings.ToLower(self.Label),
		Environment: strings.ToLower(self.Environment),
		Unit:        (*models.UnitForeignKey)(unit),
	}
	cost = &models.AwsCost{
		Ts:         ts,
		Region:     self.Region,
		Service:    self.Service,
		Date:       self.Date,
		Cost:       self.Cost,
		AwsAccount: (*models.AwsAccountForeignKey)(account),
		Unit:       (*models.UnitForeignKey)(unit),
	}
	bytes, err = json.MarshalIndent(cost, "", "  ")

	return
}
