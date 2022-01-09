package maps

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"medvil/model"
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
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
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

	fields := make([][]model.Field, sx)
	for i := range fields {
		fields[i] = make([]model.Field, sy)
	}
	m := model.Map{SX: sx, SY: sy, Fields: fields}

	LoadFields(dir, &m)
	LoadPlants(dir, &m)
	return m
}
