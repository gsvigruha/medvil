package renderer

import (
	"github.com/tfriedel6/canvas"
	"medvil/model/building"
)

type RenderedWall struct {
	X    [4]float64
	Y    [4]float64
	Wall *building.BuildingWall
}

type RenderedBuildingUnit struct {
	Walls []RenderedWall
	Unit  *building.BuildingUnit
}

func (rw RenderedWall) Draw(cv *canvas.Canvas) {
	cv.BeginPath()
	cv.LineTo(rw.X[0], rw.Y[0])
	cv.LineTo(rw.X[1], rw.Y[1])
	cv.LineTo(rw.X[2], rw.Y[2])
	cv.LineTo(rw.X[3], rw.Y[3])
	cv.ClosePath()
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
