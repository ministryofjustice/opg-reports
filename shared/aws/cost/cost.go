package cost

import (
	"opg-reports/shared/data"

	"github.com/google/uuid"
)

// Cost struct captures data for a cost explorer item in a flat structure that
// is easier to search and store
//
// UUID is generates and used as a unique index for usage in the a data store.
//
// Impliments data.IEntry
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

// UID is the unique id (UUID) for this Cost item
func (i *Cost) UID() string {
	return i.UUID
}

// Valid returns true only if all fields are present with a non-empty
// value
// Uses json marshal to a map to check this
func (i *Cost) Valid() (valid bool) {
	mapped, _ := data.ToMap(i)

	for _, val := range mapped {
		if val == "" {
			return false
		}
	}
	return true
}

// New creates a Cost with the uid passed or
// c=a new uuid if that is nil
func New(uid *string) *Cost {
	c := &Cost{}

	if uid != nil {
		c.UUID = *uid
	} else {
		c.UUID = uuid.NewString()
	}
	return c
}
