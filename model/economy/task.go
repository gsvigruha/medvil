package economy

import (
	"medvil/model/navigation"
	"medvil/model/time"
)

type Task interface {
	Complete(Calendar *time.CalendarType, tool bool) bool
	Field() *navigation.Field
	Blocked() bool
	Name() string
	Tag() string
	Expired(Calendar *time.CalendarType) bool
	Pause(bool)
	IsPaused() bool
	SetTraveller(traveller *navigation.Traveller)
	Motion() uint8
	IsFieldCenter() bool
}

type TaskBase struct {
	Paused      bool
	Traveller   *navigation.Traveller
	FieldCenter bool
}

func (t *TaskBase) Pause(paused bool) {
	t.Paused = paused
}

func (t *TaskBase) IsPaused() bool {
	return t.Paused
}

func (t *TaskBase) SetTraveller(traveller *navigation.Traveller) {
	t.Traveller = traveller
}

func (t *TaskBase) IsFieldCenter() bool {
	return t.FieldCenter
}
