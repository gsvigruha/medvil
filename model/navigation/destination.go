package navigation

import (
	"medvil/model/building"
)

type Destination interface {
	Check(PathElement) bool
}

type BuildingDestination struct {
	B  *building.Building
	ET *building.BuildingExtensionType
}

func (bd *BuildingDestination) extensionMatch(bc building.BuildingComponent) bool {
	unit, ok := bc.(*building.ExtensionUnit)
	if !ok {
		return bd.ET == nil
	} else {
		return unit.T == bd.ET
	}
}

func (bd *BuildingDestination) Check(pe PathElement) bool {
	if bpe, ok := pe.(*BuildingPathElement); ok {
		if bpe.BC.Building() == bd.B {
			return bd.extensionMatch(bpe.BC)
		}
	}
	if f, ok := pe.(*Field); ok {
		if f.Building.GetBuilding() == bd.B {
			return bd.extensionMatch(f.Building.BuildingComponents[0])
		}
	}
	return false
}

type TravellerDestination struct {
	T *Traveller
}

func (td TravellerDestination) Check(pe PathElement) bool {
	return pe.GetLocation().X == td.T.FX && pe.GetLocation().Y == td.T.FY && pe.GetLocation().Z == td.T.FZ
}
