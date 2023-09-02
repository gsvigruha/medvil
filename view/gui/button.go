package gui

import (
	"github.com/tfriedel6/canvas"
	"image/color"
)

var ButtonColorHighlight = color.RGBA{R: 192, G: 224, B: 255, A: 192}

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

type SimpleButton struct {
	ButtonGUI
	Highlight func() bool
	ClickImpl func()
}

func (b SimpleButton) Render(cv *canvas.Canvas) {
	if b.Highlight != nil && b.Highlight() {
		cv.SetFillStyle(ButtonColorHighlight)
		cv.FillRect(b.X, b.Y, b.SX, b.SY)
	}
	b.ButtonGUI.Render(cv)
}

func (b SimpleButton) Click() {
	b.ClickImpl()
}

type ImageButton struct {
	ButtonGUI
	Style     uint8
	ClickImpl func()
}

func (b *ImageButton) Render(cv *canvas.Canvas) {
	if b.Style == ImageLabelStyleHighlight {
		cv.SetFillStyle(color.RGBA{R: 224, G: 240, B: 255, A: 240})
		cv.FillRect(b.X, b.Y, b.SX, b.SY)
	}
	b.ButtonGUI.Render(cv)
	if b.Style == ImageLabelStyleDisabled {
		cv.SetFillStyle(color.RGBA{R: 0, G: 0, B: 0, A: 64})
		cv.FillRect(b.X, b.Y, b.SX, b.SY)
	}
}

func (b ImageButton) Click() {
	b.ClickImpl()
}
