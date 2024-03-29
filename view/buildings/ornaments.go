package buildings

import (
	"github.com/tfriedel6/canvas"
	"medvil/controller"
	"medvil/model/building"
	"medvil/renderer"
	"path/filepath"
	"strconv"
)

func RenderOrnaments(cv *canvas.Canvas, unit *building.BuildingUnit, rf renderer.RenderedField, rw renderer.RenderedWall) {
	if unit.B.Plan.BuildingType == building.BuildingTypeWorkshop || unit.B.Plan.BuildingType == building.BuildingTypeTownhall {
		if unit.B.Shape%2 == 0 {
			cv.SetFillStyle(filepath.FromSlash("texture/building/ornament.png"))
		} else {
			cv.SetFillStyle(filepath.FromSlash("texture/building/ornament_dark.png"))
		}
		cv.BeginPath()
		cv.LineTo(rw.X[1], rw.Y[1])
		cv.LineTo(rw.X[2], rw.Y[2])
		cv.LineTo(rw.X[2], rw.Y[2]+5)
		cv.LineTo(rw.X[1], rw.Y[1]+5)
		cv.ClosePath()
		cv.Fill()
		cv.BeginPath()
		cv.LineTo(rw.X[1], rw.Y[1])
		cv.LineTo(rw.X[1], rw.Y[1]*0.75+rw.Y[0]*0.25)
		cv.LineTo(rw.X[1]*0.8+rw.X[2]*0.2, rw.Y[1]*0.8+rw.Y[2]*0.2)
		cv.ClosePath()
		cv.Fill()
		cv.BeginPath()
		cv.LineTo(rw.X[2], rw.Y[2])
		cv.LineTo(rw.X[2], rw.Y[2]*0.75+rw.Y[3]*0.25)
		cv.LineTo(rw.X[2]*0.8+rw.X[1]*0.2, rw.Y[2]*0.8+rw.Y[1]*0.2)
		cv.ClosePath()
		cv.Fill()
	} else if unit.B.Plan.BuildingType == building.BuildingTypeFarm || unit.B.Plan.BuildingType == building.BuildingTypeMine {
		cv.SetFillStyle(filepath.FromSlash("texture/building/ornament_wood.png"))
		cv.BeginPath()
		cv.LineTo(rw.X[1], rw.Y[1])
		cv.LineTo(rw.X[2], rw.Y[2])
		cv.LineTo(rw.X[2], rw.Y[2]+5)
		cv.LineTo(rw.X[1], rw.Y[1]+5)
		cv.ClosePath()
		cv.Fill()
		if unit.B.Plan.BuildingType == building.BuildingTypeFarm {
			cv.BeginPath()
			cv.LineTo(rw.X[1], rw.Y[1]*0.7+rw.Y[0]*0.3)
			cv.LineTo(rw.X[1], rw.Y[1]*0.6+rw.Y[0]*0.4)
			cv.LineTo(rw.X[1]*0.7+rw.X[2]*0.3, rw.Y[1]*0.7+rw.Y[2]*0.3)
			cv.LineTo(rw.X[1]*0.8+rw.X[2]*0.2, rw.Y[1]*0.8+rw.Y[2]*0.2)
			cv.ClosePath()
			cv.Fill()
			cv.BeginPath()
			cv.LineTo(rw.X[2], rw.Y[2]*0.7+rw.Y[3]*0.3)
			cv.LineTo(rw.X[2], rw.Y[2]*0.6+rw.Y[3]*0.4)
			cv.LineTo(rw.X[2]*0.7+rw.X[1]*0.3, rw.Y[2]*0.7+rw.Y[1]*0.3)
			cv.LineTo(rw.X[2]*0.8+rw.X[1]*0.2, rw.Y[2]*0.8+rw.Y[1]*0.2)
			cv.ClosePath()
			cv.Fill()
		}
	} else if unit.B.Plan.BuildingType == building.BuildingTypeWall {
		cv.SetFillStyle(filepath.FromSlash("texture/building/ornament_wood_dark.png"))
		cv.BeginPath()
		cv.LineTo(rw.X[1], rw.Y[1])
		cv.LineTo(rw.X[2], rw.Y[2])
		cv.LineTo(rw.X[2], rw.Y[2]+5)
		cv.LineTo(rw.X[1], rw.Y[1]+5)
		cv.ClosePath()
		cv.Fill()
		cv.BeginPath()
		cv.LineTo(rw.X[1], rw.Y[1])
		cv.LineTo(rw.X[1], rw.Y[1]*0.75+rw.Y[0]*0.25)
		cv.LineTo(rw.X[1]*0.8+rw.X[2]*0.2, rw.Y[1]*0.8+rw.Y[2]*0.2)
		cv.ClosePath()
		cv.Fill()
		cv.BeginPath()
		cv.LineTo(rw.X[2], rw.Y[2])
		cv.LineTo(rw.X[2], rw.Y[2]*0.75+rw.Y[3]*0.25)
		cv.LineTo(rw.X[2]*0.8+rw.X[1]*0.2, rw.Y[2]*0.8+rw.Y[1]*0.2)
		cv.ClosePath()
		cv.Fill()
	}
}

