package vehicles

import (
	//"fmt"
	"github.com/tfriedel6/canvas"
	"math"
	"medvil/controller"
	"medvil/model/navigation"
	"medvil/view/animation"
	"path/filepath"
)

type BoatConfig struct {
	f1 float64
	f2 float64
	s  float64
	z  float64
	h1 float64
	h2 float64
	p  float64
}

func DrawTradingBoat(cv *canvas.Canvas, t *navigation.Traveller, x float64, y float64, c *controller.Controller) {
	bc := BoatConfig{
		f1: 24.0,
		f2: 8.0,
		s:  8.0,
		z:  6.0,
		h1: 8.0,
		h2: 6.0,
		p:  16.0,
	}
	drawBoatBody(cv, t, x, y, bc, c)

	dirIdx := (c.Perspective - t.Direction) % 4
	pm := animation.ProjectionMatrices[dirIdx]

	if dirIdx == 0 || dirIdx == 2 {
		cv.SetFillStyle(filepath.FromSlash("texture/vehicle/textile.png"))
	} else {
		cv.SetFillStyle(filepath.FromSlash("texture/vehicle/textile_flipped.png"))
	}
	cv.BeginPath()
	cv.LineTo(x+bc.f1*pm.XX-bc.h1*pm.XY+0*pm.XZ, y+bc.f1*pm.YX-bc.h1*pm.YY+0*pm.YZ)
	cv.LineTo(x+bc.f2*pm.XX-bc.h2*pm.XY-bc.s*pm.XZ, y+bc.f2*pm.YX-bc.h2*pm.YY-bc.s*pm.YZ)
	cv.LineTo(x+bc.f2*pm.XX-bc.h2*pm.XY+bc.s*pm.XZ, y+bc.f2*pm.YX-bc.h2*pm.YY+bc.s*pm.YZ)
	cv.ClosePath()
	cv.Fill()

	cv.BeginPath()
	cv.LineTo(x-bc.f1*pm.XX-bc.h1*pm.XY+0*pm.XZ, y-bc.f1*pm.YX-bc.h1*pm.YY+0*pm.YZ)
	cv.LineTo(x-bc.f2*pm.XX-bc.h2*pm.XY-bc.s*pm.XZ, y-bc.f2*pm.YX-bc.h2*pm.YY-bc.s*pm.YZ)
	cv.LineTo(x-bc.f2*pm.XX-bc.h2*pm.XY+bc.s*pm.XZ, y-bc.f2*pm.YX-bc.h2*pm.YY+bc.s*pm.YZ)
	cv.ClosePath()
	cv.Fill()
}

func DrawBoat(cv *canvas.Canvas, t *navigation.Traveller, x float64, y float64, c *controller.Controller) {
	bc := BoatConfig{
		f1: 20.0,
		f2: 6.0,
		s:  8.0,
		z:  6.0,
		h1: 6.0,
		h2: 6.0,
		p:  16.0,
	}
	drawBoatBody(cv, t, x, y, bc, c)
}

