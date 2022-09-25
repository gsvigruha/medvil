package view

import (
	"github.com/tfriedel6/canvas"
	"image/color"
	"math"
	"medvil/controller"
	"medvil/model/building"
	"medvil/model/materials"
	"medvil/renderer"
	"medvil/view/buildings"
	"strconv"
)

const BuildingAnimationMaxPhase = 24

const BuildingUnitHeight = buildings.BuildingUnitHeight

func WallMaterialName(t building.BuildingType, m *materials.Material, shape uint8) string {
	if m == materials.GetMaterial("brick") {
		if shape == 0 {
			return "painted_yellow"
		} else if shape == 1 {
			return "painted_red"
		} else if shape == 2 {
			return "painted_brown"
		} else if shape == 3 {
			return "painted_beige"
		} else if shape == 4 {
			return "painted_sand"
		}
	}
	if t == building.BuildingTypeWall && m == materials.GetMaterial("stone") {
		if shape == 0 {
			return "stone_1"
		} else if shape == 1 {
			return "stone_2"
		} else {
			return "stone"
		}
	}
	return m.Name
}

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

	cv.SetFillStyle("texture/building/gray_marble.png")
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

func RenderBuildingUnit(cv *canvas.Canvas, unit *building.BuildingUnit, rf renderer.RenderedField, k int, c *controller.Controller) renderer.RenderedBuildingUnit {
	var rws = []renderer.RenderedWall{}
	startI := 2 + c.Perspective
	for i := uint8(startI); i < 4+startI; i++ {
		wall := unit.Walls[i%4]
		if wall == nil {
			continue
		}
		rfIdx1 := (3 - (-c.Perspective + i)) % 4
		rfIdx2 := (2 - (-c.Perspective + i)) % 4
		if !unit.Construction && (rfIdx1 == 0 || rfIdx1 == 1) && unit.B.Plan.BuildingType != building.BuildingTypeGate {
			continue
		}
		var suffix = ""
		if rfIdx1%2 == 1 {
			suffix = "_flipped"
		}
		if cv != nil {
			if !unit.Construction {
				cv.SetFillStyle("texture/building/" + WallMaterialName(unit.B.Plan.BuildingType, wall.M, unit.B.Shape) + suffix + ".png")
			} else {
				cv.SetFillStyle("texture/building/construction" + suffix + ".png")
			}
		}

		z := float64(k*BuildingUnitHeight) * DZ
		rw := renderer.RenderedWall{
			X: [4]float64{rf.X[rfIdx1], rf.X[rfIdx1], rf.X[rfIdx2], rf.X[rfIdx2]},
			Y: [4]float64{
				rf.Y[rfIdx1] - rf.Z[rfIdx1] - z, rf.Y[rfIdx1] - rf.Z[rfIdx1] - z - BuildingUnitHeight*DZ,
				rf.Y[rfIdx2] - rf.Z[rfIdx2] - z - BuildingUnitHeight*DZ, rf.Y[rfIdx2] - rf.Z[rfIdx2] - z},
			Wall: wall,
		}
		rws = append(rws, rw)
		if cv != nil {
			cv.SetStrokeStyle(color.RGBA{R: 192, G: 128, B: 64, A: 24})
			cv.BeginPath()
			if wall.Arch {
				cv.LineTo(rw.X[0], rw.Y[0])
				dx := (rw.X[3]*0.8 - rw.X[0]*0.8) / 12
				dy := (rw.Y[3]*0.8 - rw.Y[0]*0.8) / 12
				for n := 0.0; n <= 12; n++ {
					zn := math.Pow((6.0-(math.Abs(n-6.0)))/6.0, 0.2 /*arch exponent*/) * 0.8
					cv.LineTo(rw.X[0]*0.9+rw.X[3]*0.1+n*dx, rw.Y[0]*0.9+rw.Y[3]*0.1+n*dy-zn*BuildingUnitHeight*DZ)
				}
				cv.LineTo(rw.X[3], rw.Y[3])
				cv.LineTo(rw.X[2], rw.Y[2])
				cv.LineTo(rw.X[1], rw.Y[1])
			} else {
				cv.LineTo(rw.X[0], rw.Y[0])
				cv.LineTo(rw.X[1], rw.Y[1])
				cv.LineTo(rw.X[2], rw.Y[2])
				cv.LineTo(rw.X[3], rw.Y[3])
			}
			cv.ClosePath()
			cv.Fill()
			cv.Stroke()

			/*
				if !wall.Arch {
					cv.SetFillStyle("texture/building/ornament" + suffix + ".png")
					cv.BeginPath()
					cv.LineTo(rw.X[0], rw.Y[0]*0.2+rw.Y[1]*0.8)
					cv.LineTo(rw.X[1], rw.Y[1])
					cv.LineTo(rw.X[2], rw.Y[2])
					cv.LineTo(rw.X[3], rw.Y[3]*0.2+rw.Y[2]*0.8)
					cv.ClosePath()
					cv.Fill()
				}
			*/

			z := math.Min(math.Min(math.Min(rf.Z[0], rf.Z[1]), rf.Z[2]), rf.Z[3]) + float64(k*BuildingUnitHeight)*DZ
			if !unit.Construction && wall.Windows != building.WindowTypeNone {
				cv.SetFillStyle("texture/building/glass_2.png")
				cv.SetStrokeStyle(color.RGBA{R: 32, G: 32, B: 0, A: 64})
				cv.SetLineWidth(2)
				RenderWindows(cv, rf, rfIdx1, rfIdx2, z, wall.Door, wall.Windows == building.WindowTypeFrench)
				if wall.Windows == building.WindowTypeBalcony {
					RenderBalcony(cv, rf, rfIdx1, rfIdx2, z, wall.Door)
				}
			}

			if wall.Door && !unit.Construction {
				cv.SetFillStyle("texture/building/door.png")
				cv.SetStrokeStyle(color.RGBA{R: 64, G: 32, B: 0, A: 64})
				cv.SetLineWidth(2)

				cv.BeginPath()
				cv.LineTo((3*rf.X[rfIdx1]+7*rf.X[rfIdx2])/10, (3*rf.Y[rfIdx1]+7*rf.Y[rfIdx2])/10-z)
				cv.LineTo((3*rf.X[rfIdx1]+7*rf.X[rfIdx2])/10, (3*rf.Y[rfIdx1]+7*rf.Y[rfIdx2])/10-z-BuildingUnitHeight*DZ*3/5)
				cv.LineTo((1*rf.X[rfIdx1]+9*rf.X[rfIdx2])/10, (1*rf.Y[rfIdx1]+9*rf.Y[rfIdx2])/10-z-BuildingUnitHeight*DZ*3/5)
				cv.LineTo((1*rf.X[rfIdx1]+9*rf.X[rfIdx2])/10, (1*rf.Y[rfIdx1]+9*rf.Y[rfIdx2])/10-z)
				cv.ClosePath()
				cv.Fill()
				cv.Stroke()

				workshop := c.ReverseReferences.BuildingToWorkshop[unit.Building()]
				if unit.NamePlate() && workshop != nil && workshop.Manufacture != nil {
					dX := float64((int(rfIdx1)%2)*2 - 1)
					cv.SetStrokeStyle(color.RGBA{R: 48, G: 32, B: 0, A: 192})
					cv.SetFillStyle("#514")
					cv.SetLineWidth(2)
					cv.BeginPath()
					xm, ym := (3*rf.X[rfIdx1]+7*rf.X[rfIdx2])/10, (3*rf.Y[rfIdx1]+7*rf.Y[rfIdx2])/10-BuildingUnitHeight*DZ*4/5
					cv.LineTo(xm, ym)
					cv.LineTo(xm+dX*18, ym+12)
					cv.LineTo(xm+dX*18, ym+28)
					cv.LineTo(xm+dX*2, ym+17)
					cv.LineTo(xm+dX*2, ym+1)
					cv.ClosePath()
					cv.Stroke()
					cv.Fill()
					cv.DrawImage("icon/gui/tasks/"+workshop.Manufacture.Name+".png", xm+dX*2, ym+5, 16*dX, 16)
				}
			}
		}
	}
	return renderer.RenderedBuildingUnit{Walls: rws, Unit: unit}
}

