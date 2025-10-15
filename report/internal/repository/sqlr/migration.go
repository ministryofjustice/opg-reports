package sqlr

// migrated tracks any migrations made upwards so they only run once
var migrated_up map[string]bool = map[string]bool{}

// MigrateUp goes up the DB_MIGRATIONS_UP commands in order to create
// the database as required and can be run at any time
func MigrateUp(r RepositoryWriter) (err error) {
	var key string

	key, err = r.ID()
	if err != nil {
		return
	}
	// check if its already been migrated
	if _, ok := migrated_up[key]; ok {
		return
	}

	for _, stmt := range DB_MIGRATIONS_UP {
		_, err = r.Exec(stmt)
		if err != nil {
			return
		}
	}

	migrated_up[key] = true

	return
}
