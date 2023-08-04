package building

import (
	"fmt"
	"medvil/model/materials"
)

type ConnectionType uint8

const ConnectionTypeNone = 0
const ConnectionTypeLowerLevel = 1
const ConnectionTypeUpperLevel = 2
const ConnectionTypeGround = 3

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

func (b *BuildingComponentBase) IsConstruction() bool {
	return b.Construction
}

type RoofUnit struct {
	BuildingComponentBase
	Roof      Roof
	Connected [4]bool
}

func (u *RoofUnit) Connection(dir uint8) ConnectionType {
	switch u.Roof.RoofType {
	case RoofTypeSplit:
		return ConnectionTypeNone
	case RoofTypeRamp:
		if dir == u.Roof.RampD {
			return ConnectionTypeUpperLevel
		}
		if OppDir(dir) == u.Roof.RampD {
			return ConnectionTypeLowerLevel
		}
	case RoofTypeFlat:
		return ConnectionTypeLowerLevel
	}
	return ConnectionTypeNone
}

func (u *RoofUnit) NamePlate() bool {
	return false
}

type WindowType uint8

const WindowTypeNone WindowType = 0
const WindowTypePlain WindowType = 1
const WindowTypeBalcony WindowType = 2
const WindowTypeFrench WindowType = 3
const WindowTypeFactory WindowType = 4

type BuildingWall struct {
	M       *materials.Material
	Windows WindowType
	Door    bool
	Arch    bool
}

type BuildingUnit struct {
	BuildingComponentBase
	Walls []*BuildingWall
}

func (u *BuildingUnit) Connection(dir uint8) ConnectionType {
	// Gates can only be passed through one direction
	if u.B.Plan.BuildingType == BuildingTypeGate {
		if dir%2 == u.B.Direction%2 {
			return ConnectionTypeGround
		}
	}
	// Towers are accessible to all neighbors
	if u.B.Plan.BuildingType == BuildingTypeTower {
		return ConnectionTypeLowerLevel
	}
	return ConnectionTypeNone
}

func (u *BuildingUnit) NamePlate() bool {
	return true
}

func (u *BuildingUnit) HasDoor() bool {
	for _, w := range u.Walls {
		if w != nil && w.Door {
			return true
		}
	}
	return false
}

type ExtensionUnit struct {
	BuildingComponentBase
	Direction uint8
	T         *BuildingExtensionType
}

func (u *ExtensionUnit) Connection(dir uint8) ConnectionType {
	return ConnectionTypeNone
}

func (u *ExtensionUnit) NamePlate() bool {
	return false
}

type BuildingComponent interface {
	Building() *Building
	SetConstruction(bool)
	IsConstruction() bool
	Connection(dir uint8) ConnectionType
	NamePlate() bool
}

func (b BuildingUnit) Walkable() bool { return false }
func (b BuildingUnit) LiftN() int8    { return 0 }
func (b BuildingUnit) LiftE() int8    { return 0 }
func (b BuildingUnit) LiftS() int8    { return 0 }
func (b BuildingUnit) LiftW() int8    { return 0 }

func (r *RoofUnit) CacheKey() string {
	return fmt.Sprintf("%v#%v#%v#%v#%v#%v#%s", r.Roof.M.Name, r.Connected, r.Roof.RoofType, r.Construction, r.B.Shape, r.B.Plan.BuildingType, r.B.Broken)
}

func (e *ExtensionUnit) CacheKey() string {
	return fmt.Sprintf("%v#%v#%v#%v", e.T, e.Direction, e.Construction, e.B.Shape)
}

func (u *BuildingUnit) CacheKey() string {
	var s = fmt.Sprintf("%v", u.Construction)
	for i := range u.Walls {
		w := u.Walls[i]
		if w != nil {
			s += fmt.Sprintf("[%v#%v#%v#%v#%v#%v#%v#%v]", w.M.Name, w.Windows, w.Door, u.B.Shape, u.B.Plan.BuildingType, w.Arch, u.B.Broken)
		} else {
			s += "[]"
		}
	}
	return s
}
