package maps

import (
	"math"
	"math/rand"
	"medvil/model"
	"medvil/model/building"
	"medvil/model/navigation"
	"medvil/model/social"
	"medvil/model/terrain"
	"medvil/model/time"
)

const HillAreaRatio = 10000
const LakeAreaRatio = 10000
const LakeLength = 150
const MaxIter = 6
const HillBranching = 2
const LakeBranching = 4
const TreeProb = 30
const ResourcesProb = 1000

type MapConfig struct {
	Size      int
	Hills     int
	Lakes     int
	Trees     int
	Resources int
}

func setupTerrain(m *model.Map, config MapConfig) {
	fields := m.Fields
	for i := range fields {
		for j := range fields[i] {
			fields[i][j] = &navigation.Field{X: uint16(i), Y: uint16(j)}
			fields[i][j].Terrain.T = terrain.Grass
			fields[i][j].Terrain.Shape = uint8(rand.Intn(4))
			if fields[i][j].Terrain.T == terrain.Grass && rand.Intn(TreeProb) < config.Trees {
				fields[i][j].Plant = &terrain.Plant{
					T:             terrain.AllTreeTypes[rand.Intn(2)],
					X:             uint16(i),
					Y:             uint16(j),
					BirthDateDays: uint32(1000*12*30 - rand.Intn(20*12*30)),
					Shape:         uint8(rand.Intn(terrain.TreeNumShapes)),
				}
			}
		}
	}

	sx, sy := config.Size, config.Size
	area := sx * sy
	for k := 0; k < area*config.Hills/HillAreaRatio; k++ {
		x, y := rand.Intn(sx), rand.Intn(sy)
		peak := rand.Intn(30) + 10
		GenerateHills(x, y, peak, 1, fields)
	}

	for k := 0; k < area*config.Lakes/LakeAreaRatio; k++ {
		size := float64(rand.Intn(20))
		x, y := rand.Intn(sx), rand.Intn(sy)
		GenerateLakes(x, y, size, 1, fields, terrain.Water)
	}

	for i := range fields {
		for j := range fields[i] {
			for k, _ := range navigation.DirectionOrthogonalXY {
				di1 := navigation.DirectionOrthogonalXY[k][0]
				dj1 := navigation.DirectionOrthogonalXY[k][1]
				di2 := navigation.DirectionOrthogonalXY[(k+1)%4][0]
				dj2 := navigation.DirectionOrthogonalXY[(k+1)%4][1]
				di3 := navigation.DirectionDiagonalXY[k][0]
				dj3 := navigation.DirectionDiagonalXY[k][1]
				if i > 0 && j > 0 && i < sx-1 && j < sy-1 {
					fields[i][j].Surroundings[k] = GetSurroundingType(fields[i][j], fields[i+di1][j+dj1], fields[i+di2][j+dj2], fields[i+di3][j+dj3])
				}
			}
			if rand.Intn(ResourcesProb) < config.Resources && fields[i][j].Plant == nil {
				if !m.Shore(uint16(i), uint16(j)) {
					if fields[i][j].Terrain.T == terrain.Grass && fields[i][j].Flat() {
						fields[i][j].Terrain.T = terrain.Mud
					} else if fields[i][j].Terrain.T == terrain.Grass && !fields[i][j].Flat() {
						if rand.Float64() < 0.5 {
							fields[i][j].Terrain.T = terrain.Rock
						} else {
							fields[i][j].Terrain.T = terrain.IronBog
						}
					}
				} else {
					fields[i][j].Terrain.T = terrain.Gold
				}
			}
		}
	}
}

func GenerateLakes(x, y int, size float64, n int, fields [][]*navigation.Field, t *terrain.TerrainType) {
	sint := int(size)
	if sint < 6 {
		return
	}
	for i := range fields {
		for j := range fields[i] {
			dist := math.Sqrt(float64((x-i)*(x-i)) + float64((y-j)*(y-j)))
			if dist < size && fields[i][j].NW == 0 && fields[i][j].SW == 0 && fields[i][j].NE == 0 && fields[i][j].SE == 0 {
				fields[i][j].Terrain.T = t
				fields[i][j].Plant = nil
			}
		}
	}

	for l := 0; l < LakeBranching; l++ {
		nx, ny := x+rand.Intn(sint)-sint/2, y+rand.Intn(sint)-sint/2
		nsize := (0.9 + rand.Float64()/10) * size
		if n < MaxIter && nsize > 0 && rand.Float64() < 0.75 {
			GenerateLakes(nx, ny, nsize, n+1, fields, t)
		}
	}
}

