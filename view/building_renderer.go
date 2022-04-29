package view

import (
	"github.com/tfriedel6/canvas"
	"image/color"
	"math"
	"medvil/controller"
	"medvil/model/building"
	"medvil/renderer"
)

const BuildingUnitHeight = 3

const BuildingAnimationMaxPhase = 24

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
		if !unit.Construction && (rfIdx1 == 0 || rfIdx1 == 1) {
			continue
		}
		var suffix = ""
		if rfIdx1%2 == 1 {
			suffix = "_flipped"
		}
		if cv != nil {
			if !unit.Construction {
				cv.SetFillStyle("texture/building/" + wall.M.Name + suffix + ".png")
			} else {
				cv.SetFillStyle("texture/building/construction" + suffix + ".png")
			}
		}
		z := math.Min(math.Min(math.Min(rf.Z[0], rf.Z[1]), rf.Z[2]), rf.Z[3]) + float64(k*BuildingUnitHeight)*DZ

		rw := renderer.RenderedWall{
			X:    [4]float64{rf.X[rfIdx1], rf.X[rfIdx1], rf.X[rfIdx2], rf.X[rfIdx2]},
			Y:    [4]float64{rf.Y[rfIdx1] - z, rf.Y[rfIdx1] - z - BuildingUnitHeight*DZ, rf.Y[rfIdx2] - z - BuildingUnitHeight*DZ, rf.Y[rfIdx2] - z},
			Wall: wall,
		}
		rws = append(rws, rw)
		if cv != nil {
			rw.Draw(cv)
			cv.Fill()

			if wall.Windows && !unit.Construction {
				cv.SetFillStyle("texture/building/glass.png")
				cv.SetStrokeStyle(color.RGBA{R: 64, G: 32, B: 0, A: 64})
				cv.SetLineWidth(2)

				cv.BeginPath()
				cv.LineTo((6*rf.X[rfIdx1]+1*rf.X[rfIdx2])/7, (6*rf.Y[rfIdx1]+1*rf.Y[rfIdx2])/7-z-BuildingUnitHeight*DZ*1/3)
				cv.LineTo((6*rf.X[rfIdx1]+1*rf.X[rfIdx2])/7, (6*rf.Y[rfIdx1]+1*rf.Y[rfIdx2])/7-z-BuildingUnitHeight*DZ*2/3)
				cv.LineTo((5*rf.X[rfIdx1]+2*rf.X[rfIdx2])/7, (5*rf.Y[rfIdx1]+2*rf.Y[rfIdx2])/7-z-BuildingUnitHeight*DZ*2/3)
				cv.LineTo((5*rf.X[rfIdx1]+2*rf.X[rfIdx2])/7, (5*rf.Y[rfIdx1]+2*rf.Y[rfIdx2])/7-z-BuildingUnitHeight*DZ*1/3)
				cv.ClosePath()
				cv.Fill()
				cv.Stroke()

				cv.BeginPath()
				cv.LineTo((4*rf.X[rfIdx1]+3*rf.X[rfIdx2])/7, (4*rf.Y[rfIdx1]+3*rf.Y[rfIdx2])/7-z-BuildingUnitHeight*DZ*1/3)
				cv.LineTo((4*rf.X[rfIdx1]+3*rf.X[rfIdx2])/7, (4*rf.Y[rfIdx1]+3*rf.Y[rfIdx2])/7-z-BuildingUnitHeight*DZ*2/3)
				cv.LineTo((3*rf.X[rfIdx1]+4*rf.X[rfIdx2])/7, (3*rf.Y[rfIdx1]+4*rf.Y[rfIdx2])/7-z-BuildingUnitHeight*DZ*2/3)
				cv.LineTo((3*rf.X[rfIdx1]+4*rf.X[rfIdx2])/7, (3*rf.Y[rfIdx1]+4*rf.Y[rfIdx2])/7-z-BuildingUnitHeight*DZ*1/3)
				cv.ClosePath()
				cv.Fill()
				cv.Stroke()

				cv.BeginPath()
				cv.LineTo((2*rf.X[rfIdx1]+5*rf.X[rfIdx2])/7, (2*rf.Y[rfIdx1]+5*rf.Y[rfIdx2])/7-z-BuildingUnitHeight*DZ*1/3)
				cv.LineTo((2*rf.X[rfIdx1]+5*rf.X[rfIdx2])/7, (2*rf.Y[rfIdx1]+5*rf.Y[rfIdx2])/7-z-BuildingUnitHeight*DZ*2/3)
				cv.LineTo((1*rf.X[rfIdx1]+6*rf.X[rfIdx2])/7, (1*rf.Y[rfIdx1]+6*rf.Y[rfIdx2])/7-z-BuildingUnitHeight*DZ*2/3)
				cv.LineTo((1*rf.X[rfIdx1]+6*rf.X[rfIdx2])/7, (1*rf.Y[rfIdx1]+6*rf.Y[rfIdx2])/7-z-BuildingUnitHeight*DZ*1/3)
				cv.ClosePath()
				cv.Fill()
				cv.Stroke()
			}

			if wall.Door && !unit.Construction {
				cv.SetFillStyle("texture/building/door.png")
				cv.SetStrokeStyle(color.RGBA{R: 64, G: 32, B: 0, A: 64})
				cv.SetLineWidth(2)

				cv.BeginPath()
				cv.LineTo((3*rf.X[rfIdx1]+4*rf.X[rfIdx2])/7, (3*rf.Y[rfIdx1]+4*rf.Y[rfIdx2])/7-z)
				cv.LineTo((3*rf.X[rfIdx1]+4*rf.X[rfIdx2])/7, (3*rf.Y[rfIdx1]+4*rf.Y[rfIdx2])/7-z-BuildingUnitHeight*DZ*2/3)
				cv.LineTo((1*rf.X[rfIdx1]+6*rf.X[rfIdx2])/7, (1*rf.Y[rfIdx1]+6*rf.Y[rfIdx2])/7-z-BuildingUnitHeight*DZ*2/3)
				cv.LineTo((1*rf.X[rfIdx1]+6*rf.X[rfIdx2])/7, (1*rf.Y[rfIdx1]+6*rf.Y[rfIdx2])/7-z)
				cv.ClosePath()
				cv.Fill()
				cv.Stroke()
			}
		}
	}
	return renderer.RenderedBuildingUnit{Walls: rws, Unit: unit}
}

