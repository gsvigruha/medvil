package gui

import (
	"github.com/tfriedel6/canvas"
	"image/color"
)

type Label interface {
	Render(cv *canvas.Canvas)
}

type TextLabel struct {
	X    float64
	Y    float64
	SX   float64
	SY   float64
	Text string
}

func (l *TextLabel) Render(cv *canvas.Canvas) {
	cv.SetFillStyle("#FED")
	cv.SetFont("texture/font/Go-Regular.ttf", 12)
	cv.FillText(l.Text, l.X, l.Y)
}

const ImageLabelStyleRegular = 0
const ImageLabelStyleHighlight = 1
const ImageLabelStyleDisabled = 2

type ImageLabel struct {
	X     float64
	Y     float64
	SX    float64
	SY    float64
	Icon  string
	Style uint8
}

func (l *ImageLabel) Render(cv *canvas.Canvas) {
	if l.Style == ImageLabelStyleHighlight {
		cv.SetFillStyle(color.RGBA{R: 224, G: 240, B: 255, A: 240})
		cv.FillRect(l.X, l.Y, l.SX, l.SY)
	}
	cv.DrawImage("icon/gui/"+l.Icon+".png", l.X, l.Y, l.SX, l.SY)
	if l.Style == ImageLabelStyleDisabled {
		cv.SetFillStyle(color.RGBA{R: 0, G: 0, B: 0, A: 64})
		cv.FillRect(l.X, l.Y, l.SX, l.SY)
	}
}

type ScaleLabel struct {
	X      float64
	Y      float64
	SX     float64
	SY     float64
	ScaleW float64
	Icon   string
	Scale  float64
}

func (l *ScaleLabel) Render(cv *canvas.Canvas) {
	cv.DrawImage("icon/gui/"+l.Icon+".png", l.X, l.Y, l.SX, l.SY)
	cv.SetFillStyle("#B00")
	cv.FillRect(l.X+l.SX, l.Y+l.SY, l.ScaleW, -l.SY*l.Scale)
}

type TextureLabel struct {
	X       float64
	Y       float64
	SX      float64
	SY      float64
	Texture string
}

func (l *TextureLabel) Render(cv *canvas.Canvas) {
	cv.SetFillStyle("texture/" + l.Texture + ".png")
	cv.FillRect(l.X, l.Y, l.SX, l.SY)
}