func GenerateHills(x, y, peak, n int, fields [][]*navigation.Field) {
	slope := rand.Float64()*2.0 + 1.0
	rad := int(float64(peak) / slope)
	for i := range fields {
		for j := range fields[i] {
			h0 := float64(peak) - slope*math.Sqrt(float64((x-i)*(x-i))+float64((y-j)*(y-j)))
			h := uint8(math.Max(math.Max(h0, 0), float64(fields[i][j].NW)))
			if i > 0 && j > 0 {
				fields[i-1][j-1].SE = uint8(h)
			}
			if j > 0 {
				fields[i][j-1].SW = uint8(h)
			}
			if i > 0 {
				fields[i-1][j].NE = uint8(h)
			}
			fields[i][j].NW = uint8(h)
		}
	}
	if n < MaxIter && peak > 4 && rad > 4 {
		for l := 0; l < HillBranching; l++ {
			nx, ny := x+rand.Intn(rad/4)-rad/2, y+rand.Intn(rad/4)-rad/2
			npeak := rand.Intn(peak/2) + peak/2
			GenerateHills(nx, ny, npeak, n+1, fields)
		}
	}
}

func GetSurroundingType(f *navigation.Field, of1 *navigation.Field, of2 *navigation.Field, of3 *navigation.Field) uint8 {
	if f.Terrain.T == terrain.Grass && of1.Terrain.T == terrain.Water && of2.Terrain.T == terrain.Water && of3.Terrain.T == terrain.Water {
		return navigation.SurroundingWater
	} else if f.Terrain.T == terrain.Water && of1.Terrain.T == terrain.Grass && of2.Terrain.T == terrain.Grass && of3.Terrain.T == terrain.Grass {
		return navigation.SurroundingGrass
	} else if f.Terrain.T == terrain.Grass && of1.Terrain.T == terrain.Grass && of2.Terrain.T == terrain.Grass && of3.Terrain.T == terrain.Grass {
		if !f.DarkSlope() && of1.DarkSlope() && of2.DarkSlope() {
			return navigation.SurroundingDarkSlope
		}
	}
	return navigation.SurroundingSame
}

func findStartingLocation(m *model.Map) (int, int) {
	var x, y = 0, 0
	var maxScore = 0
	for i := range m.Fields {
		for j := range m.Fields[i] {
			dx := float64(int(m.SX/2) - i)
			dy := float64(int(m.SY/2) - j)
			var score = int(m.SX+m.SY)/2 - int(math.Sqrt(dx*dx+dy*dy))
			var suitable = true
			for di := -10; di <= 10; di++ {
				for dj := -10; dj <= 10; dj++ {
					if i+dj >= 0 && j+dj >= 0 {
						f := m.GetField(uint16(i+di), uint16(j+dj))
						if f != nil {
							if dj >= -5 && dj <= 5 && di >= -5 && di <= 5 {
								if !f.Flat() || f.Terrain.T != terrain.Grass {
									suitable = false
								}
							}
							if f.Terrain.T == terrain.Water {
								score++
							} else if f.Terrain.T == terrain.Rock {
								score++
							} else if f.Terrain.T == terrain.Gold {
								score++
							} else if f.Terrain.T == terrain.IronBog {
								score++
							} else if f.Terrain.T == terrain.Mud {
								score++
							}
						} else {
							suitable = false
						}
					}
				}
			}
			if suitable && score > maxScore {
				maxScore = score
				x = i
				y = j
			}
		}
	}
	return x, y
}

func NewMap(config MapConfig) *model.Map {
	fields := make([][]*navigation.Field, config.Size)
	for i := range fields {
		fields[i] = make([]*navigation.Field, config.Size)
	}
	m := &model.Map{SX: uint16(config.Size), SY: uint16(config.Size), Fields: fields}
	setupTerrain(m, config)
	calendar := &time.CalendarType{
		Year:  1000,
		Month: 1,
		Day:   1,
		Hour:  0,
	}
	m.Calendar = calendar

	tx, ty := findStartingLocation(m)

	townhall := &building.Building{
		Plan: building.BuildingPlanFromJSON("samples/building/townhouse_1.building.json"),
		X:    uint16(tx - 2),
		Y:    uint16(ty),
	}
	AddBuilding(townhall, m)
	marketplace := &building.Building{
		Plan: building.BuildingPlanFromJSON("samples/building/marketplace_1.building.json"),
		X:    uint16(tx + 2),
		Y:    uint16(ty),
	}
	AddBuilding(marketplace, m)

	m.Countries = []*social.Country{&social.Country{Towns: []*social.Town{&social.Town{}}}}
	town := m.Countries[0].Towns[0]
	town.Country = m.Countries[0]
	town.Townhall = &social.Townhall{Household: &social.Household{Building: townhall, Town: town}}
	town.Marketplace = &social.Marketplace{Building: marketplace, Town: town}
	town.Init()
	town.Marketplace.Init()
	return m
}
