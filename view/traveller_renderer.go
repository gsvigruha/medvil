package view

import (
	//"fmt"
	"github.com/tfriedel6/canvas"
	"math"
	"medvil/controller"
	"medvil/model/building"
	"medvil/model/navigation"
	"medvil/renderer"
	"medvil/view/animation"
)

const MaxPX = navigation.MaxPX
const MaxPY = navigation.MaxPY

func getZByDir(bpe *navigation.BuildingPathElement, dir uint8) float64 {
	if bpe.BC.Connection(dir) == building.ConnectionTypeUpperLevel {
		return float64(bpe.GetLocation().Z) * DZ * BuildingUnitHeight
	} else if bpe.BC.Connection(dir) == building.ConnectionTypeLowerLevel {
		return float64(bpe.GetLocation().Z-1) * DZ * BuildingUnitHeight
	}
	return 0
}

func RenderTravellers(cv *canvas.Canvas, travellers []*navigation.Traveller, rf renderer.RenderedField, c *controller.Controller) {
	for i := range travellers {
		t := travellers[i]
		px := float64(t.PX)
		py := float64(t.PY)
		NEPX := rf.X[(2+c.Perspective)%4]
		NEPY := rf.Y[(2+c.Perspective)%4] - rf.Z[(2+c.Perspective)%4]
		SEPX := rf.X[(1+c.Perspective)%4]
		SEPY := rf.Y[(1+c.Perspective)%4] - rf.Z[(1+c.Perspective)%4]
		SWPX := rf.X[(0+c.Perspective)%4]
		SWPY := rf.Y[(0+c.Perspective)%4] - rf.Z[(0+c.Perspective)%4]
		NWPX := rf.X[(3+c.Perspective)%4]
		NWPY := rf.Y[(3+c.Perspective)%4] - rf.Z[(3+c.Perspective)%4]
		x := (NWPX*(MaxPX-px)*(MaxPY-py) +
			SWPX*(MaxPX-px)*py +
			NEPX*px*(MaxPY-py) +
			SEPX*px*py) / (MaxPX * MaxPY)
		y := (NWPY*(MaxPX-px)*(MaxPY-py) +
			SWPY*(MaxPX-px)*py +
			NEPY*px*(MaxPY-py) +
			SEPY*px*py) / (MaxPX * MaxPY)
		if t.PE != nil && t.PE.GetLocation().Z > 0 {
			if bpe, ok := t.PE.(*navigation.BuildingPathElement); ok {
				z1 := getZByDir(bpe, t.Direction)
				z2 := getZByDir(bpe, building.OppDir(t.Direction))
				var z = 0.0
				switch t.Direction {
				case navigation.DirectionN:
					z = (z1*(MaxPY-py) + z2*py) / MaxPY
				case navigation.DirectionS:
					z = (z1*py + z2*(MaxPY-py)) / MaxPY
				case navigation.DirectionW:
					z = (z1*(MaxPX-px) + z2*px) / MaxPX
				case navigation.DirectionE:
					z = (z1*px + z2*(MaxPX-px)) / MaxPX
				}
				DrawTraveller(cv, t, x, y-5-z, c)
			} else {
				DrawTraveller(cv, t, x, y-5, c)
			}
		} else {
			DrawTraveller(cv, t, x, y-5, c)
		}
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

func DrawTool(cv *canvas.Canvas, pm animation.ProjectionMatrix, m animation.PersonMotion, x, y float64, p uint8) {
	// Tool
	cv.SetStrokeStyle("#321")
	cv.SetLineWidth(3)

	lh := m.LeftHand[p]
	lhx := x + lh[0]*pm.XX + lh[1]*pm.XY + lh[2]*pm.XZ
	lhy := y + lh[0]*pm.YX + lh[1]*pm.YY + lh[2]*pm.YZ

	rh := m.RightHand[p]
	rhx := x + rh[0]*pm.XX + rh[1]*pm.XY + rh[2]*pm.XZ
	rhy := y + rh[0]*pm.YX + rh[1]*pm.YY + rh[2]*pm.YZ

	tx := rhx + (lhx-rhx)*2
	ty := rhy + (lhy-rhy)*2

	cv.BeginPath()
	cv.MoveTo(rhx, rhy)
	cv.LineTo(tx, ty)
	cv.ClosePath()
	cv.Stroke()

	cv.SetFillStyle("#999")
	cv.BeginPath()
	cv.LineTo(tx, ty-2)
	cv.LineTo(tx-2, ty)
	cv.LineTo(tx, ty+4)
	cv.LineTo(tx+2, ty)
	cv.ClosePath()
	cv.Fill()
}

func DrawTraveller(cv *canvas.Canvas, t *navigation.Traveller, x float64, y float64, c *controller.Controller) {
	if t.T == navigation.TravellerTypePedestrian {
		inBoat := t.Vehicle != nil && t.Vehicle.TravellerType() == navigation.TravellerTypeBoat
		DrawPerson(cv, t, x, y, inBoat, c)
	} else if t.T == navigation.TravellerTypeBoat {
		DrawBoat(cv, t, x, y, c)
	} else if t.T == navigation.TravellerTypeCart {
		DrawCart(cv, t, x, y, c)
	}
	c.AddRenderedTraveller(&renderer.RenderedTraveller{X: x, Y: y, H: 30, W: 10, Traveller: t})
}

func DrawCart(cv *canvas.Canvas, t *navigation.Traveller, x float64, y float64, c *controller.Controller) {
	if !t.Visible {
		return
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
		dx0 := math.Cos(math.Pi*2.0*i/8.0)*(h1/2.0) + (f2-f1)/2.0 + f1
		dy0 := math.Sin(math.Pi*2.0*i/8.0)*(h1/2.0) + h1/2.0
		cv.LineTo(x+dx0*pm.XX-dy0*pm.XY+z*pm.XZ*dir, y+dx0*pm.YX-dy0*pm.YY+z*pm.YZ*dir)
	}
	cv.ClosePath()
	cv.Fill()
	cv.Stroke()
}

func DrawBoat(cv *canvas.Canvas, t *navigation.Traveller, x float64, y float64, c *controller.Controller) {
	if !t.Visible {
		return
	}
	dirIdx := (c.Perspective - t.Direction) % 4
	pm := animation.ProjectionMatrices[dirIdx]
	var r = 8.0 - float64(t.Phase%16)
	if r < 0.0 {
		r = -r
	}

	f := 18.0
	s := 8.0
	z := 6.0
	h := 6.0
	p := 16.0

	cv.SetFillStyle("texture/vehicle/boat_bottom.png")
	cv.BeginPath()
	cv.LineTo(x-f*pm.XX+0*pm.XY+0*pm.XZ, y-f*pm.YX+0*pm.YY+0*pm.YZ)
	cv.LineTo(x+0*pm.XX+0*pm.XY-z*pm.XZ, y+0*pm.YX+0*pm.YY-z*pm.YZ)
	cv.LineTo(x+f*pm.XX+0*pm.XY+0*pm.XZ, y+f*pm.YX+0*pm.YY+0*pm.YZ)
	cv.LineTo(x+0*pm.XX+0*pm.XY+z*pm.XZ, y+0*pm.YX+0*pm.YY+z*pm.YZ)
	cv.ClosePath()
	cv.Fill()

	cv.SetFillStyle("texture/vehicle/boat_side.png")
	cv.BeginPath()
	cv.LineTo(x-f*pm.XX+0*pm.XY+0*pm.XZ, y-f*pm.YX+0*pm.YY+0*pm.YZ)
	cv.LineTo(x+0*pm.XX+0*pm.XY-z*pm.XZ, y+0*pm.YX+0*pm.YY-z*pm.YZ)
	cv.LineTo(x+f*pm.XX+0*pm.XY+0*pm.XZ, y+f*pm.YX+0*pm.YY+0*pm.YZ)
	cv.LineTo(x+f*pm.XX-h*pm.XY+0*pm.XZ, y+f*pm.YX-h*pm.YY+0*pm.YZ)
	cv.LineTo(x+0*pm.XX-h*pm.XY-s*pm.XZ, y+0*pm.YX-h*pm.YY-s*pm.YZ)
	cv.LineTo(x-f*pm.XX-h*pm.XY+0*pm.XZ, y-f*pm.YX-h*pm.YY+0*pm.YZ)
	cv.ClosePath()
	cv.Fill()

	cv.SetFillStyle("texture/vehicle/boat_side.png")
	cv.BeginPath()
	cv.LineTo(x-f*pm.XX+0*pm.XY+0*pm.XZ, y-f*pm.YX+0*pm.YY+0*pm.YZ)
	cv.LineTo(x+0*pm.XX+0*pm.XY+z*pm.XZ, y+0*pm.YX+0*pm.YY+z*pm.YZ)
	cv.LineTo(x+f*pm.XX+0*pm.XY+0*pm.XZ, y+f*pm.YX+0*pm.YY+0*pm.YZ)
	cv.LineTo(x+f*pm.XX-h*pm.XY+0*pm.XZ, y+f*pm.YX-h*pm.YY+0*pm.YZ)
	cv.LineTo(x+0*pm.XX-h*pm.XY+s*pm.XZ, y+0*pm.YX-h*pm.YY+s*pm.YZ)
	cv.LineTo(x-f*pm.XX-h*pm.XY+0*pm.XZ, y-f*pm.YX-h*pm.YY+0*pm.YZ)
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

}

func DrawPerson(cv *canvas.Canvas, t *navigation.Traveller, x float64, y float64, InVehicle bool, c *controller.Controller) {
	if !t.Visible {
		return
	}
	if InVehicle {
		y += 5
	}
	var m animation.PersonMotion
	switch t.Motion {
	case navigation.MotionWalk:
		m = animation.PersonMotionWalk
	case navigation.MotionFieldWork:
		m = animation.PersonMotionFieldWork
	case navigation.MotionBuild:
		m = animation.PersonMotionBuild
	case navigation.MotionMine:
		m = animation.PersonMotionMine
	}
	p := (t.Phase / 2) % 8
	dirIdx := (c.Perspective - t.Direction) % 4
	pm := animation.ProjectionMatrices[dirIdx]

	if dirIdx >= 2 {
		DrawLeftArm(cv, pm, m, x, y, p)
		if !InVehicle {
			DrawLeftLeg(cv, pm, m, x, y, p)
		}
	} else {
		DrawRightArm(cv, pm, m, x, y, p)
		if !InVehicle {
			DrawRightLeg(cv, pm, m, x, y, p)
		}
	}
	if dirIdx == 1 || dirIdx == 2 {
		if m.Tool {
			DrawTool(cv, pm, m, x, y, p)
		}
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
		if !InVehicle {
			DrawRightLeg(cv, pm, m, x, y, p)
		}
		DrawRightArm(cv, pm, m, x, y, p)
	} else {
		if !InVehicle {
			DrawLeftLeg(cv, pm, m, x, y, p)
		}
		DrawLeftArm(cv, pm, m, x, y, p)
	}
	if dirIdx == 0 || dirIdx == 3 {
		if m.Tool {
			DrawTool(cv, pm, m, x, y, p)
		}
	}
}
