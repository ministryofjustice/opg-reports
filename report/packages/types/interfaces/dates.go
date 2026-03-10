package interfaces

import "time"

// DateTypes used in date helpers to limit types
type DateTypes interface {
	time.Time | string
}
