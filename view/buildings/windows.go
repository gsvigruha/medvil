package buildings

import (
	"github.com/tfriedel6/canvas"
	"image/color"
	"medvil/renderer"
	"path/filepath"
)

func RenderWindows(cv *canvas.Canvas, rf renderer.RenderedField, rfIdx1, rfIdx2 uint8, z float64, door, french bool) {
	cv.BeginPath()
	cv.LineTo((6*rf.X[rfIdx1]+1*rf.X[rfIdx2])/7, (6*rf.Y[rfIdx1]+1*rf.Y[rfIdx2])/7-z-BuildingUnitHeight*DZ*1/3)
	cv.LineTo((6*rf.X[rfIdx1]+1*rf.X[rfIdx2])/7, (6*rf.Y[rfIdx1]+1*rf.Y[rfIdx2])/7-z-BuildingUnitHeight*DZ*2/3)
	cv.LineTo((5*rf.X[rfIdx1]+2*rf.X[rfIdx2])/7, (5*rf.Y[rfIdx1]+2*rf.Y[rfIdx2])/7-z-BuildingUnitHeight*DZ*2/3)
	cv.LineTo((5*rf.X[rfIdx1]+2*rf.X[rfIdx2])/7, (5*rf.Y[rfIdx1]+2*rf.Y[rfIdx2])/7-z-BuildingUnitHeight*DZ*1/3)
	cv.ClosePath()
	cv.Fill()
	cv.Stroke()

	if !door {
		cv.BeginPath()
		cv.LineTo((2*rf.X[rfIdx1]+5*rf.X[rfIdx2])/7, (2*rf.Y[rfIdx1]+5*rf.Y[rfIdx2])/7-z-BuildingUnitHeight*DZ*1/3)
		cv.LineTo((2*rf.X[rfIdx1]+5*rf.X[rfIdx2])/7, (2*rf.Y[rfIdx1]+5*rf.Y[rfIdx2])/7-z-BuildingUnitHeight*DZ*2/3)
		cv.LineTo((1*rf.X[rfIdx1]+6*rf.X[rfIdx2])/7, (1*rf.Y[rfIdx1]+6*rf.Y[rfIdx2])/7-z-BuildingUnitHeight*DZ*2/3)
		cv.LineTo((1*rf.X[rfIdx1]+6*rf.X[rfIdx2])/7, (1*rf.Y[rfIdx1]+6*rf.Y[rfIdx2])/7-z-BuildingUnitHeight*DZ*1/3)
		cv.ClosePath()
		cv.Fill()
		cv.Stroke()
	}

	if french {
		x1 := (4*rf.X[rfIdx1] + 3*rf.X[rfIdx2]) / 7
		x2 := (3*rf.X[rfIdx1] + 4*rf.X[rfIdx2]) / 7
		y1 := (4*rf.Y[rfIdx1] + 3*rf.Y[rfIdx2]) / 7
		y2 := (3*rf.Y[rfIdx1] + 4*rf.Y[rfIdx2]) / 7
		var dx1, dy1, dx2, dy2 float64
		if rfIdx1 == 1 || rfIdx1 == 3 {
			dx1 = (x1 - x2) / 3.0
			dy1 = (y2 - y1) / 3.0
			dx2 = -dx1
			dy2 = dy1
		} else {
			dx1 = (x2 - x1) / 3.0
			dy1 = (y1 - y2) / 3.0
			dx2 = dx1
			dy2 = -dy1
		}

		cv.BeginPath()
		cv.LineTo(x1-dx2, y1-dy2-z-BuildingUnitHeight*DZ*1/3)
		cv.LineTo(x1-dx2, y1-dy2-z-BuildingUnitHeight*DZ*2/3)
		cv.LineTo(x1+dx1, y1+dy1-z-BuildingUnitHeight*DZ*2/3)
		cv.LineTo(x1+dx1, y1+dy1-z-BuildingUnitHeight*DZ*1/3)
		cv.ClosePath()
		cv.Fill()
		cv.Stroke()

		cv.BeginPath()
		cv.LineTo(x1+dx1, y1+dy1-z-BuildingUnitHeight*DZ*1/3)
		cv.LineTo(x1+dx1, y1+dy1-z-BuildingUnitHeight*DZ*2/3)
		cv.LineTo(x2+dx1, y2+dy1-z-BuildingUnitHeight*DZ*2/3)
		cv.LineTo(x2+dx1, y2+dy1-z-BuildingUnitHeight*DZ*1/3)
		cv.ClosePath()
		cv.Fill()
		cv.Stroke()

		cv.BeginPath()
		cv.LineTo(x2+dx2, y2+dy2-z-BuildingUnitHeight*DZ*1/3)
		cv.LineTo(x2+dx2, y2+dy2-z-BuildingUnitHeight*DZ*2/3)
		cv.LineTo(x2+dx1, y2+dy1-z-BuildingUnitHeight*DZ*2/3)
		cv.LineTo(x2+dx1, y2+dy1-z-BuildingUnitHeight*DZ*1/3)
		cv.ClosePath()
		cv.Fill()
		cv.Stroke()

		cv.SetFillStyle(color.RGBA{R: 32, G: 32, B: 32, A: 64})
		cv.BeginPath()
		cv.LineTo(x1-dx2, y1-dy2-z-BuildingUnitHeight*DZ*2/3)
		cv.LineTo(x1+dx1, y1+dy1-z-BuildingUnitHeight*DZ*2/3)
		cv.LineTo(x2+dx1, y2+dy1-z-BuildingUnitHeight*DZ*2/3)
		cv.LineTo(x2+dx2, y2+dy2-z-BuildingUnitHeight*DZ*2/3)
		cv.ClosePath()
		cv.Fill()
	} else {
		cv.BeginPath()
		cv.LineTo((4*rf.X[rfIdx1]+3*rf.X[rfIdx2])/7, (4*rf.Y[rfIdx1]+3*rf.Y[rfIdx2])/7-z-BuildingUnitHeight*DZ*1/3)
		cv.LineTo((4*rf.X[rfIdx1]+3*rf.X[rfIdx2])/7, (4*rf.Y[rfIdx1]+3*rf.Y[rfIdx2])/7-z-BuildingUnitHeight*DZ*2/3)
		cv.LineTo((3*rf.X[rfIdx1]+4*rf.X[rfIdx2])/7, (3*rf.Y[rfIdx1]+4*rf.Y[rfIdx2])/7-z-BuildingUnitHeight*DZ*2/3)
		cv.LineTo((3*rf.X[rfIdx1]+4*rf.X[rfIdx2])/7, (3*rf.Y[rfIdx1]+4*rf.Y[rfIdx2])/7-z-BuildingUnitHeight*DZ*1/3)
		cv.ClosePath()
		cv.Fill()
		cv.Stroke()
	}

	cv.SetStrokeStyle(color.RGBA{R: 128, G: 64, B: 32, A: 32})
	cv.SetLineWidth(3)
	cv.BeginPath()
	cv.LineTo(rf.X[rfIdx1], rf.Y[rfIdx1]-z-BuildingUnitHeight*DZ*1/3+2)
	cv.LineTo(rf.X[rfIdx2], rf.Y[rfIdx2]-z-BuildingUnitHeight*DZ*1/3+2)
	cv.ClosePath()
	cv.Stroke()
}

