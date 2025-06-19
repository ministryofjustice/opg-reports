package interfaces

type Service interface {
	Import() (err error)
}
