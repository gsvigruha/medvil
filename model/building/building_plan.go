package building

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"medvil/model/artifacts"
	"medvil/model/materials"
	"os"
)

const BuildingBaseMaxSize = 5
const MaxFloors = 4

const DirectionN uint8 = 0
const DirectionE uint8 = 1
const DirectionS uint8 = 2
const DirectionW uint8 = 3

var FloorMaterials = []*materials.Material{
	materials.GetMaterial("wood"),
	materials.GetMaterial("stone"),
	materials.GetMaterial("sandstone"),
	materials.GetMaterial("brick"),
	materials.GetMaterial("whitewash"),
}

var RoofMaterials = []*materials.Material{
	materials.GetMaterial("hay"),
	materials.GetMaterial("tile"),
}

var FlatRoofMaterials = []*materials.Material{
	materials.GetMaterial("stone"),
	materials.GetMaterial("sandstone"),
}

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

func (b BuildingPlan) BaseArea() uint16 {
	var baseArea uint16 = 0
	for i := 0; i < BuildingBaseMaxSize; i++ {
		for j := 0; j < BuildingBaseMaxSize; j++ {
			if b.BaseShape[i][j] {
				baseArea += 1
			}
		}
	}
	return baseArea
}

func (b BuildingPlan) Area() uint16 {
	baseArea := b.BaseArea()
	area := baseArea * uint16(len(b.Floors))
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

var cube = artifacts.GetArtifact("cube")
var board = artifacts.GetArtifact("board")

func (b BuildingPlan) ConstructionCost() []artifacts.Artifacts {
	var cubes uint16 = 0
	var boards uint16 = 0
	baseArea := b.BaseArea()
	for _, floor := range b.Floors {
		switch floor.M {
		case materials.GetMaterial("wood"):
			boards += baseArea * 3
		case materials.GetMaterial("sandstone"):
			boards += baseArea
			cubes += baseArea * 2
		case materials.GetMaterial("stone"):
			boards += baseArea
			cubes += baseArea * 2
		case materials.GetMaterial("brick"):
			boards += baseArea
			cubes += baseArea * 2
		case materials.GetMaterial("whitewash"):
			boards += baseArea
			cubes += baseArea * 2
		}
	}
	return []artifacts.Artifacts{
		artifacts.Artifacts{A: cube, Quantity: cubes},
		artifacts.Artifacts{A: board, Quantity: boards},
	}
}

func (b BuildingPlan) IsComplete() bool {
	if len(b.Floors) == 0 {
		return false
	}
	if b.Roof.M == nil {
		return false
	}
	for i := 0; i < BuildingBaseMaxSize-1; i++ {
		for j := 0; j < BuildingBaseMaxSize-1; j++ {
			if b.BaseShape[i][j] {
				return true
			}
		}
	}
	return false
}
