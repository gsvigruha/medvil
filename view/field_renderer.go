package view

import (
	"github.com/tfriedel6/canvas"
	"medvil/controller"
	"medvil/model"
	"medvil/model/navigation"
	"medvil/renderer"
	//"fmt"
)

func RenderField(ic *ImageCache, cv *canvas.Canvas, rf renderer.RenderedField, t, l, b, r uint8, m model.Map, f *navigation.Field, c *controller.Controller) {

	offsetX, offsetY := rf.Offset()
	fimg := ic.Fic.RenderFieldOnBuffer(f, rf)
	cv.DrawImage(fimg, offsetX, offsetY, BufferW, BufferH)

	units := f.Building.BuildingUnits
	for k := 0; k < len(units); k++ {
		rbu := RenderBuildingUnit(cv, &units[k], rf, k, c)
		c.AddRenderedBuildingUnit(&rbu)
	}
	roof := f.Building.RoofUnit
	RenderBuildingRoof(cv, roof, rf, len(units), c)
	if f.Road.T != nil {
		cv.SetFillStyle("texture/infra/" + f.Road.T.Name + ".png")
		rf.Draw(cv)
		cv.Fill()
	}
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
