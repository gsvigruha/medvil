package controller

import (
	"github.com/tfriedel6/canvas"
	"medvil/model/navigation"
	"medvil/model/social"
	"medvil/renderer"
)

func DemolishToControlPanel(cp *ControlPanel, th *social.Townhall) {
	if th == nil {
		return
	}
	dc := &DemolishController{th: th}

	cp.SetDynamicPanel(dc)
	cp.C.ClickHandler = dc
}

type DemolishController struct {
	th *social.Townhall
}

func (dc *DemolishController) HandleClick(c *Controller, rf *renderer.RenderedField) bool {
	if rf.F.Building.GetBuilding() != nil || rf.F.Road != nil {
		dc.th.Household.Town.CreateDemolishTask(rf.F.Building.GetBuilding(), rf.F.Road, rf.F, c.Map)
		return true
	}
	return false
}

func (dc *DemolishController) GetActiveFields(c *Controller, rf *renderer.RenderedField) []navigation.FieldWithContext {
	var fields []navigation.FieldWithContext
	for _, h := range dc.th.Household.Town.GetHouseholds() {
		for _, coords := range h.Building.GetBuildingXYs(true) {
			fields = append(fields, c.Map.GetField(coords[0], coords[1]))
		}
	}
	for _, w := range dc.th.Household.Town.Walls {
		fields = append(fields, c.Map.GetField(w.Building.X, w.Building.Y))
	}
	for _, r := range dc.th.Household.Town.Roads {
		fields = append(fields, c.Map.GetField(r.X, r.Y))
	}
	if dc.th.Household.Town.Marketplace != nil {
		for _, coords := range dc.th.Household.Town.Marketplace.Building.GetBuildingXYs(true) {
			fields = append(fields, c.Map.GetField(coords[0], coords[1]))
		}
	}
	return fields
}

func (dc *DemolishController) CaptureMove(x, y float64) {}

func (dc *DemolishController) CaptureClick(x, y float64) {}

func (dc *DemolishController) Render(cv *canvas.Canvas) {}

func (dc *DemolishController) Clear() {}

func (dc *DemolishController) Refresh() {}
