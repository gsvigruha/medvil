package gui

import (
	"github.com/tfriedel6/canvas"
)

type DropDown struct {
	X        float64
	Y        float64
	SX       float64
	SY       float64
	Options  []string
	Selected int
	Open     bool
}

func (d *DropDown) Render(cv *canvas.Canvas) {
	cv.SetFillStyle("#321")
	if d.Open {
		cv.FillRect(d.X, d.Y, d.SX, d.SY*float64(len(d.Options)))
	} else {
		cv.FillRect(d.X, d.Y, d.SX, d.SY)
	}
	cv.SetFillStyle("#FED")
	cv.SetFont("texture/font/Go-Regular.ttf", 12)
	if d.Open {
		for i, t := range d.Options {
			cv.FillText(t, d.X, d.Y+float64(i)*d.SY+d.SY-4)
		}
	} else {
		cv.FillText(d.Options[d.Selected], d.X, d.Y+d.SY-4)
	}
}

func (d *DropDown) CaptureClick(x float64, y float64) bool {
	if d.Open {
		if x >= d.X && x < d.X+d.SX && y >= d.Y && y < d.Y+d.SY*float64(len(d.Options)) {
			d.Selected = int((y - d.Y) / d.SY)
			d.Open = false
			return true
		}
	} else {
		if x >= d.X && x < d.X+d.SX && y >= d.Y && y < d.Y+d.SY {
			d.Open = true
			return true
		}
	}
	return false
}
