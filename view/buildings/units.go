package buildings

import (
	"github.com/tfriedel6/canvas"
	"image/color"
	"math"
	"medvil/controller"
	"medvil/model/building"
	"medvil/model/materials"
	"medvil/renderer"
	"path/filepath"
)

func BrickMaterialName(shape uint8) string {
	switch shape % 5 {
	case 0:
		return "painted_yellow"
	case 1:
		return "painted_red"
	case 2:
		return "painted_brown"
	case 3:
		return "painted_beige"
	case 4:
		return "painted_sand"
	}
	return "painted"
}

func WallMaterialName(t building.BuildingType, m *materials.Material, shape uint8, broken bool) string {
	if m == materials.GetMaterial("brick") {
		return BrickMaterialName(shape)
	}
	if t == building.BuildingTypeWorkshop && m == materials.GetMaterial("stone") {
		return "stone_dark"
	}
	if t == building.BuildingTypeWall && m == materials.GetMaterial("stone") {
		if broken {
			return "stone_broken"
		}
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
				cv.SetFillStyle(filepath.FromSlash("texture/building/" + WallMaterialName(unit.B.Plan.BuildingType, wall.M, unit.B.Shape, unit.B.Broken) + suffix + ".png"))
			} else {
				cv.SetFillStyle(filepath.FromSlash("texture/building/construction" + suffix + ".png"))
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

			if !wall.Arch && !unit.Construction {
				RenderOrnaments(cv, unit, rf, rw)
			}

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
				cv.SetLineWidth(2)
				if wall.Windows == building.WindowTypeFactory {
					cv.SetFillStyle(filepath.FromSlash("texture/building/glass_3.png"))
					cv.SetStrokeStyle(color.RGBA{R: 32, G: 32, B: 0, A: 192})
					RenderFactoryWindows(cv, rf, rfIdx1, rfIdx2, z, wall.Door)
				} else {
					cv.SetFillStyle(filepath.FromSlash("texture/building/glass_2.png"))
					cv.SetStrokeStyle(color.RGBA{R: 32, G: 32, B: 0, A: 64})
					RenderWindows(cv, rf, rfIdx1, rfIdx2, z, wall.Door, wall.Windows == building.WindowTypeFrench)
					if wall.Windows == building.WindowTypeBalcony {
						RenderBalcony(cv, rf, rfIdx1, rfIdx2, z, wall.Door)
					}
				}
			}

			if wall.Door && !unit.Construction {
				cv.SetFillStyle(filepath.FromSlash("texture/building/door.png"))
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
					cv.DrawImage(filepath.FromSlash("icon/gui/tasks/"+workshop.Manufacture.Name+".png"), xm+dX*2, ym+5, 16*dX, 16)
				}
			}
		}
	}
	return renderer.RenderedBuildingUnit{Walls: rws, Unit: unit}
}
