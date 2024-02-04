package controller

import (
	"medvil/model/artifacts"
	"medvil/model/building"
	"medvil/view/gui"
)

const ConstructionControllerTop = 0.175

var ArtifactOrder = []*artifacts.Artifact{
	building.Cube,
	building.Board,
	building.Brick,
	building.Thatch,
	building.Tile,
	building.Textile,
	building.Paper,
}

func ConstructionToControlPanel(cp *ControlPanel, c *building.Construction) {
	p := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: HouseholdControllerSY}
	p.AddTextLabel("Construction", 24, ControlPanelSY*0.15)
	top := ConstructionControllerTop * ControlPanelSY
	p.AddScaleLabel("tasks/building", 24, top, IconS, IconS, 4, float64(c.Progress)/float64(c.MaxProgress), false,
		func(scaleStr string) {
			cp.HelperMessage("Building completion: "+scaleStr, false)
		})
	p.AddButton(&gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "demolish", X: float64(24 + IconW*7), Y: top, SX: IconS, SY: IconS, OnHoover: func() {
			cp.HelperMessage("Demolish construction", true)
		}},
		ClickImpl: func() {
			c.Delete()
		}})
	var i = 1
	p.AddLabel(&gui.ImageLabel{Icon: "coin", X: 24, Y: top + float64(IconH*2), SX: IconS, SY: IconS, Style: gui.ImageLabelStyleRegular, OnHoover: func() {
		cp.HelperMessage("Total materials needed for the building", false)
	}})
	for _, a := range c.Cost {
		ArtifactsToControlPanel(cp, p, i, a.A, a.Quantity, top+float64(IconH*2))
		i++
	}
	i = 1
	p.AddLabel(&gui.ImageLabel{Icon: "barrel", X: 24, Y: top + float64(IconH*4), SX: IconS, SY: IconS, Style: gui.ImageLabelStyleRegular, OnHoover: func() {
		cp.HelperMessage("Total materials stored at this construction", false)
	}})
	for _, a := range ArtifactOrder {
		if q, ok := c.Storage.Artifacts[a]; ok {
			ArtifactsToControlPanel(cp, p, i, a, q, top+float64(IconH*4))
			i++
		}
	}
	cp.SelectedHelperMessage("Construction in progress")
	cp.SetDynamicPanel(p)
}
