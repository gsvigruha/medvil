package buildings

import (
	"github.com/tfriedel6/canvas"
	"medvil/model/building"
	"medvil/renderer"
)

func RenderOrnaments(cv *canvas.Canvas, unit *building.BuildingUnit, rf renderer.RenderedField, rw renderer.RenderedWall) {
	if unit.B.Plan.BuildingType == building.BuildingTypeWorkshop || unit.B.Plan.BuildingType == building.BuildingTypeTownhall {
		cv.SetFillStyle("texture/building/ornament.png")
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
	} else if unit.B.Plan.BuildingType == building.BuildingTypeFarm {
		cv.SetFillStyle("texture/building/ornament_wood.png")
		cv.BeginPath()
		cv.LineTo(rw.X[1], rw.Y[1])
		cv.LineTo(rw.X[2], rw.Y[2])
		cv.LineTo(rw.X[2], rw.Y[2]+5)
		cv.LineTo(rw.X[1], rw.Y[1]+5)
		cv.ClosePath()
		cv.Fill()
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
}
