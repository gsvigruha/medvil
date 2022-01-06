package building

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"medvil/model/materials"
	"os"
)

const BuildingBaseMaxSize = 5

const DirectionN uint8 = 0
const DirectionE uint8 = 1
const DirectionS uint8 = 2
const DirectionW uint8 = 3

type Floor struct {
	M *materials.Material `json:"material"`
}

func (f *Floor) UnmarshalJSON(data []byte) error {
	var j map[string]string
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	f.M = materials.GetMaterial(j["material"])
	return nil
}

type Roof struct {
	M    *materials.Material `json:"material"`
	Flat bool                `json:"flat"`
}

func (r *Roof) UnmarshalJSON(data []byte) error {
	var j map[string]string
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	r.M = materials.GetMaterial(j["material"])
	r.Flat = j["flat"] == "true"
	return nil
}

type BuildingPlan struct {
	BaseShape        [BuildingBaseMaxSize][BuildingBaseMaxSize]bool `json:"baseShape"`
	WindowStartFloor [4]uint8                                       `json:"windowStartFloor"`
	Floors           []Floor                                        `json:"floors"`
	Roof             Roof                                           `json:"roof"`
	DoorX            uint8                                          `json:"doorX"`
	DoorY            uint8                                          `json:"doorY"`
	DoorD            uint8                                          `json:"doorD"`
}

func (b BuildingPlan) Area() uint16 {
	baseArea := 0
	for i := 0; i < BuildingBaseMaxSize; i++ {
		for j := 0; j < BuildingBaseMaxSize; j++ {
			if b.BaseShape[i][j] {
				baseArea += 1
			}
		}
	}
	area := baseArea * len(b.Floors)
	if !b.Roof.Flat {
		area += baseArea / 2
	}
	return uint16(area)
}

func (b BuildingPlan) Perimeter() uint16 {
	perimeter := 0
	for i := 0; i < BuildingBaseMaxSize; i++ {
		if b.BaseShape[i][0] {
			perimeter += 1
		}
		if b.BaseShape[i][BuildingBaseMaxSize-1] {
			perimeter += 1
		}
		if b.BaseShape[0][i] {
			perimeter += 1
		}
		if b.BaseShape[BuildingBaseMaxSize-1][i] {
			perimeter += 1
		}
	}
	for i := 0; i < BuildingBaseMaxSize-1; i++ {
		for j := 0; j < BuildingBaseMaxSize-1; j++ {
			if b.BaseShape[i][j] != b.BaseShape[i+1][j] {
				perimeter += 1
			}
			if b.BaseShape[i][j] != b.BaseShape[i][j+1] {
				perimeter += 1
			}
		}
	}
	return uint16(perimeter)
}

func (b BuildingPlan) Windows() uint16 {
	windows := 0
	for i := 0; i < BuildingBaseMaxSize; i++ {
		if b.BaseShape[i][0] {
			windows += len(b.Floors) - int(b.WindowStartFloor[DirectionN])
		}
		if b.BaseShape[i][BuildingBaseMaxSize-1] {
			windows += len(b.Floors) - int(b.WindowStartFloor[DirectionS])
		}
		if b.BaseShape[0][i] {
			windows += len(b.Floors) - int(b.WindowStartFloor[DirectionW])
		}
		if b.BaseShape[BuildingBaseMaxSize-1][i] {
			windows += len(b.Floors) - int(b.WindowStartFloor[DirectionE])
		}
	}
	for i := 0; i < BuildingBaseMaxSize-1; i++ {
		for j := 0; j < BuildingBaseMaxSize-1; j++ {
			if b.BaseShape[i][j] != b.BaseShape[i+1][j] {
				windows += len(b.Floors)
			}
			if b.BaseShape[i][j] != b.BaseShape[i][j+1] {
				windows += len(b.Floors)
			}
		}
	}
	return uint16(windows)
}

func (b BuildingPlan) GetRoof(x uint8, y uint8) *RoofUnit {
	if !b.BaseShape[x][y] {
		return nil
	}
	return &RoofUnit{
		Roof: b.Roof,
		Elevated: [4]bool{
			y > 0 && b.BaseShape[x][y-1],
			x < BuildingBaseMaxSize-1 && b.BaseShape[x+1][y],
			y < BuildingBaseMaxSize-1 && b.BaseShape[x][y+1],
			x > 0 && b.BaseShape[x-1][y]}}
}

func (b BuildingPlan) ToBuildingUnits(x uint8, y uint8) []BuildingUnit {
	if !b.BaseShape[x][y] {
		return []BuildingUnit{}
	}
	numFloors := uint8(len(b.Floors))
	units := make([]BuildingUnit, numFloors)
	for i := uint8(0); i < numFloors; i++ {
		unitDoor := (i == 0 && b.DoorX == x && b.DoorY == y)
		var n *BuildingWall
		if y == 0 || !b.BaseShape[x][y-1] {
			door := (unitDoor && b.DoorD == DirectionN)
			n = &BuildingWall{M: b.Floors[i].M, Windows: !door && b.WindowStartFloor[DirectionN] <= i, Door: door}
		}
		var e *BuildingWall
		if x == BuildingBaseMaxSize-1 || !b.BaseShape[x+1][y] {
			door := (unitDoor && b.DoorD == DirectionE)
			e = &BuildingWall{M: b.Floors[i].M, Windows: !door && b.WindowStartFloor[DirectionE] <= i, Door: door}
		}
		var s *BuildingWall
		if y == BuildingBaseMaxSize-1 || !b.BaseShape[x][y+1] {
			door := (unitDoor && b.DoorD == DirectionS)
			s = &BuildingWall{M: b.Floors[i].M, Windows: !door && b.WindowStartFloor[DirectionS] <= i, Door: door}
		}
		var w *BuildingWall
		if x == 0 || !b.BaseShape[x-1][y] {
			door := (unitDoor && b.DoorD == DirectionW)
			w = &BuildingWall{M: b.Floors[i].M, Windows: !door && b.WindowStartFloor[DirectionW] <= i, Door: door}
		}
		units[i].Walls = []*BuildingWall{n, e, s, w}
	}
	return units
}

func BuildingPlanFromJSON(fileName string) BuildingPlan {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}
	var plan BuildingPlan
	json.Unmarshal(byteValue, &plan)
	return plan
}
