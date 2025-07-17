package sqlr

// BoundStatement is used to handle sql statments that use named parameters
type BoundStatement struct {
	Statement string
	Data      interface{}
	Returned  interface{}
}