func drawBoatBody(cv *canvas.Canvas, t *navigation.Traveller, x float64, y float64, c BoatConfig, ctrl *controller.Controller) {
	if !t.Visible {
		return
	}
	dirIdx := (ctrl.Perspective - t.Direction) % 4
	pm := animation.ProjectionMatrices[dirIdx]
	var r = 8.0 - float64(t.Phase%16)
	if r < 0.0 {
		r = -r
	}

	cv.SetFillStyle(filepath.FromSlash("texture/vehicle/boat_bottom.png"))
	cv.BeginPath()
	cv.LineTo(x-c.f1*pm.XX+0*pm.XY+0*pm.XZ, y-c.f1*pm.YX+0*pm.YY+0*pm.YZ)
	cv.LineTo(x-c.f2*pm.XX+0*pm.XY-c.z*pm.XZ, y-c.f2*pm.YX+0*pm.YY-c.z*pm.YZ)
	cv.LineTo(x+c.f2*pm.XX+0*pm.XY-c.z*pm.XZ, y+c.f2*pm.YX+0*pm.YY-c.z*pm.YZ)
	cv.LineTo(x+c.f1*pm.XX+0*pm.XY+0*pm.XZ, y+c.f1*pm.YX+0*pm.YY+0*pm.YZ)
	cv.LineTo(x+c.f2*pm.XX+0*pm.XY+c.z*pm.XZ, y+c.f2*pm.YX+0*pm.YY+c.z*pm.YZ)
	cv.LineTo(x-c.f2*pm.XX+0*pm.XY+c.z*pm.XZ, y-c.f2*pm.YX+0*pm.YY+c.z*pm.YZ)
	cv.ClosePath()
	cv.Fill()

	cv.SetFillStyle(filepath.FromSlash("texture/vehicle/boat_side.png"))
	cv.BeginPath()
	cv.LineTo(x-c.f1*pm.XX+0*pm.XY+0*pm.XZ, y-c.f1*pm.YX+0*pm.YY+0*pm.YZ)
	cv.LineTo(x-c.f2*pm.XX+0*pm.XY-c.z*pm.XZ, y-c.f2*pm.YX+0*pm.YY-c.z*pm.YZ)
	cv.LineTo(x+c.f2*pm.XX+0*pm.XY-c.z*pm.XZ, y+c.f2*pm.YX+0*pm.YY-c.z*pm.YZ)
	cv.LineTo(x+c.f1*pm.XX+0*pm.XY+0*pm.XZ, y+c.f1*pm.YX+0*pm.YY+0*pm.YZ)
	cv.LineTo(x+c.f1*pm.XX-c.h1*pm.XY+0*pm.XZ, y+c.f1*pm.YX-c.h1*pm.YY+0*pm.YZ)
	cv.LineTo(x+c.f2*pm.XX-c.h2*pm.XY-c.s*pm.XZ, y+c.f2*pm.YX-c.h2*pm.YY-c.s*pm.YZ)
	cv.LineTo(x-c.f2*pm.XX-c.h2*pm.XY-c.s*pm.XZ, y-c.f2*pm.YX-c.h2*pm.YY-c.s*pm.YZ)
	cv.LineTo(x-c.f1*pm.XX-c.h1*pm.XY+0*pm.XZ, y-c.f1*pm.YX-c.h1*pm.YY+0*pm.YZ)
	cv.ClosePath()
	cv.Fill()

	cv.SetFillStyle(filepath.FromSlash("texture/vehicle/boat_side.png"))
	cv.BeginPath()
	cv.LineTo(x-c.f1*pm.XX+0*pm.XY+0*pm.XZ, y-c.f1*pm.YX+0*pm.YY+0*pm.YZ)
	cv.LineTo(x-c.f2*pm.XX+0*pm.XY+c.z*pm.XZ, y-c.f2*pm.YX+0*pm.YY+c.z*pm.YZ)
	cv.LineTo(x+c.f2*pm.XX+0*pm.XY+c.z*pm.XZ, y+c.f2*pm.YX+0*pm.YY+c.z*pm.YZ)
	cv.LineTo(x+c.f1*pm.XX+0*pm.XY+0*pm.XZ, y+c.f1*pm.YX+0*pm.YY+0*pm.YZ)
	cv.LineTo(x+c.f1*pm.XX-c.h1*pm.XY+0*pm.XZ, y+c.f1*pm.YX-c.h1*pm.YY+0*pm.YZ)
	cv.LineTo(x+c.f2*pm.XX-c.h2*pm.XY+c.s*pm.XZ, y+c.f2*pm.YX-c.h2*pm.YY+c.s*pm.YZ)
	cv.LineTo(x-c.f2*pm.XX-c.h2*pm.XY+c.s*pm.XZ, y-c.f2*pm.YX-c.h2*pm.YY+c.s*pm.YZ)
	cv.LineTo(x-c.f1*pm.XX-c.h1*pm.XY+0*pm.XZ, y-c.f1*pm.YX-c.h1*pm.YY+0*pm.YZ)
	cv.ClosePath()
	cv.Fill()

	// Paddles
	cv.SetStrokeStyle("#321")
	cv.SetLineWidth(2)
	cv.BeginPath()
	cv.MoveTo(x+0*pm.XX-c.h2*pm.XY+c.s*pm.XZ, y+0*pm.YX-c.h2*pm.YY+c.s*pm.YZ)
	cv.LineTo(x+r*pm.XX-0*pm.XY+c.p*pm.XZ, y+r*pm.YX-0*pm.YY+c.p*pm.YZ)
	cv.ClosePath()
	cv.Stroke()
	cv.BeginPath()
	cv.MoveTo(x+0*pm.XX-c.h2*pm.XY-c.s*pm.XZ, y+0*pm.YX-c.h2*pm.YY-c.s*pm.YZ)
	cv.LineTo(x+r*pm.XX-0*pm.XY-c.p*pm.XZ, y+r*pm.YX-0*pm.YY-c.p*pm.YZ)
	cv.ClosePath()
	cv.Stroke()
}

