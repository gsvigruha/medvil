package economy

import (
	"medvil/model/navigation"
	"medvil/model/time"
)

type Task interface {
	Complete(Calendar *time.CalendarType) bool
	Field() *navigation.Field
	Blocked() bool
	Name() string
	Tag() string
}
