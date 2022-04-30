package building

import (
	"encoding/json"
	"math/rand"
)

type Building struct {
	Plan BuildingPlan
	X    uint16
	Y    uint16
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
	windows := (b.Plan.BuildingType != BuildingTypeWall)
	for i := uint8(0); i < numFloors; i++ {
		var n *BuildingWall
		if y == 0 || !b.Plan.HasUnit(x, y-1, i) {
			n = &BuildingWall{M: p.Floors[i].M, Windows: windows && !b.Plan.HasUnitOrRoof(x, y-1, i), Door: false}
		}
		var e *BuildingWall
		if x == BuildingBaseMaxSize-1 || !b.Plan.HasUnit(x+1, y, i) {
			e = &BuildingWall{M: p.Floors[i].M, Windows: windows && !b.Plan.HasUnitOrRoof(x+1, y, i), Door: false}
		}
		var s *BuildingWall
		if y == BuildingBaseMaxSize-1 || !b.Plan.HasUnit(x, y+1, i) {
			s = &BuildingWall{M: p.Floors[i].M, Windows: windows && !b.Plan.HasUnitOrRoof(x, y+1, i), Door: false}
		}
		var w *BuildingWall
		if x == 0 || !b.Plan.HasUnit(x-1, y, i) {
			w = &BuildingWall{M: p.Floors[i].M, Windows: windows && !b.Plan.HasUnitOrRoof(x-1, y, i), Door: false}
		}
		units[i] = &BuildingUnit{
			BuildingComponentBase: BuildingComponentBase{B: b, Construction: construction},
			Walls:                 []*BuildingWall{n, e, s, w},
		}
	}
	units = append(units, b.getRoof(x, y, construction))
	return units
}

func (b *Building) GetRandomBuildingXY() (uint16, uint16) {
	var fields [][2]uint16
	for i := uint16(0); i < 5; i++ {
		for j := uint16(0); j < 5; j++ {
			bx := uint16(b.X+i) - 2
			by := uint16(b.Y+j) - 2
			if b.Plan.BaseShape[i][j] != nil {
				fields = append(fields, [2]uint16{bx, by})
			}
		}
	}
	idx := rand.Intn(len(fields))
	return fields[idx][0], fields[idx][1]
}

func (b *Building) GetExtension() *BuildingExtension {
	for i := uint16(0); i < 5; i++ {
		for j := uint16(0); j < 5; j++ {
			if b.Plan.BaseShape[i][j] != nil && b.Plan.BaseShape[i][j].Extension != nil {
				return b.Plan.BaseShape[i][j].Extension
			}
		}
	}
	return nil
}
