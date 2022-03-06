package view

import (
	"github.com/tfriedel6/canvas"
	"medvil/controller"
	"medvil/model"
	"medvil/model/navigation"
	"medvil/renderer"
	//"fmt"
)

func RenderField(ic *ImageCache, cv *canvas.Canvas, rf renderer.RenderedField, m model.Map, f *navigation.Field, c *controller.Controller) {
	xMin, yMin, _, _ := rf.BoundingBox()
	fimg := ic.Fic.RenderFieldOnBuffer(f, rf, c)
	cv.DrawImage(fimg, xMin, yMin, float64(fimg.Width()), float64(fimg.Height()))

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
		img := ic.Pic.RenderPlantOnBuffer(f.Plant, rf.Move(-tx, -ty), c)
		cv.DrawImage(img, tx, ty, BufferW, BufferH)
	}
	if f.Terrain.Resources.HasRealArtifacts() {
		cv.DrawImage("texture/terrain/barrel.png", rf.X[1]+44, rf.Y[2]-64, 32, 32)
	}
}