func DrawExpeditionBoat(cv *canvas.Canvas, t *navigation.Traveller, x float64, y float64, c *controller.Controller) {
	if !t.Visible {
		return
	}
	p := t.DrawingPhase()
	dirIdx := (c.Perspective - t.Direction) % 4
	pm := animation.ProjectionMatrices[dirIdx]

	f1 := 48.0
	f2 := 16.0
	s := 18.0
	z := 14.0
	h1 := 15.0
	h2 := 24.0

	cv.SetFillStyle(filepath.FromSlash("texture/vehicle/boat_bottom.png"))
	cv.BeginPath()
	cv.LineTo(x-f1*pm.XX+0*pm.XY+0*pm.XZ, y-f1*pm.YX+0*pm.YY+0*pm.YZ)
	cv.LineTo(x-f2*pm.XX+0*pm.XY-z*pm.XZ, y-f2*pm.YX+0*pm.YY-z*pm.YZ)
	cv.LineTo(x+f2*pm.XX+0*pm.XY-z*pm.XZ, y+f2*pm.YX+0*pm.YY-z*pm.YZ)
	cv.LineTo(x+f1*pm.XX+0*pm.XY+0*pm.XZ, y+f1*pm.YX+0*pm.YY+0*pm.YZ)
	cv.LineTo(x+f2*pm.XX+0*pm.XY+z*pm.XZ, y+f2*pm.YX+0*pm.YY+z*pm.YZ)
	cv.LineTo(x-f2*pm.XX+0*pm.XY+z*pm.XZ, y-f2*pm.YX+0*pm.YY+z*pm.YZ)
	cv.ClosePath()
	cv.Fill()

	cv.SetFillStyle(filepath.FromSlash("texture/vehicle/boat_side.png"))
	cv.BeginPath()
	cv.LineTo(x-f1*pm.XX+0*pm.XY+0*pm.XZ, y-f1*pm.YX+0*pm.YY+0*pm.YZ)
	cv.LineTo(x-f2*pm.XX+0*pm.XY-z*pm.XZ, y-f2*pm.YX+0*pm.YY-z*pm.YZ)
	cv.LineTo(x+f2*pm.XX+0*pm.XY-z*pm.XZ, y+f2*pm.YX+0*pm.YY-z*pm.YZ)
	cv.LineTo(x+f1*pm.XX+0*pm.XY+0*pm.XZ, y+f1*pm.YX+0*pm.YY+0*pm.YZ)
	cv.LineTo(x+f1*pm.XX-h2*pm.XY+0*pm.XZ, y+f1*pm.YX-h2*pm.YY+0*pm.YZ)
	cv.LineTo(x+f2*pm.XX-h1*pm.XY-s*pm.XZ, y+f2*pm.YX-h1*pm.YY-s*pm.YZ)
	cv.LineTo(x-f2*pm.XX-h1*pm.XY-s*pm.XZ, y-f2*pm.YX-h1*pm.YY-s*pm.YZ)
	cv.LineTo(x-f1*pm.XX-h2*pm.XY+0*pm.XZ, y-f1*pm.YX-h2*pm.YY+0*pm.YZ)
	cv.ClosePath()
	cv.Fill()

	cv.SetFillStyle(filepath.FromSlash("texture/vehicle/boat_side.png"))
	cv.BeginPath()
	cv.LineTo(x-f1*pm.XX+0*pm.XY+0*pm.XZ, y-f1*pm.YX+0*pm.YY+0*pm.YZ)
	cv.LineTo(x-f2*pm.XX+0*pm.XY+z*pm.XZ, y-f2*pm.YX+0*pm.YY+z*pm.YZ)
	cv.LineTo(x+f2*pm.XX+0*pm.XY+z*pm.XZ, y+f2*pm.YX+0*pm.YY+z*pm.YZ)
	cv.LineTo(x+f1*pm.XX+0*pm.XY+0*pm.XZ, y+f1*pm.YX+0*pm.YY+0*pm.YZ)
	cv.LineTo(x+f1*pm.XX-h2*pm.XY+0*pm.XZ, y+f1*pm.YX-h2*pm.YY+0*pm.YZ)
	cv.LineTo(x+f2*pm.XX-h1*pm.XY+s*pm.XZ, y+f2*pm.YX-h1*pm.YY+s*pm.YZ)
	cv.LineTo(x-f2*pm.XX-h1*pm.XY+s*pm.XZ, y-f2*pm.YX-h1*pm.YY+s*pm.YZ)
	cv.LineTo(x-f1*pm.XX-h2*pm.XY+0*pm.XZ, y-f1*pm.YX-h2*pm.YY+0*pm.YZ)
	cv.ClosePath()
	cv.Fill()

	mw1 := 3.0
	mw2 := 2.0
	mh := 64.0
	sh1 := 20.0
	sh2 := 52.0
	sw := 18.0

	if dirIdx == 1 || dirIdx == 2 {
		drawSail(cv, pm, p, x, y, sh1, sh2, sw)
		drawMast(cv, pm, x, y, mw1, mw2, mh, sh1, sh2, sw)
	} else {
		drawMast(cv, pm, x, y, mw1, mw2, mh, sh1, sh2, sw)
		drawSail(cv, pm, p, x, y, sh1, sh2, sw)
	}
}

