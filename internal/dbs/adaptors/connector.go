package adaptors

import (
	"fmt"
)

// Connection is a generic connection struct used for
// creating a new database connection with
// sqlx.ConnectContext
//
// Implements dbs.Connector interface
type Connection struct {
	Driver     string
	Path       string
	Parameters string
}

// GetConnectionString returns the full connection string used to connect to database
func (self *Connection) String() string {
	return fmt.Sprintf("%s%s", self.Path, self.Parameters)
}

// GetDriverName returns the driver name to sued for this connection
func (self *Connection) DriverName() string {
	return self.Driver
}
