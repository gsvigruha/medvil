package util

import (
	"github.com/tfriedel6/canvas"
	"medvil/renderer"
)

func RenderPolygon(cv *canvas.Canvas, polygon renderer.Polygon, stroke bool) {
	cv.BeginPath()
	for _, p := range polygon.Points {
		cv.LineTo(p.X, p.Y)
	}
	cv.ClosePath()
	cv.Fill()
	if stroke {
		cv.Stroke()
	}
}
