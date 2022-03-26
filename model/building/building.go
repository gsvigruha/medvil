package building

import (
	"encoding/json"
	"math/rand"
)

type Building struct {
	Plan *BuildingPlan
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
	b.Plan = &plan
	return nil
}

func (b *Building) GetRoof(x uint8, y uint8) *RoofUnit {
	p := b.Plan
	if p.BaseShape[x][y] == nil {
		return nil
	}
	z := uint8(len(p.BaseShape[x][y].Floors))
	return &RoofUnit{
		B:    b,
		Roof: *(p.BaseShape[x][y].Roof),
		Elevated: [4]bool{
			y > 0 && p.HasUnitOrRoof(x, y-1, z),
			x < BuildingBaseMaxSize-1 && p.HasUnitOrRoof(x+1, y, z),
			y < BuildingBaseMaxSize-1 && p.HasUnitOrRoof(x, y+1, z),
			x > 0 && p.HasUnitOrRoof(x-1, y, z)}}
}

func (b *Building) ToBuildingUnits(x uint8, y uint8, construction bool) []BuildingUnit {
	if b.Plan.BaseShape[x][y] == nil {
		return []BuildingUnit{}
	}
	p := b.Plan.BaseShape[x][y]
	numFloors := uint8(len(p.Floors))
	units := make([]BuildingUnit, numFloors)
	windows := (b.Plan.WindowType == WindowTypeRegular)
	for i := uint8(0); i < numFloors; i++ {
		var n *BuildingWall
		if y == 0 || !b.Plan.HasUnit(x, y-1, i) {
			n = &BuildingWall{M: p.Floors[i].M, Windows: windows && !b.Plan.HasUnitOrRoof(x, y-1, i), Door: false, B: b, Construction: construction}
		}
		var e *BuildingWall
		if x == BuildingBaseMaxSize-1 || !b.Plan.HasUnit(x+1, y, i) {
			e = &BuildingWall{M: p.Floors[i].M, Windows: windows && !b.Plan.HasUnitOrRoof(x+1, y, i), Door: false, B: b, Construction: construction}
		}
		var s *BuildingWall
		if y == BuildingBaseMaxSize-1 || !b.Plan.HasUnit(x, y+1, i) {
			s = &BuildingWall{M: p.Floors[i].M, Windows: windows && !b.Plan.HasUnitOrRoof(x, y+1, i), Door: false, B: b, Construction: construction}
		}
		var w *BuildingWall
		if x == 0 || !b.Plan.HasUnit(x-1, y, i) {
			w = &BuildingWall{M: p.Floors[i].M, Windows: windows && !b.Plan.HasUnitOrRoof(x-1, y, i), Door: false, B: b, Construction: construction}
		}
		units[i].Walls = []*BuildingWall{n, e, s, w}
		units[i].B = b
	}
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
