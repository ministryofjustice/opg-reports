package repository

const SCHEMA string = `
CREATE TABLE IF NOT EXISTS owner (
	id INTEGER PRIMARY KEY,
	created_at TEXT,
	name TEXT NOT NULL UNIQUE
) STRICT;

CREATE INDEX IF NOT EXISTS owner_name ON owner(name)
`
