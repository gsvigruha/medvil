package navigation

import (
	"medvil/model/building"
)

type FieldObject interface {
	Walkable() bool
	LiftN() int8
	LiftE() int8
	LiftS() int8
	LiftW() int8
}

type FieldBuildingObjects struct {
	BuildingComponents []building.BuildingComponent
}

func (o FieldBuildingObjects) Empty() bool {
	return len(o.BuildingComponents) == 0
}

func (o FieldBuildingObjects) GetBuildingComponent(z uint8) building.BuildingComponent {
	if o.Empty() {
		return nil
	}
	if len(o.BuildingComponents) <= int(z) {
		return nil
	}
	return o.BuildingComponents[z]
}

func (o FieldBuildingObjects) GetBuilding() *building.Building {
	if o.Empty() {
		return nil
	}
	return o.BuildingComponents[0].Building()
}

func (o FieldBuildingObjects) IsBuildingExtension() bool {
	if o.Empty() {
		return false
	}
	_, ok := o.BuildingComponents[0].(*building.ExtensionUnit)
	return ok
}

type BuildingPathElement struct {
	BC building.BuildingComponent
	L  Location
}

func (bpe *BuildingPathElement) GetLocation() Location {
	return bpe.L
}

func (bpe *BuildingPathElement) GetNeighbors(m IMap) []PathElement {
	f := m.GetField(bpe.L.X, bpe.L.Y)
	var n = []PathElement{}
	for dir, coordDelta := range building.CoordDeltaByDirection {
		x, y := uint16(coordDelta[0]+int(bpe.L.X)), uint16(coordDelta[1]+int(bpe.L.Y))
		nf := m.GetField(x, y)
		if nf != nil {
			if nf.Building.Empty() {
				if bpe.BC.Connection(uint8(dir)) == building.ConnectionTypeLowerLevel && bpe.L.Z == 1 {
					n = append(n, nf)
				}
			} else {
				oppDir := uint8((dir + 2) % 4)
				// The 0th building unit is the 1st Z, the 0th Z is the field level
				nbcBelow := nf.Building.GetBuildingComponent(bpe.L.Z - 2)
				nbcSame := nf.Building.GetBuildingComponent(bpe.L.Z - 1)
				nbcAbove := nf.Building.GetBuildingComponent(bpe.L.Z)
				if bpe.BC.Connection(uint8(dir)) == building.ConnectionTypeUpperLevel && nbcSame != nil && nbcSame.Connection(oppDir) == building.ConnectionTypeUpperLevel {
					n = append(n, &BuildingPathElement{BC: nbcSame, L: Location{X: x, Y: y, Z: bpe.L.Z}})
				} else if bpe.BC.Connection(uint8(dir)) == building.ConnectionTypeUpperLevel && nbcAbove != nil && nbcAbove.Connection(oppDir) == building.ConnectionTypeLowerLevel {
					n = append(n, &BuildingPathElement{BC: nbcAbove, L: Location{X: x, Y: y, Z: bpe.L.Z + 1}})
				} else if bpe.BC.Connection(uint8(dir)) == building.ConnectionTypeLowerLevel && nbcSame != nil && nbcSame.Connection(oppDir) == building.ConnectionTypeLowerLevel {
					n = append(n, &BuildingPathElement{BC: nbcSame, L: Location{X: x, Y: y, Z: bpe.L.Z}})
				} else if bpe.BC.Connection(uint8(dir)) == building.ConnectionTypeLowerLevel && nbcBelow != nil && nbcBelow.Connection(oppDir) == building.ConnectionTypeUpperLevel {
					n = append(n, &BuildingPathElement{BC: nbcBelow, L: Location{X: x, Y: y, Z: bpe.L.Z - 1}})
				}
			}
		}
	}
	// Towers allow vertical movement
	if bpe.BC.Building().Plan.BuildingType == building.BuildingTypeTower {
		for l := uint8(0); l < uint8(len(f.Building.BuildingComponents)); l++ {
			bc := f.Building.GetBuildingComponent(l)
			if bc != nil && !bc.IsConstruction() {
				n = append(n, &BuildingPathElement{BC: bc, L: Location{X: f.X, Y: f.Y, Z: l + 1}})
			}
		}
	}
	return n
}

func (bpe *BuildingPathElement) GetSpeed() float64 {
	return 1.0
}

func (bpe *BuildingPathElement) Walkable() bool {
	if bpe.BC.IsConstruction() {
		return false
	}
	return (bpe.BC.Building().Plan.BuildingType == building.BuildingTypeWall ||
		bpe.BC.Building().Plan.BuildingType == building.BuildingTypeGate ||
		bpe.BC.Building().Plan.BuildingType == building.BuildingTypeTower)
}

func (bpe *BuildingPathElement) Sailable() bool {
	if bpe.BC.Building().Plan.BuildingType == building.BuildingTypeGate {
		return true
	}
	return false
}

func (bpe *BuildingPathElement) TravellerVisible() bool {
	return bpe.BC.Building().Plan.BuildingType != building.BuildingTypeTower
}
