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
	var maxLength = 0.0
	for _, line := range lines {
		l := float64(len(line))
		if l > maxLength {
			maxLength = l
		}
	}

	var dh = 0.0
	if s.Y < 200 {
		dh = iconD
	}
	tw := maxLength * FontSize * 0.5
	th := math.Max(float64(len(lines))*FontSize, iconD)

	cv.SetFillStyle("texture/wood.png")
	cv.SetStrokeStyle("#DDD")
	cv.SetLineWidth(2)
	cv.BeginPath()
	cv.LineTo(s.X, s.Y)
	cv.LineTo(s.X+120, s.Y-20+dh)
	cv.LineTo(s.X+120, s.Y-20-th/2.0+dh)
	cv.LineTo(s.X+120+tw+iconD, s.Y-20-th/2.0+dh)
	cv.LineTo(s.X+120+tw+iconD, s.Y+20+th/2.0+dh)
	cv.LineTo(s.X+120, s.Y+20+th/2.0+dh)
	cv.LineTo(s.X+120, s.Y+20+dh)
	cv.ClosePath()
	cv.Fill()
	cv.Stroke()
	cv.DrawImage("icon/gui/"+s.Icon+".png", s.X+130, s.Y-iconS/2.0+dh, iconS, iconS)
	cv.SetFillStyle("#FED")
	cv.SetFont("texture/font/Go-Regular.ttf", FontSize)
	for i, line := range lines {
		cv.FillText(line, s.X+130+iconD, s.Y-float64(len(lines)-2)*FontSize/2.0-4+float64(i)*FontSize+dh)
	}
}
