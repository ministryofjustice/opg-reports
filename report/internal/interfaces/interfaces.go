package interfaces

type Model interface{}
type Repository interface{}

type Service interface {
	Close() (err error)
}