func RenderRoofFence(cv *canvas.Canvas, roof *building.RoofUnit, rp1 renderer.Polygon, c *controller.Controller) {
	if roof.B.Plan.BuildingType == building.BuildingTypeWorkshop || roof.B.Plan.BuildingType == building.BuildingTypeFactory {
		cv.SetFillStyle(filepath.FromSlash("texture/building/ornament_2.png"))
		for i := uint8(0); i < 4; i++ {
			if !roof.Connected[i] {
				rfIdx1 := (2 - (-c.Perspective + i)) % 4
				rfIdx2 := (3 - (-c.Perspective + i)) % 4
				dx := (rp1.Points[rfIdx2].X - rp1.Points[rfIdx1].X) / 5.0
				dy := (rp1.Points[rfIdx2].Y - rp1.Points[rfIdx1].Y) / 5.0
				for j := float64(0); j <= 5; j++ {
					cv.BeginPath()
					cv.LineTo(rp1.Points[rfIdx1].X+dx*j-3, rp1.Points[rfIdx1].Y+dy*j)
					cv.LineTo(rp1.Points[rfIdx1].X+dx*j, rp1.Points[rfIdx1].Y+dy*j+2)
					cv.LineTo(rp1.Points[rfIdx1].X+dx*j+3, rp1.Points[rfIdx1].Y+dy*j)
					cv.LineTo(rp1.Points[rfIdx1].X+dx*j+3, rp1.Points[rfIdx1].Y+dy*j-15)
					cv.LineTo(rp1.Points[rfIdx1].X+dx*j-3, rp1.Points[rfIdx1].Y+dy*j-15)
					cv.ClosePath()
					cv.Fill()
				}
				cv.BeginPath()
				cv.LineTo(rp1.Points[rfIdx1].X, rp1.Points[rfIdx1].Y-12)
				cv.LineTo(rp1.Points[rfIdx1].X, rp1.Points[rfIdx1].Y-17)
				cv.LineTo(rp1.Points[rfIdx2].X, rp1.Points[rfIdx2].Y-17)
				cv.LineTo(rp1.Points[rfIdx2].X, rp1.Points[rfIdx2].Y-12)
				cv.ClosePath()
				cv.Fill()
			}
		}
	}
}

func RenderWallOrnaments(cv *canvas.Canvas, unit *building.BuildingUnit, rf renderer.RenderedField, rw renderer.RenderedWall) {
	if unit.B.Plan.BuildingType == building.BuildingTypeWall {
		sfx := int(unit.B.Shape) % 3
		if unit.B.Shape%4 == 1 {
			cv.SetFillStyle(filepath.FromSlash("texture/building/wall_rose_" + strconv.Itoa(sfx) + ".png"))
		} else if unit.B.Shape%4 == 2 {
			cv.SetFillStyle(filepath.FromSlash("texture/building/wall_green_" + strconv.Itoa(sfx) + ".png"))
		} else {
			return
		}
		cv.BeginPath()
		cv.LineTo(rw.X[0], rw.Y[0])
		cv.LineTo(rw.X[1], rw.Y[1])
		cv.LineTo(rw.X[2], rw.Y[2])
		cv.LineTo(rw.X[3], rw.Y[3])
		cv.ClosePath()
		cv.Fill()
	}
}
