package renderer

import (
	"github.com/tfriedel6/canvas"
	"math"
	"medvil/model/navigation"
	"strconv"
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
	if x < math.Min(math.Min(math.Min(rf.X[0], rf.X[1]), rf.X[2]), rf.X[3]) ||
		x > math.Max(math.Max(math.Max(rf.X[0], rf.X[1]), rf.X[2]), rf.X[3]) ||
		y < math.Min(math.Min(math.Min(rf.Y[0]-rf.Z[0], rf.Y[1]-rf.Z[1]), rf.Y[2]-rf.Z[2]), rf.X[3]-rf.Z[3]) ||
		y > math.Max(math.Max(math.Max(rf.Y[0]-rf.Z[0], rf.Y[1]-rf.Z[1]), rf.Y[2]-rf.Z[2]), rf.X[3]-rf.Z[3]) {
		return false
	}
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

func (rf RenderedField) Move(dx, dy float64) RenderedField {
	return RenderedField{
		X: MoveVector(rf.X, dx),
		Y: MoveVector(rf.Y, dy),
		Z: rf.Z,
		F: rf.F,
	}
}

func (rf RenderedField) BoundingBox() (float64, float64, float64, float64) {
	var xMin float64 = math.Inf(1)
	var yMin float64 = math.Inf(1)
	var xMax float64 = math.Inf(-1)
	var yMax float64 = math.Inf(-1)

	for i := 0; i < 4; i++ {
		if rf.X[i] < xMin {
			xMin = rf.X[i]
		}
		if rf.Y[i]-rf.Z[i] < yMin {
			yMin = rf.Y[i] - rf.Z[i]
		}
		if rf.X[i] > xMax {
			xMax = rf.X[i]
		}
		if rf.Y[i]-rf.Z[i] > yMax {
			yMax = rf.Y[i] - rf.Z[i]
		}
	}
	return xMin, yMin, xMax, yMax
}

func (rf RenderedField) CacheKey() string {
	return (strconv.Itoa(int(rf.Z[0])) + "#" +
		strconv.Itoa(int(rf.Z[1])) + "#" +
		strconv.Itoa(int(rf.Z[2])) + "#" +
		strconv.Itoa(int(rf.Z[3])))
}

func (rf RenderedField) MidPoint() (float64, float64) {
	midX := (rf.X[0] + rf.X[1] + rf.X[2] + rf.X[3]) / 4
	midY := (rf.Y[0] + rf.Y[1] + rf.Y[2] + rf.Y[3]) / 4
	return midX, midY
}