func RenderBalcony(cv *canvas.Canvas, rf renderer.RenderedField, rfIdx1, rfIdx2 uint8, z float64, door bool) {
	x1 := (7*rf.X[rfIdx1] + 4*rf.X[rfIdx2]) / 11
	x2 := (4*rf.X[rfIdx1] + 7*rf.X[rfIdx2]) / 11
	y1 := (7*rf.Y[rfIdx1]+4*rf.Y[rfIdx2])/11 + 1
	y2 := (4*rf.Y[rfIdx1]+7*rf.Y[rfIdx2])/11 + 1
	var dx, dy float64
	if rfIdx1 == 1 || rfIdx1 == 3 {
		dx = (x1 - x2) / 3.0
		dy = (y2 - y1) / 3.0
	} else {
		dx = (x2 - x1) / 3.0
		dy = (y1 - y2) / 3.0
	}

	cv.SetFillStyle(filepath.FromSlash("texture/building/gray_marble.png"))
	cv.BeginPath()
	cv.LineTo(x1, y1-z-BuildingUnitHeight*DZ*1/3)
	cv.LineTo(x1+dx, y1+dy-z-BuildingUnitHeight*DZ*1/3)
	cv.LineTo(x2+dx, y2+dy-z-BuildingUnitHeight*DZ*1/3)
	cv.LineTo(x2, y2-z-BuildingUnitHeight*DZ*1/3)
	cv.ClosePath()
	cv.Fill()

	cv.SetStrokeStyle(color.RGBA{R: 16, G: 16, B: 32, A: 192})
	cv.SetLineWidth(2)
	cv.BeginPath()
	cv.LineTo(x1+dx, y1+dy-z-BuildingUnitHeight*DZ*1/3)
	cv.LineTo(x1+dx, y1+dy-z-BuildingUnitHeight*DZ*1/2)
	cv.LineTo(x2+dx, y2+dy-z-BuildingUnitHeight*DZ*1/2)
	cv.LineTo(x2+dx, y2+dy-z-BuildingUnitHeight*DZ*1/3)
	cv.ClosePath()
	cv.Stroke()

	cv.BeginPath()
	cv.LineTo(x1, y1-z-BuildingUnitHeight*DZ*1/3)
	cv.LineTo(x1, y1-z-BuildingUnitHeight*DZ*1/2)
	cv.LineTo(x1+dx, y1+dy-z-BuildingUnitHeight*DZ*1/2)
	cv.LineTo(x1+dx, y1+dy-z-BuildingUnitHeight*DZ*1/3)
	cv.ClosePath()
	cv.Stroke()

	cv.BeginPath()
	cv.LineTo(x2, y2-z-BuildingUnitHeight*DZ*1/3)
	cv.LineTo(x2, y2-z-BuildingUnitHeight*DZ*1/2)
	cv.LineTo(x2+dx, y2+dy-z-BuildingUnitHeight*DZ*1/2)
	cv.LineTo(x2+dx, y2+dy-z-BuildingUnitHeight*DZ*1/3)
	cv.ClosePath()
	cv.Stroke()

	cv.BeginPath()
	cv.LineTo(x2+dx/2, y2+dy/2-z-BuildingUnitHeight*DZ*1/3)
	cv.LineTo(x2+dx/2, y2+dy/2-z-BuildingUnitHeight*DZ*1/2)
	cv.ClosePath()
	cv.Stroke()

	cv.BeginPath()
	cv.LineTo(x1+dx/2, y1+dy/2-z-BuildingUnitHeight*DZ*1/3)
	cv.LineTo(x1+dx/2, y1+dy/2-z-BuildingUnitHeight*DZ*1/2)
	cv.ClosePath()
	cv.Stroke()

	ddx := ((x1 + dx) - (x2 + dx)) / 5
	ddy := ((y1 + dy) - (y2 + dy)) / 5
	for i := 0.0; i < 5; i++ {
		cv.BeginPath()
		cv.LineTo(x2+dx+ddx*i, y2+dy+ddy*i-z-BuildingUnitHeight*DZ*1/3)
		cv.LineTo(x2+dx+ddx*i, y2+dy+ddy*i-z-BuildingUnitHeight*DZ*1/2)
		cv.ClosePath()
		cv.Stroke()
	}
}

