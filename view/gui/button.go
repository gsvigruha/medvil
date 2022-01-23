package gui

import (
	"github.com/tfriedel6/canvas"
)

type Button interface {
	Click()
	Render(cv *canvas.Canvas)
	Contains(x float64, y float64) bool
}

type ButtonGUI struct {
	X       float64
	Y       float64
	SX      float64
	SY      float64
	Icon    string
	Texture string
}

func (b ButtonGUI) Render(cv *canvas.Canvas) {
	if b.Texture != "" {
		cv.SetFillStyle("texture/" + b.Texture + ".png")
		cv.FillRect(b.X, b.Y, b.SX, b.SY)
	}
	if b.Icon != "" {
		cv.DrawImage("icon/gui/"+b.Icon+".png", b.X, b.Y, b.SX, b.SY)
	}
}

func (b ButtonGUI) Contains(x float64, y float64) bool {
	return b.X <= x && b.X+b.SX >= x && b.Y <= y && b.Y+b.SY >= y
}
