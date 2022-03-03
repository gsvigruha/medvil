package navigation

import (
	"medvil/model/building"
)

type IMap interface {
	GetField(x uint16, y uint16) *Field
	ShortPath(sx, sy, ex, ey uint16, travellerType uint8) *Path
	FindDest(sx, sy uint16, dest Destination, travellerType uint8) *Field
	SetBuildingUnits(b *building.Building, construction bool)
}
