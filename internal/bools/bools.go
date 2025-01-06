package bools

// Int returns a unit8 representation of a boolean:
//   - True = 1
//   - False = 0
//
// Used as sqlite doesnt have a boolean type by default
func Int(b bool) (i uint8) {
	i = 0
	if b {
		i = 1
	}
	return
}
