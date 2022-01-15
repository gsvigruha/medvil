package navigation

import (
	"medvil/model/terrain"
	//"fmt"
)

type Field struct {
	NE uint8
	SE uint8
	SW uint8
	NW uint8

	Terrain    terrain.Terrain
	Building   FieldBuildingObjects
	Plant      *terrain.Plant
	Travellers []*Traveller
}

func (f Field) Walkable() bool {
	return f.Terrain.T.Walkable && ((f.NE == f.NW && f.SE == f.SW) || (f.NE == f.SE && f.NW == f.SW))
}

func (f *Field) RegisterTraveller(t *Traveller) {
	f.Travellers = append(f.Travellers, t)
}

func (f *Field) UnregisterTraveller(t *Traveller) {
	for i := range f.Travellers {
		if f.Travellers[i] == t {
			f.Travellers = append(f.Travellers[0:i], f.Travellers[i+1:]...)
			return
		}
	}
}
