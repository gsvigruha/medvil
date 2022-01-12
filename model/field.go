package model

import (
	"medvil/model/social"
	"medvil/model/terrain"
)

// Implements navigation.IField
type Field struct {
	NE uint8
	SE uint8
	SW uint8
	NW uint8

	Terrain  terrain.Terrain
	Building FieldBuildingObjects
	Plant    *terrain.Plant
	Farm     *social.Farm
}

func (f Field) Walkable() bool {
	return f.Terrain.T.Walkable && ((f.NE == f.NW && f.SE == f.SW) || (f.NE == f.SE && f.NW == f.SW))
}
