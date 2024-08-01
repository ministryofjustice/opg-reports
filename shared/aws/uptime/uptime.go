package uptime

import (
	"log/slog"
	"opg-reports/shared/data"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/google/uuid"
)

// Uptime struct captures data for an aws uptime item
//
// UUID is generates and used as a unique index for usage in the a data store.
//
// Impliments data.IEntry
type Uptime struct {
	UUID      string    `json:"uuid"`
	Timestamp time.Time `json:"time"`

	Average  float64   `json:"average"`
	Unit     string    `json:"unit"`
	DateTime time.Time `json:"date_time"`

	AccountUnit  string `json:"account_unit"`
	AccountLabel string `json:"account_label"`
}

func (i *Uptime) TS() time.Time {
	return i.Timestamp
}

// UID is the unique id (UUID) for this Cost item
func (i *Uptime) UID() string {
	slog.Debug("[aws/uptime] UID()", slog.String("UID", i.UUID))
	return i.UUID
}

// Valid returns true only if all fields are present with a non-empty
// value
// Uses json marshal to a map to check this
func (i *Uptime) Valid() (valid bool) {

	mapped, _ := data.ToMap(i)

	for k, val := range mapped {
		if val == "" {
			slog.Debug("[aws/uptime] invalid",
				slog.String("UID", i.UID()),
				slog.String(k, val.(string)))
			return false
		}
	}
	slog.Debug("[aws/uptime] valid", slog.String("UID", i.UID()))
	return true
}

var _ data.IEntry = &Uptime{}

// New creates an Uptime with the uid passed or
// c=a new uuid if that is nil
func New(uid *string) *Uptime {
	c := &Uptime{Timestamp: time.Now().UTC()}
	if uid != nil {
		c.UUID = *uid
	} else {
		c.UUID = uuid.NewString()
	}
	slog.Debug("[aws/uptime] new", slog.String("UID", c.UUID))
	return c
}

func NewFromDatapoint(uid *string, dp *cloudwatch.Datapoint) *Uptime {
	up := New(uid)
	up.Average = *dp.Average
	up.Unit = *dp.Unit
	up.DateTime = *dp.Timestamp
	return up
}
