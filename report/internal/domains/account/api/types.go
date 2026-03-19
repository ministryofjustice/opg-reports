package api

// Result is wrapper for the api data results
type Result struct {
	Accounts []*Account `json:"accounts"` // all accounts
}

// Account is used in the api result and contain
// the result of the database select query
type Account struct {
	ID          string `json:"id"`
	Name        string `json:"name" `
	Label       string `json:"label" `
	Environment string `json:"environment" `
	TeamName    string `json:"team"`
}

// Sequence is used to return the columns in the order they are selected.
func (self *Account) Sequence() []any {
	return []any{
		&self.ID,
		&self.Name,
		&self.Label,
		&self.Environment,
		&self.TeamName,
	}
}
