package view

import (
	"github.com/tfriedel6/canvas"
)

type RenderedField struct {
	X [4]float64
	Y [4]float64
	Z [4]float64
}

func (rf RenderedField) Draw(cv *canvas.Canvas) {
	cv.BeginPath()
	cv.LineTo(rf.X[0], rf.Y[0]-rf.Z[0])
	cv.LineTo(rf.X[1], rf.Y[1]-rf.Z[1])
	cv.LineTo(rf.X[2], rf.Y[2]-rf.Z[2])
	cv.LineTo(rf.X[3], rf.Y[3]-rf.Z[3])
	cv.ClosePath()

}
