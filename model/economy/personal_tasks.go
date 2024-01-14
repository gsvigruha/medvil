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
	D navigation.Destination
	P Person
}

type GoToTask struct {
	TaskBase
	F *navigation.Field
}

func (t *EatTask) Destination() navigation.Destination {
	return t.F
}

func (t *EatTask) Complete(m navigation.IMap, tool bool) bool {
	t.P.Eat()
	return true
}

func (t *EatTask) Blocked() bool {
	return !t.P.HasFood()
}

func (t *EatTask) Name() string {
	return "eat"
}

func (t *EatTask) Tags() Tags {
	return EmptyTags
}

func (t *EatTask) Expired(Calendar *time.CalendarType) bool {
	return false
}

func (t *EatTask) Motion() uint8 {
	return navigation.MotionStand
}

func (t *EatTask) Description() string {
	return "Eat food"
}

func (t *DrinkTask) Destination() navigation.Destination {
	return t.F
}

func (t *DrinkTask) Complete(m navigation.IMap, tool bool) bool {
	t.P.Drink()
	return true
}

func (t *DrinkTask) Blocked() bool {
	return !t.P.HasDrink()
}

func (t *DrinkTask) Name() string {
	return "drink"
}

func (t *DrinkTask) Tags() Tags {
	return EmptyTags
}

func (t *DrinkTask) Expired(Calendar *time.CalendarType) bool {
	return false
}

func (t *DrinkTask) Motion() uint8 {
	return navigation.MotionStand
}

func (t *DrinkTask) Description() string {
	return "Drink"
}

func (t *HealTask) Destination() navigation.Destination {
	return t.F
}

func (t *HealTask) Complete(m navigation.IMap, tool bool) bool {
	t.P.Heal()
	return true
}

func (t *HealTask) Blocked() bool {
	return !t.P.HasMedicine()
}

func (t *HealTask) Name() string {
	return "heal"
}

func (t *HealTask) Tags() Tags {
	return EmptyTags
}

func (t *HealTask) Expired(Calendar *time.CalendarType) bool {
	return false
}

func (t *HealTask) Motion() uint8 {
	return navigation.MotionStand
}

func (t *HealTask) Description() string {
	return "Heal"
}

func (t *RelaxTask) Destination() navigation.Destination {
	return t.F
}

func (t *RelaxTask) Complete(m navigation.IMap, tool bool) bool {
	t.P.DrinkBeer()
	return true
}

func (t *RelaxTask) Blocked() bool {
	return !t.P.HasBeer()
}

func (t *RelaxTask) Name() string {
	return "relax"
}

func (t *RelaxTask) Tags() Tags {
	return EmptyTags
}

func (t *RelaxTask) Expired(Calendar *time.CalendarType) bool {
	return false
}

func (t *RelaxTask) Motion() uint8 {
	return navigation.MotionStand
}

func (t *RelaxTask) Description() string {
	return "Relax"
}

func (t *GoHomeTask) Destination() navigation.Destination {
	return t.D
}

func (t *GoHomeTask) Complete(m navigation.IMap, tool bool) bool {
	t.P.SetHome()
	return true
}

func (t *GoHomeTask) Blocked() bool {
	return false
}

func (t *GoHomeTask) Name() string {
	return "gohome"
}

func (t *GoHomeTask) Tags() Tags {
	return EmptyTags
}

func (t *GoHomeTask) Expired(Calendar *time.CalendarType) bool {
	return false
}

func (t *GoHomeTask) Motion() uint8 {
	return navigation.MotionStand
}

func (t *GoHomeTask) Description() string {
	return "Go home"
}

func (t *GoToTask) Destination() navigation.Destination {
	return t.F
}

func (t *GoToTask) Complete(m navigation.IMap, tool bool) bool {
	return true
}

func (t *GoToTask) Blocked() bool {
	return false
}

func (t *GoToTask) Name() string {
	return "goto"
}

func (t *GoToTask) Tags() Tags {
	return EmptyTags
}

func (t *GoToTask) Expired(Calendar *time.CalendarType) bool {
	return false
}

func (t *GoToTask) Motion() uint8 {
	return navigation.MotionStand
}

func (t *GoToTask) Description() string {
	return "Go to place"
}
