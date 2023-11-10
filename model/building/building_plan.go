package building

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
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

func RandomBuildingDir() uint8 {
	return uint8(rand.Intn(4))
}

func OppDir(dir uint8) uint8 {
	return uint8((dir + 2) % 4)
}

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
	Floors    []Floor
	Roof      *Roof
	Extension *BuildingExtension
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
						m := materials.GetMaterial(shape[i][j][k])
						b.BaseShape[i][j].Roof = &Roof{M: m, RoofType: GetRoofType(m)}
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
			if b.BaseShape[i][j] != nil && b.BaseShape[i][j].Roof != nil && !b.BaseShape[i][j].Roof.Flat() {
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

var Cube = artifacts.GetArtifact("cube")
var Board = artifacts.GetArtifact("board")
var Brick = artifacts.GetArtifact("brick")
var Thatch = artifacts.GetArtifact("thatch")
var Tile = artifacts.GetArtifact("tile")
var Textile = artifacts.GetArtifact("textile")
var Paper = artifacts.GetArtifact("paper")

func (b BuildingPlan) RepairCost() []artifacts.Artifacts {
	cc := b.ConstructionCost()
	var rc = make([]artifacts.Artifacts, len(cc))
	for i, as := range cc {
		rc[i] = artifacts.Artifacts{A: as.A, Quantity: uint16(math.Ceil(float64(as.Quantity) / 2))}
	}
	return artifacts.Filter(rc)
}

func (b BuildingPlan) ConstructionCost() []artifacts.Artifacts {
	var cubes uint16 = 0
	var boards uint16 = 0
	var bricks uint16 = 0
	var thatches uint16 = 0
	var tiles uint16 = 0
	var textiles uint16 = 0
	var papers uint16 = 0
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
					case materials.GetMaterial("marble"):
						boards += 1
						cubes += 2
					case materials.GetMaterial("stone"):
						boards += 1
						cubes += 2
					case materials.GetMaterial("brick"):
						boards += 1
						bricks += 2
					case materials.GetMaterial("whitewash"):
						boards += 1
						bricks += 2
					}
					if b.BuildingType == BuildingTypeGate {
						papers += 2
					}
				}
				if b.BaseShape[i][j].Roof != nil && !b.BaseShape[i][j].Roof.Flat() {
					switch b.BaseShape[i][j].Roof.M {
					case materials.GetMaterial("tile"):
						tiles += 1
						boards += 1
					case materials.GetMaterial("reed"):
						thatches += 1
						boards += 1
					case materials.GetMaterial("stone"):
						cubes += 1
					case materials.GetMaterial("textile"):
						textiles += 1
					case materials.GetMaterial("copper"):
						tiles += 1
						boards += 1
					}
				}
				if b.BaseShape[i][j].Extension != nil {
					switch b.BaseShape[i][j].Extension.T {
					case WaterMillWheel:
						boards += 2
					case Forge:
						cubes += 2
						tiles += 1
					case Kiln:
						bricks += 2
						tiles += 1
					case Cooker:
						bricks += 1
					case Workshop:
						boards += 1
					}
				}
			}
		}
	}

	return artifacts.Filter([]artifacts.Artifacts{
		artifacts.Artifacts{A: Cube, Quantity: cubes},
		artifacts.Artifacts{A: Board, Quantity: boards},
		artifacts.Artifacts{A: Brick, Quantity: bricks},
		artifacts.Artifacts{A: Thatch, Quantity: thatches},
		artifacts.Artifacts{A: Tile, Quantity: tiles},
		artifacts.Artifacts{A: Textile, Quantity: textiles},
		artifacts.Artifacts{A: Paper, Quantity: papers},
	})
}

func (b BuildingPlan) IsComplete() bool {
	for i := 0; i < BuildingBaseMaxSize; i++ {
		for j := 0; j < BuildingBaseMaxSize; j++ {
			if b.BaseShape[i][j] != nil {
				return true
			}
		}
	}
	return false
}

func (b BuildingPlan) HasNeighborUnit(x, y, z uint8) bool {
	return b.HasUnit(x, y-1, z) || b.HasUnit(x+1, y, z) || b.HasUnit(x, y+1, z) || b.HasUnit(x-1, y, z)
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
	if b.BaseShape[x][y].Extension != nil && !b.BaseShape[x][y].Extension.T.InUnit {
		return z == 0
	}
	var oz = len(b.BaseShape[x][y].Floors)
	if b.BaseShape[x][y].Roof != nil && !b.BaseShape[x][y].Roof.Flat() {
		oz++
	}
	return oz > int(z)
}

func (b *BuildingPlan) Copy() *BuildingPlan {
	var bs [BuildingBaseMaxSize][BuildingBaseMaxSize]*PlanUnits
	for i := 0; i < BuildingBaseMaxSize; i++ {
		for j := 0; j < BuildingBaseMaxSize; j++ {
			if b.BaseShape[i][j] != nil {
				var extension *BuildingExtension = nil
				var roof *Roof = nil
				if b.BaseShape[i][j].Extension != nil {
					e := *b.BaseShape[i][j].Extension
					extension = &e
				}
				if b.BaseShape[i][j].Roof != nil {
					r := *b.BaseShape[i][j].Roof
					roof = &r
				}
				bs[i][j] = &PlanUnits{
					Floors:    b.BaseShape[i][j].Floors,
					Roof:      roof,
					Extension: extension,
				}
			}
		}
	}
	return &BuildingPlan{
		BaseShape:    bs,
		BuildingType: b.BuildingType,
	}
}

func (b *BuildingPlan) GetExtension(et *BuildingExtensionType) *BuildingExtension {
	es := b.GetExtensionsWithCoords(et)
	if len(es) > 0 {
		return es[0].E
	}
	return nil
}

type ExtensionWithCoords struct {
	E *BuildingExtension
	X uint16
	Y uint16
}

func (b *BuildingPlan) GetExtensionsWithCoords(et *BuildingExtensionType) []ExtensionWithCoords {
	var result []ExtensionWithCoords
	for i := uint16(0); i < BuildingBaseMaxSize; i++ {
		for j := uint16(0); j < BuildingBaseMaxSize; j++ {
			if b.BaseShape[i][j] != nil && b.BaseShape[i][j].Extension != nil && b.BaseShape[i][j].Extension.T == et {
				result = append(result, ExtensionWithCoords{E: b.BaseShape[i][j].Extension, X: i, Y: j})
			}
		}
	}
	return result
}

func (b *BuildingPlan) GetExtensions() []*BuildingExtension {
	var es []*BuildingExtension
	for i := uint16(0); i < BuildingBaseMaxSize; i++ {
		for j := uint16(0); j < BuildingBaseMaxSize; j++ {
			if b.BaseShape[i][j] != nil && b.BaseShape[i][j].Extension != nil {
				es = append(es, b.BaseShape[i][j].Extension)
			}
		}
	}
	return es
}

func GetRoofType(m *materials.Material) RoofType {
	if m == materials.GetMaterial("textile") {
		return RoofTypeFlat
	}
	return RoofTypeSplit
}
