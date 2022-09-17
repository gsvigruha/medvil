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
	HasBeer() bool
	HasDrink() bool
	HasFood() bool
	HasMedicine() bool
	Field(navigation.IMap) *navigation.Field
	RandomField(navigation.IMap) *navigation.Field
	NextTask(navigation.IMap, economy.Equipment) economy.Task
	GetResources() *artifacts.Resources
	GetBuilding() *building.Building
	GetHeating() float64
	HasEnoughTextile() bool
	AddVehicle(*vehicles.Vehicle)
	GetVehicle() *vehicles.Vehicle
	NumTasks(name string, tag string) int
}

func FindWaterTask(h Home, numP uint16, m navigation.IMap) {
	water := artifacts.GetArtifact("water")
	if h.GetResources().Get(water) < economy.MinFoodOrDrinkPerPerson*numP &&
		NumBatchesSimple(economy.MaxFoodOrDrinkPerPerson*numP, WaterTransportQuantity) > h.NumTasks("transport", "water") {
		hf := h.RandomField(m)
		if hf != nil {
			dest := m.FindDest(navigation.Location{X: hf.X, Y: hf.Y, Z: 0}, economy.WaterDestination{}, navigation.TravellerTypePedestrian)
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
