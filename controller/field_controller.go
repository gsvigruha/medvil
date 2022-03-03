package controller

import (
	"medvil/model/navigation"
	"medvil/view/gui"
)

func FieldToControlPanel(cp *ControlPanel, f *navigation.Field) {
	p := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: ControlPanelDynamicPanelSY}
	p.AddTextureLabel("terrain/"+f.Terrain.T.Name, 10, 70, 32, 32)
	var aI = 0
	for a, q := range f.Terrain.Resources.Artifacts {
		ArtifactsToControlPanel(p, aI, a, q, ArtifactsGUIY)
		aI++
	}
	cp.SetDynamicPanel(p)
}
