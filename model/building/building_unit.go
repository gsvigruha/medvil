package building

import (
	"medvil/model/materials"
)

type RoofUnit struct {
	Roof     Roof
	Elevated [4]bool
}

type BuildingWall struct {
	M       *materials.Material
	Windows bool
	Door    bool
}

type BuildingUnit struct {
	Walls []*BuildingWall
}

func (b BuildingUnit) Walkable() bool { return false }
func (b BuildingUnit) LiftN() int8    { return 0 }
func (b BuildingUnit) LiftE() int8    { return 0 }
func (b BuildingUnit) LiftS() int8    { return 0 }
func (b BuildingUnit) LiftW() int8    { return 0 }
