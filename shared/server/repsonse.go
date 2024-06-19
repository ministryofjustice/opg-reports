package server

type IResponse interface {
	Body() []byte
	Status() int
	Errors() []error
}

type IApiResponse interface {
	Start()
	End()
	Set(results interface{}, status int, errors []error)
	Results() interface{}
	Status() int
	Errors() []error
	Body() []byte
}
