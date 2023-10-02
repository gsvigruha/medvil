package vehicles

import (
	//"fmt"
	"github.com/tfriedel6/canvas"
	"math"
	"medvil/controller"
	"medvil/model/navigation"
	"medvil/view/animation"
)

func DrawTradingCart(cv *canvas.Canvas, t *navigation.Traveller, x float64, y float64, c *controller.Controller) {
	f1, f2, z, _, h2 := DrawCart(cv, t, x, y, c)

	dirIdx := (c.Perspective - t.Direction) % 4
	pm := animation.ProjectionMatrices[dirIdx]

	h3 := 21.0

	cv.SetStrokeStyle("#321")
	cv.SetLineWidth(1)

	cv.BeginPath()
	cv.LineTo(x+f1*pm.XX-h2*pm.XY-z*pm.XZ, y+f1*pm.YX-h2*pm.YY-z*pm.YZ)
	cv.LineTo(x+f1*pm.XX-h3*pm.XY-z*pm.XZ, y+f1*pm.YX-h3*pm.YY-z*pm.YZ)
	cv.ClosePath()
	cv.Stroke()

	cv.BeginPath()
	cv.LineTo(x+f1*pm.XX-h2*pm.XY+z*pm.XZ, y+f1*pm.YX-h2*pm.YY+z*pm.YZ)
	cv.LineTo(x+f1*pm.XX-h3*pm.XY+z*pm.XZ, y+f1*pm.YX-h3*pm.YY+z*pm.YZ)
	cv.ClosePath()
	cv.Stroke()

	cv.BeginPath()
	cv.LineTo(x+f2*pm.XX-h2*pm.XY+z*pm.XZ, y+f2*pm.YX-h2*pm.YY+z*pm.YZ)
	cv.LineTo(x+f2*pm.XX-h3*pm.XY+z*pm.XZ, y+f2*pm.YX-h3*pm.YY+z*pm.YZ)
	cv.ClosePath()
	cv.Stroke()

	cv.BeginPath()
	cv.LineTo(x+f2*pm.XX-h2*pm.XY-z*pm.XZ, y+f2*pm.YX-h2*pm.YY-z*pm.YZ)
	cv.LineTo(x+f2*pm.XX-h3*pm.XY-z*pm.XZ, y+f2*pm.YX-h3*pm.YY-z*pm.YZ)
	cv.ClosePath()
	cv.Stroke()

	if dirIdx == 1 || dirIdx == 3 {
		cv.SetFillStyle("texture/vehicle/textile.png")
	} else {
		cv.SetFillStyle("texture/vehicle/textile_flipped.png")
	}
	cv.BeginPath()
	cv.LineTo(x+f1*pm.XX-h3*pm.XY-z*pm.XZ, y+f1*pm.YX-h3*pm.YY-z*pm.YZ)
	cv.LineTo(x+f1*pm.XX-h3*pm.XY+z*pm.XZ, y+f1*pm.YX-h3*pm.YY+z*pm.YZ)
	cv.LineTo(x+f2*pm.XX-h3*pm.XY+z*pm.XZ, y+f2*pm.YX-h3*pm.YY+z*pm.YZ)
	cv.LineTo(x+f2*pm.XX-h3*pm.XY-z*pm.XZ, y+f2*pm.YX-h3*pm.YY-z*pm.YZ)
	cv.ClosePath()
	cv.Fill()
}

