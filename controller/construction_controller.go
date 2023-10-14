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
	top := ConstructionControllerTop * ControlPanelSY
	p.AddScaleLabel("tasks/building", 24, top, IconS, IconS, 4, float64(c.Progress)/float64(c.MaxProgress), false)
	var i = 1
	p.AddImageLabel("coin", 24, top+float64(IconH*2), IconS, IconS, gui.ImageLabelStyleRegular)
	for _, a := range c.Cost {
		ArtifactsToControlPanel(cp, p, i, a.A, a.Quantity, top+float64(IconH*2))
		i++
	}
	i = 1
	p.AddImageLabel("barrel", 24, top+float64(IconH*4), IconS, IconS, gui.ImageLabelStyleRegular)
	for _, a := range ArtifactOrder {
		if q, ok := c.Storage.Artifacts[a]; ok {
			ArtifactsToControlPanel(cp, p, i, a, q, top+float64(IconH*4))
			i++
		}
	}
	cp.HelperMessage("Construction in progress")
	cp.SetDynamicPanel(p)
}
