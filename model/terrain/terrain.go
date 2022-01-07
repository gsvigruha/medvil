package terrain

import (
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

var Water = TerrainType{Walkable: false, Arable: false, Pasture: false, Water: true, Buildable: false, Name: "water"}
var Grass = TerrainType{Walkable: true, Arable: true, Pasture: true, Water: false, Buildable: true, Name: "grass"}
var Sand = TerrainType{Walkable: true, Arable: false, Pasture: false, Water: false, Buildable: false, Name: "sand"}
var Dirt = TerrainType{Walkable: true, Arable: true, Pasture: false, Water: false, Buildable: true, Name: "dirt"}
var Rock = TerrainType{Walkable: true, Arable: false, Pasture: false, Water: false, Buildable: true, Name: "rock"}

type Terrain struct {
	T         TerrainType
	Resources []artifacts.RawResource
}
