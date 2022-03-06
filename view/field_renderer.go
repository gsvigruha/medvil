package view

import (
	"github.com/tfriedel6/canvas"
	"medvil/controller"
	"medvil/model/navigation"
	"medvil/renderer"
	//"fmt"
)

func RenderField(ic *ImageCache, cv *canvas.Canvas, rf renderer.RenderedField, f *navigation.Field, c *controller.Controller) {
	xMin, yMin, xMax, yMax := rf.BoundingBox()
	fieldImg := ic.Fic.RenderFieldOnBuffer(f, rf, c)
	cv.DrawImage(fieldImg, xMin, yMin, xMax-xMin, yMax-yMin)

	units := f.Building.BuildingUnits
	if len(units) > 0 {
		for k := 0; k < len(units); k++ {
			//rbu := RenderBuildingUnit(cv, &units[k], rf, k, c)
			z := float64((k+1)*BuildingUnitHeight) * DZ
			unitImg, rbu := ic.Bic.RenderBuildingUnitOnBuffer(&units[k], rf, k, c)
			cv.DrawImage(unitImg, xMin, yMin-z, float64(unitImg.Width()), float64(unitImg.Height()))
			c.AddRenderedBuildingUnit(&rbu)
		}
		if f.Building.RoofUnit != nil {
			numUnits := len(units)
			z := float64((numUnits+1)*BuildingUnitHeight) * DZ
			roofImg := ic.Bic.RenderBuildingRoofOnBuffer(f.Building.RoofUnit, rf, numUnits, c)
			cv.DrawImage(roofImg, xMin, yMin-z, float64(roofImg.Width()), float64(roofImg.Height()))
		}
	}

	if f.Road.T != nil {
		cv.SetFillStyle("texture/infra/" + f.Road.T.Name + ".png")
		rf.Draw(cv)
		cv.Fill()
	}
	if f.Plant != nil {
		tx := rf.X[0] - BufferW/2
		ty := rf.Y[2] - BufferH
		img := ic.Pic.RenderPlantOnBuffer(f.Plant, rf.Move(-tx, -ty), c)
		cv.DrawImage(img, tx, ty, BufferW, BufferH)
	}
	if f.Terrain.Resources.HasRealArtifacts() {
		cv.DrawImage("texture/terrain/barrel.png", rf.X[1]+44, rf.Y[2]-64, 32, 32)
	}
}