func renderFactoryWindow(cv *canvas.Canvas, x1, y1, z1, x2, y2, z2 float64) {
	cv.BeginPath()
	cv.LineTo(x1, y1-z1)
	cv.LineTo(x1, y1-z2)
	cv.LineTo(x2, y2-z2)
	cv.LineTo(x2, y2-z1)
	cv.ClosePath()
	cv.Fill()
	cv.Stroke()

	dx := (x2 - x1) / 3.0
	dy := (y2 - y1) / 3.0
	dz := (z2 - z1) / 3.0
	for i := 0.0; i < 3.0; i++ {
		cv.BeginPath()
		cv.MoveTo(x1+dx*i, y1+dy*i-z1)
		cv.LineTo(x1+dx*i, y1+dy*i-z2)
		cv.ClosePath()
		cv.Stroke()
		cv.BeginPath()
		cv.MoveTo(x1, y1-z1-dz*i)
		cv.LineTo(x2, y2-z1-dz*i)
		cv.ClosePath()
		cv.Stroke()
	}
}

func RenderFactoryWindows(cv *canvas.Canvas, rf renderer.RenderedField, rfIdx1, rfIdx2 uint8, z float64, door bool) {
	renderFactoryWindow(cv,
		(6*rf.X[rfIdx1]+1*rf.X[rfIdx2])/7,
		(6*rf.Y[rfIdx1]+1*rf.Y[rfIdx2])/7,
		z+BuildingUnitHeight*DZ*1/3,
		(4*rf.X[rfIdx1]+3*rf.X[rfIdx2])/7,
		(4*rf.Y[rfIdx1]+3*rf.Y[rfIdx2])/7,
		z+BuildingUnitHeight*DZ*2/3,
	)

	if !door {
		renderFactoryWindow(cv,
			(3*rf.X[rfIdx1]+4*rf.X[rfIdx2])/7,
			(3*rf.Y[rfIdx1]+4*rf.Y[rfIdx2])/7,
			z+BuildingUnitHeight*DZ*1/3,
			(1*rf.X[rfIdx1]+6*rf.X[rfIdx2])/7,
			(1*rf.Y[rfIdx1]+6*rf.Y[rfIdx2])/7,
			z+BuildingUnitHeight*DZ*2/3,
		)
	}
}
