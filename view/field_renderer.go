package view

import (
	"github.com/tfriedel6/canvas"
	"medvil/controller"
	"medvil/model/building"
	"medvil/model/navigation"
	"medvil/renderer"
	"sort"
)

func renderPlant(ic *ImageCache, cv *canvas.Canvas, rf renderer.RenderedField, f *navigation.Field, c *controller.Controller) {
	tx := rf.X[0] - PlantBufferW/2
	ty := rf.Y[2] - PlantBufferH
	img := ic.Pic.RenderPlantOnBuffer(f.Plant, rf.Move(-tx, -ty), c)
	cv.DrawImage(img, tx, ty, PlantBufferW, PlantBufferH)
}

func RenderField(ic *ImageCache, cv *canvas.Canvas, rf renderer.RenderedField, f *navigation.Field, c *controller.Controller) {
	xMin, yMin, xMax, yMax := rf.BoundingBox()
	fieldImg := ic.Fic.RenderFieldOnBuffer(f, rf, c)
	cv.DrawImage(fieldImg, xMin, yMin, xMax-xMin, yMax-yMin)

	if f.Construction || f.Road != nil {
		RenderRoad(cv, rf, f, c)
	}

	if f.Plant != nil && !f.Plant.IsTree() {
		renderPlant(ic, cv, rf, f, c)
	}

	_, midY := rf.MidPoint()
	if f.Travellers != nil {
		sort.Slice(f.Travellers, func(i, j int) bool { return GetScreenY(f.Travellers[i], rf, c) < GetScreenY(f.Travellers[j], rf, c) })
		RenderTravellers(cv, f.Travellers, 0, midY, rf, c)
	}

	components := f.Building.BuildingComponents
	if len(components) > 0 {
		for k := 0; k < len(components); k++ {
			if unit, ok := components[k].(*building.BuildingUnit); ok {
				unitImg, rbu, x, y := ic.Bic.RenderBuildingUnitOnBuffer(unit, rf, k, c)
				cv.DrawImage(unitImg, x, y, float64(unitImg.Width()), float64(unitImg.Height()))
				c.AddRenderedBuildingPart(&rbu)
			} else if roof, ok := components[k].(*building.RoofUnit); ok {
				roofImg, rbr, x, y := ic.Bic.RenderBuildingRoofOnBuffer(roof, rf, k, c)
				cv.DrawImage(roofImg, x, y, float64(roofImg.Width()), float64(roofImg.Height()))
				c.AddRenderedBuildingPart(&rbr)
			} else if extension, ok := components[k].(*building.ExtensionUnit); ok {
				extensionImg, x, y := ic.Bic.RenderBuildingExtensionOnBuffer(extension, rf, k, c)
				cv.DrawImage(extensionImg, x, y, float64(extensionImg.Width()), float64(extensionImg.Height()))
			}
		}
	}

	if f.Terrain.Resources.HasRealArtifacts() {
		cv.DrawImage("texture/terrain/barrel.png", rf.X[1]+40, rf.Y[2]-72, 32, 32)
	}
	if f.Plant != nil && f.Plant.IsTree() {
		renderPlant(ic, cv, rf, f, c)
	}
	if f.Animal != nil && !f.Animal.Corralled {
		cv.DrawImage("texture/terrain/"+f.Animal.T.Name+".png", rf.X[0]-32, rf.Y[2]-64, 64, 64)
	}
	if f.Travellers != nil {
		RenderTravellers(cv, f.Travellers, midY, float64(c.H), rf, c)
	}
}
