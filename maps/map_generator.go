package maps

import (
	//"math"
	//"fmt"
	"math/rand"
	"medvil/model"
	"medvil/model/building"
	"medvil/model/terrain"
	"medvil/util"
)

func setupTerrain(fields [][]model.Field) {
	for i := range fields {
		for j := range fields[i] {
			switch rand.Intn(5) {
			case 0:
				fields[i][j].Terrain.T = terrain.Grass
			case 1:
				fields[i][j].Terrain.T = terrain.Sand
			case 2:
				fields[i][j].Terrain.T = terrain.Rock
			case 3:
				fields[i][j].Terrain.T = terrain.Dirt
			case 4:
				fields[i][j].Terrain.T = terrain.Water
			}
			if i > 0 && j > 0 {
				max := util.Max(int(fields[i-1][j-1].NE)+2, int(fields[i-1][j-1].SW)+2)
				min := util.Max(util.Min(int(fields[i-1][j-1].NE)-2, int(fields[i-1][j-1].SW)-2), 0)
				r := max - min
				var h = uint8(0)
				if r > 0 {
					//	h = uint8(rand.Intn(r) + min)
				}
				fields[i-1][j-1].SE = h
				fields[i][j-1].SW = h
				fields[i][j].NW = h
				fields[i-1][j].NE = h
			}
			if fields[i][j].Terrain.T == terrain.Grass && rand.Intn(2) == 0 {
				fields[i][j].Plant = &terrain.Plant{T: &terrain.AllPlantTypes[0], X: uint16(i), Y: uint16(j), Age: uint8(rand.Intn(20)), Shape: uint8(rand.Intn(10))}
			}
		}
	}
}

func addHouse(name string, x int, y int, m model.Map) {
	bp := building.BuildingPlanFromJSON(name)
	for bx := 0; bx < 5; bx++ {
		for by := 0; by < 5; by++ {
			m.Fields[x+bx][y+by].Building.BuildingUnits = bp.ToBuildingUnits(uint8(bx), uint8(by))
			m.Fields[x+bx][y+by].Building.RoofUnit = bp.GetRoof(uint8(bx), uint8(by))
		}
	}
}

func NewMap(sizeX uint16, sizeY uint16) model.Map {
	fields := make([][]model.Field, sizeX)
	for i := range fields {
		fields[i] = make([]model.Field, sizeY)
	}
	setupTerrain(fields)
	m := model.Map{SX: sizeX, SY: sizeY, Fields: fields}
	addHouse("samples/building/townhouse_1.building.json", 2, 2, m)
	addHouse("samples/building/rural_1.building.json", 6, 6, m)
	return m
}
