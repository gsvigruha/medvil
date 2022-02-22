package economy

import (
	"medvil/model/artifacts"
	"medvil/model/navigation"
)

var fruit = artifacts.GetArtifact("fruit")
var vegetable = artifacts.GetArtifact("vegetable")
var bread = artifacts.GetArtifact("bread")
var meat = artifacts.GetArtifact("meat")

var water = artifacts.GetArtifact("water")
var wine = artifacts.GetArtifact("wine")
var beer = artifacts.GetArtifact("beer")

var Foods = []*artifacts.Artifact{fruit, vegetable, bread, meat}
var Drinks = []*artifacts.Artifact{water, wine, beer}

const MinFoodOrDrinkPerPerson uint16 = 2
const MaxFoodOrDrinkPerPerson uint16 = 5

func BuyFoodOrDrinkPerPerson() uint16 {
	return MaxFoodOrDrinkPerPerson - MinFoodOrDrinkPerPerson
}

func HasFood(r artifacts.Resources) bool {
	return AvailableFood(r) != nil
}

func HasDrink(r artifacts.Resources) bool {
	return AvailableDrink(r) != nil
}

func AvailableFood(r artifacts.Resources) []*artifacts.Artifact {
	var available []*artifacts.Artifact = nil
	for _, a := range Foods {
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
	for _, a := range Drinks {
		if q, ok := r.Artifacts[a]; ok {
			if q > 0 {
				available = append(available, a)
			}
		}
	}
	return available
}

func IsFoodOrDrink(a *artifacts.Artifact) bool {
	for _, a2 := range Foods {
		if a == a2 {
			return true
		}
	}
	for _, a2 := range Drinks {
		if a == a2 {
			return true
		}
	}
	return false
}

type PersonStateChange struct {
	Food      uint8
	Water     uint8
	Happiness uint8
	Health    uint8
}

var ArtifactToPersonState = map[*artifacts.Artifact]PersonStateChange{
	fruit:     PersonStateChange{Food: 50, Water: 50, Happiness: 25, Health: 25},
	vegetable: PersonStateChange{Food: 100, Water: 0, Happiness: 0, Health: 50},
	bread:     PersonStateChange{Food: 200, Water: 0, Happiness: 0, Health: 0},
	meat:      PersonStateChange{Food: 200, Water: 0, Happiness: 50, Health: 0},
	water:     PersonStateChange{Food: 0, Water: 200, Happiness: 0, Health: 0},
	wine:      PersonStateChange{Food: 0, Water: 150, Happiness: 100, Health: 50},
	beer:      PersonStateChange{Food: 0, Water: 150, Happiness: 100, Health: 0},
}

type WaterDestination struct{}

func (d WaterDestination) Check(f *navigation.Field) bool {
	return f.Terrain.T.Water
}