package view

import (
	"github.com/tfriedel6/canvas"
	"image/color"
	"math"
	"medvil/controller"
	"medvil/model/building"
	"medvil/model/materials"
	"medvil/renderer"
	"strconv"
)

const BuildingUnitHeight = 3

const BuildingAnimationMaxPhase = 24

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

func RoofMaterialName(m *materials.Material, shape uint8) string {
	if m == materials.GetMaterial("tile") {
		if shape == 0 {
			return "tile_red"
		} else if shape == 1 {
			return "tile_darkred"
		} else if shape == 2 {
			return "tile_darkred"
		} else if shape == 3 {
			return "tile_red"
		} else if shape == 4 {
			return "tile_darkred"
		}
	}
	return m.Name
}

func RenderWindows(cv *canvas.Canvas, rf renderer.RenderedField, rfIdx1, rfIdx2 uint8, z float64, door bool) {
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

	cv.SetStrokeStyle(color.RGBA{R: 128, G: 64, B: 32, A: 32})
	cv.SetLineWidth(3)
	cv.BeginPath()
	cv.LineTo(rf.X[rfIdx1], rf.Y[rfIdx1]-z-BuildingUnitHeight*DZ*1/3+2)
	cv.LineTo(rf.X[rfIdx2], rf.Y[rfIdx2]-z-BuildingUnitHeight*DZ*1/3+2)
	cv.ClosePath()
	cv.Stroke()
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

			z := math.Min(math.Min(math.Min(rf.Z[0], rf.Z[1]), rf.Z[2]), rf.Z[3]) + float64(k*BuildingUnitHeight)*DZ
			if wall.Windows && !unit.Construction {
				cv.SetFillStyle("texture/building/glass_2.png")
				cv.SetStrokeStyle(color.RGBA{R: 64, G: 32, B: 0, A: 64})
				cv.SetLineWidth(2)
				RenderWindows(cv, rf, rfIdx1, rfIdx2, z, wall.Door)
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
					cv.SetStrokeStyle("#320")
					cv.SetFillStyle("#320")
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

func RenderPolygon(cv *canvas.Canvas, polygon renderer.Polygon, stroke bool) {
	cv.BeginPath()
	for _, p := range polygon.Points {
		cv.LineTo(p.X, p.Y)
	}
	cv.ClosePath()
	cv.Fill()
	if stroke {
		cv.Stroke()
	}
}

func RenderBuildingRoof(cv *canvas.Canvas, roof *building.RoofUnit, rf renderer.RenderedField, k int, c *controller.Controller) *renderer.RenderedBuildingRoof {
	if roof == nil {
		return nil
	}
	var roofPolygons []renderer.Polygon
	startL := 2 + c.Perspective
	if roof.Roof.RoofType == building.RoofTypeSplit {
		z := math.Min(math.Min(math.Min(rf.Z[0], rf.Z[1]), rf.Z[2]), rf.Z[3]) + float64(k*BuildingUnitHeight)*DZ
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

				sideMidX := (rf.X[rfIdx1] + rf.X[rfIdx2]) / 2
				sideMidY := (rf.Y[rfIdx1] + rf.Y[rfIdx2]) / 2
				rp1 := renderer.Polygon{Points: []renderer.Point{
					renderer.Point{X: rf.X[rfIdx1], Y: rf.Y[rfIdx1] - z},
					renderer.Point{X: sideMidX, Y: sideMidY - z - BuildingUnitHeight*DZ},
					renderer.Point{X: midX, Y: midY - z - BuildingUnitHeight*DZ},
				}}
				rp2 := renderer.Polygon{Points: []renderer.Point{
					renderer.Point{X: rf.X[rfIdx2], Y: rf.Y[rfIdx2] - z},
					renderer.Point{X: sideMidX, Y: sideMidY - z - BuildingUnitHeight*DZ},
					renderer.Point{X: midX, Y: midY - z - BuildingUnitHeight*DZ},
				}}
				roofPolygons = append(roofPolygons, rp1, rp2)

				if cv != nil {
					if !roof.Construction {
						cv.SetFillStyle("texture/building/" + RoofMaterialName(roof.Roof.M, roof.B.Shape) + suffix + ".png")
					} else {
						cv.SetFillStyle("texture/building/construction" + suffix + ".png")
					}

					cv.SetStrokeStyle(color.RGBA{R: 192, G: 128, B: 64, A: 32})
					cv.SetLineWidth(3)

					RenderPolygon(cv, rp1, true)
					RenderPolygon(cv, rp2, true)
				}
			} else {
				var suffix = ""
				if rfIdx1%2 == 1 {
					suffix = "_flipped"
				}

				rp1 := renderer.Polygon{Points: []renderer.Point{
					renderer.Point{X: rf.X[rfIdx1], Y: rf.Y[rfIdx1] - z},
					renderer.Point{X: rf.X[rfIdx2], Y: rf.Y[rfIdx2] - z},
					renderer.Point{X: midX, Y: midY - z - BuildingUnitHeight*DZ},
				}}
				roofPolygons = append(roofPolygons, rp1)

				if cv != nil {
					if !roof.Construction {
						cv.SetFillStyle("texture/building/" + RoofMaterialName(roof.Roof.M, roof.B.Shape) + suffix + ".png")
					} else {
						cv.SetFillStyle("texture/building/construction" + suffix + ".png")
					}
					cv.SetStrokeStyle(color.RGBA{R: 64, G: 32, B: 0, A: 32})
					cv.SetLineWidth(3)
					RenderPolygon(cv, rp1, true)
				}
			}
		}
	} else if roof.Roof.RoofType == building.RoofTypeFlat {
		z := float64(k*BuildingUnitHeight) * DZ
		if !roof.Construction {
			rp1 := renderer.Polygon{Points: []renderer.Point{
				renderer.Point{X: rf.X[0], Y: rf.Y[0] - rf.Z[0] - z},
				renderer.Point{X: rf.X[1], Y: rf.Y[1] - rf.Z[1] - z},
				renderer.Point{X: rf.X[2], Y: rf.Y[2] - rf.Z[2] - z},
				renderer.Point{X: rf.X[3], Y: rf.Y[3] - rf.Z[3] - z},
			}}
			roofPolygons = append(roofPolygons, rp1)
			if cv != nil {
				cv.SetFillStyle("texture/building/" + RoofMaterialName(roof.Roof.M, roof.B.Shape) + "_flat.png")
				RenderPolygon(cv, rp1, false)
			}
		}
	} else if roof.Roof.RoofType == building.RoofTypeRamp {
		z := math.Min(math.Min(math.Min(rf.Z[0], rf.Z[1]), rf.Z[2]), rf.Z[3]) + float64(k*BuildingUnitHeight)*DZ
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
				rp1 := renderer.Polygon{Points: []renderer.Point{
					renderer.Point{X: rf.X[rfIdx1], Y: rf.Y[rfIdx1] - z - BuildingUnitHeight*DZ},
					renderer.Point{X: rf.X[rfIdx1], Y: rf.Y[rfIdx1] - z},
					renderer.Point{X: rf.X[rfIdx4], Y: rf.Y[rfIdx4] - z},
				}}
				rp2 := renderer.Polygon{Points: []renderer.Point{
					renderer.Point{X: rf.X[rfIdx2], Y: rf.Y[rfIdx2] - z - BuildingUnitHeight*DZ},
					renderer.Point{X: rf.X[rfIdx2], Y: rf.Y[rfIdx2] - z},
					renderer.Point{X: rf.X[rfIdx3], Y: rf.Y[rfIdx3] - z},
				}}
				rp3 := renderer.Polygon{Points: []renderer.Point{
					renderer.Point{X: rf.X[rfIdx1], Y: rf.Y[rfIdx1] - z - BuildingUnitHeight*DZ},
					renderer.Point{X: rf.X[rfIdx2], Y: rf.Y[rfIdx2] - z - BuildingUnitHeight*DZ},
					renderer.Point{X: rf.X[rfIdx3], Y: rf.Y[rfIdx3] - z},
					renderer.Point{X: rf.X[rfIdx4], Y: rf.Y[rfIdx4] - z},
				}}
				roofPolygons = append(roofPolygons, rp1, rp2, rp3)

				if cv != nil {
					if !roof.Construction {
						cv.SetFillStyle("texture/building/" + RoofMaterialName(roof.Roof.M, roof.B.Shape) + suffix + ".png")
					} else {
						cv.SetFillStyle("texture/building/construction" + suffix + ".png")
					}

					cv.SetStrokeStyle(color.RGBA{R: 192, G: 128, B: 64, A: 32})
					cv.SetLineWidth(3)

					RenderPolygon(cv, rp1, true)
					RenderPolygon(cv, rp2, true)

					if !roof.Construction {
						cv.SetFillStyle("texture/building/" + RoofMaterialName(roof.Roof.M, roof.B.Shape) + "_flat.png")
					} else {
						cv.SetFillStyle("texture/building/construction" + suffix + ".png")
					}

					RenderPolygon(cv, rp3, true)
				}
			}
		}
	}
	return &renderer.RenderedBuildingRoof{B: roof.Building(), Ps: roofPolygons}
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
		rfIdx1 := (3 - (-c.Perspective + extension.Direction)) % 4
		rfIdx2 := (2 - (-c.Perspective + extension.Direction)) % 4
		rfIdx3 := (1 - (-c.Perspective + extension.Direction)) % 4
		rfIdx4 := (0 - (-c.Perspective + extension.Direction)) % 4
		if !extension.Construction {
			cv.SetFillStyle("texture/building/fire_" + strconv.Itoa(int(phase/3)) + ".png")
			if rfIdx1 == 2 || rfIdx1 == 3 {
				RenderWindows(cv, rf, rfIdx1, rfIdx2, 0, false)
			} else {
				RenderWindows(cv, rf, rfIdx3, rfIdx4, 0, false)
			}
		}
		RenderBuildingRoof(cv, building.ForgeBuildingRoof(extension.B, materials.GetMaterial("tile"), extension.Construction), rf, 1, c)
		if !extension.Construction {
			RenderChimey(cv, rf, 1)
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

func RenderChimey(cv *canvas.Canvas, rf renderer.RenderedField, k int) {
	z := math.Min(math.Min(math.Min(rf.Z[0], rf.Z[1]), rf.Z[2]), rf.Z[3]) + float64(k*BuildingUnitHeight)*DZ
	midX := (rf.X[0] + rf.X[2]) / 2
	midY := (rf.Y[0] + rf.Y[2]) / 2
	h := 8.0
	rp1 := renderer.Polygon{Points: []renderer.Point{
		renderer.Point{X: midX, Y: midY - z - BuildingUnitHeight*DZ + 12},
		renderer.Point{X: midX, Y: midY - z - BuildingUnitHeight*DZ - h},
		renderer.Point{X: midX - 9, Y: midY - z - BuildingUnitHeight*DZ - h - 6},
		renderer.Point{X: midX - 9, Y: midY - z - BuildingUnitHeight*DZ + 6},
	}}
	cv.SetFillStyle("texture/building/stone.png")
	RenderPolygon(cv, rp1, true)

	rp2 := renderer.Polygon{Points: []renderer.Point{
		renderer.Point{X: midX, Y: midY - z - BuildingUnitHeight*DZ + 12},
		renderer.Point{X: midX, Y: midY - z - BuildingUnitHeight*DZ - h},
		renderer.Point{X: midX + 9, Y: midY - z - BuildingUnitHeight*DZ - h - 6},
		renderer.Point{X: midX + 9, Y: midY - z - BuildingUnitHeight*DZ + 6},
	}}
	cv.SetFillStyle("texture/building/stone_flipped.png")
	RenderPolygon(cv, rp2, true)

	rp3 := renderer.Polygon{Points: []renderer.Point{
		renderer.Point{X: midX + 9, Y: midY - z - BuildingUnitHeight*DZ - h - 6},
		renderer.Point{X: midX, Y: midY - z - BuildingUnitHeight*DZ - h},
		renderer.Point{X: midX - 9, Y: midY - z - BuildingUnitHeight*DZ - h - 6},
		renderer.Point{X: midX, Y: midY - z - BuildingUnitHeight*DZ - h - 12},
	}}
	cv.SetFillStyle("texture/building/stone_flat.png")
	RenderPolygon(cv, rp3, true)
	cv.SetFillStyle(color.RGBA{R: 0, G: 0, B: 0, A: 224})
	RenderPolygon(cv, rp3, true)
}
