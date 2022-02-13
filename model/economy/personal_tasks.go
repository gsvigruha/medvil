package economy

import (
	"medvil/model/artifacts"
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

var fruit = artifacts.GetArtifact("fruit")
var vegetable = artifacts.GetArtifact("vegetable")
var bread = artifacts.GetArtifact("bread")
var meat = artifacts.GetArtifact("meat")

var water = artifacts.GetArtifact("water")
var wine = artifacts.GetArtifact("wine")
var beer = artifacts.GetArtifact("beer")

func HasFood(r artifacts.Resources) bool {
	return AvailableFood(r) != nil
}

func HasDrink(r artifacts.Resources) bool {
	return AvailableDrink(r) != nil
}

func AvailableFood(r artifacts.Resources) []*artifacts.Artifact {
	var available []*artifacts.Artifact = nil
	for _, a := range []*artifacts.Artifact{fruit, vegetable, bread, meat} {
		if q, ok := r.Artifacts[a]; ok {
			if q > 0 {
				available = append(available, a)
			}
		}
	}
	return available
}

func AvailableDrink(r artifacts.Resources) []*artifacts.Artifact {
	var available []*artifacts.Artifact = nil
	for _, a := range []*artifacts.Artifact{water, wine, beer} {
		if q, ok := r.Artifacts[a]; ok {
			if q > 0 {
				available = append(available, a)
			}
		}
	}
	return available
}

type PersonStateChange struct {
	Food      uint8
	Water     uint8
	Happiness uint8
	Health    uint8
}

var ArtifactToPersonState = map[*artifacts.Artifact]PersonStateChange{
	fruit:     PersonStateChange{Food: 100, Water: 50, Happiness: 50, Health: 50},
	vegetable: PersonStateChange{Food: 150, Water: 0, Happiness: 0, Health: 50},
	bread:     PersonStateChange{Food: 200, Water: 0, Happiness: 0, Health: 0},
	meat:      PersonStateChange{Food: 200, Water: 0, Happiness: 50, Health: 0},
	water:     PersonStateChange{Food: 0, Water: 200, Happiness: 0, Health: 0},
	wine:      PersonStateChange{Food: 0, Water: 150, Happiness: 100, Health: 50},
	beer:      PersonStateChange{Food: 0, Water: 150, Happiness: 100, Health: 0},
}
