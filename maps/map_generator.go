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

type MapConfig struct {
	SizeX, SizeY int
	NumHills     int
	HillSize     int
}

func setupTerrain(fields [][]*navigation.Field, config MapConfig) {
	for i := range fields {
		for j := range fields[i] {
			fields[i][j] = &navigation.Field{X: uint16(i), Y: uint16(j)}
			fields[i][j].Terrain.T = terrain.Grass
			fields[i][j].Terrain.Shape = uint8(rand.Intn(4))
			if fields[i][j].Terrain.T == terrain.Grass && rand.Intn(5) == 0 {
				fields[i][j].Plant = &terrain.Plant{
					T:             terrain.AllTreeTypes[rand.Intn(2)],
					X:             uint16(i),
					Y:             uint16(j),
					BirthDateDays: uint32(1000*12*30 - rand.Intn(20*12*30)),
					Shape:         uint8(rand.Intn(10)),
				}
			}
		}
	}

	sx, sy := config.SizeX, config.SizeY
	for k := 0; k < config.NumHills; k++ {
		peak := rand.Intn((sx + sy) / 5)
		x, y := rand.Intn(sx), rand.Intn(sy)
		for l := 0; l < config.HillSize; l++ {
			x, y := x+rand.Intn(sx/3)-sx/6, y+rand.Intn(sy/3)-sy/6
			peak := peak + rand.Intn(10) - 5
			slope := rand.Float64() + 1.0
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
		}
	}
}

func NewMap(config MapConfig) *model.Map {
	fields := make([][]*navigation.Field, config.SizeX)
	for i := range fields {
		fields[i] = make([]*navigation.Field, config.SizeY)
	}
	setupTerrain(fields, config)
	m := &model.Map{SX: uint16(config.SizeX), SY: uint16(config.SizeY), Fields: fields}
	calendar := &time.CalendarType{
		Year:  1000,
		Month: 1,
		Day:   1,
		Hour:  0,
	}
	m.Calendar = calendar
	townhall := &building.Building{
		Plan: building.BuildingPlanFromJSON("samples/building/townhouse_1.building.json"),
		X:    m.SX / 2,
		Y:    m.SY / 2,
	}
	AddBuilding(townhall, m)
	marketplace := &building.Building{
		Plan: building.BuildingPlanFromJSON("samples/building/marketplace_1.building.json"),
		X:    m.SX/2 + 5,
		Y:    m.SY / 2,
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
