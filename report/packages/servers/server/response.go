package server

import "opg-reports/report/packages/types"

type Response struct {
	// context
	ctx types.ContextLogger
	// w is the writer to send data back
	w types.Writer
}

// Ctx returns the context logger
func (self *Response) Ctx() types.ContextLogger {
	return self.ctx
}

// SetWritter
func (self *Response) SetWriter(w types.Writer) {
	self.w = w
}

// Writer
func (self *Response) Writer() types.Writer {
	return self.w
}

// Response is a thing proxy from handler to render the result ... expect
// templates to be set on the writer already, then writer data gets reset and
// the new data attached then writen
//  - reset would clear any attached data on this and the writer
//  - ctx & templates are left as is

// Writer would be configured in main server command to have templates attached

// During server mux handler wrapper func the Response gets a Writer attached

// During the handler itself, the response will have its data reset at the start
// to avoid an old values rendering.
// At the end, the data is set and written - maybe in one?

// Data would be a map[string][]interface{} - generated from a struct -
// will need a type alias ( ResponseData ?)

//
