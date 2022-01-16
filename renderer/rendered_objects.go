package renderer

import (
	"github.com/tfriedel6/canvas"
	"medvil/model/navigation"
	"math"
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

func RayIntersects(x float64, y float64, lx1 float64, ly1 float64, lx2 float64, ly2 float64) bool {
	if lx1 != lx2 {
		a := (ly2 - ly1) / (lx2 - lx1)
		b := ly1 - a * lx1
		xi := (y - b) / a
		return xi >= math.Min(lx1, lx2) && xi <= math.Max(lx1, lx2) && xi >= x
	} else {
		xi := lx1
		yi := y
		return yi >= math.Min(ly1, ly2) && yi <= math.Max(ly1, ly2) && xi >= x
	}
}

func BtoI(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

func (rf *RenderedField) Contains(x float64, y float64) bool {
	return (
		BtoI(RayIntersects(x, y, rf.X[0], rf.Y[0]-rf.Z[0], rf.X[1], rf.Y[1]-rf.Z[1])) + 
		BtoI(RayIntersects(x, y, rf.X[1], rf.Y[1]-rf.Z[1], rf.X[2], rf.Y[2]-rf.Z[2])) + 
		BtoI(RayIntersects(x, y, rf.X[2], rf.Y[2]-rf.Z[2], rf.X[3], rf.Y[3]-rf.Z[3])) + 
		BtoI(RayIntersects(x, y, rf.X[3], rf.Y[3]-rf.Z[3], rf.X[0], rf.Y[0]-rf.Z[0]))) % 2 == 1
}
