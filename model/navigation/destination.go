package navigation

import (
	"medvil/model/building"
)

func (l Location) Check(pe PathElement) bool {
	return l == pe.GetLocation()
}

type Destination interface {
	Check(PathElement) bool
}

type BuildingDestination struct {
	B *building.Building
}

func (bd BuildingDestination) Check(pe PathElement) bool {
	if bpe, ok := pe.(*BuildingPathElement); ok {
		return bpe.BC.Building() == bd.B
	}
	return false
}
