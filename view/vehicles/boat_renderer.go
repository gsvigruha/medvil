package vehicles

import (
	//"fmt"
	"github.com/tfriedel6/canvas"
	"math"
	"medvil/controller"
	"medvil/model/navigation"
	"medvil/view/animation"
)

func DrawTradingBoat(cv *canvas.Canvas, t *navigation.Traveller, x float64, y float64, c *controller.Controller) {
	f1, f2, s, _, h, _ := DrawBoat(cv, t, x, y, c)

	dirIdx := (c.Perspective - t.Direction) % 4
	pm := animation.ProjectionMatrices[dirIdx]

	if dirIdx == 0 || dirIdx == 2 {
		cv.SetFillStyle("texture/vehicle/textile.png")
	} else {
		cv.SetFillStyle("texture/vehicle/textile_flipped.png")
	}
	cv.BeginPath()
	cv.LineTo(x+f1*pm.XX-h*pm.XY+0*pm.XZ, y+f1*pm.YX-h*pm.YY+0*pm.YZ)
	cv.LineTo(x+f2*pm.XX-h*pm.XY-s*pm.XZ, y+f2*pm.YX-h*pm.YY-s*pm.YZ)
	cv.LineTo(x+f2*pm.XX-h*pm.XY+s*pm.XZ, y+f2*pm.YX-h*pm.YY+s*pm.YZ)
	cv.ClosePath()
	cv.Fill()
}

func DrawBoat(cv *canvas.Canvas, t *navigation.Traveller, x float64, y float64, c *controller.Controller) (float64, float64, float64, float64, float64, float64) {
	if !t.Visible {
		return 0, 0, 0, 0, 0, 0
	}
	dirIdx := (c.Perspective - t.Direction) % 4
	pm := animation.ProjectionMatrices[dirIdx]
	var r = 8.0 - float64(t.Phase%16)
	if r < 0.0 {
		r = -r
	}

	f1 := 20.0
	f2 := 6.0
	s := 8.0
	z := 6.0
	h := 6.0
	p := 16.0

	cv.SetFillStyle("texture/vehicle/boat_bottom.png")
	cv.BeginPath()
	cv.LineTo(x-f1*pm.XX+0*pm.XY+0*pm.XZ, y-f1*pm.YX+0*pm.YY+0*pm.YZ)
	cv.LineTo(x-f2*pm.XX+0*pm.XY-z*pm.XZ, y-f2*pm.YX+0*pm.YY-z*pm.YZ)
	cv.LineTo(x+f2*pm.XX+0*pm.XY-z*pm.XZ, y+f2*pm.YX+0*pm.YY-z*pm.YZ)
	cv.LineTo(x+f1*pm.XX+0*pm.XY+0*pm.XZ, y+f1*pm.YX+0*pm.YY+0*pm.YZ)
	cv.LineTo(x+f2*pm.XX+0*pm.XY+z*pm.XZ, y+f2*pm.YX+0*pm.YY+z*pm.YZ)
	cv.LineTo(x-f2*pm.XX+0*pm.XY+z*pm.XZ, y-f2*pm.YX+0*pm.YY+z*pm.YZ)
	cv.ClosePath()
	cv.Fill()

	cv.SetFillStyle("texture/vehicle/boat_side.png")
	cv.BeginPath()
	cv.LineTo(x-f1*pm.XX+0*pm.XY+0*pm.XZ, y-f1*pm.YX+0*pm.YY+0*pm.YZ)
	cv.LineTo(x-f2*pm.XX+0*pm.XY-z*pm.XZ, y-f2*pm.YX+0*pm.YY-z*pm.YZ)
	cv.LineTo(x+f2*pm.XX+0*pm.XY-z*pm.XZ, y+f2*pm.YX+0*pm.YY-z*pm.YZ)
	cv.LineTo(x+f1*pm.XX+0*pm.XY+0*pm.XZ, y+f1*pm.YX+0*pm.YY+0*pm.YZ)
	cv.LineTo(x+f1*pm.XX-h*pm.XY+0*pm.XZ, y+f1*pm.YX-h*pm.YY+0*pm.YZ)
	cv.LineTo(x+f2*pm.XX-h*pm.XY-s*pm.XZ, y+f2*pm.YX-h*pm.YY-s*pm.YZ)
	cv.LineTo(x-f2*pm.XX-h*pm.XY-s*pm.XZ, y-f2*pm.YX-h*pm.YY-s*pm.YZ)
	cv.LineTo(x-f1*pm.XX-h*pm.XY+0*pm.XZ, y-f1*pm.YX-h*pm.YY+0*pm.YZ)
	cv.ClosePath()
	cv.Fill()

	cv.SetFillStyle("texture/vehicle/boat_side.png")
	cv.BeginPath()
	cv.LineTo(x-f1*pm.XX+0*pm.XY+0*pm.XZ, y-f1*pm.YX+0*pm.YY+0*pm.YZ)
	cv.LineTo(x-f2*pm.XX+0*pm.XY+z*pm.XZ, y-f2*pm.YX+0*pm.YY+z*pm.YZ)
	cv.LineTo(x+f2*pm.XX+0*pm.XY+z*pm.XZ, y+f2*pm.YX+0*pm.YY+z*pm.YZ)
	cv.LineTo(x+f1*pm.XX+0*pm.XY+0*pm.XZ, y+f1*pm.YX+0*pm.YY+0*pm.YZ)
	cv.LineTo(x+f1*pm.XX-h*pm.XY+0*pm.XZ, y+f1*pm.YX-h*pm.YY+0*pm.YZ)
	cv.LineTo(x+f2*pm.XX-h*pm.XY+s*pm.XZ, y+f2*pm.YX-h*pm.YY+s*pm.YZ)
	cv.LineTo(x-f2*pm.XX-h*pm.XY+s*pm.XZ, y-f2*pm.YX-h*pm.YY+s*pm.YZ)
	cv.LineTo(x-f1*pm.XX-h*pm.XY+0*pm.XZ, y-f1*pm.YX-h*pm.YY+0*pm.YZ)
	cv.ClosePath()
	cv.Fill()

	cv.SetStrokeStyle("#321")
	cv.SetLineWidth(2)
	cv.BeginPath()
	cv.MoveTo(x+0*pm.XX-h*pm.XY+s*pm.XZ, y+0*pm.YX-h*pm.YY+s*pm.YZ)
	cv.LineTo(x+r*pm.XX-0*pm.XY+p*pm.XZ, y+r*pm.YX-0*pm.YY+p*pm.YZ)
	cv.ClosePath()
	cv.Stroke()
	cv.BeginPath()
	cv.MoveTo(x+0*pm.XX-h*pm.XY-s*pm.XZ, y+0*pm.YX-h*pm.YY-s*pm.YZ)
	cv.LineTo(x+r*pm.XX-0*pm.XY-p*pm.XZ, y+r*pm.YX-0*pm.YY-p*pm.YZ)
	cv.ClosePath()
	cv.Stroke()

	return f1, f2, s, z, h, p
}

