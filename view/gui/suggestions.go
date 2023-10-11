package gui

import (
	"github.com/tfriedel6/canvas"
	"math"
	"strings"
)

type Suggestion struct {
	Message string
	Icon    string
	X       float64
	Y       float64
}

func (s *Suggestion) Render(cv *canvas.Canvas, iconS, iconD float64) {
	lines := strings.Split(s.Message, "\n")
	var mw = 0.0
	for _, line := range lines {
		w := EstimateWidth(line) * FontSize
		if mw < w {
			mw = w
		}
	}

	var dh = 0.0
	if s.Y <= iconD*3 {
		dh = iconD
	}
	th := math.Max(float64(len(lines))*(FontSize+4), iconD)
	p := iconD / 3.0
	dw := iconD * 2.0

	cv.SetFillStyle("texture/wood.png")
	cv.SetStrokeStyle("#DDD")
	cv.SetLineWidth(2)
	cv.BeginPath()
	cv.LineTo(s.X, s.Y)
	cv.LineTo(s.X+dw, s.Y-p+dh)
	cv.LineTo(s.X+dw, s.Y-p-th/2.0+dh)
	cv.LineTo(s.X+dw+mw+iconD+p*3, s.Y-p-th/2.0+dh)
	cv.LineTo(s.X+dw+mw+iconD+p*3, s.Y+p+th/2.0+dh)
	cv.LineTo(s.X+dw, s.Y+p+th/2.0+dh)
	cv.LineTo(s.X+dw, s.Y+p+dh)
	cv.ClosePath()
	cv.Fill()
	cv.Stroke()
	cv.DrawImage("icon/gui/"+s.Icon+".png", s.X+dw+p, s.Y-iconS/2.0+dh, iconS, iconS)
	cv.SetFillStyle("#FED")
	cv.SetFont("texture/font/Go-Regular.ttf", FontSize)
	for i, line := range lines {
		cv.FillText(line, s.X+dw+p+iconD, s.Y-float64(len(lines)-2)*(FontSize+4)/2.0-4+float64(i)*(FontSize+4)+dh)
	}
}
