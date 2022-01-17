package view

import (
	"github.com/tfriedel6/canvas"
	"medvil/controller"
	"medvil/model/navigation"
	"medvil/renderer"
	"math"
	"medvil/view/animation"
)

const MaxPX = navigation.MaxPX
const MaxPY = navigation.MaxPY

func RenderTravellers(cv *canvas.Canvas, travellers []*navigation.Traveller, rf renderer.RenderedField, c *controller.Controller) {
	for i := range travellers {
		t := travellers[i]
		px := float64(t.PX)
		py := float64(t.PY)
		NEPX := rf.X[(2+c.Perspective)%4]
		NEPY := rf.Y[(2+c.Perspective)%4]
		SEPX := rf.X[(1+c.Perspective)%4]
		SEPY := rf.Y[(1+c.Perspective)%4]
		SWPX := rf.X[(0+c.Perspective)%4]
		SWPY := rf.Y[(0+c.Perspective)%4]
		NWPX := rf.X[(3+c.Perspective)%4]
		NWPY := rf.Y[(3+c.Perspective)%4]
		x := (NWPX*(MaxPX-px)*(MaxPY-py) +
			SWPX*(MaxPX-px)*py +
			NEPX*px*(MaxPY-py) +
			SEPX*px*py) / (MaxPX * MaxPY)
		y := (NWPY*(MaxPX-px)*(MaxPY-py) +
			SWPY*(MaxPX-px)*py +
			NEPY*px*(MaxPY-py) +
			SEPY*px*py) / (MaxPX * MaxPY)
		DrawPerson(cv, t, x, y, c)
	}
}

func DrawLimb(cv *canvas.Canvas, pm animation.ProjectionMatrix, x, y, cx1, cy1, cz1, w1, cx2, cy2, cz2, w2 float64) {
	cv.BeginPath()
	pcx1 := x + cx1 * pm.XX + cy1 * pm.XY + cz1 * pm.XZ
	pcy1 := y + cx1 * pm.YX + cy1 * pm.YY + cz1 * pm.YZ
	pcx2 := x + cx2 * pm.XX + cy2 * pm.XY + cz2 * pm.XZ
	pcy2 := y + cx2 * pm.YX + cy2 * pm.YY + cz2 * pm.YZ
	cv.LineTo(pcx1 - w1, pcy1)
	cv.LineTo(pcx1 + w1, pcy1)
	cv.LineTo(pcx2 + w2, pcy2)
	cv.LineTo(pcx2 - w2, pcy2)
	cv.ClosePath()
	cv.Fill()
}

func DrawPerson(cv *canvas.Canvas, t *navigation.Traveller, x float64, y float64, c *controller.Controller) {
	m := animation.PersonMotionWalk
	p := (t.Phase / 4) % 8
	var pm = animation.ProjectionMatrixNE
	switch c.Perspective {
	case controller.PerspectiveNE:
		pm = animation.ProjectionMatrixNE
	}

	// Arm
	cv.SetFillStyle("#974")
	// LeftElbow
	DrawLimb(cv, pm, x, y, 0, -24, -4, 2, m.LeftElbow[p][0], -18+m.LeftElbow[p][1], -4, 2)
	// LeftHand
	DrawLimb(cv, pm, x, y, m.LeftElbow[p][0], -18+m.LeftElbow[p][1], -4, 2, m.LeftKnee[p][0], -12+m.LeftKnee[p][1], -4, 2)

	// Body
	cv.SetFillStyle("#A84")
	cv.FillRect(x-2, y-28, 4, 3)
	cv.FillRect(x-4, y-25, 8, 10)
	cv.SetFillStyle("#840")
	// Head
	cv.BeginPath()
	cv.Arc(x, y-30, 3, 0, math.Pi*2, false)
	cv.ClosePath()
	cv.Fill()
	// Legs
	cv.SetFillStyle("#420")
	// LeftKnee
	DrawLimb(cv, pm, x, y, 0, -15, -3, 3, m.LeftKnee[p][0], -8+m.LeftKnee[p][1], -3, 2)
	// LeftFoot
	DrawLimb(cv, pm, x, y, m.LeftKnee[p][0], -8+m.LeftKnee[p][1], -3, 2, m.LeftFoot[p][0], m.LeftFoot[p][1], -3, 2)

	// RightKnee
	DrawLimb(cv, pm, x, y, 0, -15, 3, 3, m.RightKnee[p][0], -8+m.RightKnee[p][1], 3, 2)
	// RightFoot
	DrawLimb(cv, pm, x, y, m.RightKnee[p][0], -8+m.RightKnee[p][1], 3, 2, m.RightFoot[p][0], m.RightFoot[p][1], 3, 2)

	// Arm
	cv.SetFillStyle("#974")
	// RightElbow
	DrawLimb(cv, pm, x, y, 0, -24, 4, 2, m.RightElbow[p][0], -18+m.RightElbow[p][1], 4, 2)
	// RightHand
	DrawLimb(cv, pm, x, y, m.RightElbow[p][0], -18+m.RightElbow[p][1], 4, 2, m.RightKnee[p][0], -12+m.RightKnee[p][1], 4, 2)
}
