package building

import (
	"encoding/json"
)

type Building struct {
	Plan      BuildingPlan
	X         uint16
	Y         uint16
	Shape     uint8
	Direction uint8
	Broken    bool
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
	var connected [4]bool
	if roof.RoofType == RoofTypeSplit {
		connected = [4]bool{
			y > 0 && p.HasUnitOrRoof(x, y-1, z),
			x < BuildingBaseMaxSize-1 && p.HasUnitOrRoof(x+1, y, z),
			y < BuildingBaseMaxSize-1 && p.HasUnitOrRoof(x, y+1, z),
			x > 0 && p.HasUnitOrRoof(x-1, y, z)}
	} else if roof.RoofType == RoofTypeFlat {
		connected = [4]bool{
			y > 0 && p.HasUnit(x, y-1, z-1),
			x < BuildingBaseMaxSize-1 && p.HasUnit(x+1, y, z-1),
			y < BuildingBaseMaxSize-1 && p.HasUnit(x, y+1, z-1),
			x > 0 && p.HasUnit(x-1, y, z-1)}
	} else if roof.RoofType == RoofTypeRamp {
		if roof.RampD == DirectionN {
			connected = [4]bool{true, false, false, false}
		} else if roof.RampD == DirectionE {
			connected = [4]bool{false, true, false, false}
		} else if roof.RampD == DirectionS {
			connected = [4]bool{false, false, true, false}
		} else if roof.RampD == DirectionW {
			connected = [4]bool{false, false, false, true}
		}
	}
	return &RoofUnit{
		BuildingComponentBase: BuildingComponentBase{B: b, Construction: construction},
		Roof:                  *roof,
		Connected:             connected}
}

func (b *Building) hasArch(d uint8) bool {
	if b.Plan.BuildingType == BuildingTypeMarket {
		return true
	}
	if b.Plan.BuildingType != BuildingTypeGate {
		return false
	}
	return b.Direction%2 == d%2
}

func (b *Building) hasDoor(d uint8, floor uint8, open bool) bool {
	if !open {
		return false
	}
	if b.Plan.BuildingType == BuildingTypeGate || b.Plan.BuildingType == BuildingTypeWall || b.Plan.BuildingType == BuildingTypeMarket {
		return false
	}
	if b.Direction != d {
		return false
	}
	if floor != 0 {
		return false
	}
	return true
}

func (b *Building) getWindowType(open bool, floor uint8) WindowType {
	if !open {
		return WindowTypeNone
	}
	if b.Plan.BuildingType == BuildingTypeWall || b.Plan.BuildingType == BuildingTypeGate || b.Plan.BuildingType == BuildingTypeTower || b.Plan.BuildingType == BuildingTypeMarket {
		return WindowTypeNone
	}
	if b.Plan.BuildingType == BuildingTypeFactory {
		return WindowTypeFactory
	}
	if b.Plan.BuildingType == BuildingTypeWorkshop || b.Plan.BuildingType == BuildingTypeTownhall {
		switch b.Shape % 3 {
		case 0:
			if floor == 0 {
				return WindowTypePlain
			} else {
				return WindowTypeBalcony
			}
		case 1:
			return WindowTypePlain
		case 2:
			return WindowTypeFrench
		}
	}
	return WindowTypePlain
}

func (b *Building) ToBuildingUnits(x uint8, y uint8, construction bool) []BuildingComponent {
	if b.Plan.BaseShape[x][y] == nil {
		return []BuildingComponent{}
	}
	p := b.Plan.BaseShape[x][y]
	numFloors := uint8(len(p.Floors))
	units := make([]BuildingComponent, numFloors)
	if p.Extension != nil && !p.Extension.T.InUnit {
		units = append(units, &ExtensionUnit{
			T:                     p.Extension.T,
			Direction:             GetExtensionDirection(p.Extension.T, x, y, b.Plan),
			BuildingComponentBase: BuildingComponentBase{B: b, Construction: construction},
		})
		return units
	}
	for i := uint8(0); i < numFloors; i++ {
		var n *BuildingWall
		if y == 0 || !b.Plan.HasUnit(x, y-1, i) {
			open := !b.Plan.HasUnitOrRoof(x, y-1, i)
			n = &BuildingWall{M: p.Floors[i].M, Windows: b.getWindowType(open, i), Door: b.hasDoor(0, i, open), Arch: b.hasArch(0)}
		}
		var e *BuildingWall
		if x == BuildingBaseMaxSize-1 || !b.Plan.HasUnit(x+1, y, i) {
			open := !b.Plan.HasUnitOrRoof(x+1, y, i)
			e = &BuildingWall{M: p.Floors[i].M, Windows: b.getWindowType(open, i), Door: b.hasDoor(1, i, open), Arch: b.hasArch(1)}
		}
		var s *BuildingWall
		if y == BuildingBaseMaxSize-1 || !b.Plan.HasUnit(x, y+1, i) {
			open := !b.Plan.HasUnitOrRoof(x, y+1, i)
			s = &BuildingWall{M: p.Floors[i].M, Windows: b.getWindowType(open, i), Door: b.hasDoor(2, i, open), Arch: b.hasArch(2)}
		}
		var w *BuildingWall
		if x == 0 || !b.Plan.HasUnit(x-1, y, i) {
			open := !b.Plan.HasUnitOrRoof(x-1, y, i)
			w = &BuildingWall{M: p.Floors[i].M, Windows: b.getWindowType(open, i), Door: b.hasDoor(3, i, open), Arch: b.hasArch(3)}
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
			if b.Plan.BaseShape[i][j] != nil && (includeExtensions || b.Plan.BaseShape[i][j].Extension == nil || len(b.Plan.BaseShape[i][j].Floors) > 0) {
				fields = append(fields, [2]uint16{bx, by})
			}
		}
	}
	return fields
}

func (b *Building) HasExtension(et *BuildingExtensionType) bool {
	for _, e := range b.Plan.GetExtensions() {
		if e.T == et {
			return true
		}
	}
	return false
}

func (b *Building) GetExtensionWithCoords(et *BuildingExtensionType) (*BuildingExtension, uint16, uint16) {
	e, i, j := b.Plan.GetExtensionWithCoords(et)
	if e != nil {
		return e, uint16(b.X+i) - 2, uint16(b.Y+j) - 2
	}
	return nil, 0, 0
}
