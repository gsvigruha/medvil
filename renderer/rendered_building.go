package renderer

import (
	"medvil/model/building"
)

type RenderedBuildingPart interface {
	Contains(x float64, y float64) bool
	GetBuilding() *building.Building
}

type RenderedWall struct {
	X    [4]float64
	Y    [4]float64
	Wall *building.BuildingWall
}

type RenderedBuildingUnit struct {
	Walls []RenderedWall
	Unit  *building.BuildingUnit
}

func (rw *RenderedWall) Contains(x float64, y float64) bool {
	return (BtoI(RayIntersects(x, y, rw.X[0], rw.Y[0], rw.X[1], rw.Y[1]))+
		BtoI(RayIntersects(x, y, rw.X[1], rw.Y[1], rw.X[2], rw.Y[2]))+
		BtoI(RayIntersects(x, y, rw.X[2], rw.Y[2], rw.X[3], rw.Y[3]))+
		BtoI(RayIntersects(x, y, rw.X[3], rw.Y[3], rw.X[0], rw.Y[0])))%2 == 1
}

func (rbu *RenderedBuildingUnit) Contains(x float64, y float64) bool {
	for i := range rbu.Walls {
		if rbu.Walls[i].Contains(x, y) {
			return true
		}
	}
	return false
}

func (rbu *RenderedBuildingUnit) GetBuilding() *building.Building {
	return rbu.Unit.B
}

func (rbu RenderedBuildingUnit) Move(dx, dy float64) RenderedBuildingUnit {
	var walls []RenderedWall
	for i := range rbu.Walls {
		walls = append(walls, RenderedWall{
			X:    MoveVector(rbu.Walls[i].X, dx),
			Y:    MoveVector(rbu.Walls[i].Y, dy),
			Wall: rbu.Walls[i].Wall,
		})
	}
	return RenderedBuildingUnit{
		Walls: walls,
		Unit:  rbu.Unit,
	}
}

type RenderedBuildingRoof struct {
	Ps []Polygon
	B  *building.Building
}

func (rbr *RenderedBuildingRoof) Contains(x float64, y float64) bool {
	for _, p := range rbr.Ps {
		if p.Contains(x, y) {
			return true
		}
	}
	return false
}

func (rbr *RenderedBuildingRoof) GetBuilding() *building.Building {
	return rbr.B
}

func (rbr RenderedBuildingRoof) Move(dx, dy float64) RenderedBuildingRoof {
	ps := make([]Polygon, len(rbr.Ps))
	for i, p := range rbr.Ps {
		ps[i] = p.Move(dx, dy)
	}
	return RenderedBuildingRoof{
		Ps: ps,
		B:  rbr.B,
	}
}
