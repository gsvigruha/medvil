package maps

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"medvil/model"
	"medvil/model/building"
	"medvil/model/navigation"
	"medvil/model/social"
	"medvil/model/terrain"
	"os"
	"strconv"
)

func LoadPlants(dir string, m *model.Map) {
	file, err := os.Open(dir + "/plants.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for j := uint16(0); j < m.SY; j++ {
		scanner.Scan()
		plants := scanner.Text()
		for i := uint16(0); i < m.SX; i++ {
			switch plants[i+0 : i+1] {
			case "G":
				m.Fields[i][j].Plant = &terrain.Plant{
					T:             &terrain.AllCropTypes[0],
					X:             uint16(i),
					Y:             uint16(j),
					BirthDateDays: uint32(1000*12*30 - rand.Intn(20*12*30)),
					Shape:         uint8(rand.Intn(10)),
				}

			case "O":
				m.Fields[i][j].Plant = &terrain.Plant{
					T:             &terrain.AllTreeTypes[0],
					X:             uint16(i),
					Y:             uint16(j),
					BirthDateDays: uint32(1000*12*30 - rand.Intn(20*12*30)),
					Shape:         uint8(rand.Intn(10)),
				}
			case "A":
				m.Fields[i][j].Plant = &terrain.Plant{
					T:             &terrain.AllTreeTypes[1],
					X:             uint16(i),
					Y:             uint16(j),
					BirthDateDays: uint32(1000*12*30 - rand.Intn(20*12*30)),
					Shape:         uint8(rand.Intn(10)),
				}
			}
		}
	}
}

func LoadFields(dir string, m *model.Map) {
	file, err := os.Open(dir + "/fields.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for j := uint16(0); j < m.SY; j++ {
		scanner.Scan()
		heights := scanner.Text()
		scanner.Scan()
		fieldTypes := scanner.Text()
		for i := uint16(0); i < m.SX; i++ {
			switch fieldTypes[i*2+1 : i*2+2] {
			case "G":
				m.Fields[i][j].Terrain.T = terrain.Grass
			case "S":
				m.Fields[i][j].Terrain.T = terrain.Sand
			case "R":
				m.Fields[i][j].Terrain.T = terrain.Rock
			case "D":
				m.Fields[i][j].Terrain.T = terrain.Dirt
			case "W":
				m.Fields[i][j].Terrain.T = terrain.Water
			}
			h, err := strconv.Atoi(heights[i*2 : i*2+1])
			if err != nil {
				fmt.Println(err)
			}
			if i > 0 && j > 0 {
				m.Fields[i-1][j-1].SE = uint8(h)
			}
			if j > 0 {
				m.Fields[i][j-1].SW = uint8(h)
			}
			if i > 0 {
				m.Fields[i-1][j].NE = uint8(h)
			}
			m.Fields[i][j].NW = uint8(h)
			m.Fields[i][j].X = i
			m.Fields[i][j].Y = j
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func AddBuilding(b *building.Building, m *model.Map) {
	for bx := uint16(0); bx < 5; bx++ {
		for by := uint16(0); by < 5; by++ {
			m.Fields[b.X+bx][b.Y+by].Building.BuildingUnits = b.ToBuildingUnits(uint8(bx), uint8(by))
			m.Fields[b.X+bx][b.Y+by].Building.RoofUnit = b.GetRoof(uint8(bx), uint8(by))
		}
	}
}

func LoadSociety(dir string, m *model.Map) {
	jsonFile, err := os.Open(dir + "/society.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}

	var cm map[string][]social.Country
	if err := json.Unmarshal(byteValue, &cm); err != nil {
		fmt.Println(err)
	}
	countries := cm["countries"]
	for i := range countries {
		country := countries[i]
		for j := range countries[i].Towns {
			town := countries[i].Towns[j]
			town.Country = &country
			town.Townhall.Household.People = make([]*social.Person, 10)
			town.Townhall.Household.TargetNumPeople = 10
			for i := range town.Townhall.Household.People {
				town.Townhall.Household.People[i] = town.Townhall.Household.NewPerson()
			}
			AddBuilding(town.Townhall.Household.Building, m)
			for k := range town.Farms {
				farm := town.Farms[k]
				farm.Household.Town = town
				AddBuilding(farm.Household.Building, m)
				for l := range farm.Land {
					land := farm.Land[l]
					farm.Land[l].F = &m.Fields[land.X][land.Y]
					farm.Land[l].F.Allocated = true
				}
				farm.Household.People = []*social.Person{farm.Household.NewPerson(), farm.Household.NewPerson()}
				farm.Household.TargetNumPeople = 2
				for _, p := range farm.Household.People {
					m.GetField(farm.Household.Building.X, farm.Household.Building.Y).RegisterTraveller(p.Traveller)
				}
			}
		}
	}
	m.Countries = countries
}

func LoadMap(dir string) model.Map {
	jsonFile, err := os.Open(dir + "/meta.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}

	var meta map[string]uint16
	if err := json.Unmarshal(byteValue, &meta); err != nil {
		fmt.Println(err)
	}
	sx := meta["SX"]
	sy := meta["SY"]

	fields := make([][]navigation.Field, sx)
	for i := range fields {
		fields[i] = make([]navigation.Field, sy)
	}
	m := model.Map{SX: sx, SY: sy, Fields: fields}

	LoadFields(dir, &m)
	LoadPlants(dir, &m)
	LoadSociety(dir, &m)
	return m
}
