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
			unitImg, rbu, x, y := ic.Bic.RenderBuildingUnitOnBuffer(&units[k], rf, k, c)
			cv.DrawImage(unitImg, x, y, float64(unitImg.Width()), float64(unitImg.Height()))
			c.AddRenderedBuildingUnit(&rbu)
		}
		if f.Building.RoofUnit != nil {
			roofImg, x, y := ic.Bic.RenderBuildingRoofOnBuffer(f.Building.RoofUnit, rf, len(units), c)
			cv.DrawImage(roofImg, x, y, float64(roofImg.Width()), float64(roofImg.Height()))
		}
		workshop := c.ReverseReferences.BuildingToWorkshop[units[0].B]
		if workshop != nil {
			if workshop.Manufacture != nil {
				cv.DrawImage("icon/gui/tasks/"+workshop.Manufacture.Name+".png", rf.X[2]-8, rf.Y[2]-DZ*3-8, 16, 16)
			}
		}
	}

	if f.Road != nil {
		if f.Road.Construction {
			cv.SetFillStyle("texture/infra/construction.png")
		} else {
			cv.SetFillStyle("texture/infra/" + f.Road.T.Name + ".png")
		}
		rf.Draw(cv)
		cv.Fill()
	}
	if f.Plant != nil {
		tx := rf.X[0] - PlantBufferW/2
		ty := rf.Y[2] - PlantBufferH
		img := ic.Pic.RenderPlantOnBuffer(f.Plant, rf.Move(-tx, -ty), c)
		cv.DrawImage(img, tx, ty, PlantBufferW, PlantBufferH)
	}
	if f.Terrain.Resources.HasRealArtifacts() {
		cv.DrawImage("texture/terrain/barrel.png", rf.X[1]+44, rf.Y[2]-64, 32, 32)
	}
}
