package navigation

import (
	"medvil/model/terrain"
	//"fmt"
)

type Destination interface {
	Check(*Field) bool
}

type FieldWithContext interface {
	Field() *Field
	Context() string
}

type Field struct {
	X uint16
	Y uint16

	NE uint8
	SE uint8
	SW uint8
	NW uint8

	Terrain    terrain.Terrain
	Building   FieldBuildingObjects
	Plant      *terrain.Plant
	Travellers []*Traveller
	Allocated  bool
}

func (f *Field) Field() *Field {
	return f
}

func (f *Field) Context() string {
	return ""
}

func (f Field) Walkable() bool {
	if !f.Building.Empty() {
		return false
	}
	return f.Terrain.T.Walkable && ((f.NE == f.NW && f.SE == f.SW) || (f.NE == f.SE && f.NW == f.SW))
}

func (f Field) Buildable() bool {
	if !f.Building.Empty() {
		return false
	}
	if f.Plant != nil {
		return false
	}
	if f.Allocated {
		return false
	}
	return f.Terrain.T.Buildable && f.NE == f.NW && f.SE == f.SW && f.NE == f.SE && f.NW == f.SW
}

func (f Field) Arable() bool {
	if !f.Building.Empty() {
		return false
	}
	return f.Terrain.T.Arable
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

func (f *Field) Check(f2 *Field) bool {
	return f == f2
}
