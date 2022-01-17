package view

import (
	//"fmt"
	"github.com/tfriedel6/canvas"
	"math"
	"medvil/controller"
	"medvil/model/navigation"
	"medvil/renderer"
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
		DrawPerson(cv, t, x, y-5, c)
	}
}

func DrawLimb(cv *canvas.Canvas, pm animation.ProjectionMatrix, x, y, w1, w2 float64, c1, c2 [3]float64) {
	cv.BeginPath()
	pcx1 := x + c1[0]*pm.XX + c1[1]*pm.XY + c1[2]*pm.XZ
	pcy1 := y + c1[0]*pm.YX + c1[1]*pm.YY + c1[2]*pm.YZ
	pcx2 := x + c2[0]*pm.XX + c2[1]*pm.XY + c2[2]*pm.XZ
	pcy2 := y + c2[0]*pm.YX + c2[1]*pm.YY + c2[2]*pm.YZ
	a := math.Tanh((pcy2-pcy1)/(pcx2-pcx1)) + math.Pi/2
	dx1 := w1 * math.Cos(a)
	dy1 := w1 * math.Sin(a)
	dx2 := w2 * math.Cos(a)
	dy2 := w2 * math.Sin(a)
	cv.LineTo(pcx1-dx1, pcy1-dy1)
	cv.LineTo(pcx1+dx1, pcy1+dy1)
	cv.LineTo(pcx2+dx2, pcy2+dy2)
	cv.LineTo(pcx2-dx2, pcy2-dy2)
	cv.ClosePath()
	cv.Fill()
}

func DrawLeftArm(cv *canvas.Canvas, pm animation.ProjectionMatrix, m animation.PersonMotion, x, y float64, p uint8) {
	// Arm
	cv.SetFillStyle("#762")
	// LeftElbow
	DrawLimb(cv, pm, x, y, 1, 2, m.LeftShoulder, m.LeftElbow[p])
	// LeftHand
	DrawLimb(cv, pm, x, y, 2, 1, m.LeftElbow[p], m.LeftHand[p])
}

func DrawLeftLeg(cv *canvas.Canvas, pm animation.ProjectionMatrix, m animation.PersonMotion, x, y float64, p uint8) {
	// Legs
	cv.SetFillStyle("#420")
	// LeftKnee
	DrawLimb(cv, pm, x, y, 3, 2, m.LeftHip, m.LeftKnee[p])
	// LeftFoot
	DrawLimb(cv, pm, x, y, 2, 2, m.LeftKnee[p], m.LeftFoot[p])
}

func DrawRightLeg(cv *canvas.Canvas, pm animation.ProjectionMatrix, m animation.PersonMotion, x, y float64, p uint8) {
	// Legs
	cv.SetFillStyle("#420")
	// RightKnee
	DrawLimb(cv, pm, x, y, 3, 2, m.RightHip, m.RightKnee[p])
	// LeftFoot
	DrawLimb(cv, pm, x, y, 2, 2, m.RightKnee[p], m.RightFoot[p])
}

func DrawRightArm(cv *canvas.Canvas, pm animation.ProjectionMatrix, m animation.PersonMotion, x, y float64, p uint8) {
	// Arm
	cv.SetFillStyle("#762")
	// RightElbow
	DrawLimb(cv, pm, x, y, 1, 2, m.RightShoulder, m.RightElbow[p])
	// LeftHand
	DrawLimb(cv, pm, x, y, 2, 1, m.RightElbow[p], m.RightHand[p])
}

func DrawPerson(cv *canvas.Canvas, t *navigation.Traveller, x float64, y float64, c *controller.Controller) {
	m := animation.PersonMotionWalk
	p := (t.Phase / 4) % 8
	dirIdx := (c.Perspective - t.Direction) % 4
	pm := animation.ProjectionMatrices[dirIdx]

	if dirIdx >= 2 {
		DrawLeftArm(cv, pm, m, x, y, p)
		DrawLeftLeg(cv, pm, m, x, y, p)
	} else {
		DrawRightArm(cv, pm, m, x, y, p)
		DrawRightLeg(cv, pm, m, x, y, p)
	}

	// Body
	cv.SetFillStyle("#BA6")
	cv.FillRect(x-2, y-28, 4, 3)
	cv.FillRect(x-4, y-25, 8, 10)
	cv.SetFillStyle("#840")
	// Head
	cv.BeginPath()
	cv.Arc(x, y-30, 3, 0, math.Pi*2, false)
	cv.ClosePath()
	cv.Fill()

	if dirIdx >= 2 {
		DrawRightLeg(cv, pm, m, x, y, p)
		DrawRightArm(cv, pm, m, x, y, p)
	} else {
		DrawLeftLeg(cv, pm, m, x, y, p)
		DrawLeftArm(cv, pm, m, x, y, p)
	}
}
