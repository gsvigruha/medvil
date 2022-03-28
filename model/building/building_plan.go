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
const DirectionNone uint8 = 255

var CoordDeltaByDirection = [4][2]int{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}

type Floor struct {
	M *materials.Material `json:"material"`
}

type RoofType uint8

const RoofTypeFlat = 1
const RoofTypeSplit = 2
const RoofTypeRamp = 3

type Roof struct {
	M        *materials.Material `json:"material"`
	RoofType RoofType
	RampD    uint8
}

func (r Roof) Flat() bool {
	return r.RoofType == RoofTypeFlat
}

type PlanUnits struct {
	Floors []Floor
	Roof   *Roof
}

type BuildingPlan struct {
	BaseShape    [BuildingBaseMaxSize][BuildingBaseMaxSize]*PlanUnits
	BuildingType BuildingType
}

func (b *BuildingPlan) UnmarshalJSON(data []byte) error {
	var j map[string][][][]string
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	shape := j["baseShape"]
	for i := range shape {
		for j := range shape[i] {
			if shape[i][j] != nil {
				b.BaseShape[i][j] = &PlanUnits{}
				for k := range shape[i][j] {
					if k < len(shape[i][j])-1 {
						b.BaseShape[i][j].Floors = append(b.BaseShape[i][j].Floors, Floor{M: materials.GetMaterial(shape[i][j][k])})
					} else {
						b.BaseShape[i][j].Roof = &Roof{M: materials.GetMaterial(shape[i][j][k]), RoofType: RoofTypeSplit}
					}
				}
			}
		}
	}
	return nil
}

func (b BuildingPlan) BaseArea() uint16 {
	var baseArea uint16 = 0
	for i := 0; i < BuildingBaseMaxSize; i++ {
		for j := 0; j < BuildingBaseMaxSize; j++ {
			if b.BaseShape[i][j] != nil {
				baseArea += 1
			}
		}
	}
	return baseArea
}

func (b BuildingPlan) Area() uint16 {
	var area = uint16(0)
	for i := 0; i < BuildingBaseMaxSize; i++ {
		for j := 0; j < BuildingBaseMaxSize; j++ {
			if b.BaseShape[i][j] != nil {
				area += uint16(len(b.BaseShape[i][j].Floors))
			}
		}
	}
	return area
}

func (b BuildingPlan) RoofArea() uint16 {
	var area = uint16(0)
	for i := 0; i < BuildingBaseMaxSize; i++ {
		for j := 0; j < BuildingBaseMaxSize; j++ {
			if b.BaseShape[i][j] != nil && !b.BaseShape[i][j].Roof.Flat() {
				area += 1
			}
		}
	}
	return area
}

func (b BuildingPlan) Perimeter() uint16 {
	perimeter := 0
	for i := 0; i < BuildingBaseMaxSize; i++ {
		if b.BaseShape[i][0] != nil {
			perimeter += 1
		}
		if b.BaseShape[i][BuildingBaseMaxSize-1] != nil {
			perimeter += 1
		}
		if b.BaseShape[0][i] != nil {
			perimeter += 1
		}
		if b.BaseShape[BuildingBaseMaxSize-1][i] != nil {
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
var brick = artifacts.GetArtifact("brick")
var thatch = artifacts.GetArtifact("thatch")
var tile = artifacts.GetArtifact("tile")

func (b BuildingPlan) ConstructionCost() []artifacts.Artifacts {
	var cubes uint16 = 0
	var boards uint16 = 0
	var bricks uint16 = 0
	var thatches uint16 = 0
	var tiles uint16 = 0
	for i := 0; i < BuildingBaseMaxSize; i++ {
		for j := 0; j < BuildingBaseMaxSize; j++ {
			if b.BaseShape[i][j] != nil {
				for _, floor := range b.BaseShape[i][j].Floors {
					switch floor.M {
					case materials.GetMaterial("wood"):
						boards += 3
					case materials.GetMaterial("sandstone"):
						boards += 1
						cubes += 2
					case materials.GetMaterial("stone"):
						boards += 1
						cubes += 2
					case materials.GetMaterial("brick"):
						boards += 1
						cubes += 2
					case materials.GetMaterial("whitewash"):
						boards += 1
						cubes += 2
					}
				}
				if !b.BaseShape[i][j].Roof.Flat() {
					switch b.BaseShape[i][j].Roof.M {
					case materials.GetMaterial("tile"):
						tiles += 1
						boards += 1
					case materials.GetMaterial("hay"):
						thatches += 1
						boards += 1
					case materials.GetMaterial("stone"):
						cubes += 1
					}
				}
			}
		}
	}

	return []artifacts.Artifacts{
		artifacts.Artifacts{A: cube, Quantity: cubes},
		artifacts.Artifacts{A: board, Quantity: boards},
		artifacts.Artifacts{A: brick, Quantity: bricks},
		artifacts.Artifacts{A: thatch, Quantity: thatches},
		artifacts.Artifacts{A: tile, Quantity: tiles},
	}
}

func (b BuildingPlan) IsComplete() bool {
	for i := 0; i < BuildingBaseMaxSize-1; i++ {
		for j := 0; j < BuildingBaseMaxSize-1; j++ {
			if b.BaseShape[i][j] != nil {
				return true
			}
		}
	}
	return false
}

func (b BuildingPlan) HasUnit(x, y, z uint8) bool {
	if x >= BuildingBaseMaxSize || y >= BuildingBaseMaxSize {
		return false
	}
	if b.BaseShape[x][y] == nil {
		return false
	}
	return len(b.BaseShape[x][y].Floors) > int(z)
}

func (b BuildingPlan) HasUnitOrRoof(x, y, z uint8) bool {
	if x >= BuildingBaseMaxSize || y >= BuildingBaseMaxSize {
		return false
	}
	if b.BaseShape[x][y] == nil {
		return false
	}
	var oz = len(b.BaseShape[x][y].Floors)
	if !b.BaseShape[x][y].Roof.Flat() {
		oz++
	}
	return oz > int(z)
}

func (b *BuildingPlan) Copy() *BuildingPlan {
	var bs [BuildingBaseMaxSize][BuildingBaseMaxSize]*PlanUnits
	for i := 0; i < BuildingBaseMaxSize-1; i++ {
		for j := 0; j < BuildingBaseMaxSize-1; j++ {
			if b.BaseShape[i][j] != nil {
				roof := *b.BaseShape[i][j].Roof
				bs[i][j] = &PlanUnits{
					Floors: b.BaseShape[i][j].Floors,
					Roof:   &roof,
				}
			}
		}
	}
	return &BuildingPlan{
		BaseShape:    bs,
		BuildingType: b.BuildingType,
	}
}
