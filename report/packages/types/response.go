package types

import (
	"net/http"
	"text/template"
)

type Templater interface {
	// returns the template name
	Template() string
	// list of template files
	Files() []string
	// the functions to use in the template generation
	Funcs() template.FuncMap
}

// Writer is the interface to send data back
// in the response - should be json / html variant
type Writer interface {
	http.ResponseWriter
	// uses context
	Contexter
	// SetData attaches the values to itself before sending
	SetData(data any)
	// SetTemplates configures template data on the instance
	//  - generally used on html only
	SetTemplates(cfg Templater)
	// Respond is the main method to use instead of write and
	// applies extra parsing, data marshaling etc
	Respond()
}

// Writeable use to attach response writer
// to the struct which will be called at the
// end to send data back
type Writerable interface {
	SetWriter(w Writer)
	Writer() Writer
}

type Responder interface {
	// uses context
	Contexter
	// Needs to be reset
	Resetable
	// has a writer to send data back
	Writerable
}