func DrawCart(cv *canvas.Canvas, t *navigation.Traveller, x float64, y float64, c *controller.Controller) (float64, float64, float64, float64, float64) {
	if !t.Visible {
		return 0, 0, 0, 0, 0
	}
	dirIdx := (c.Perspective - t.Direction) % 4
	pm := animation.ProjectionMatrices[dirIdx]
	var r = 8.0 - float64(t.Phase%16)
	if r < 0.0 {
		r = -r
	}
	var dir = 1.0
	if dirIdx == 0 || dirIdx == 1 {
		dir = -1.0
	}

	f1 := 3.0
	f2 := 17.0
	z := 6.0
	h1 := 8.0
	h2 := 12.0

	cv.SetFillStyle("texture/vehicle/boat_bottom.png")
	cv.SetStrokeStyle("#321")
	cv.SetLineWidth(1)
	cv.BeginPath()
	for i := 0.0; i < 8; i++ {
		dx0 := math.Cos(math.Pi*2.0*i/8.0)*(h1/2.0) + (f2-f1)/2.0
		dy0 := math.Sin(math.Pi*2.0*i/8.0)*(h1/2.0) + h1/2.0
		cv.LineTo(x+dx0*pm.XX-dy0*pm.XY-z*pm.XZ*dir, y+dx0*pm.YX-dy0*pm.YY-z*pm.YZ*dir)
	}
	cv.ClosePath()
	cv.Fill()
	cv.Stroke()

	cv.SetFillStyle("texture/vehicle/boat_bottom.png")
	cv.BeginPath()
	cv.LineTo(x+f1*pm.XX-h1*pm.XY-z*pm.XZ, y+f1*pm.YX-h1*pm.YY-z*pm.YZ)
	cv.LineTo(x+f1*pm.XX-h1*pm.XY+z*pm.XZ, y+f1*pm.YX-h1*pm.YY+z*pm.YZ)
	cv.LineTo(x+f2*pm.XX-h1*pm.XY+z*pm.XZ, y+f2*pm.YX-h1*pm.YY+z*pm.YZ)
	cv.LineTo(x+f2*pm.XX-h1*pm.XY-z*pm.XZ, y+f2*pm.YX-h1*pm.YY-z*pm.YZ)
	cv.ClosePath()
	cv.Fill()

	cv.SetStrokeStyle("#321")
	cv.SetLineWidth(2)
	cv.BeginPath()
	cv.LineTo(x+f1*pm.XX-h2*pm.XY-z*pm.XZ, y+f1*pm.YX-h2*pm.YY-z*pm.YZ)
	cv.LineTo(x+f1*pm.XX-h2*pm.XY+z*pm.XZ, y+f1*pm.YX-h2*pm.YY+z*pm.YZ)
	cv.LineTo(x+f2*pm.XX-h2*pm.XY+z*pm.XZ, y+f2*pm.YX-h2*pm.YY+z*pm.YZ)
	cv.LineTo(x+f2*pm.XX-h2*pm.XY-z*pm.XZ, y+f2*pm.YX-h2*pm.YY-z*pm.YZ)
	cv.ClosePath()
	cv.Stroke()

	cv.SetFillStyle("texture/vehicle/boat_bottom.png")
	cv.SetStrokeStyle("#321")
	cv.SetLineWidth(1)
	cv.BeginPath()
	for i := 0.0; i < 8; i++ {
		dx0 := math.Cos(math.Pi*2.0*i/8.0)*(h1/2.0) + (f2-f1)/2.0
		dy0 := math.Sin(math.Pi*2.0*i/8.0)*(h1/2.0) + h1/2.0
		cv.LineTo(x+dx0*pm.XX-dy0*pm.XY+z*pm.XZ*dir, y+dx0*pm.YX-dy0*pm.YY+z*pm.YZ*dir)
	}
	cv.ClosePath()
	cv.Fill()
	cv.Stroke()

	return f1, f2, z, h1, h2
}

func drawWheel(cv *canvas.Canvas, pm animation.ProjectionMatrix, phase uint8, x, y, f, h1, z, wr, dir float64) {
	cv.SetStrokeStyle("#321")
	cv.SetLineWidth(3)
	cv.BeginPath()
	for i := 0.0; i < 8; i++ {
		p := (i*4 + float64(phase)) / 32.0
		dx0 := math.Cos(math.Pi*2.0*p)*wr + f
		dy0 := math.Sin(math.Pi*2.0*p)*wr + h1/2.0
		cv.LineTo(x+dx0*pm.XX-dy0*pm.XY-z*pm.XZ*dir, y+dx0*pm.YX-dy0*pm.YY-z*pm.YZ*dir)
	}
	cv.ClosePath()
	cv.Stroke()

	cv.SetLineWidth(1)
	for i := 0.0; i < 8; i++ {
		p := (i*4 + float64(phase)) / 32.0
		dx0 := math.Cos(math.Pi*2.0*p)*wr + f
		dy0 := math.Sin(math.Pi*2.0*p)*wr + h1/2.0
		cv.BeginPath()
		cv.MoveTo(x+f*pm.XX-h1/2.0*pm.XY-z*pm.XZ*dir, y+f*pm.YX-h1/2.0*pm.YY-z*pm.YZ*dir)
		cv.LineTo(x+dx0*pm.XX-dy0*pm.XY-z*pm.XZ*dir, y+dx0*pm.YX-dy0*pm.YY-z*pm.YZ*dir)
		cv.ClosePath()
		cv.Stroke()
	}
}

