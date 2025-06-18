package interfaces

type Service interface {
	Seed() (err error)
}
