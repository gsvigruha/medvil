package view

import (
	"github.com/tfriedel6/canvas"
	"medvil/controller"
	"medvil/model/building"
	"medvil/model/navigation"
	"medvil/renderer"
)

func RenderField(ic *ImageCache, cv *canvas.Canvas, rf renderer.RenderedField, f *navigation.Field, c *controller.Controller) {
	xMin, yMin, xMax, yMax := rf.BoundingBox()
	fieldImg := ic.Fic.RenderFieldOnBuffer(f, rf, c)
	cv.DrawImage(fieldImg, xMin, yMin, xMax-xMin, yMax-yMin)

	components := f.Building.BuildingComponents
	if len(components) > 0 {
		for k := 0; k < len(components); k++ {
			if unit, ok := components[k].(*building.BuildingUnit); ok {
				unitImg, rbu, x, y := ic.Bic.RenderBuildingUnitOnBuffer(unit, rf, k, c)
				cv.DrawImage(unitImg, x, y, float64(unitImg.Width()), float64(unitImg.Height()))
				c.AddRenderedBuildingUnit(&rbu)
			} else if roof, ok := components[k].(*building.RoofUnit); ok {
				roofImg, x, y := ic.Bic.RenderBuildingRoofOnBuffer(roof, rf, k, c)
				cv.DrawImage(roofImg, x, y, float64(roofImg.Width()), float64(roofImg.Height()))
			}
		}
		workshop := c.ReverseReferences.BuildingToWorkshop[components[0].Building()]
		if workshop != nil {
			if workshop.Manufacture != nil {
				cv.SetFillStyle("#320")
				cv.FillRect(rf.X[2]-10-30, rf.Y[2]-DZ*3-2, 20, 20)
				cv.FillRect(rf.X[2]-10-30, rf.Y[2]-DZ*3-2, 30, 2)
				cv.DrawImage("icon/gui/tasks/"+workshop.Manufacture.Name+".png", rf.X[2]-8-30, rf.Y[2]-DZ*3, 16, 16)
			}
		}
	}

	if f.Construction || f.Road != nil {
		if f.Construction || f.Road.Construction {
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
	if f.Animal != nil && !f.Animal.Corralled {
		cv.DrawImage("texture/terrain/"+f.Animal.T.Name+".png", rf.X[0]-32, rf.Y[2]-64, 64, 64)
	}
	if f.Terrain.Resources.HasRealArtifacts() {
		cv.DrawImage("texture/terrain/barrel.png", rf.X[1]+44, rf.Y[2]-64, 32, 32)
	}
}
