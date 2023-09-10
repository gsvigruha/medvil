package economy

import (
	"medvil/model/navigation"
	"medvil/model/time"
)

type Task interface {
	Complete(Calendar *time.CalendarType, tool bool) bool
	Destination() navigation.Destination
	Blocked() bool
	Name() string
	IconName() string
	Tag() string
	Expired(Calendar *time.CalendarType) bool
	Pause(bool)
	IsPaused() bool
	SetUp(traveller *navigation.Traveller, household Household, person Person)
	Motion() uint8
	IsFieldCenter() bool
	Equipped(*EquipmentType) bool
}

type TaskBase struct {
	Paused      bool
	Traveller   *navigation.Traveller
	FieldCenter bool
	Household   Household
	Person      Person
}

func (t *TaskBase) Pause(paused bool) {
	t.Paused = paused
}

func (t *TaskBase) IsPaused() bool {
	return t.Paused
}

func (t *TaskBase) SetUp(traveller *navigation.Traveller, household Household, person Person) {
	t.Traveller = traveller
	t.Household = household
	t.Person = person
}

func (t *TaskBase) IsFieldCenter() bool {
	return t.FieldCenter
}

func (t *TaskBase) Equipped(*EquipmentType) bool {
	return true
}

func (t *TaskBase) IconName() string {
	return ""
}

func IconName(t Task) string {
	if t.IconName() != "" {
		return t.IconName()
	}
	return t.Name()
}

func HasPersonalTask(tasks []Task) bool {
	for _, task := range tasks {
		if IsPersonalTask(task.Name()) {
			return true
		}
	}
	return false
}
