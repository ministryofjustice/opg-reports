package sqlr

// MigrateUp goes up the DB_MIGRATIONS commands in order
func MigrateUp(r *Repository) (err error) {
	for _, stmt := range DB_MIGRATIONS_UP {
		_, err = r.Exec(stmt)
		if err != nil {
			return
		}
	}
	return
}
