package view

import (
	//"fmt"
	"github.com/tfriedel6/canvas"
	"math"
	"medvil/controller"
	"medvil/model/building"
	"medvil/model/navigation"
	"medvil/model/social"
	"medvil/renderer"
	"medvil/view/animation"
	"medvil/view/vehicles"
	"path/filepath"
)

const MaxPX = navigation.MaxPX
const MaxPY = navigation.MaxPY

const ClothesRed = 1
const ClothesBlue = 2
const ClothesPurple = 3
const ClothesYellow = 4
const ClothesMetal = 5
const ClothesBrown = 6

func Move(v1, v2 [3]float64) [3]float64 {
	var v3 [3]float64
	for i, _ := range v1 {
		v3[i] = v1[i] + v2[i]
	}
	return v3
}

func GetScreenY(t *navigation.Traveller, rf renderer.RenderedField, c *controller.Controller) float64 {
	_, y := GetScreenXY(t, rf, c)
	return y
}

func GetScreenXY(t *navigation.Traveller, rf renderer.RenderedField, c *controller.Controller) (float64, float64) {
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
	return x, y
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

func DrawLeftArm(cv *canvas.Canvas, pm animation.ProjectionMatrix, m animation.PersonMotion, x, y float64, p uint8, color int) {
	// Arm
	setClothesColor(cv, color, true)
	// LeftElbow
	DrawLimb(cv, pm, x, y, 1, 2, m.LeftShoulder, m.LeftElbow[p])
	// LeftHand
	DrawLimb(cv, pm, x, y, 2, 1, m.LeftElbow[p], m.LeftHand[p])
}

func DrawLeftLeg(cv *canvas.Canvas, pm animation.ProjectionMatrix, m animation.PersonMotion, x, y float64, p uint8) {
	// Legs
	cv.SetFillStyle(filepath.FromSlash("texture/people/leather.png"))
	// LeftKnee
	DrawLimb(cv, pm, x, y, 3, 2, m.LeftHip, m.LeftKnee[p])
	// LeftShin
	DrawLimb(cv, pm, x, y, 2, 2, m.LeftKnee[p], m.LeftFoot[p])
	// Left Foot
	DrawLimb(cv, pm, x, y, 1, 1, Move(m.LeftFoot[p], [3]float64{-1.0, 0.0, 0.0}), Move(m.LeftFoot[p], [3]float64{4.0, 0.0, 0.0}))
}

func DrawRightLeg(cv *canvas.Canvas, pm animation.ProjectionMatrix, m animation.PersonMotion, x, y float64, p uint8) {
	// Legs
	cv.SetFillStyle(filepath.FromSlash("texture/people/leather.png"))
	// RightKnee
	DrawLimb(cv, pm, x, y, 3, 2, m.RightHip, m.RightKnee[p])
	// LeftShin
	DrawLimb(cv, pm, x, y, 2, 2, m.RightKnee[p], m.RightFoot[p])
	// Right Foot
	DrawLimb(cv, pm, x, y, 1, 1, Move(m.RightFoot[p], [3]float64{-1.0, 0.0, 0.0}), Move(m.RightFoot[p], [3]float64{4.0, 0.0, 0.0}))
}

func DrawRightArm(cv *canvas.Canvas, pm animation.ProjectionMatrix, m animation.PersonMotion, x, y float64, p uint8, color int) {
	// Arm
	setClothesColor(cv, color, true)
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

func tallPlant(f *navigation.Field) bool {
	return f.Plant != nil && f.Plant.T.Tall
}

func DrawTraveller(cv *canvas.Canvas, t *navigation.Traveller, x float64, y float64, f *navigation.Field, c *controller.Controller) {
	if cv != nil {
		if t.T == navigation.TravellerTypePedestrianM || t.T == navigation.TravellerTypePedestrianF {
			inBoat := t.Vehicle != nil && t.Vehicle.Water()
			if inBoat {
				y += 5
			}
			DrawPerson(cv, t, x, y, !inBoat && !tallPlant(f), c)
		} else if t.T == navigation.TravellerTypeBoat {
			vehicles.DrawBoat(cv, t, x, y, c)
		} else if t.T == navigation.TravellerTypeTradingBoat {
			vehicles.DrawTradingBoat(cv, t, x, y, c)
		} else if t.T == navigation.TravellerTypeExpeditionBoat {
			vehicles.DrawExpeditionBoat(cv, t, x, y, c)
		} else if t.T == navigation.TravellerTypeCart {
			vehicles.DrawCart(cv, t, x, y, c)
		} else if t.T == navigation.TravellerTypeTradingCart {
			vehicles.DrawTradingCart(cv, t, x, y, c)
		} else if t.T == navigation.TravellerTypeExpeditionCart {
			vehicles.DrawExpeditionCart(cv, t, x, y, c)
		}
	}
}

func drawHair(cv *canvas.Canvas, t *navigation.Traveller, x float64, y float64) {
	cv.SetFillStyle("#630")
	cv.BeginPath()
	if t.T == navigation.TravellerTypePedestrianM {
		cv.Ellipse(x, y-32, 3, 4, 0, 0, math.Pi*2, false)
	} else if t.T == navigation.TravellerTypePedestrianF {
		cv.Ellipse(x, y-29, 4, 6, 0, 0, math.Pi*2, false)
	}
	cv.ClosePath()
	cv.Fill()
}

func drawHead(cv *canvas.Canvas, t *navigation.Traveller, x float64, y float64) {
	cv.SetFillStyle("#A64")
	cv.BeginPath()
	if t.T == navigation.TravellerTypePedestrianM {
		cv.Arc(x, y-30, 3, 0, math.Pi*2, false)
	} else if t.T == navigation.TravellerTypePedestrianF {
		cv.Arc(x, y-28, 3, 0, math.Pi*2, false)
	}
	cv.ClosePath()
	cv.Fill()
}

func setClothesColor(cv *canvas.Canvas, color int, dark bool) {
	var sfx = ""
	if dark {
		sfx = "_dark"
	}
	switch color {
	case ClothesMetal:
		cv.SetFillStyle(filepath.FromSlash("texture/people/metal.png"))
	case ClothesBrown:
		cv.SetFillStyle("#952")
	case ClothesYellow:
		cv.SetFillStyle(filepath.FromSlash("texture/people/textile_yellow" + sfx + ".png"))
	case ClothesRed:
		cv.SetFillStyle(filepath.FromSlash("texture/people/textile_red" + sfx + ".png"))
	case ClothesPurple:
		cv.SetFillStyle(filepath.FromSlash("texture/people/textile_purple" + sfx + ".png"))
	case ClothesBlue:
		cv.SetFillStyle(filepath.FromSlash("texture/people/textile_blue" + sfx + ".png"))
	}
}

func DrawPerson(cv *canvas.Canvas, t *navigation.Traveller, x float64, y float64, drawLeg bool, c *controller.Controller) {
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
	case navigation.MotionCut:
		m = animation.PersonMotionCut
	}
	p := t.DrawingPhase()
	dirIdx := (c.Perspective - t.Direction) % 4
	pm := animation.ProjectionMatrices[dirIdx]

	var color int
	person := c.ReverseReferences.TravellerToPerson[t]
	if person != nil {
		if person.Home.GetTown().Country.T == social.CountryTypeOutlaw {
			color = ClothesBrown
		} else if person.Equipment.Weapon {
			color = ClothesMetal
		} else if person.Home.GetBuilding() != nil {
			switch person.Home.GetBuilding().Plan.BuildingType {
			case building.BuildingTypeFarm:
				color = ClothesYellow
			case building.BuildingTypeMine:
				color = ClothesRed
			case building.BuildingTypeWorkshop:
				color = ClothesPurple
			case building.BuildingTypeFactory:
				color = ClothesPurple
			case building.BuildingTypeTownhall:
				color = ClothesBlue
			}
		}
	}

	if dirIdx >= 2 {
		DrawLeftArm(cv, pm, m, x, y, p, color)
		if drawLeg {
			DrawLeftLeg(cv, pm, m, x, y, p)
		}
	} else {
		DrawRightArm(cv, pm, m, x, y, p, color)
		if drawLeg {
			DrawRightLeg(cv, pm, m, x, y, p)
		}
	}
	if dirIdx == 1 || dirIdx == 2 {
		if m.Tool {
			DrawTool(cv, pm, m, x, y, p)
		}
	}

	if drawLeg {
		if dirIdx >= 2 {
			DrawRightLeg(cv, pm, m, x, y, p)
		} else {
			DrawLeftLeg(cv, pm, m, x, y, p)
		}
	}

	if dirIdx == 0 || dirIdx == 3 {
		drawHair(cv, t, x, y)
	} else {
		drawHead(cv, t, x, y)
	}

	// Body
	setClothesColor(cv, color, false)
	if t.T == navigation.TravellerTypePedestrianM {
		cv.FillRect(x-2, y-28, 4, 3)
		cv.FillRect(x-4, y-25, 8, 11)
	} else if t.T == navigation.TravellerTypePedestrianF {
		cv.BeginPath()
		cv.LineTo(x-2, y-28)
		cv.LineTo(x-4, y-25)
		cv.LineTo(x-5, y-9)
		cv.LineTo(x+5, y-9)
		cv.LineTo(x+4, y-25)
		cv.LineTo(x+2, y-28)
		cv.ClosePath()
		cv.Fill()
	}

	if dirIdx == 0 || dirIdx == 3 {
		drawHead(cv, t, x, y)
	} else {
		drawHair(cv, t, x, y)
	}

	if dirIdx >= 2 {
		DrawRightArm(cv, pm, m, x, y, p, color)
	} else {
		DrawLeftArm(cv, pm, m, x, y, p, color)
	}
	if dirIdx == 0 || dirIdx == 3 {
		if m.Tool {
			DrawTool(cv, pm, m, x, y, p)
		}
	}
}
