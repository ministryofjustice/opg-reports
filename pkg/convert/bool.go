package convert

import "strconv"

// IntToBool helper used with sql conversion as sqlite has no
// boolean type, they are stored as 1 (true) or 0
func BoolToInt(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

// BoolStringToInt helper to deal with get param bools
// that convert over to 1 | 0 for the db
func BoolStringToInt(s string) int {
	b, err := strconv.ParseBool(s)
	if err == nil && b {
		return 1
	}
	return 0
}
