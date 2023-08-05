package gui

import (
	"github.com/tfriedel6/canvas"
	"image/color"
)

var FontSize = 12.0

type Label interface {
	Render(cv *canvas.Canvas)
	CaptureClick(x float64, y float64)
}

type TextLabel struct {
	X     float64
	Y     float64
	SX    float64
	SY    float64
	Text  string
	Large bool
}

func (l *TextLabel) Render(cv *canvas.Canvas) {
	cv.SetFillStyle("#FED")
	if l.Large {
		cv.SetFont("texture/font/Go-Regular.ttf", FontSize*1.5)
	} else {
		cv.SetFont("texture/font/Go-Regular.ttf", FontSize)
	}
	cv.FillText(l.Text, l.X, l.Y)
}

func (l *TextLabel) CaptureClick(x float64, y float64) {}

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

func (l *ImageLabel) CaptureClick(x float64, y float64) {}

type DoubleImageLabel struct {
	X       float64
	Y       float64
	SX      float64
	SY      float64
	Icon    string
	SubIcon string
	Style   uint8
}

func (l *DoubleImageLabel) Render(cv *canvas.Canvas) {
	if l.Style == ImageLabelStyleHighlight {
		cv.SetFillStyle(color.RGBA{R: 224, G: 240, B: 255, A: 240})
		cv.FillRect(l.X, l.Y, l.SX, l.SY)
	}
	cv.DrawImage("icon/gui/"+l.Icon+".png", l.X, l.Y, l.SX, l.SY)
	cv.DrawImage("icon/gui/"+l.SubIcon+".png", l.X, l.Y, l.SX/2, l.SY/2)
	if l.Style == ImageLabelStyleDisabled {
		cv.SetFillStyle(color.RGBA{R: 0, G: 0, B: 0, A: 64})
		cv.FillRect(l.X, l.Y, l.SX, l.SY)
	}
}

func (l *DoubleImageLabel) CaptureClick(x float64, y float64) {}

type ScaleLabel struct {
	X       float64
	Y       float64
	SX      float64
	SY      float64
	ScaleW  float64
	Icon    string
	Scale   float64
	Stacked bool
}

func (l *ScaleLabel) Render(cv *canvas.Canvas) {
	if l.Stacked {
		iconTop := l.Y + l.SY - l.SX
		cv.DrawImage("icon/gui/"+l.Icon+".png", l.X, iconTop, l.SX, l.SX)
		cv.SetFillStyle("#B00")
		var s = l.Scale
		if s >= 1.0 {
			s = 1.0
		}
		cv.FillRect(l.X+l.SX/2-l.ScaleW/2, iconTop, l.ScaleW, -(l.SY-l.SX)*s)
	} else {
		cv.DrawImage("icon/gui/"+l.Icon+".png", l.X, l.Y, l.SX, l.SY)
		cv.SetFillStyle("#B00")
		var s = l.Scale
		if s >= 1.0 {
			s = 1.0
		}
		cv.FillRect(l.X+l.SX, l.Y+l.SY, l.ScaleW, -l.SY*s)
	}
}

func (l *ScaleLabel) CaptureClick(x float64, y float64) {}

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

func (l *TextureLabel) CaptureClick(x float64, y float64) {}

type DynamicImageLabel struct {
	X    float64
	Y    float64
	SX   float64
	SY   float64
	Icon func() string
}

func (l *DynamicImageLabel) Render(cv *canvas.Canvas) {
	cv.DrawImage("icon/gui/"+l.Icon()+".png", l.X, l.Y, l.SX, l.SY)
}

func (l *DynamicImageLabel) CaptureClick(x float64, y float64) {}
