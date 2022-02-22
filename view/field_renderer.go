package view

import (
	"github.com/tfriedel6/canvas"
	"image/color"
	"medvil/controller"
	"medvil/model"
	"medvil/model/navigation"
	"medvil/renderer"
	//"fmt"
)

func RenderField(ic *ImageCache, cv *canvas.Canvas, rf renderer.RenderedField, t, l, b, r uint8, m model.Map, f *navigation.Field, c *controller.Controller) {

	cv.SetFillStyle("texture/terrain/" + f.Terrain.T.Name + ".png")

	rf.Draw(cv)
	cv.Fill()

	if (f.SE + f.SW) > (f.NE + f.NW) {
		slope := (f.SE + f.SW) - (f.NE + f.NW)
		cv.SetFillStyle(color.RGBA{R: 255, G: 255, B: 255, A: slope * 4})
		rf.Draw(cv)
		cv.Fill()
	} else if (f.SE + f.SW) < (f.NE + f.NW) {
		slope := (f.NE + f.NW) - (f.SE + f.SW)
		cv.SetFillStyle(color.RGBA{R: 0, G: 0, B: 0, A: slope * 16})
		rf.Draw(cv)
		cv.Fill()
	}

	units := f.Building.BuildingUnits
	for k := 0; k < len(units); k++ {
		rbu := RenderBuildingUnit(cv, &units[k], rf, k, c)
		c.AddRenderedBuildingUnit(&rbu)
	}
	roof := f.Building.RoofUnit
	RenderBuildingRoof(cv, roof, rf, len(units), c)
	if f.Plant != nil {
		//RenderPlant(cv, f.Plant, rf, c)
		tx := rf.X[0] - BufferW/2
		ty := rf.Y[2] - BufferH
		img := ic.RenderPlantOnBuffer(f.Plant, rf.Move(-tx, -ty), c)
		cv.DrawImage(img, tx, ty, BufferW, BufferH)
	}
	if f.Terrain.Resources.HasRealArtifacts() {
		cv.DrawImage("texture/terrain/barrel.png", rf.X[1]+44, rf.Y[2]-64, 32, 32)
	}
}
