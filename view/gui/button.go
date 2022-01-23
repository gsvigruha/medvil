package gui

import (
	"github.com/tfriedel6/canvas"
)

type Button struct {
	X        float64
	Y        float64
	SX       float64
	SY       float64
	Icon     string
	Callback func(interface{})
}

func (b *Button) Render(cv *canvas.Canvas) {
	cv.DrawImage("icon/gui/"+b.Icon+".png", b.X, b.Y, b.SX, b.SY)
}

func (b *Button) Contains(x float64, y float64) bool {
	return b.X <= x && b.X+b.SX >= x && b.Y <= y && b.Y+b.SY >= y
}
