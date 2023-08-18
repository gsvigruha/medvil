package economy

import (
	"medvil/model/navigation"
	"medvil/model/time"
)

type Person interface {
	Eat()
	Drink()
	Heal()
	DrinkBeer()
	SetHome()
	HasFood() bool
	HasDrink() bool
	HasMedicine() bool
	HasBeer() bool
}

var PersonalTasks = []string{"eat", "drink", "heal", "relax", "gohome", "goto"}

func IsPersonalTask(n string) bool {
	for _, t := range PersonalTasks {
		if t == n {
			return true
		}
	}
	return false
}

type EatTask struct {
	TaskBase
	F *navigation.Field
	P Person
}

type DrinkTask struct {
	TaskBase
	F *navigation.Field
	P Person
}

type HealTask struct {
	TaskBase
	F *navigation.Field
	P Person
}

type RelaxTask struct {
	TaskBase
	F *navigation.Field
	P Person
}

type GoHomeTask struct {
	TaskBase
	F *navigation.Field
	P Person
}

type GoToTask struct {
	TaskBase
	F *navigation.Field
}

func (t *EatTask) Destination() navigation.Destination {
	return t.F
}

func (t *EatTask) Complete(Calendar *time.CalendarType, tool bool) bool {
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

func (t *EatTask) Expired(Calendar *time.CalendarType) bool {
	return false
}

func (t *EatTask) Motion() uint8 {
	return navigation.MotionStand
}

func (t *DrinkTask) Destination() navigation.Destination {
	return t.F
}

func (t *DrinkTask) Complete(Calendar *time.CalendarType, tool bool) bool {
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

func (t *DrinkTask) Expired(Calendar *time.CalendarType) bool {
	return false
}

func (t *DrinkTask) Motion() uint8 {
	return navigation.MotionStand
}

func (t *HealTask) Destination() navigation.Destination {
	return t.F
}

func (t *HealTask) Complete(Calendar *time.CalendarType, tool bool) bool {
	t.P.Heal()
	return true
}

func (t *HealTask) Blocked() bool {
	return !t.P.HasMedicine()
}

func (t *HealTask) Name() string {
	return "heal"
}

func (t *HealTask) Tag() string {
	return ""
}

func (t *HealTask) Expired(Calendar *time.CalendarType) bool {
	return false
}

func (t *HealTask) Motion() uint8 {
	return navigation.MotionStand
}

func (t *RelaxTask) Destination() navigation.Destination {
	return t.F
}

func (t *RelaxTask) Complete(Calendar *time.CalendarType, tool bool) bool {
	t.P.DrinkBeer()
	return true
}

func (t *RelaxTask) Blocked() bool {
	return !t.P.HasBeer()
}

func (t *RelaxTask) Name() string {
	return "relax"
}

func (t *RelaxTask) Tag() string {
	return ""
}

func (t *RelaxTask) Expired(Calendar *time.CalendarType) bool {
	return false
}

func (t *RelaxTask) Motion() uint8 {
	return navigation.MotionStand
}

func (t *GoHomeTask) Destination() navigation.Destination {
	return t.F
}

func (t *GoHomeTask) Complete(Calendar *time.CalendarType, tool bool) bool {
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

func (t *GoHomeTask) Expired(Calendar *time.CalendarType) bool {
	return false
}

func (t *GoHomeTask) Motion() uint8 {
	return navigation.MotionStand
}

func (t *GoToTask) Destination() navigation.Destination {
	return t.F
}

func (t *GoToTask) Complete(Calendar *time.CalendarType, tool bool) bool {
	return true
}

func (t *GoToTask) Blocked() bool {
	return false
}

func (t *GoToTask) Name() string {
	return "goto"
}

func (t *GoToTask) Tag() string {
	return ""
}

func (t *GoToTask) Expired(Calendar *time.CalendarType) bool {
	return false
}

func (t *GoToTask) Motion() uint8 {
	return navigation.MotionStand
}