func drawMast(cv *canvas.Canvas, pm animation.ProjectionMatrix, x, y, mw1, mw2, mh, sh1, sh2, sw float64) {
	cv.SetFillStyle(filepath.FromSlash("texture/vehicle/boat_side.png"))
	cv.BeginPath()
	cv.LineTo(x-mw1, y)
	cv.LineTo(x+mw1, y)
	cv.LineTo(x+mw2, y-mh)
	cv.LineTo(x-mw2, y-mh)
	cv.ClosePath()
	cv.Fill()

	h1 := sh1 + 4.0
	h2 := sh2 - 4.0
	w := sw + 4.0
	cv.BeginPath()
	cv.LineTo(x+0*pm.XX-h1*pm.XY+w*pm.XZ, y+0*pm.YX-h1*pm.YY+w*pm.YZ+1)
	cv.LineTo(x+0*pm.XX-h1*pm.XY+w*pm.XZ, y+0*pm.YX-h1*pm.YY+w*pm.YZ-2)
	cv.LineTo(x+0*pm.XX-h1*pm.XY-w*pm.XZ, y+0*pm.YX-h1*pm.YY-w*pm.YZ-2)
	cv.LineTo(x+0*pm.XX-h1*pm.XY-w*pm.XZ, y+0*pm.YX-h1*pm.YY-w*pm.YZ+1)
	cv.ClosePath()
	cv.Fill()

	cv.BeginPath()
	cv.LineTo(x+0*pm.XX-h2*pm.XY+w*pm.XZ, y+0*pm.YX-h2*pm.YY+w*pm.YZ+1)
	cv.LineTo(x+0*pm.XX-h2*pm.XY+w*pm.XZ, y+0*pm.YX-h2*pm.YY+w*pm.YZ-2)
	cv.LineTo(x+0*pm.XX-h2*pm.XY-w*pm.XZ, y+0*pm.YX-h2*pm.YY-w*pm.YZ-2)
	cv.LineTo(x+0*pm.XX-h2*pm.XY-w*pm.XZ, y+0*pm.YX-h2*pm.YY-w*pm.YZ+1)
	cv.ClosePath()
	cv.Fill()
}

func drawSail(cv *canvas.Canvas, pm animation.ProjectionMatrix, p uint8, x, y, sh1, sh2, sw float64) {
	p2 := (p + 1) % 8
	cv.SetFillStyle(filepath.FromSlash("texture/vehicle/boat_sail.png"))
	cv.BeginPath()
	for i := 0.0; i <= 6.0; i++ {
		h := (sh1*i + sh2*(6-i)) / 6
		dx := math.Sin(i*math.Pi/6) * (math.Sin(float64(p)*math.Pi/8)*2.0 + 2.0)
		cv.LineTo(x+dx*pm.XX-h*pm.XY+sw*pm.XZ, y+dx*pm.YX-h*pm.YY+sw*pm.YZ)
	}
	for i := 6.0; i >= 0.0; i-- {
		h := (sh1*i + sh2*(6-i)) / 6
		dx := math.Sin(i*math.Pi/6) * (math.Sin(float64(p2)*math.Pi/8)*2.0 + 2.0)
		cv.LineTo(x+dx*pm.XX-h*pm.XY-sw*pm.XZ, y+dx*pm.YX-h*pm.YY-sw*pm.YZ)
	}
	cv.ClosePath()
	cv.Fill()
}
