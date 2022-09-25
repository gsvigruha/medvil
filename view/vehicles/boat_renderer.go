package vehicles

import (
	//"fmt"
	"github.com/tfriedel6/canvas"
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
