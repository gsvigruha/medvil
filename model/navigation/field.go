package navigation

import (
	"medvil/model/building"
	"medvil/model/terrain"
	"strconv"
)

type Location struct {
	X uint16
	Y uint16
	Z uint8
}

type Destination interface {
	Check(PathElement) bool
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

	Terrain      terrain.Terrain
	Building     FieldBuildingObjects
	Plant        *terrain.Plant
	Road         *building.Road
	Travellers   []*Traveller
	Allocated    bool
	Construction bool
}

func (f *Field) GetLocation() Location {
	return Location{X: f.X, Y: f.Y, Z: 0}
}

func (f *Field) GetNeighbors(m IMap) []PathElement {
	var n = []PathElement{}
	nextCoords := [][]uint16{{f.X + 1, f.Y}, {f.X - 1, f.Y}, {f.X, f.Y + 1}, {f.X, f.Y - 1}}
	for idx := range nextCoords {
		x, y := nextCoords[idx][0], nextCoords[idx][1]
		f := m.GetField(x, y)
		if f != nil {
			n = append(n, f)
		}
	}
	return n
}

func (f *Field) GetSpeed() float64 {
	if f.Road != nil && !f.Road.Construction {
		return f.Road.T.Speed
	}
	return 1.0
}

func (f *Field) Field() *Field {
	return f
}

func (f *Field) Context() string {
	return ""
}

func (f Field) Empty() bool {
	if !f.Building.Empty() {
		return false
	}
	if f.Plant != nil {
		return false
	}
	if f.Road != nil {
		return false
	}
	return true
}

func (f Field) Walkable() bool {
	if !f.Building.Empty() {
		return false
	}
	if f.Road != nil && !f.Road.Construction {
		return true
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
	if f.Road != nil {
		return false
	}
	return f.Terrain.T.Buildable && f.NE == f.NW && f.SE == f.SW && f.NE == f.SE && f.NW == f.SW
}

func (f Field) Arable() bool {
	if !f.Building.Empty() {
		return false
	}
	if f.Road != nil {
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

func (f *Field) Check(pe PathElement) bool {
	if f2, ok := pe.(*Field); ok {
		return f2 == f
	}
	return false
}

func (f *Field) CacheKey() string {
	return (strconv.Itoa(int(f.NE)) + "#" +
		strconv.Itoa(int(f.SE)) + "#" +
		strconv.Itoa(int(f.SW)) + "#" +
		strconv.Itoa(int(f.NW)) + "#" +
		f.Terrain.T.Name)
}
