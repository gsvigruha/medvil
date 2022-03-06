package building

import (
	"fmt"
	"medvil/model/materials"
)

type RoofUnit struct {
	Roof     Roof
	Elevated [4]bool
	B        *Building
}

type BuildingWall struct {
	M            *materials.Material
	Windows      bool
	Door         bool
	B            *Building
	Construction bool
}

type BuildingUnit struct {
	Walls []*BuildingWall
	B     *Building
}

func (b BuildingUnit) Walkable() bool { return false }
func (b BuildingUnit) LiftN() int8    { return 0 }
func (b BuildingUnit) LiftE() int8    { return 0 }
func (b BuildingUnit) LiftS() int8    { return 0 }
func (b BuildingUnit) LiftW() int8    { return 0 }

func (r *RoofUnit) CacheKey() string {
	return fmt.Sprintf("%v#%v", r.Roof.M.Name, r.Elevated)
}

func (r *BuildingUnit) CacheKey() string {
	var s = ""
	for i := range r.Walls {
		w := r.Walls[i]
		if w != nil {
			s += fmt.Sprintf("[%v#%v#%v#%v]", w.M.Name, w.Windows, w.Door, w.Construction)
		} else {
			s += "[]"
		}
	}
	return s
}
