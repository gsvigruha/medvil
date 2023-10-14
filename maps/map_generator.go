package maps

import (
	"math"
	"math/rand"
	"medvil/model"
	"medvil/model/artifacts"
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
			fields[i][j].Terrain.Resources = &artifacts.Resources{}
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
		angle := rand.Float64() * math.Pi * 2
		GenerateHills(x, y, peak, 1, angle, fields)
	}

	for k := 0; k < area*config.Lakes/LakeAreaRatio; k++ {
		size := float64(rand.Intn(20))
		x, y := rand.Intn(sx), rand.Intn(sy)
		GenerateLakes(x, y, size, 1, fields, terrain.Water)
	}

	for i := range fields {
		for j := range fields[i] {
			navigation.SetSurroundingTypes(m, fields[i][j])
			if !m.Shore(uint16(i), uint16(j)) {
				if rand.Intn(ResourcesProb) < config.Resources && fields[i][j].Plant == nil && fields[i][j].Terrain.T == terrain.Grass {
					if fields[i][j].Flat() {
						if rand.Float64() < float64(i)/float64(config.Size) {
							fields[i][j].Deposit = &terrain.Deposit{T: terrain.Mud, Q: artifacts.InfiniteQuantity}
						} else {
							fields[i][j].Deposit = &terrain.Deposit{T: terrain.Rock, Q: artifacts.InfiniteQuantity}
						}
					} else {
						if rand.Float64() < float64(j)/float64(config.Size) {
							fields[i][j].Deposit = &terrain.Deposit{T: terrain.Rock, Q: artifacts.InfiniteQuantity}
						} else {
							fields[i][j].Deposit = &terrain.Deposit{T: terrain.IronBog, Q: uint16((rand.Intn(5) + 1) * 1000)}
						}
					}
				}
			} else {
				if rand.Intn(ResourcesProb) < config.Resources*2 {
					fields[i][j].Deposit = &terrain.Deposit{T: terrain.Gold, Q: uint16((rand.Intn(5) + 1) * 1000)}
				}
			}
			if fields[i][j].Terrain.T == terrain.Water {
				fields[i][j].Terrain.Resources.Add(artifacts.GetArtifact("water"), artifacts.InfiniteQuantity)
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

func GenerateHills(x, y, peak, n int, angle float64, fields [][]*navigation.Field) {
	slope := rand.Float64()*2.0 + 1.0
	rad := float64(peak) / slope
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
	if n < MaxIter && peak > 4 {
		radX := int(math.Abs(rad*math.Cos(angle))) + 1
		radY := int(math.Abs(rad*math.Sin(angle))) + 1
		for l := 0; l < HillBranching; l++ {
			nx := x + rand.Intn(radX*2) - radX
			ny := y + rand.Intn(radY*2) - radY
			npeak := rand.Intn(peak/2) + peak/2
			GenerateHills(nx, ny, npeak, n+1, angle, fields)
		}
	}
}

func findStartingLocation(m *model.Map) (int, int) {
	var x, y = 0, 0
	var maxScore = 0
	for i := range m.Fields {
		for j := range m.Fields[i] {
			if i < 5 || j < 5 || i > len(m.Fields)-5 || j > len(m.Fields[i])-5 {
				continue
			}
			dx := float64(int(m.SX/2) - i)
			dy := float64(int(m.SY/2) - j)
			var score = int(m.SX+m.SY)/4 - int(math.Sqrt(dx*dx+dy*dy))
			var suitable = true
			var water = false
			var rock = 0
			var gold = 0
			var iron = 0
			var mud = 0
			for di := -20; di <= 20; di++ {
				for dj := -20; dj <= 20; dj++ {
					if i+dj >= 0 && j+dj >= 0 {
						f := m.GetField(uint16(i+di), uint16(j+dj))
						if dj >= -6 && dj <= 6 && di >= -6 && di <= 6 {
							if f == nil || !f.Flat() || f.Terrain.T != terrain.Grass {
								suitable = false
							}
						}
						if f != nil {
							if f.Building.GetBuilding() != nil {
								suitable = false
							}
							if dj >= -10 && dj <= 10 && di >= -10 && di <= 10 {
								if f.Terrain.T == terrain.Water {
									water = true
								}
							}
							if f.Deposit != nil {
								if f.Deposit.T == terrain.Rock {
									rock++
								} else if f.Deposit.T == terrain.Gold {
									gold++
								} else if f.Deposit.T == terrain.IronBog {
									iron++
								} else if f.Deposit.T == terrain.Mud {
									mud++
								}
							}
						}
					}
					if !suitable {
						break
					}
				}
				if !suitable {
					break
				}
			}
			score += 10*int(math.Log2(float64(rock+1))) + 10*int(math.Log2(float64(gold+1))) + 10*int(math.Log2(float64(iron+1))) + 10*int(math.Log2(float64(mud+1)))
			if water && suitable && score > maxScore {
				maxScore = score
				x = i
				y = j
			}
		}
	}
	return x, y
}

func NewMap(config MapConfig) *model.Map {
	for {
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

		var success = true
		success = success && GenerateCountry(social.CountryTypePlayer, m)
		for i := 0; i < (config.Size-50)/50; i++ {
			success = success && GenerateCountry(social.CountryTypeOutlaw, m)
		}

		if !success {
			continue
		}

		return m
	}
}