func RenderBuildingRoof(cv *canvas.Canvas, roof *building.RoofUnit, rf renderer.RenderedField, k int, c *controller.Controller) {
	if roof == nil {
		return
	}
	startL := 2 + c.Perspective
	z := math.Min(math.Min(math.Min(rf.Z[0], rf.Z[1]), rf.Z[2]), rf.Z[3]) + float64(k*BuildingUnitHeight)*DZ
	if roof.Roof.RoofType == building.RoofTypeSplit {
		midX := (rf.X[0] + rf.X[2]) / 2
		midY := (rf.Y[0] + rf.Y[2]) / 2

		for l := uint8(startL); l < 4+startL; l++ {
			rfIdx1 := (3 - (-c.Perspective + l)) % 4
			rfIdx2 := (2 - (-c.Perspective + l)) % 4
			if roof.Elevated[l%4] {
				var suffix = ""
				if rfIdx1%2 == 0 {
					suffix = "_flipped"
				}
				if !roof.Construction {
					cv.SetFillStyle("texture/building/" + roof.Roof.M.Name + suffix + ".png")
				} else {
					cv.SetFillStyle("texture/building/construction" + suffix + ".png")
				}

				cv.SetStrokeStyle(color.RGBA{R: 192, G: 128, B: 64, A: 32})
				cv.SetLineWidth(3)

				sideMidX := (rf.X[rfIdx1] + rf.X[rfIdx2]) / 2
				sideMidY := (rf.Y[rfIdx1] + rf.Y[rfIdx2]) / 2
				cv.BeginPath()
				cv.LineTo(rf.X[rfIdx1], rf.Y[rfIdx1]-z)
				cv.LineTo(sideMidX, sideMidY-z-BuildingUnitHeight*DZ)
				cv.LineTo(midX, midY-z-BuildingUnitHeight*DZ)
				cv.ClosePath()
				cv.Fill()
				cv.Stroke()

				cv.BeginPath()
				cv.LineTo(rf.X[rfIdx2], rf.Y[rfIdx2]-z)
				cv.LineTo(sideMidX, sideMidY-z-BuildingUnitHeight*DZ)
				cv.LineTo(midX, midY-z-BuildingUnitHeight*DZ)
				cv.ClosePath()
				cv.Fill()
				cv.Stroke()
			} else {
				var suffix = ""
				if rfIdx1%2 == 1 {
					suffix = "_flipped"
				}
				if !roof.Construction {
					cv.SetFillStyle("texture/building/" + roof.Roof.M.Name + suffix + ".png")
				} else {
					cv.SetFillStyle("texture/building/construction" + suffix + ".png")
				}

				cv.SetStrokeStyle(color.RGBA{R: 64, G: 32, B: 0, A: 32})
				cv.SetLineWidth(3)

				cv.BeginPath()
				cv.LineTo(rf.X[rfIdx1], rf.Y[rfIdx1]-z)
				cv.LineTo(rf.X[rfIdx2], rf.Y[rfIdx2]-z)
				cv.LineTo(midX, midY-z-BuildingUnitHeight*DZ)
				cv.ClosePath()
				cv.Fill()
				cv.Stroke()
			}
		}
	} else if roof.Roof.RoofType == building.RoofTypeFlat {
		if !roof.Construction {
			cv.SetFillStyle("texture/building/" + roof.Roof.M.Name + "_flat.png")
			cv.BeginPath()
			cv.LineTo(rf.X[0], rf.Y[0]-z)
			cv.LineTo(rf.X[1], rf.Y[1]-z)
			cv.LineTo(rf.X[2], rf.Y[2]-z)
			cv.LineTo(rf.X[3], rf.Y[3]-z)
			cv.ClosePath()
			cv.Fill()
		}
	} else if roof.Roof.RoofType == building.RoofTypeRamp {
		for l := uint8(startL); l < 4+startL; l++ {
			rfIdx1 := (3 - (-c.Perspective + l)) % 4
			rfIdx2 := (2 - (-c.Perspective + l)) % 4
			rfIdx3 := (1 - (-c.Perspective + l)) % 4
			rfIdx4 := (0 - (-c.Perspective + l)) % 4
			if roof.Elevated[l%4] {
				var suffix = ""
				if rfIdx1%2 == 0 {
					suffix = "_flipped"
				}

				if !roof.Construction {
					cv.SetFillStyle("texture/building/" + roof.Roof.M.Name + suffix + ".png")
				} else {
					cv.SetFillStyle("texture/building/construction" + suffix + ".png")
				}

				cv.SetStrokeStyle(color.RGBA{R: 192, G: 128, B: 64, A: 32})
				cv.SetLineWidth(3)

				cv.BeginPath()
				cv.LineTo(rf.X[rfIdx1], rf.Y[rfIdx1]-z-BuildingUnitHeight*DZ)
				cv.LineTo(rf.X[rfIdx1], rf.Y[rfIdx1]-z)
				cv.LineTo(rf.X[rfIdx4], rf.Y[rfIdx4]-z)
				cv.ClosePath()
				cv.Fill()
				cv.Stroke()

				cv.BeginPath()
				cv.LineTo(rf.X[rfIdx2], rf.Y[rfIdx2]-z-BuildingUnitHeight*DZ)
				cv.LineTo(rf.X[rfIdx2], rf.Y[rfIdx2]-z)
				cv.LineTo(rf.X[rfIdx3], rf.Y[rfIdx3]-z)
				cv.ClosePath()
				cv.Fill()
				cv.Stroke()

				if !roof.Construction {
					cv.SetFillStyle("texture/building/" + roof.Roof.M.Name + "_flat.png")
				} else {
					cv.SetFillStyle("texture/building/construction" + suffix + ".png")
				}

				cv.BeginPath()
				cv.LineTo(rf.X[rfIdx1], rf.Y[rfIdx1]-z-BuildingUnitHeight*DZ)
				cv.LineTo(rf.X[rfIdx2], rf.Y[rfIdx2]-z-BuildingUnitHeight*DZ)
				cv.LineTo(rf.X[rfIdx3], rf.Y[rfIdx3]-z)
				cv.LineTo(rf.X[rfIdx4], rf.Y[rfIdx4]-z)
				cv.ClosePath()
				cv.Fill()
				cv.Stroke()
			}
		}
	}
}

