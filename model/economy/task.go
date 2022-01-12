package economy

import (
	"medvil/model/navigation"
)

type Location struct {
	X uint16
	Y uint16
	F navigation.IField
}

type Task interface {
	Tick()
	Location() Location
}
