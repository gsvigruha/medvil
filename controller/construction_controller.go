package controller

import (
	"medvil/model/building"
	"medvil/view/gui"
)

const ConstructionControllerTop = 110

func ConstructionToControlPanel(cp *ControlPanel, c *building.Construction) {
	p := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: HouseholdControllerSY}
	p.AddScaleLabel("tasks/building", 10, ConstructionControllerTop, 32, 32, 4, float64(c.Progress)/float64(c.MaxProgress), false)
	var i = 0
	for _, a := range c.Cost {
		ArtifactsToControlPanel(cp, p, i, a.A, a.Quantity, ConstructionControllerTop+50)
		i++
	}
	i = 0
	for a, q := range c.Storage.Artifacts {
		ArtifactsToControlPanel(cp, p, i, a, q, ConstructionControllerTop+250)
		i++
	}
	cp.SetDynamicPanel(p)
}
