package building

import (
	"fmt"
	"medvil/model/materials"
)

type BuildingComponentBase struct {
	B            *Building
	Construction bool
}

func (b BuildingComponentBase) Building() *Building {
	return b.B
}

type RoofUnit struct {
	BuildingComponentBase
	Roof     Roof
	Elevated [4]bool
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

type BuildingComponent interface {
	Building() *Building
}

func (b BuildingUnit) Walkable() bool { return false }
func (b BuildingUnit) LiftN() int8    { return 0 }
func (b BuildingUnit) LiftE() int8    { return 0 }
func (b BuildingUnit) LiftS() int8    { return 0 }
func (b BuildingUnit) LiftW() int8    { return 0 }

func (r *RoofUnit) CacheKey() string {
	return fmt.Sprintf("%v#%v#%v", r.Roof.M.Name, r.Elevated, r.Construction)
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
