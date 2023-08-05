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

func (l *Location) Check(pe PathElement) bool {
	return *l == pe.GetLocation()
}

type FieldWithContext interface {
	Field() *Field
	Context() string
}

const SurroundingSame uint8 = 0
const SurroundingWater uint8 = 1
const SurroundingGrass uint8 = 2
const SurroundingDarkSlope uint8 = 3

type Field struct {
	X uint16
	Y uint16

	NE uint8
	SE uint8
	SW uint8
	NW uint8

	Surroundings [4]uint8

	Terrain      terrain.Terrain
	Building     FieldBuildingObjects `json:"-"`
	Plant        *terrain.Plant
	Animal       *terrain.Animal
	Road         *building.Road
	Travellers   []*Traveller `json:"-"`
	Allocated    bool
	Construction bool
}

func (f *Field) GetLocation() Location {
	return Location{X: f.X, Y: f.Y, Z: 0}
}

func (f *Field) GetPathElement(z uint8) PathElement {
	if z == 0 {
		return f
	}
	bc := f.Building.GetBuildingComponent(z - 1)
	if bc != nil {
		return &BuildingPathElement{BC: bc, L: Location{X: f.X, Y: f.Y, Z: z}}
	}
	return nil
}

func (f *Field) GetNeighbors(m IMap) []PathElement {
	var n = []PathElement{}
	// Connecting field to other field or building component neighbors
	for dir, coordDelta := range building.CoordDeltaByDirection {
		x, y := uint16(coordDelta[0]+int(f.X)), uint16(coordDelta[1]+int(f.Y))
		nf := m.GetField(x, y)
		if nf != nil {
			if nf.Building.Empty() {
				n = append(n, nf)
			} else {
				nbc := nf.Building.GetBuildingComponent(0)
				// Ground level connections
				if nbc == nil || nbc.IsConstruction() {
					n = append(n, nf)
				} else if nbc.Building().Plan.BuildingType != building.BuildingTypeGate {
					// Regular (not gate) buildings can be final ground destinations
					n = append(n, nf)
				} else if nbc != nil && nbc.Connection(building.OppDir(uint8(dir))) == building.ConnectionTypeGround {
					// Some buildings (gate) passable through the ground
					n = append(n, nf)
				}
				// Upper level (building type) connections
				if nbc != nil && nbc.Connection(building.OppDir(uint8(dir))) == building.ConnectionTypeLowerLevel {
					n = append(n, &BuildingPathElement{BC: nbc, L: Location{X: nf.X, Y: nf.Y, Z: 1}})
				}
			}
		}
	}
	// Towers allow vertical movement
	if !f.Building.Empty() && f.Building.GetBuilding().Plan.BuildingType == building.BuildingTypeTower {
		bc := f.Building.GetBuildingComponent(0)
		if bc != nil && !bc.IsConstruction() {
			n = append(n, &BuildingPathElement{BC: bc, L: Location{X: f.X, Y: f.Y, Z: 1}})
		}
	}
	return n
}

func (f *Field) GetSpeed() float64 {
	if f.Road != nil && !f.Road.Construction && !f.Road.Broken {
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
	if f.Animal != nil {
		return false
	}
	return true
}

func (f Field) Walkable() bool {
	if !f.Building.Empty() && f.Building.GetBuilding().Plan.BuildingType != building.BuildingTypeGate {
		return false
	}
	if f.Road != nil && !f.Road.Construction {
		return true
	}
	return f.Terrain.T.Walkable
}

func (f Field) BuildingNonExtension() bool {
	if f.Building.Empty() {
		return false
	}
	_, ok := f.Building.BuildingComponents[0].(*building.ExtensionUnit)
	return !ok
}

func (f Field) Sailable() bool {
	if f.Road != nil && !f.Road.Construction {
		return false
	}
	return f.Terrain.T == terrain.Water
}

func (f Field) BoatDestination() bool {
	if f.Building.Empty() {
		return false
	}
	unit, ok := f.Building.BuildingComponents[0].(*building.ExtensionUnit)
	if !ok {
		return false
	}
	return f.Terrain.T == terrain.Water && unit.T == building.Deck
}

func (f Field) Buildable() bool {
	if !f.Empty() {
		return false
	}
	if f.Allocated {
		return false
	}
	return f.Terrain.T.Buildable && f.Flat()
}

func (f Field) Flat() bool {
	return f.NE == f.NW && f.SE == f.SW && f.NE == f.SE && f.NW == f.SW
}

func (f Field) RoadCompatible() bool {
	if !f.Empty() {
		return false
	}
	if f.Allocated {
		return false
	}
	if !f.Walkable() {
		return false
	}
	if !((f.NE == f.NW && f.SE == f.SW) || (f.NE == f.SE && f.NW == f.SW)) {
		return false
	}
	return f.Terrain.T.Buildable
}

func (f Field) Arable() bool {
	if !f.Building.Empty() {
		return false
	}
	if f.Road != nil {
		return false
	}
	if !f.Flat() {
		return false
	}
	return f.Terrain.T.Arable
}

func (f Field) Plantable() bool {
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

func (f *Field) DarkSlope() bool {
	return (f.SE + f.SW) < (f.NE + f.NW)
}

func (f *Field) LightSlope() bool {
	return (f.SE + f.SW) > (f.NE + f.NW)
}

func (f *Field) Check(pe PathElement) bool {
	if f2, ok := pe.(*Field); ok {
		return f2 == f
	}
	return false
}

func min(x, y uint8) uint8 {
	if x < y {
		return x
	} else {
		return y
	}
}

func (f *Field) CacheKey() string {
	base := min(f.SW, min(f.NW, min(f.NE, f.SE)))
	return (strconv.Itoa(int(f.NE-base)) + "#" +
		strconv.Itoa(int(f.SE-base)) + "#" +
		strconv.Itoa(int(f.SW-base)) + "#" +
		strconv.Itoa(int(f.NW-base)) + "#" +
		strconv.Itoa(int(f.Surroundings[0])) + "#" +
		strconv.Itoa(int(f.Surroundings[1])) + "#" +
		strconv.Itoa(int(f.Surroundings[2])) + "#" +
		strconv.Itoa(int(f.Surroundings[3])) + "#" +
		f.Terrain.T.Name + "#" +
		strconv.Itoa(int(f.Terrain.Shape)))
}

func (f *Field) TravellerVisible() bool {
	return true
}

func (f *Field) TopLocation() *Location {
	return &Location{X: f.X, Y: f.Y, Z: GetZForField(f)}
}
