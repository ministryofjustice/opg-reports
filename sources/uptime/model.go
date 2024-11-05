package uptime

import "fmt"

type Uptime struct {
	ID      int     `json:"id,omitempty" db:"id" faker:"unique, boundary_start=1, boundary_end=2000000" doc:"Database primary key."`               // ID is a generated primary key
	Ts      string  `json:"ts,omitempty" db:"ts" faker:"time_string" doc:"Time the record was created."`                                           // TS is timestamp when the record was created
	Unit    string  `json:"unit,omitempty" db:"unit" faker:"oneof: unitA, unitB, unitC" doc:"The name of the unit / team that owns this account."` // Unit is the team that owns this account, passed directly
	Date    string  `json:"date,omitempty" db:"date" faker:"date_string" doc:"Date this cost was generated."`                                      // The interval date for when this uptime was logged
	Average float64 `json:"average,omitempty" db:"average" doc:"Percentage uptime average for this record period."`
}

// UID returns a unique uid
// Record interface
func (self *Uptime) UID() string {
	return fmt.Sprintf("%s-%d", "uptime", self.ID)
}

// TDate
// Transformable interface
func (self *Uptime) TDate() string {
	return self.Date
}

// TValue
// Transformable interface
func (self *Uptime) TValue() string {
	return fmt.Sprintf("%f", self.Average)
}
