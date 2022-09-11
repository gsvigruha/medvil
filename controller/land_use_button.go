package controller

import (
	"github.com/tfriedel6/canvas"
	"image/color"
	"medvil/view/gui"
)

type LandUseController interface {
	SetUseType(uint8)
	GetUseType() uint8
}

type LandUseButton struct {
	b       gui.ButtonGUI
	luc     LandUseController
	useType uint8
}

func (b LandUseButton) Click() {
	b.luc.SetUseType(b.useType)
}

func (b LandUseButton) Render(cv *canvas.Canvas) {
	b.b.Render(cv)
	if b.luc.GetUseType() != b.useType {
		cv.SetFillStyle(color.RGBA{R: 64, G: 0, B: 0, A: 128})
		cv.FillRect(b.b.X, b.b.Y, IconS, IconS)
	}
}

func (b LandUseButton) Contains(x float64, y float64) bool {
	return b.b.Contains(x, y)
}

func (b LandUseButton) Enabled() bool {
	return true
}
