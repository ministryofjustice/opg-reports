package types

// Requestable exposes methods to attach the
// modfied request.
//
// Used by handlers
type Requestable interface {
	SetRequest(r ServerRequester)
	Request() ServerRequester
}
