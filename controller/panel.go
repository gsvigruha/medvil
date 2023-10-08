package controller

import (
	"github.com/tfriedel6/canvas"
)

type Panel interface {
	Clear()
	CaptureClick(x float64, y float64)
	CaptureMove(x float64, y float64)
	Render(cv *canvas.Canvas)
	Refresh()
}
