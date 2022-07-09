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
	SetUp(traveller *navigation.Traveller, household Household)
	Motion() uint8
	IsFieldCenter() bool
}

type TaskBase struct {
	Paused      bool
	Traveller   *navigation.Traveller
	FieldCenter bool
	Household   Household
}

func (t *TaskBase) Pause(paused bool) {
	t.Paused = paused
}

func (t *TaskBase) IsPaused() bool {
	return t.Paused
}

func (t *TaskBase) SetUp(traveller *navigation.Traveller, household Household) {
	t.Traveller = traveller
	t.Household = household
}

func (t *TaskBase) IsFieldCenter() bool {
	return t.FieldCenter
}
