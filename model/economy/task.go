package economy

import (
	"medvil/model/navigation"
	"medvil/model/time"
)

type Task interface {
	Complete(Calendar *time.CalendarType) bool
	Location() navigation.Location
	Blocked() bool
	Name() string
}
