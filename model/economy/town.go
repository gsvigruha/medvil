package economy

import (
	"medvil/model/building"
	"medvil/model/navigation"
)

type ITown interface {
	DestroyBuilding(building *building.Building, m navigation.IMap)
	DestroyRoad(road *building.Road, m navigation.IMap)
}
