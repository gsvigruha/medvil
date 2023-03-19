package economy

import (
	"medvil/model/building"
	"medvil/model/navigation"
)

type ITown interface {
	DestroyMine(building *building.Building, m navigation.IMap)
	DestroyFarm(building *building.Building, m navigation.IMap)
}
