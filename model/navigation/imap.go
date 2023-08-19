package navigation

import (
	"medvil/model/building"
)

type IMap interface {
	GetField(x uint16, y uint16) *Field
	GetNField(x uint16, dx int, y uint16, dy int) *Field
	ShortPath(start Location, dest Destination, pathType PathType) *Path
	FindDest(start Location, dest Destination, pathType PathType) *Field
	SetBuildingUnits(b *building.Building, construction bool)
	Shore(x, y uint16) bool
	RandomSpot(x, y uint16, r int) *Field
}
