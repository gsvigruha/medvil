package building

import (
	"encoding/json"
	"math/rand"
)

const NumShapes = 5

type Building struct {
	Plan      BuildingPlan
	X         uint16
	Y         uint16
	Shape     uint8
	Direction uint8
}

func (b *Building) UnmarshalJSON(data []byte) error {
	var j map[string]json.RawMessage
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	if err := json.Unmarshal(j["x"], &b.X); err != nil {
		return err
	}
	if err := json.Unmarshal(j["y"], &b.Y); err != nil {
		return err
	}
	var planName string
	if err := json.Unmarshal(j["plan"], &planName); err != nil {
		return err
	}
	plan := BuildingPlanFromJSON("samples/building/" + planName + ".building.json")
	b.Plan = plan
	return nil
}

func (b *Building) getRoof(x uint8, y uint8, construction bool) *RoofUnit {
	p := b.Plan
	if p.BaseShape[x][y] == nil || p.BaseShape[x][y].Roof == nil {
		return nil
	}
	z := uint8(len(p.BaseShape[x][y].Floors))
	roof := p.BaseShape[x][y].Roof
	var elevated [4]bool
	if roof.RoofType == RoofTypeSplit {
		elevated = [4]bool{
			y > 0 && p.HasUnitOrRoof(x, y-1, z),
			x < BuildingBaseMaxSize-1 && p.HasUnitOrRoof(x+1, y, z),
			y < BuildingBaseMaxSize-1 && p.HasUnitOrRoof(x, y+1, z),
			x > 0 && p.HasUnitOrRoof(x-1, y, z)}
	} else if roof.RoofType == RoofTypeFlat {
		elevated = [4]bool{false, false, false, false}
	} else if roof.RoofType == RoofTypeRamp {
		if roof.RampD == DirectionN {
			elevated = [4]bool{true, false, false, false}
		} else if roof.RampD == DirectionE {
			elevated = [4]bool{false, true, false, false}
		} else if roof.RampD == DirectionS {
			elevated = [4]bool{false, false, true, false}
		} else if roof.RampD == DirectionW {
			elevated = [4]bool{false, false, false, true}
		}
	}
	return &RoofUnit{
		BuildingComponentBase: BuildingComponentBase{B: b, Construction: construction},
		Roof:                  *roof,
		Elevated:              elevated}
}

func (b *Building) hasWall(d uint8) bool {
	if b.Plan.BuildingType != BuildingTypeGate {
		return true
	}
	return b.Direction%2 != d%2
}

func (b *Building) hasDoor(d uint8, floor uint8) bool {
	if b.Direction != d {
		return false
	}
	if floor != 0 {
		return false
	}
	return true
}

func (b *Building) ToBuildingUnits(x uint8, y uint8, construction bool) []BuildingComponent {
	if b.Plan.BaseShape[x][y] == nil {
		return []BuildingComponent{}
	}
	p := b.Plan.BaseShape[x][y]
	numFloors := uint8(len(p.Floors))
	units := make([]BuildingComponent, numFloors)
	if p.Extension != nil {
		units = append(units, &ExtensionUnit{
			T:                     p.Extension.T,
			Direction:             GetExtensionDirection(p.Extension.T, x, y, b.Plan),
			BuildingComponentBase: BuildingComponentBase{B: b, Construction: construction},
		})
		return units
	}
	windows := (b.Plan.BuildingType != BuildingTypeWall && b.Plan.BuildingType != BuildingTypeGate)
	for i := uint8(0); i < numFloors; i++ {
		var n *BuildingWall
		if y == 0 || (!b.Plan.HasUnit(x, y-1, i) && b.hasWall(0)) {
			n = &BuildingWall{M: p.Floors[i].M, Windows: windows && !b.Plan.HasUnitOrRoof(x, y-1, i), Door: b.hasDoor(0, i)}
		}
		var e *BuildingWall
		if x == BuildingBaseMaxSize-1 || (!b.Plan.HasUnit(x+1, y, i) && b.hasWall(1)) {
			e = &BuildingWall{M: p.Floors[i].M, Windows: windows && !b.Plan.HasUnitOrRoof(x+1, y, i), Door: b.hasDoor(1, i)}
		}
		var s *BuildingWall
		if y == BuildingBaseMaxSize-1 || (!b.Plan.HasUnit(x, y+1, i) && b.hasWall(2)) {
			s = &BuildingWall{M: p.Floors[i].M, Windows: windows && !b.Plan.HasUnitOrRoof(x, y+1, i), Door: b.hasDoor(2, i)}
		}
		var w *BuildingWall
		if x == 0 || (!b.Plan.HasUnit(x-1, y, i) && b.hasWall(3)) {
			w = &BuildingWall{M: p.Floors[i].M, Windows: windows && !b.Plan.HasUnitOrRoof(x-1, y, i), Door: b.hasDoor(3, i)}
		}
		units[i] = &BuildingUnit{
			BuildingComponentBase: BuildingComponentBase{B: b, Construction: construction},
			Walls:                 []*BuildingWall{n, e, s, w},
		}
	}
	units = append(units, b.getRoof(x, y, construction))
	return units
}

func (b *Building) GetBuildingXYs(includeExtensions bool) [][2]uint16 {
	var fields [][2]uint16 = nil
	for i := uint16(0); i < 5; i++ {
		for j := uint16(0); j < 5; j++ {
			bx := uint16(b.X+i) - 2
			by := uint16(b.Y+j) - 2
			if b.Plan.BaseShape[i][j] != nil && (includeExtensions || b.Plan.BaseShape[i][j].Extension == nil) {
				fields = append(fields, [2]uint16{bx, by})
			}
		}
	}
	return fields
}

func (b *Building) GetRandomBuildingXY() (uint16, uint16, bool) {
	fields := b.GetBuildingXYs(true)
	if fields == nil {
		return 0, 0, false
	}
	idx := rand.Intn(len(fields))
	return fields[idx][0], fields[idx][1], true
}

func (b *Building) HasExtension(et *BuildingExtensionType) bool {
	e, _, _ := b.GetExtensionWithCoords()
	return e != nil && e.T == et
}

func (b *Building) GetExtensionWithCoords() (*BuildingExtension, uint16, uint16) {
	e, i, j := b.Plan.GetExtensionWithCoords()
	if e != nil {
		return e, uint16(b.X+i) - 2, uint16(b.Y+j) - 2
	}
	return nil, 0, 0
}