func DrawExpeditionBoat(cv *canvas.Canvas, t *navigation.Traveller, x float64, y float64, c *controller.Controller) {
	if !t.Visible {
		return
	}
	p := t.DrawingPhase()
	dirIdx := (c.Perspective - t.Direction) % 4
	pm := animation.ProjectionMatrices[dirIdx]

	f1 := 40.0
	f2 := 12.0
	s := 16.0
	z := 12.0
	h1 := 12.0
	h2 := 18.0

	cv.SetFillStyle("texture/vehicle/boat_bottom.png")
	cv.BeginPath()
	cv.LineTo(x-f1*pm.XX+0*pm.XY+0*pm.XZ, y-f1*pm.YX+0*pm.YY+0*pm.YZ)
	cv.LineTo(x-f2*pm.XX+0*pm.XY-z*pm.XZ, y-f2*pm.YX+0*pm.YY-z*pm.YZ)
	cv.LineTo(x+f2*pm.XX+0*pm.XY-z*pm.XZ, y+f2*pm.YX+0*pm.YY-z*pm.YZ)
	cv.LineTo(x+f1*pm.XX+0*pm.XY+0*pm.XZ, y+f1*pm.YX+0*pm.YY+0*pm.YZ)
	cv.LineTo(x+f2*pm.XX+0*pm.XY+z*pm.XZ, y+f2*pm.YX+0*pm.YY+z*pm.YZ)
	cv.LineTo(x-f2*pm.XX+0*pm.XY+z*pm.XZ, y-f2*pm.YX+0*pm.YY+z*pm.YZ)
	cv.ClosePath()
	cv.Fill()

	cv.SetFillStyle("texture/vehicle/boat_side.png")
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

	cv.SetFillStyle("texture/vehicle/boat_side.png")
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
	mh := 54.0
	sh1 := 20.0
	sh2 := 48.0
	sw := 16.0

	if dirIdx == 1 || dirIdx == 2 {
		drawSail(cv, pm, p, x, y, sh1, sh2, sw)
		drawMast(cv, pm, x, y, mw1, mw2, mh)
	} else {
		drawMast(cv, pm, x, y, mw1, mw2, mh)
		drawSail(cv, pm, p, x, y, sh1, sh2, sw)
	}
}

func drawMast(cv *canvas.Canvas, pm animation.ProjectionMatrix, x, y, mw1, mw2, mh float64) {
	cv.SetFillStyle("texture/vehicle/boat_side.png")
	cv.BeginPath()
	cv.LineTo(x-mw1, y)
	cv.LineTo(x+mw1, y)
	cv.LineTo(x+mw2, y-mh)
	cv.LineTo(x-mw2, y-mh)
	cv.ClosePath()
	cv.Fill()
}

func drawSail(cv *canvas.Canvas, pm animation.ProjectionMatrix, p uint8, x, y, sh1, sh2, sw float64) {
	p2 := (p + 1) % 8
	cv.SetFillStyle("texture/vehicle/boat_sail.png")
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
