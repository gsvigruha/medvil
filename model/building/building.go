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

func (b *Building) Windows() uint16 {
	p := b.Plan
	windows := 0
	for i := 0; i < BuildingBaseMaxSize; i++ {
		if p.BaseShape[i][0] {
			windows += len(p.Floors) - int(p.WindowStartFloor[DirectionN])
		}
		if p.BaseShape[i][BuildingBaseMaxSize-1] {
			windows += len(p.Floors) - int(p.WindowStartFloor[DirectionS])
		}
		if p.BaseShape[0][i] {
			windows += len(p.Floors) - int(p.WindowStartFloor[DirectionW])
		}
		if p.BaseShape[BuildingBaseMaxSize-1][i] {
			windows += len(p.Floors) - int(p.WindowStartFloor[DirectionE])
		}
	}
	for i := 0; i < BuildingBaseMaxSize-1; i++ {
		for j := 0; j < BuildingBaseMaxSize-1; j++ {
			if p.BaseShape[i][j] != p.BaseShape[i+1][j] {
				windows += len(p.Floors)
			}
			if p.BaseShape[i][j] != p.BaseShape[i][j+1] {
				windows += len(p.Floors)
			}
		}
	}
	return uint16(windows)
}

func (b *Building) GetRoof(x uint8, y uint8) *RoofUnit {
	p := b.Plan
	if !p.BaseShape[x][y] {
		return nil
	}
	return &RoofUnit{
		B:    b,
		Roof: p.Roof,
		Elevated: [4]bool{
			y > 0 && p.BaseShape[x][y-1],
			x < BuildingBaseMaxSize-1 && p.BaseShape[x+1][y],
			y < BuildingBaseMaxSize-1 && p.BaseShape[x][y+1],
			x > 0 && p.BaseShape[x-1][y]}}
}

func (b *Building) ToBuildingUnits(x uint8, y uint8, construction bool) []BuildingUnit {
	p := b.Plan
	if !p.BaseShape[x][y] {
		return []BuildingUnit{}
	}
	numFloors := uint8(len(p.Floors))
	units := make([]BuildingUnit, numFloors)
	for i := uint8(0); i < numFloors; i++ {
		unitDoor := (i == 0 && p.DoorX == x && p.DoorY == y)
		var n *BuildingWall
		if y == 0 || !p.BaseShape[x][y-1] {
			door := (unitDoor && p.DoorD == DirectionN)
			n = &BuildingWall{M: p.Floors[i].M, Windows: !door && p.WindowStartFloor[DirectionN] <= i, Door: door, B: b, Construction: construction}
		}
		var e *BuildingWall
		if x == BuildingBaseMaxSize-1 || !p.BaseShape[x+1][y] {
			door := (unitDoor && p.DoorD == DirectionE)
			e = &BuildingWall{M: p.Floors[i].M, Windows: !door && p.WindowStartFloor[DirectionE] <= i, Door: door, B: b, Construction: construction}
		}
		var s *BuildingWall
		if y == BuildingBaseMaxSize-1 || !p.BaseShape[x][y+1] {
			door := (unitDoor && p.DoorD == DirectionS)
			s = &BuildingWall{M: p.Floors[i].M, Windows: !door && p.WindowStartFloor[DirectionS] <= i, Door: door, B: b, Construction: construction}
		}
		var w *BuildingWall
		if x == 0 || !p.BaseShape[x-1][y] {
			door := (unitDoor && p.DoorD == DirectionW)
			w = &BuildingWall{M: p.Floors[i].M, Windows: !door && p.WindowStartFloor[DirectionW] <= i, Door: door, B: b, Construction: construction}
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
			if b.Plan.BaseShape[i][j] {
				fields = append(fields, [2]uint16{bx, by})
			}
		}
	}
	idx := rand.Intn(len(fields))
	return fields[idx][0], fields[idx][1]
}
