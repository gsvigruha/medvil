package buildings

import (
	"github.com/tfriedel6/canvas"
	"math"
	"medvil/controller"
	"medvil/model/building"
	"medvil/model/materials"
	"medvil/renderer"
	"strconv"
)

const BuildingAnimationMaxPhase = 24

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
		RenderBuildingRoof(cv, building.ForgeBuildingRoof(extension.B, materials.GetMaterial("tile"), extension.Construction), rf, 1, c)
		if !extension.Construction {
			RenderChimney(cv, rf, 1, phase)
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
