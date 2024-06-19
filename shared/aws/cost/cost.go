package cost

import (
	"opg-reports/shared/data"

	"github.com/google/uuid"
)

type Cost struct {
	UUID                string `json:"uuid"`
	AccountOrganisation string `json:"account_organsiation"`
	AccountId           string `json:"account_id"`
	AccountEnvironment  string `json:"account_environment"`
	AccountName         string `json:"account_name"`
	AccountUnit         string `json:"account_unit"`
	AccountLabel        string `json:"account_label"`
	Service             string `json:"service"`
	Region              string `json:"region"`
	Date                string `json:"date"`
	Cost                string `json:"cost"`
}

func (i *Cost) UID() string {
	return i.UUID
}

func (i *Cost) Valid() (valid bool) {
	mapped, _ := data.ToMap(i)

	for _, val := range mapped {
		if val == "" {
			return false
		}
	}
	return true
}

func New(uid *string) *Cost {
	c := &Cost{}

	if uid != nil {
		c.UUID = *uid
	} else {
		c.UUID = uuid.NewString()
	}
	return c
}
