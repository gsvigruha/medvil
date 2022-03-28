package building

import (
	"fmt"
	"medvil/model/materials"
)

type ConnectionType uint8

const ConnectionTypeNone = 0
const ConnectionTypeUpperLevel = 1
const ConnectionTypeLowerLevel = 2

type BuildingComponentBase struct {
	B            *Building
	Construction bool
}

func (b BuildingComponentBase) Building() *Building {
	return b.B
}

func (b *BuildingComponentBase) SetConstruction(c bool) {
	b.Construction = c
}

type RoofUnit struct {
	BuildingComponentBase
	Roof     Roof
	Elevated [4]bool
}

func (u *RoofUnit) Connection(dir uint8) ConnectionType {
	switch u.Roof.RoofType {
	case RoofTypeSplit:
		return ConnectionTypeNone
	case RoofTypeRamp:
		if dir == u.Roof.RampD {
			return ConnectionTypeUpperLevel
		}
		oppDir := uint8((dir + 2) % 4)
		if oppDir == u.Roof.RampD {
			return ConnectionTypeLowerLevel
		}
	case RoofTypeFlat:
		return ConnectionTypeLowerLevel
	}
	return ConnectionTypeNone
}

type BuildingWall struct {
	M       *materials.Material
	Windows bool
	Door    bool
}

type BuildingUnit struct {
	BuildingComponentBase
	Walls []*BuildingWall
}

func (u *BuildingUnit) Connection(dir uint8) ConnectionType {
	return ConnectionTypeNone
}

type BuildingComponent interface {
	Building() *Building
	SetConstruction(bool)
	Connection(dir uint8) ConnectionType
}

func (b BuildingUnit) Walkable() bool { return false }
func (b BuildingUnit) LiftN() int8    { return 0 }
func (b BuildingUnit) LiftE() int8    { return 0 }
func (b BuildingUnit) LiftS() int8    { return 0 }
func (b BuildingUnit) LiftW() int8    { return 0 }

func (r *RoofUnit) CacheKey() string {
	return fmt.Sprintf("%v#%v#%v#%v", r.Roof.M.Name, r.Elevated, r.Roof.RoofType, r.Construction)
}

func (r *BuildingUnit) CacheKey() string {
	var s = fmt.Sprintf("%v", r.Construction)
	for i := range r.Walls {
		w := r.Walls[i]
		if w != nil {
			s += fmt.Sprintf("[%v#%v#%v#%v]", w.M.Name, w.Windows, w.Door)
		} else {
			s += "[]"
		}
	}
	return s
}
