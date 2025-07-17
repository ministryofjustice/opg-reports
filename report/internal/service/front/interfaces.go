package front

type Response interface{}
type Result interface{}
type Closer interface {
	Close() (err error)
}
