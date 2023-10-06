package renderer

import (
	"github.com/tfriedel6/canvas"
	"image/color"
	"medvil/model/navigation"
)

type RenderedTraveller struct {
	X         float64
	Y         float64
	W         float64
	H         float64
	Traveller *navigation.Traveller
}

func (rt *RenderedTraveller) Contains(x float64, y float64) bool {
	return x >= rt.X && x <= rt.X+rt.W && y >= rt.Y && y <= rt.Y+rt.H
}

func (rt *RenderedTraveller) Draw(cv *canvas.Canvas) {
	cv.SetStrokeStyle(color.RGBA{R: 0, G: 192, B: 0, A: 255})
	cv.SetLineWidth(2)
	cv.StrokeRect(rt.X, rt.Y, rt.W, rt.H)
}
