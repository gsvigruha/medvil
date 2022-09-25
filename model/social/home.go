package social

import (
	"medvil/model/artifacts"
	"medvil/model/building"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/vehicles"
)

type Home interface {
	AddTask(economy.Task)
	AddPriorityTask(economy.Task)
	GetTasks() []economy.Task
	SetTasks([]economy.Task)
	HasBeer() bool
	HasDrink() bool
	HasFood() bool
	HasMedicine() bool
	Field(navigation.IMap) *navigation.Field
	RandomField(navigation.IMap, func(navigation.Field) bool) *navigation.Field
	NextTask(navigation.IMap, economy.Equipment) economy.Task
	GetResources() *artifacts.Resources
	GetBuilding() *building.Building
	GetHeating() float64
	HasEnoughTextile() bool
	AddVehicle(*vehicles.Vehicle)
	GetVehicle() *vehicles.Vehicle
	NumTasks(name string, tag string) int
	GetMoney() *uint32
}

var water = artifacts.GetArtifact("water")

func needsWater(h Home, numP uint16) bool {
	if h.GetResources().Get(water) < economy.MinFoodOrDrinkPerPerson*numP &&
		NumBatchesSimple(economy.MaxFoodOrDrinkPerPerson*numP, WaterTransportQuantity) > h.NumTasks("transport", "water") {
		return true
	}
	if h.GetResources().Get(water) == 0 && h.NumTasks("transport", "water") == 0 {
		return true
	}
	return false
}

func FindWaterTask(h Home, numP uint16, m navigation.IMap) {
	if needsWater(h, numP) {
		hf := h.RandomField(m, navigation.Field.BuildingNonExtension)
		if hf != nil {
			dest := m.FindDest(navigation.Location{X: hf.X, Y: hf.Y, Z: 0}, economy.WaterDestination{}, navigation.PathTypePedestrian)
			if dest != nil {
				h.AddPriorityTask(&economy.TransportTask{
					PickupF:  dest,
					DropoffF: hf,
					PickupR:  &dest.Terrain.Resources,
					DropoffR: h.GetResources(),
					A:        water,
					Quantity: WaterTransportQuantity,
				})
			}
		}
	}
}

func needsFood(h Home, numP uint16, a *artifacts.Artifact) bool {
	tag := "food_shopping#" + a.Name
	if NumBatchesSimple(economy.BuyFoodOrDrinkPerPerson()*numP, FoodTransportQuantity) > h.NumTasks("exchange", tag) {
		return true
	}
	if h.GetResources().Get(a) == 0 && h.NumTasks("exchange", tag) == 0 {
		return true
	}
	return false
}

func GetFoodTasks(h Home, numP uint16, mp *Marketplace) {
	var numFoodBatchesNeeded = 0
	for _, a := range economy.Foods {
		if h.GetResources().Get(a) < economy.MinFoodOrDrinkPerPerson*numP {
			numFoodBatchesNeeded += NumBatchesSimple(economy.BuyFoodOrDrinkPerPerson()*numP, FoodTransportQuantity)
		}
	}
	if numFoodBatchesNeeded == 0 {
		numFoodBatchesNeeded = 1
	}
	for _, a := range economy.Foods {
		if h.GetResources().Get(a) < economy.MinFoodOrDrinkPerPerson*numP {
			if needsFood(h, numP, a) {
				needs := []artifacts.Artifacts{artifacts.Artifacts{A: a, Quantity: FoodTransportQuantity}}
				var maxPrice = *h.GetMoney() / uint32(numFoodBatchesNeeded)
				if maxPrice > mp.Price(needs)*2 {
					maxPrice = mp.Price(needs) * 2
				}
				if *h.GetMoney() >= mp.Price(needs) && mp.HasTraded(a) {
					h.AddPriorityTask(&economy.BuyTask{
						Exchange:       mp,
						HouseholdMoney: h.GetMoney(),
						Goods:          needs,
						MaxPrice:       maxPrice,
						TaskTag:        "food_shopping#" + a.Name,
					})
				}
			}
		}
	}
}
