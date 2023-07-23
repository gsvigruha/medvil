package maps

import (
	//"math"
	//"fmt"
	"math"
	"math/rand"
	"medvil/model"
	"medvil/model/building"
	"medvil/model/navigation"
	"medvil/model/social"
	"medvil/model/terrain"
	"medvil/model/time"
	//"medvil/util"
	"fmt"
)

func setupTerrain(fields [][]*navigation.Field, sizeX uint16, sizeY uint16) {
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

	sx, sy := int(sizeX), int(sizeY)
	for k := 0; k < 5; k++ {
		peak := rand.Intn(int((sizeX + sizeY) / 5))
		x, y := rand.Intn(sx), rand.Intn(sy)
		for l := 0; l < 10; l++ {
			x, y := x+rand.Intn(sx/2)-sx/4, y+rand.Intn(sy/2)-sy/4
			peak := peak + rand.Intn(10) - 5
			slope := rand.Float64() + 1.0
			fmt.Println(x, y, peak, slope)
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

func NewMap(sizeX uint16, sizeY uint16) *model.Map {
	fields := make([][]*navigation.Field, sizeX)
	for i := range fields {
		fields[i] = make([]*navigation.Field, sizeY)
	}
	setupTerrain(fields, sizeX, sizeY)
	m := &model.Map{SX: sizeX, SY: sizeY, Fields: fields}
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
