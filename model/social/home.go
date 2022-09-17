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
	HasBeer() bool
	HasDrink() bool
	HasFood() bool
	HasMedicine() bool
	Field(navigation.IMap) *navigation.Field
	NextTask(navigation.IMap, economy.Equipment) economy.Task
	GetResources() *artifacts.Resources
	GetBuilding() *building.Building
	GetHeating() float64
	HasEnoughTextile() bool
	AddVehicle(*vehicles.Vehicle)
	GetVehicle() *vehicles.Vehicle
}
