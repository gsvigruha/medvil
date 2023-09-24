package terrain

import (
	"bytes"
	"encoding/json"
	"medvil/model/artifacts"
)

type TerrainType struct {
	Walkable  bool
	Arable    bool
	Pasture   bool
	Water     bool
	Buildable bool
	Name      string
}

func (tt *TerrainType) MarshalJSON() ([]byte, error) {
	return json.Marshal(tt.Name)
}

var Water = &TerrainType{Walkable: false, Arable: false, Pasture: false, Water: true, Buildable: false, Name: "water"}
var Grass = &TerrainType{Walkable: true, Arable: true, Pasture: true, Water: false, Buildable: true, Name: "grass"}
var Sand = &TerrainType{Walkable: true, Arable: false, Pasture: false, Water: false, Buildable: false, Name: "sand"}
var Dirt = &TerrainType{Walkable: true, Arable: true, Pasture: false, Water: false, Buildable: true, Name: "dirt"}
var Canal = &TerrainType{Walkable: true, Arable: false, Pasture: false, Water: false, Buildable: false, Name: "canal"}

var TerrainTypes = [...]*TerrainType{
	Water,
	Grass,
	Sand,
	Dirt,
	Canal,
}

func GetTerrainType(name string) *TerrainType {
	for _, t := range TerrainTypes {
		if t.Name == name {
			return t
		}
	}
	return nil
}

type Terrain struct {
	T         *TerrainType
	Resources *artifacts.Resources
	Shape     uint8
}

func (t *Terrain) UnmarshalJSON(data []byte) error {
	var j map[string]json.RawMessage
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	if err := json.Unmarshal(j["Resources"], &t.Resources); err != nil {
		return err
	}
	s := bytes.NewBuffer(j["T"]).String()
	switch s {
	case "Water":
		t.T = Water
	case "Grass":
		t.T = Grass
	case "Sand":
		t.T = Sand
	case "Dirt":
		t.T = Dirt
	case "Canal":
		t.T = Canal
	}
	return nil
}
