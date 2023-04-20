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
	if f, ok := pe.(*Field); ok {
		return f.Building.GetBuilding() == bd.B
	}
	return false
}

type TravellerDestination struct {
	T *Traveller
}

func (td TravellerDestination) Check(pe PathElement) bool {
	return pe.GetLocation().X == td.T.FX && pe.GetLocation().Y == td.T.FY && pe.GetLocation().Z == td.T.FZ
}
