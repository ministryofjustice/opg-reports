package seeder

type Seed struct {
	Table  string
	Label  string
	DB     string
	Schema string
	Source []string
	Dummy  []string
}
