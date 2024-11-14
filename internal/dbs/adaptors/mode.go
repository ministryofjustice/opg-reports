package adaptors

// ReadOnly
// Implements dbs.Moder
type ReadOnly struct{}

func (self *ReadOnly) Read() bool {
	return true
}

func (self *ReadOnly) Write() bool {
	return false
}

// ReadWrite
// Implements dbs.Moder
type ReadWrite struct{}

func (self *ReadWrite) Read() bool {
	return true
}

func (self *ReadWrite) Write() bool {
	return true
}
