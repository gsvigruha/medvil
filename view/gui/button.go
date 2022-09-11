package gui

import (
	"github.com/tfriedel6/canvas"
	"image/color"
)

type Button interface {
	Click()
	Render(cv *canvas.Canvas)
	Contains(x float64, y float64) bool
	Enabled() bool
}

type ButtonGUI struct {
	X        float64
	Y        float64
	SX       float64
	SY       float64
	Icon     string
	Texture  string
	Disabled func() bool
}

func (b ButtonGUI) Render(cv *canvas.Canvas) {
	if b.Texture != "" {
		cv.SetFillStyle("texture/" + b.Texture + ".png")
		cv.FillRect(b.X, b.Y, b.SX, b.SY)
	}
	if b.Icon != "" {
		cv.DrawImage("icon/gui/"+b.Icon+".png", b.X, b.Y, b.SX, b.SY)
	}
	if !b.Enabled() {
		cv.SetFillStyle(color.RGBA{R: 0, G: 0, B: 0, A: 128})
		cv.FillRect(b.X, b.Y, b.SX, b.SY)
	}
}

func (b ButtonGUI) Contains(x float64, y float64) bool {
	return b.X <= x && b.X+b.SX >= x && b.Y <= y && b.Y+b.SY >= y
}

func (b ButtonGUI) Enabled() bool {
	return b.Disabled == nil || !b.Disabled()
}
