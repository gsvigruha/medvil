package model

import (
	"medvil/model/terrain"
)

type Field struct {
	NE uint8
	SE uint8
	SW uint8
	NW uint8

	Terrain  terrain.Terrain
	Building FieldBuildingObjects
}

func (f Field) Walkable() bool {
	return f.Terrain.T.Walkable && ((f.NE == f.NW && f.SE == f.SW) || (f.NE == f.SE && f.NW == f.SW))
}
