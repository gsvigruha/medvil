package view

import (
	//"fmt"
	"github.com/tfriedel6/canvas"
	"image/color"
	"math"
	"medvil/controller"
	"medvil/model"
)

const BuildingUnitHeight = 3

func RenderBuildingUnit(cv *canvas.Canvas, unit model.BuildingUnit, rf RenderedField, k int) {
	startI := 2 + controller.Perspective
	for i := uint8(startI); i < 4+startI; i++ {
		wall := unit.Walls[i%4]
		if wall == nil {
			continue
		}
		rfIdx1 := (3 - (-controller.Perspective + i)) % 4
		rfIdx2 := (2 - (-controller.Perspective + i)) % 4
		var suffix = ""
		if rfIdx1%2 == 1 {
			suffix = "_flipped"
		}
		cv.SetFillStyle("texture/building/" + wall.M.Name + suffix + ".png")
		cv.SetStrokeStyle(color.RGBA{R: 192, G: 192, B: 192, A: 24})
		cv.SetLineWidth(2)
		z := math.Min(math.Min(math.Min(rf.Z[0], rf.Z[1]), rf.Z[2]), rf.Z[3]) + float64(k*BuildingUnitHeight)*DZ

		cv.BeginPath()
		cv.LineTo(rf.X[rfIdx1], rf.Y[rfIdx1]-z)
		cv.LineTo(rf.X[rfIdx1], rf.Y[rfIdx1]-z-BuildingUnitHeight*DZ)
		cv.LineTo(rf.X[rfIdx2], rf.Y[rfIdx2]-z-BuildingUnitHeight*DZ)
		cv.LineTo(rf.X[rfIdx2], rf.Y[rfIdx2]-z)
		cv.ClosePath()
		cv.Fill()
		cv.Stroke()

		if wall.Windows {
			cv.SetFillStyle("texture/building/glass.png")
			cv.BeginPath()
			cv.LineTo((6*rf.X[rfIdx1]+1*rf.X[rfIdx2])/7, (6*rf.Y[rfIdx1]+1*rf.Y[rfIdx2])/7-z-BuildingUnitHeight*DZ*1/3)
			cv.LineTo((6*rf.X[rfIdx1]+1*rf.X[rfIdx2])/7, (6*rf.Y[rfIdx1]+1*rf.Y[rfIdx2])/7-z-BuildingUnitHeight*DZ*2/3)
			cv.LineTo((5*rf.X[rfIdx1]+2*rf.X[rfIdx2])/7, (5*rf.Y[rfIdx1]+2*rf.Y[rfIdx2])/7-z-BuildingUnitHeight*DZ*2/3)
			cv.LineTo((5*rf.X[rfIdx1]+2*rf.X[rfIdx2])/7, (5*rf.Y[rfIdx1]+2*rf.Y[rfIdx2])/7-z-BuildingUnitHeight*DZ*1/3)
			cv.ClosePath()
			cv.Fill()

			cv.BeginPath()
			cv.LineTo((4*rf.X[rfIdx1]+3*rf.X[rfIdx2])/7, (4*rf.Y[rfIdx1]+3*rf.Y[rfIdx2])/7-z-BuildingUnitHeight*DZ*1/3)
			cv.LineTo((4*rf.X[rfIdx1]+3*rf.X[rfIdx2])/7, (4*rf.Y[rfIdx1]+3*rf.Y[rfIdx2])/7-z-BuildingUnitHeight*DZ*2/3)
			cv.LineTo((3*rf.X[rfIdx1]+4*rf.X[rfIdx2])/7, (3*rf.Y[rfIdx1]+4*rf.Y[rfIdx2])/7-z-BuildingUnitHeight*DZ*2/3)
			cv.LineTo((3*rf.X[rfIdx1]+4*rf.X[rfIdx2])/7, (3*rf.Y[rfIdx1]+4*rf.Y[rfIdx2])/7-z-BuildingUnitHeight*DZ*1/3)
			cv.ClosePath()
			cv.Fill()

			cv.BeginPath()
			cv.LineTo((2*rf.X[rfIdx1]+5*rf.X[rfIdx2])/7, (2*rf.Y[rfIdx1]+5*rf.Y[rfIdx2])/7-z-BuildingUnitHeight*DZ*1/3)
			cv.LineTo((2*rf.X[rfIdx1]+5*rf.X[rfIdx2])/7, (2*rf.Y[rfIdx1]+5*rf.Y[rfIdx2])/7-z-BuildingUnitHeight*DZ*2/3)
			cv.LineTo((1*rf.X[rfIdx1]+6*rf.X[rfIdx2])/7, (1*rf.Y[rfIdx1]+6*rf.Y[rfIdx2])/7-z-BuildingUnitHeight*DZ*2/3)
			cv.LineTo((1*rf.X[rfIdx1]+6*rf.X[rfIdx2])/7, (1*rf.Y[rfIdx1]+6*rf.Y[rfIdx2])/7-z-BuildingUnitHeight*DZ*1/3)
			cv.ClosePath()
			cv.Fill()
		}

		if wall.Door {
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

func RenderBuildingRoof(cv *canvas.Canvas, roof *model.Roof, rf RenderedField, k int) {
	if roof == nil {
		return
	}
	z := math.Min(math.Min(math.Min(rf.Z[0], rf.Z[1]), rf.Z[2]), rf.Z[3]) + float64(k*BuildingUnitHeight)*DZ
	if !roof.Flat {
		cv.SetFillStyle("texture/building/" + roof.M.Name + ".png")
		cv.BeginPath()
		cv.LineTo(rf.X[0], rf.Y[0]-z)
		cv.LineTo(rf.X[1], rf.Y[1]-z)
		cv.LineTo(rf.X[2], rf.Y[2]-z)
		cv.LineTo(rf.X[3], rf.Y[3]-z)
		cv.ClosePath()
		cv.Fill()
	} else {
		cv.SetFillStyle("texture/building/" + roof.M.Name + ".png")
		cv.BeginPath()
		cv.LineTo(rf.X[0], rf.Y[0]-z)
		cv.LineTo(rf.X[1], rf.Y[1]-z)
		cv.LineTo(rf.X[2], rf.Y[2]-z)
		cv.LineTo(rf.X[3], rf.Y[3]-z)
		cv.ClosePath()
		cv.Fill()
	}
}
