package types

import "database/sql"

// Semver returns the semantic version number
type Semver interface {
	// Semver returns the semantic version number
	Semver() string
}

// Hash returns the git details
type Hasher interface {
	// Hash returns the git commit sha
	Hash() string
}

// Versioner exposes methods to get version details
type Versioner interface {
	Semver
	Hasher
}

// DBDetailer exposes methods for database
// information
type DBPather interface {
	SetPath(path string)
}

// DBConnector exposes method to connect to the database
type DBConnector interface {
	Connection() (*sql.DB, error)
}

// DBer is the main DB interface
type DBer interface {
	DBPather
	DBConnector
}

// Hoster is used for fron / api host details
type Hoster interface {
	Label() string
	Host() string
}

type ServerConfigure interface {
	Database() DBer
	Versions() Versioner
	Host() Hoster
}
