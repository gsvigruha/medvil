package social

import (
	"medvil/model/artifacts"
	"medvil/model/building"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/vehicles"
	"sort"
)

type Home interface {
	economy.Wallet
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
	NextTask(navigation.IMap, *economy.EquipmentType) economy.Task
	GetResources() *artifacts.Resources
	GetBuilding() *building.Building
	GetHeating() uint8
	HasEnoughClothes() bool
	AddVehicle(*vehicles.Vehicle)
	AllocateVehicle(waterOk bool) *vehicles.Vehicle
	NumTasks(name string, tag string) int
	Destination(extensionType *building.BuildingExtensionType) navigation.Destination
	PendingCosts() uint32
	Broken() bool
	GetTown() *Town
	GetPeople() []*Person
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
					PickupD:        dest,
					DropoffD:       h.Destination(building.NonExtension),
					PickupR:        dest.Terrain.Resources,
					DropoffR:       h.GetResources(),
					A:              water,
					TargetQuantity: WaterTransportQuantity,
				})
			}
		}
	}
}

func numFoodBatchesNeeded(h Home, numP uint16, a *artifacts.Artifact) int {
	tag := "food_shopping#" + a.Name
	has := uint16(h.NumTasks("exchange", tag)*FoodTransportQuantity) + h.GetResources().Get(a)
	needs := economy.BuyFoodOrDrinkPerPerson() * numP
	if needs > has {
		return NumBatches(needs-has, 0, FoodTransportQuantity)
	}
	return 0
}

func GetFoodTasks(h Home, numP uint16, mp *Marketplace) {
	budget := int(h.GetMoney()) - int(h.PendingCosts())
	var foodBatches []*artifacts.Artifact
	for _, a := range economy.Foods {
		batch := []artifacts.Artifacts{artifacts.Artifacts{A: a, Quantity: FoodTransportQuantity}}
		if int(mp.Price(batch)) < budget && mp.HasTraded(a) {
			for i := 0; i < numFoodBatchesNeeded(h, numP, a); i++ {
				foodBatches = append(foodBatches, a)
			}
		}
	}
	sort.Slice(foodBatches, func(i, j int) bool { return mp.Prices[foodBatches[i]] < mp.Prices[foodBatches[j]] })

	var totalCost = 0
	for _, a := range foodBatches {
		needs := []artifacts.Artifacts{artifacts.Artifacts{A: a, Quantity: FoodTransportQuantity}}
		price := mp.Price(needs)
		if budget > totalCost+int(price) {
			h.AddPriorityTask(&economy.BuyTask{
				Exchange:        mp,
				HouseholdWallet: h,
				Goods:           needs,
				MaxPrice:        price * 2,
				TaskTag:         "food_shopping#" + a.Name,
			})
			totalCost += int(price)
		} else {
			break
		}
	}
}

func PendingCosts(tasks []economy.Task) uint32 {
	var costs uint32
	for _, task := range tasks {
		if buyTask, ok := task.(*economy.BuyTask); ok {
			costs += buyTask.Exchange.Price(buyTask.Goods)
		}
		if exchangeTask, ok := task.(*economy.ExchangeTask); ok {
			costs += exchangeTask.Exchange.Price(exchangeTask.GoodsToBuy)
		}
	}
	return costs
}
