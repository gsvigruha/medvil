package renderer

import (
	"github.com/tfriedel6/canvas"
	"medvil/model/navigation"
)

type RenderedField struct {
	X [4]float64
	Y [4]float64
	Z [4]float64
	F *navigation.Field
}

func (rf RenderedField) Draw(cv *canvas.Canvas) {
	cv.BeginPath()
	cv.LineTo(rf.X[0], rf.Y[0]-rf.Z[0])
	cv.LineTo(rf.X[1], rf.Y[1]-rf.Z[1])
	cv.LineTo(rf.X[2], rf.Y[2]-rf.Z[2])
	cv.LineTo(rf.X[3], rf.Y[3]-rf.Z[3])
	cv.ClosePath()
}

func (rf *RenderedField) Contains(x float64, y float64) bool {
	return (BtoI(RayIntersects(x, y, rf.X[0], rf.Y[0]-rf.Z[0], rf.X[1], rf.Y[1]-rf.Z[1]))+
		BtoI(RayIntersects(x, y, rf.X[1], rf.Y[1]-rf.Z[1], rf.X[2], rf.Y[2]-rf.Z[2]))+
		BtoI(RayIntersects(x, y, rf.X[2], rf.Y[2]-rf.Z[2], rf.X[3], rf.Y[3]-rf.Z[3]))+
		BtoI(RayIntersects(x, y, rf.X[3], rf.Y[3]-rf.Z[3], rf.X[0], rf.Y[0]-rf.Z[0])))%2 == 1
}

func MoveVector(v [4]float64, d float64) [4]float64 {
	var r [4]float64
	for i := range r {
		r[i] = v[i] + d
	}
	return r
}

func (rf RenderedField) Move(dx, dy float64) (RenderedField) {
	return RenderedField{
		X: MoveVector(rf.X, dx),
		Y: MoveVector(rf.Y, dy),
		Z: rf.Z,
		F: rf.F,
	}
}