func DrawExpeditionCart(cv *canvas.Canvas, t *navigation.Traveller, x float64, y float64, c *controller.Controller) {
	if !t.Visible {
		return
	}
	p := t.DrawingPhase()
	dirIdx := (c.Perspective - t.Direction) % 4
	pm := animation.ProjectionMatrices[dirIdx]
	var r = 8.0 - float64(t.Phase%16)
	if r < 0.0 {
		r = -r
	}
	var dir = 1.0
	if dirIdx == 0 || dirIdx == 1 {
		dir = -1.0
	}

	f1 := 24.0
	f2 := -24.0
	f3 := 12.0
	z := 16.0
	z3 := 12.0
	h1 := 12.0
	h2 := 24.0
	h3 := 40.0
	h4 := 44.0

	wr := h1 * 3 / 4

	drawWheel(cv, pm, p, x, y, f2*0.8+f1*0.2, h1, z, wr, dir)
	drawWheel(cv, pm, p, x, y, f2*0.2+f1*0.8, h1, z, wr, dir)

	cv.SetFillStyle("texture/vehicle/boat_bottom.png")
	cv.BeginPath()
	cv.LineTo(x+f1*pm.XX-h1*pm.XY-z*pm.XZ, y+f1*pm.YX-h1*pm.YY-z*pm.YZ)
	cv.LineTo(x+f1*pm.XX-h1*pm.XY+z*pm.XZ, y+f1*pm.YX-h1*pm.YY+z*pm.YZ)
	cv.LineTo(x+f2*pm.XX-h1*pm.XY+z*pm.XZ, y+f2*pm.YX-h1*pm.YY+z*pm.YZ)
	cv.LineTo(x+f2*pm.XX-h1*pm.XY-z*pm.XZ, y+f2*pm.YX-h1*pm.YY-z*pm.YZ)
	cv.ClosePath()
	cv.Fill()

	cv.SetFillStyle("texture/vehicle/boat_side.png")
	cv.BeginPath()
	cv.LineTo(x+f1*pm.XX-h1*pm.XY-z*pm.XZ, y+f1*pm.YX-h1*pm.YY-z*pm.YZ)
	cv.LineTo(x+f1*pm.XX-h1*pm.XY+z*pm.XZ, y+f1*pm.YX-h1*pm.YY+z*pm.YZ)
	cv.LineTo(x+f1*pm.XX-h2*pm.XY+z*pm.XZ, y+f1*pm.YX-h2*pm.YY+z*pm.YZ)
	cv.LineTo(x+f1*pm.XX-h2*pm.XY-z*pm.XZ, y+f1*pm.YX-h2*pm.YY-z*pm.YZ)
	cv.ClosePath()
	cv.Fill()
	cv.BeginPath()
	cv.LineTo(x+f2*pm.XX-h1*pm.XY-z*pm.XZ, y+f2*pm.YX-h1*pm.YY-z*pm.YZ)
	cv.LineTo(x+f2*pm.XX-h1*pm.XY+z*pm.XZ, y+f2*pm.YX-h1*pm.YY+z*pm.YZ)
	cv.LineTo(x+f2*pm.XX-h2*pm.XY+z*pm.XZ, y+f2*pm.YX-h2*pm.YY+z*pm.YZ)
	cv.LineTo(x+f2*pm.XX-h2*pm.XY-z*pm.XZ, y+f2*pm.YX-h2*pm.YY-z*pm.YZ)
	cv.ClosePath()
	cv.Fill()
	cv.BeginPath()
	cv.LineTo(x+f2*pm.XX-h1*pm.XY+z*pm.XZ, y+f2*pm.YX-h1*pm.YY+z*pm.YZ)
	cv.LineTo(x+f1*pm.XX-h1*pm.XY+z*pm.XZ, y+f1*pm.YX-h1*pm.YY+z*pm.YZ)
	cv.LineTo(x+f1*pm.XX-h2*pm.XY+z*pm.XZ, y+f1*pm.YX-h2*pm.YY+z*pm.YZ)
	cv.LineTo(x+f2*pm.XX-h2*pm.XY+z*pm.XZ, y+f2*pm.YX-h2*pm.YY+z*pm.YZ)
	cv.ClosePath()
	cv.Fill()
	cv.BeginPath()
	cv.LineTo(x+f2*pm.XX-h1*pm.XY-z*pm.XZ, y+f2*pm.YX-h1*pm.YY-z*pm.YZ)
	cv.LineTo(x+f1*pm.XX-h1*pm.XY-z*pm.XZ, y+f1*pm.YX-h1*pm.YY-z*pm.YZ)
	cv.LineTo(x+f1*pm.XX-h2*pm.XY-z*pm.XZ, y+f1*pm.YX-h2*pm.YY-z*pm.YZ)
	cv.LineTo(x+f2*pm.XX-h2*pm.XY-z*pm.XZ, y+f2*pm.YX-h2*pm.YY-z*pm.YZ)
	cv.ClosePath()
	cv.Fill()

	cv.SetStrokeStyle("#321")
	cv.SetLineWidth(2)
	cv.BeginPath()
	cv.LineTo(x+f1*pm.XX-h2*pm.XY-z*pm.XZ, y+f1*pm.YX-h2*pm.YY-z*pm.YZ)
	cv.LineTo(x+f1*pm.XX-h2*pm.XY+z*pm.XZ, y+f1*pm.YX-h2*pm.YY+z*pm.YZ)
	cv.LineTo(x+f2*pm.XX-h2*pm.XY+z*pm.XZ, y+f2*pm.YX-h2*pm.YY+z*pm.YZ)
	cv.LineTo(x+f2*pm.XX-h2*pm.XY-z*pm.XZ, y+f2*pm.YX-h2*pm.YY-z*pm.YZ)
	cv.ClosePath()
	cv.Stroke()

	cv.SetFillStyle("texture/vehicle/leather.png")
	cv.SetStrokeStyle("#432")
	cv.SetLineWidth(2)
	// side
	cv.BeginPath()
	cv.LineTo(x+f2*pm.XX-h2*pm.XY-z*pm.XZ*dir, y+f2*pm.YX-h2*pm.YY-z*pm.YZ*dir)
	cv.LineTo(x+f3*pm.XX-h2*pm.XY-z*pm.XZ*dir, y+f3*pm.YX-h2*pm.YY-z*pm.YZ*dir)
	cv.LineTo(x+f3*pm.XX-h3*pm.XY-z*pm.XZ*dir, y+f3*pm.YX-h3*pm.YY-z*pm.YZ*dir)
	cv.LineTo(x+f2*pm.XX-h3*pm.XY-z*pm.XZ*dir, y+f2*pm.YX-h3*pm.YY-z*pm.YZ*dir)
	cv.ClosePath()
	cv.Fill()
	cv.Stroke()
	cv.BeginPath()
	cv.LineTo(x+f3*pm.XX-h3*pm.XY-z*pm.XZ*dir, y+f3*pm.YX-h3*pm.YY-z*pm.YZ*dir)
	cv.LineTo(x+f2*pm.XX-h3*pm.XY-z*pm.XZ*dir, y+f2*pm.YX-h3*pm.YY-z*pm.YZ*dir)
	cv.LineTo(x+f2*pm.XX-h4*pm.XY-z3*pm.XZ*dir, y+f2*pm.YX-h4*pm.YY-z3*pm.YZ*dir)
	cv.LineTo(x+f3*pm.XX-h4*pm.XY-z3*pm.XZ*dir, y+f3*pm.YX-h4*pm.YY-z3*pm.YZ*dir)
	cv.ClosePath()
	cv.Fill()
	cv.Stroke()
	// back
	cv.BeginPath()
	cv.LineTo(x+f2*pm.XX-h4*pm.XY-z3*pm.XZ*dir, y+f2*pm.YX-h4*pm.YY-z3*pm.YZ*dir)
	cv.LineTo(x+f2*pm.XX-h3*pm.XY-z*pm.XZ*dir, y+f2*pm.YX-h3*pm.YY-z*pm.YZ*dir)
	cv.LineTo(x+f2*pm.XX-h2*pm.XY-z*pm.XZ*dir, y+f2*pm.YX-h2*pm.YY-z*pm.YZ*dir)
	cv.LineTo(x+f2*pm.XX-h2*pm.XY+z*pm.XZ*dir, y+f2*pm.YX-h2*pm.YY+z*pm.YZ*dir)
	cv.LineTo(x+f2*pm.XX-h3*pm.XY+z*pm.XZ*dir, y+f2*pm.YX-h3*pm.YY+z*pm.YZ*dir)
	cv.LineTo(x+f2*pm.XX-h4*pm.XY+z3*pm.XZ*dir, y+f2*pm.YX-h4*pm.YY+z3*pm.YZ*dir)
	cv.ClosePath()
	cv.Fill()
	cv.Stroke()
	// side
	cv.BeginPath()
	cv.LineTo(x+f3*pm.XX-h3*pm.XY+z*pm.XZ*dir, y+f3*pm.YX-h3*pm.YY+z*pm.YZ*dir)
	cv.LineTo(x+f2*pm.XX-h3*pm.XY+z*pm.XZ*dir, y+f2*pm.YX-h3*pm.YY+z*pm.YZ*dir)
	cv.LineTo(x+f2*pm.XX-h4*pm.XY+z3*pm.XZ*dir, y+f2*pm.YX-h4*pm.YY+z3*pm.YZ*dir)
	cv.LineTo(x+f3*pm.XX-h4*pm.XY+z3*pm.XZ*dir, y+f3*pm.YX-h4*pm.YY+z3*pm.YZ*dir)
	cv.ClosePath()
	cv.Fill()
	cv.Stroke()
	cv.BeginPath()
	cv.LineTo(x+f2*pm.XX-h2*pm.XY+z*pm.XZ*dir, y+f2*pm.YX-h2*pm.YY+z*pm.YZ*dir)
	cv.LineTo(x+f3*pm.XX-h2*pm.XY+z*pm.XZ*dir, y+f3*pm.YX-h2*pm.YY+z*pm.YZ*dir)
	cv.LineTo(x+f3*pm.XX-h3*pm.XY+z*pm.XZ*dir, y+f3*pm.YX-h3*pm.YY+z*pm.YZ*dir)
	cv.LineTo(x+f2*pm.XX-h3*pm.XY+z*pm.XZ*dir, y+f2*pm.YX-h3*pm.YY+z*pm.YZ*dir)
	cv.ClosePath()
	cv.Fill()

	// roof
	cv.Stroke()
	cv.BeginPath()
	cv.LineTo(x+f2*pm.XX-h4*pm.XY-z3*pm.XZ*dir, y+f2*pm.YX-h4*pm.YY-z3*pm.YZ*dir)
	cv.LineTo(x+f3*pm.XX-h4*pm.XY-z3*pm.XZ*dir, y+f3*pm.YX-h4*pm.YY-z3*pm.YZ*dir)
	cv.LineTo(x+f3*pm.XX-h4*pm.XY+z3*pm.XZ*dir, y+f3*pm.YX-h4*pm.YY+z3*pm.YZ*dir)
	cv.LineTo(x+f2*pm.XX-h4*pm.XY+z3*pm.XZ*dir, y+f2*pm.YX-h4*pm.YY+z3*pm.YZ*dir)
	cv.ClosePath()
	cv.Fill()
	cv.Stroke()

	drawWheel(cv, pm, p, x, y, f2*0.8+f1*0.2, h1, z, wr, -dir)
	drawWheel(cv, pm, p, x, y, f2*0.2+f1*0.8, h1, z, wr, -dir)
}