func RenderBuildingExtension(cv *canvas.Canvas, extension *building.ExtensionUnit, rf renderer.RenderedField, k int, phase uint8, c *controller.Controller) {
	if extension == nil {
		return
	}
	if extension.T == building.WaterMillWheel {
		if extension.IsConstruction() {
			return
		}
		var dir = 1.0
		var phi float64
		phi = float64(phase) * math.Pi * 2.0 / 12.0 / float64(BuildingAnimationMaxPhase)

		rfIdx1 := (3 - (-c.Perspective + extension.Direction)) % 4
		rfIdx2 := (2 - (-c.Perspective + extension.Direction)) % 4
		rfIdx3 := (1 - (-c.Perspective + extension.Direction)) % 4
		rfIdx4 := (0 - (-c.Perspective + extension.Direction)) % 4

		xo := (rf.X[rfIdx1] + rf.X[rfIdx2]) / 2.0
		yo := (rf.Y[rfIdx1] + rf.Y[rfIdx2]) / 2.0
		xi := (rf.X[rfIdx1]*4 + rf.X[rfIdx2]*4 + rf.X[rfIdx3] + rf.X[rfIdx4]) / 10
		yi := (rf.Y[rfIdx1]*4 + rf.Y[rfIdx2]*4 + rf.Y[rfIdx3] + rf.Y[rfIdx4]) / 10

		if rfIdx1%2 == 1 {
			dir = -1.0
		}
		var (
			x1, y1, x2, y2 float64
		)
		if rfIdx1 < 2 {
			x1, y1, x2, y2 = xo, yo, xi, yi
		} else {
			x1, y1, x2, y2 = xi, yi, xo, yo
		}

		cv.SetFillStyle("texture/building/waterwheel_wood.png")
		cv.SetStrokeStyle("#310")
		cv.SetLineWidth(2)

		r1 := 1.5
		r2 := 2.0
		r3 := 2.5
		for i := 0.0; i < 12; i++ {
			dx0 := math.Cos(math.Pi*2.0*i/12.0 + phi)
			dy0 := math.Sin(math.Pi*2.0*i/12.0 + phi)
			dx1 := math.Cos(math.Pi*2.0*(i+1)/12.0 + phi)
			dy1 := math.Sin(math.Pi*2.0*(i+1)/12.0 + phi)

			cv.BeginPath()
			cv.LineTo(x1+dx0*DZ*r2*2/3, y1-DZ*2+dy0*DZ*r2+dx0*DZ*r2*1/3*dir)
			cv.LineTo(x1+dx0*DZ*r1*2/3, y1-DZ*2+dy0*DZ*r1+dx0*DZ*r1*1/3*dir)
			cv.LineTo(x1+dx1*DZ*r1*2/3, y1-DZ*2+dy1*DZ*r1+dx1*DZ*r1*1/3*dir)
			cv.LineTo(x1+dx1*DZ*r2*2/3, y1-DZ*2+dy1*DZ*r2+dx1*DZ*r2*1/3*dir)
			cv.ClosePath()
			cv.Fill()

			cv.BeginPath()
			cv.MoveTo(x1+dx0*DZ*r1*2/3, y1-DZ*2+dy0*DZ*r1+dx0*DZ*r1*1/3*dir)
			cv.LineTo(x1-dx0*DZ*r1*2/3, y1-DZ*2-dy0*DZ*r1-dx0*DZ*r1*1/3*dir)
			cv.ClosePath()
			cv.Stroke()
		}

		cv.SetFillStyle("texture/building/waterwheel_wood_2.png")
		for i := 0.0; i < 12; i++ {
			dx0 := math.Cos(math.Pi*2.0*i/12.0 + phi)
			dy0 := math.Sin(math.Pi*2.0*i/12.0 + phi)
			dx1 := math.Cos(math.Pi*2.0*(i+1)/12.0 + phi)
			dy1 := math.Sin(math.Pi*2.0*(i+1)/12.0 + phi)

			cv.BeginPath()
			cv.LineTo(x1+dx0*DZ*r1*2/3, y1-DZ*2+dy0*DZ*r1+dx0*DZ*r1*1/3*dir)
			cv.LineTo(x1+dx0*DZ*r3*2/3, y1-DZ*2+dy0*DZ*r3+dx0*DZ*r3*1/3*dir)
			cv.LineTo(x2+dx0*DZ*r3*2/3, y2-DZ*2+dy0*DZ*r3+dx0*DZ*r3*1/3*dir)
			cv.LineTo(x2+dx0*DZ*r1*2/3, y2-DZ*2+dy0*DZ*r1+dx0*DZ*r1*1/3*dir)
			cv.ClosePath()
			cv.Fill()

			cv.BeginPath()
			cv.LineTo(x1+dx0*DZ*r1*2/3, y1-DZ*2+dy0*DZ*r1+dx0*DZ*r1*1/3*dir)
			cv.LineTo(x1+dx1*DZ*r1*2/3, y1-DZ*2+dy1*DZ*r1+dx1*DZ*r1*1/3*dir)
			cv.LineTo(x2+dx1*DZ*r1*2/3, y2-DZ*2+dy1*DZ*r1+dx1*DZ*r1*1/3*dir)
			cv.LineTo(x2+dx0*DZ*r1*2/3, y2-DZ*2+dy0*DZ*r1+dx0*DZ*r1*1/3*dir)
			cv.ClosePath()
			cv.Fill()
		}

		cv.SetFillStyle("texture/building/waterwheel_wood.png")
		for i := 0.0; i < 12; i++ {
			dx0 := math.Cos(math.Pi*2.0*i/12.0 + phi)
			dy0 := math.Sin(math.Pi*2.0*i/12.0 + phi)
			dx1 := math.Cos(math.Pi*2.0*(i+1)/12.0 + phi)
			dy1 := math.Sin(math.Pi*2.0*(i+1)/12.0 + phi)

			cv.BeginPath()
			cv.LineTo(x2+dx0*DZ*r2*2/3, y2-DZ*2+dy0*DZ*r2+dx0*DZ*r2*1/3*dir)
			cv.LineTo(x2+dx0*DZ*r1*2/3, y2-DZ*2+dy0*DZ*r1+dx0*DZ*r1*1/3*dir)
			cv.LineTo(x2+dx1*DZ*r1*2/3, y2-DZ*2+dy1*DZ*r1+dx1*DZ*r1*1/3*dir)
			cv.LineTo(x2+dx1*DZ*r2*2/3, y2-DZ*2+dy1*DZ*r2+dx1*DZ*r2*1/3*dir)
			cv.ClosePath()
			cv.Fill()

			cv.BeginPath()
			cv.MoveTo(x2+dx0*DZ*r1*2/3, y2-DZ*2+dy0*DZ*r1+dx0*DZ*r1*1/3*dir)
			cv.LineTo(x2-dx0*DZ*r1*2/3, y2-DZ*2-dy0*DZ*r1-dx0*DZ*r1*1/3*dir)
			cv.ClosePath()
			cv.Stroke()
		}
	} else if extension.T == building.Forge {
		RenderBuildingUnit(cv, building.ForgeBuildingUnit(extension.B, materials.GetMaterial("stone"), extension.Construction), rf, 0, c)
		if !extension.Construction {
			cv.SetFillStyle("texture/building/fire_" + strconv.Itoa(int(phase/3)) + ".png")
			RenderWindows(cv, rf, 1, 2, 0, false, false)
			RenderWindows(cv, rf, 2, 3, 0, false, false)
		}
		buildings.RenderBuildingRoof(cv, building.ForgeBuildingRoof(extension.B, materials.GetMaterial("tile"), extension.Construction), rf, 1, c)
		if !extension.Construction {
			buildings.RenderChimney(cv, rf, 1, phase)
		}
	} else if extension.T == building.Deck {
		if extension.IsConstruction() {
			return
		}

		rfIdx1 := (3 - (-c.Perspective + extension.Direction)) % 4
		rfIdx2 := (2 - (-c.Perspective + extension.Direction)) % 4
		rfIdx3 := (1 - (-c.Perspective + extension.Direction)) % 4
		rfIdx4 := (0 - (-c.Perspective + extension.Direction)) % 4

		xw1 := (rf.X[rfIdx1]*2.0 + rf.X[rfIdx4]*1.0) / 3.0
		yw1 := (rf.Y[rfIdx1]*2.0 + rf.Y[rfIdx4]*1.0) / 3.0
		zw1 := (rf.Z[rfIdx1]*2.0 + rf.Z[rfIdx4]*1.0) / 3.0
		xw2 := (rf.X[rfIdx2]*2.0 + rf.X[rfIdx3]*1.0) / 3.0
		yw2 := (rf.Y[rfIdx2]*2.0 + rf.Y[rfIdx3]*1.0) / 3.0
		zw2 := (rf.Z[rfIdx2]*2.0 + rf.Z[rfIdx3]*1.0) / 3.0

		cv.SetFillStyle("texture/building/deck.png")
		cv.BeginPath()
		cv.LineTo(rf.X[rfIdx1], rf.Y[rfIdx1]-rf.Z[rfIdx1])
		cv.LineTo(rf.X[rfIdx2], rf.Y[rfIdx2]-rf.Z[rfIdx2])
		cv.LineTo(xw2, yw2-zw2)
		cv.LineTo(xw1, yw1-zw1)
		cv.ClosePath()
		cv.Fill()
	}
}
