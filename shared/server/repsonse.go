package server

type IResponse interface {
	Body() []byte
	Status() int
	Errors() []error
}

type IApiResponse interface {
	Start()
	End()

	SetResults(results interface{})
	GetResults() interface{}

	SetStatus(status int)
	GetStatus() int

	SetErrors(errors []error)
	AddError(err error)
	AddStatusError(status int, err error)
	GetErrors() []error

	Body() []byte
}
