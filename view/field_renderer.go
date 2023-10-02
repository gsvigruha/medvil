package view

import (
	"github.com/tfriedel6/canvas"
	"medvil/controller"
	"medvil/model/building"
	"medvil/model/navigation"
	"medvil/renderer"
	"sort"
)

func RenderField(phase int, ic *ImageCache, cv *canvas.Canvas, rf renderer.RenderedField, f *navigation.Field, c *controller.Controller) {
	xMin, yMin, _, _ := rf.BoundingBox()
	if phase == RenderPhaseField {
		fieldImg := ic.Fic.RenderFieldOnBuffer(f, rf, c)
		cv.DrawImage(fieldImg, xMin, yMin)

		if f.Deposit != nil {
			xMid, yMid := rf.MidScreenPoint()
			cv.DrawImage("texture/terrain/"+f.Deposit.T.Name+".png", xMid-60, yMid-40, 120, 80)
		}

		if f.Construction || f.Road != nil {
			RenderRoad(cv, rf, f, c)
		}

		if f.Building.GetBuilding() != nil && f.Building.GetBuilding().Plan.BuildingType == building.BuildingTypeMarket {
			if _, ok := f.Building.BuildingComponents[0].(*building.ExtensionUnit); !ok {
				cv.SetFillStyle("texture/building/market.png")
				rf.Draw(cv)
				cv.Fill()
			}
		}
	}

	if phase == RenderPhaseObjects {

		if f.Plant != nil && !f.Plant.IsTree() {
			RenderPlantOnBuffer(ic, cv, rf, f, c)
		}

		_, midY := rf.MidPoint()
		if f.Travellers != nil {
			// Travellers on the ground behind other ground objects
			sort.Slice(f.Travellers, func(i, j int) bool { return GetScreenY(f.Travellers[i], rf, c) < GetScreenY(f.Travellers[j], rf, c) })
			show := func(t *navigation.Traveller) bool { _, y := GetScreenXY(t, rf, c); return y < midY && t.FZ == 0 }
			RenderTravellers(ic, cv, f.Travellers, show, rf, c)
		}
		if f.Terrain.Resources.HasRealArtifacts() {
			cv.DrawImage("texture/terrain/barrel.png", rf.X[1]+40, rf.Y[2]-72, 32, 32)
		}
		if f.Plant != nil && f.Plant.IsTree() {
			RenderPlantOnBuffer(ic, cv, rf, f, c)
		}
		if f.Animal != nil && !f.Animal.Corralled {
			RenderAnimal(cv, rf, f, c)
		}
		if f.Statue != nil && !f.Statue.Construction {
			cv.DrawImage("icon/gui/infra/"+f.Statue.T.Name+".png", rf.X[0]-32, rf.Y[2]-80, 64, 64)
		}
		if f.Travellers != nil {
			// Travellers on the ground ahead other ground objects
			show := func(t *navigation.Traveller) bool { _, y := GetScreenXY(t, rf, c); return y >= midY && t.FZ == 0 }
			RenderTravellers(ic, cv, f.Travellers, show, rf, c)
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
			if c.ViewSettings.ShowHouseIcons {
				DrawHouseholdIcons(cv, rf, f, len(components), c)
			}
			if c.ViewSettings.ShowLabels {
				DrawLabels(cv, rf, f, len(components), c)
			}
		}

		if f.Travellers != nil {
			// Travellers on top of buildings
			show := func(t *navigation.Traveller) bool { return t.FZ > 0 }
			RenderTravellers(ic, cv, f.Travellers, show, rf, c)
		}

		if c.ViewSettings.ShowAllocatedFields {
			if f.Allocated {
				cv.DrawImage("icon/gui/flag.png", rf.X[0]-32, rf.Y[2]-80, 64, 64)
			}
		}
	}
}
