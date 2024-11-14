package adaptors

// Seed captures info on if the database
// can be used for seeding or not
type Seed struct {
	seedable bool
}

// Seedable returns bool to determin if the database is
// in a position to be used for seeding data
func (self *Seed) Seedable() bool {
	return self.seedable
}

// Seeded sets the database table as having been seeded
// and no longer suitable for seeding
func (self *Seed) Seeded() {
	self.seedable = false
}
