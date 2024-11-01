package convert

// IntToBool helper used with sql conversion as sqlite has no
// boolean type, they are stored as 1 (true) or 0, this maps them back to
// a bool
func IntToBool(i uint8) bool {
	if i == 1 {
		return true
	}
	return false
}
