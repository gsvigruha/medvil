package maps

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"medvil/model"
	"medvil/model/artifacts"
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
					T:             terrain.AllCropTypes[0],
					X:             uint16(i),
					Y:             uint16(j),
					BirthDateDays: uint32(1000*12*30 - rand.Intn(20*12*30)),
					Shape:         uint8(rand.Intn(10)),
				}
			case "O":
				m.Fields[i][j].Plant = &terrain.Plant{
					T:             terrain.AllTreeTypes[0],
					X:             uint16(i),
					Y:             uint16(j),
					BirthDateDays: uint32(1000*12*30 - rand.Intn(20*12*30)),
					Shape:         uint8(rand.Intn(10)),
				}
			case "A":
				m.Fields[i][j].Plant = &terrain.Plant{
					T:             terrain.AllTreeTypes[1],
					X:             uint16(i),
					Y:             uint16(j),
					BirthDateDays: uint32(1000*12*30 - rand.Intn(20*12*30)),
					Shape:         uint8(rand.Intn(10)),
				}
			case "R":
				m.Fields[i][j].Plant = &terrain.Plant{
					T:             terrain.AllCropTypes[2],
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
			m.Fields[i][j] = &navigation.Field{}
			switch fieldTypes[i*2+1 : i*2+2] {
			case "G":
				m.Fields[i][j].Terrain.T = terrain.Grass
			case "S":
				m.Fields[i][j].Terrain.T = terrain.Sand
			case "R":
				m.Fields[i][j].Terrain.T = terrain.Rock
			case "I":
				m.Fields[i][j].Terrain.T = terrain.IronBog
			case "D":
				m.Fields[i][j].Terrain.T = terrain.Dirt
			case "M":
				m.Fields[i][j].Terrain.T = terrain.Mud
			case "W":
				m.Fields[i][j].Terrain.T = terrain.Water
				m.Fields[i][j].Terrain.Resources.Add(artifacts.GetArtifact("water"), artifacts.InfiniteQuantity)
			case "N":
				m.Fields[i][j].Terrain.T = terrain.Gold
			}
			if m.Fields[i][j].Terrain.T == terrain.Grass {
				m.Fields[i][j].Terrain.Shape = uint8(rand.Intn(4))
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
	for i := uint16(0); i < 5; i++ {
		for j := uint16(0); j < 5; j++ {
			bx := int(b.X+i) - 2
			by := int(b.Y+j) - 2
			if b.Plan.BaseShape[i][j] != nil {
				m.Fields[bx][by].Building.BuildingComponents = b.ToBuildingUnits(uint8(i), uint8(j), false)
				m.Fields[bx][by].Plant = nil
			}
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

	var cm map[string][]*social.Country
	if err := json.Unmarshal(byteValue, &cm); err != nil {
		fmt.Println(err)
	}
	countries := cm["countries"]
	for i := range countries {
		country := countries[i]
		for j := range countries[i].Towns {
			town := countries[i].Towns[j]
			town.Country = country
			town.Init(0)
			town.Townhall.Household.People = make([]*social.Person, 5)
			town.Townhall.Household.TargetNumPeople = 5
			town.Townhall.Household.Town = town
			town.Townhall.Household.Resources.VolumeCapacity = town.Townhall.Household.Building.Plan.Area() * social.StoragePerArea
			town.Townhall.Household.Building.Plan.BuildingType = building.BuildingTypeTownhall
			for i := range town.Townhall.Household.People {
				town.Townhall.Household.People[i] = town.Townhall.Household.NewPerson(m)
			}
			AddBuilding(town.Townhall.Household.Building, m)

			town.Marketplace.Town = town
			town.Marketplace.Building.Plan.BuildingType = building.BuildingTypeMarket
			AddBuilding(town.Marketplace.Building, m)
			town.Marketplace.Init()
			for k := range town.Farms {
				farm := town.Farms[k]
				farm.Household.Town = town
				farm.Household.Resources.VolumeCapacity = farm.Household.Building.Plan.Area() * social.StoragePerArea
				AddBuilding(farm.Household.Building, m)
				for l := range farm.Land {
					land := farm.Land[l]
					farm.Land[l].F = m.Fields[land.X][land.Y]
					farm.Land[l].F.Allocated = true
				}
				farm.Household.People = []*social.Person{farm.Household.NewPerson(m), farm.Household.NewPerson(m)}
				farm.Household.TargetNumPeople = 2
				farm.Household.Building.Plan.BuildingType = building.BuildingTypeFarm
				for _, p := range farm.Household.People {
					m.GetField(p.Traveller.FX, p.Traveller.FY).RegisterTraveller(p.Traveller)
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

	fields := make([][]*navigation.Field, sx)
	for i := range fields {
		fields[i] = make([]*navigation.Field, sy)
	}
	m := model.Map{SX: sx, SY: sy, Fields: fields}

	LoadFields(dir, &m)
	LoadPlants(dir, &m)
	LoadSociety(dir, &m)

	return m
}
