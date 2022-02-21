package economy

import (
	"medvil/model/navigation"
	"medvil/model/time"
)

type Person interface {
	Eat()
	Drink()
	SetHome()
	HasFood() bool
	HasDrink() bool
}

type EatTask struct {
	F *navigation.Field
	P Person
}

type DrinkTask struct {
	F *navigation.Field
	P Person
}

type GoHomeTask struct {
	F *navigation.Field
	P Person
}

func (t *EatTask) Field() *navigation.Field {
	return t.F
}

func (t *EatTask) Complete(Calendar *time.CalendarType) bool {
	t.P.Eat()
	return true
}

func (t *EatTask) Blocked() bool {
	return !t.P.HasFood()
}

func (t *EatTask) Name() string {
	return "eat"
}

func (t *EatTask) Tag() string {
	return ""
}

func (t *DrinkTask) Field() *navigation.Field {
	return t.F
}

func (t *DrinkTask) Complete(Calendar *time.CalendarType) bool {
	t.P.Drink()
	return true
}

func (t *DrinkTask) Blocked() bool {
	return !t.P.HasDrink()
}

func (t *DrinkTask) Name() string {
	return "drink"
}

func (t *DrinkTask) Tag() string {
	return ""
}

func (t *GoHomeTask) Field() *navigation.Field {
	return t.F
}

func (t *GoHomeTask) Complete(Calendar *time.CalendarType) bool {
	t.P.SetHome()
	return true
}

func (t *GoHomeTask) Blocked() bool {
	return false
}

func (t *GoHomeTask) Name() string {
	return "gohome"
}

func (t *GoHomeTask) Tag() string {
	return ""
}