func RenderBuildingExtension(cv *canvas.Canvas, extension *building.ExtensionUnit, rf renderer.RenderedField, k int, phase uint8, c *controller.Controller) {
	if extension == nil {
		return
	}
	if extension.IsConstruction() {
		return
	}
	if extension.T == building.WaterMillWheel {
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
		rfIdx1 := (3 - (-c.Perspective + extension.Direction)) % 4
		rfIdx2 := (2 - (-c.Perspective + extension.Direction)) % 4
		rfIdx3 := (1 - (-c.Perspective + extension.Direction)) % 4
		rfIdx4 := (0 - (-c.Perspective + extension.Direction)) % 4
		var suffix1 = ""
		var suffix2 = "_flipped"
		if rfIdx1%2 == 1 {
			suffix1 = "_flipped"
			suffix2 = ""
		}

		x1 := (rf.X[rfIdx1]*3.0 + rf.X[rfIdx2]*1.0) / 4.0
		y1 := ((rf.Y[rfIdx1]-rf.Z[rfIdx1])*3.0 + (rf.Y[rfIdx2]-rf.Z[rfIdx1])*1.0) / 4.0
		x2 := (rf.X[rfIdx1]*1.0 + rf.X[rfIdx2]*3.0) / 4.0
		y2 := ((rf.Y[rfIdx1]-rf.Z[rfIdx1])*1.0 + (rf.Y[rfIdx2] - -rf.Z[rfIdx2])*3.0) / 4.0

		x3 := (rf.X[rfIdx3]*3.0 + rf.X[rfIdx4]*1.0) / 4.0
		y3 := ((rf.Y[rfIdx3]-rf.Z[rfIdx3])*3.0 + (rf.Y[rfIdx4]-rf.Z[rfIdx4])*1.0) / 4.0
		x4 := (rf.X[rfIdx3]*1.0 + rf.X[rfIdx4]*3.0) / 4.0
		y4 := ((rf.Y[rfIdx3]-rf.Z[rfIdx3])*1.0 + (rf.Y[rfIdx4]-rf.Z[rfIdx4])*3.0) / 4.0

		x5 := (x1 + x4) / 2.0
		y5 := (y1 + y4) / 2.0
		x6 := (x2 + x3) / 2.0
		y6 := (y2 + y3) / 2.0

		x7 := (x1 + x5) / 2.0
		y7 := (y1 + y5) / 2.0
		x8 := (x2 + x6) / 2.0
		y8 := (y2 + y6) / 2.0

		cv.SetFillStyle("texture/building/stone" + suffix1 + ".png")

		cv.BeginPath()
		cv.LineTo(x1, y1)
		cv.LineTo(x1, y1-DZ*2)
		cv.LineTo(x2, y2-DZ*2)
		cv.LineTo(x2, y2)
		cv.ClosePath()
		cv.Fill()

		cv.BeginPath()
		cv.LineTo(x5, y5)
		cv.LineTo(x5, y5-DZ)
		cv.LineTo(x6, y6-DZ)
		cv.LineTo(x6, y6)
		cv.ClosePath()
		cv.Fill()

		cv.BeginPath()
		cv.LineTo(x8, y8-DZ)
		cv.LineTo(x8, y8-DZ*2)
		cv.LineTo(x7, y7-DZ*2)
		cv.LineTo(x7, y7-DZ)
		cv.ClosePath()
		cv.Fill()

		cv.SetFillStyle("texture/building/stone" + suffix2 + ".png")

		cv.BeginPath()
		cv.LineTo(x1, y1)
		cv.LineTo(x1, y1-DZ*2)
		cv.LineTo(x7, y7-DZ*2)
		cv.LineTo(x7, y7-DZ)
		cv.LineTo(x5, y5-DZ)
		cv.LineTo(x5, y5)
		cv.ClosePath()
		cv.Fill()

		cv.BeginPath()
		cv.LineTo(x2, y2)
		cv.LineTo(x2, y2-DZ*2)
		cv.LineTo(x8, y8-DZ*2)
		cv.LineTo(x8, y8-DZ)
		cv.LineTo(x6, y6-DZ)
		cv.LineTo(x6, y6)
		cv.ClosePath()
		cv.Fill()

		cv.SetFillStyle("texture/building/stone_flat.png")

		cv.BeginPath()
		cv.LineTo(x8, y8-DZ)
		cv.LineTo(x7, y7-DZ)
		cv.LineTo(x5, y5-DZ)
		cv.LineTo(x6, y6-DZ)
		cv.ClosePath()
		cv.Fill()

		cv.BeginPath()
		cv.LineTo(x8, y8-DZ*2)
		cv.LineTo(x7, y7-DZ*2)
		cv.LineTo(x1, y1-DZ*2)
		cv.LineTo(x2, y2-DZ*2)
		cv.ClosePath()
		cv.Fill()
	}
}
