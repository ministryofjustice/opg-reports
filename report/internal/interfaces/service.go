package interfaces

type Service interface {
	// Seed is used to populate the service with valid test data
	Seed() (err error)

	// Import is used to populate the database with existing data - either
	// from previous versions of the service or by generating accurate new
	// information
	Import(repository Repository) (err error)
}
