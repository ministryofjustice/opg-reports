package cost

import (
	"log/slog"
	"opg-reports/shared/data"
	"strconv"

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
	slog.Debug("[aws/cost] UID()", slog.String("UID", i.UUID))
	return i.UUID
}

// Valid returns true only if all fields are present with a non-empty
// value
// Uses json marshal to a map to check this
func (i *Cost) Valid() (valid bool) {

	mapped, _ := data.ToMap(i)

	for k, val := range mapped {
		if val == "" {
			slog.Debug("[aws/cost] invalid", slog.String("UID", i.UID()), slog.String(k, val))
			return false
		}
	}
	slog.Debug("[aws/cost] valid", slog.String("UID", i.UID()))
	return true
}

func Total(items []*Cost) (total float64) {
	total = 0.0
	for _, c := range items {
		if v, err := strconv.ParseFloat(c.Cost, 10); err == nil {
			total += v
		}
	}
	return
}

var _ data.IEntry = &Cost{}

// New creates a Cost with the uid passed or
// c=a new uuid if that is nil
func New(uid *string) *Cost {
	c := &Cost{}

	if uid != nil {
		c.UUID = *uid
	} else {
		c.UUID = uuid.NewString()
	}
	slog.Debug("[aws/cost] new", slog.String("UID", c.UUID))
	return c
}
