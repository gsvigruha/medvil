package controller

import (
	"medvil/model/navigation"
	"medvil/view/gui"
)

var FieldGUIY = 0.15

func FieldToControlPanel(cp *ControlPanel, f *navigation.Field) {
	p := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: ControlPanelDynamicPanelSY}
	p.AddTextureLabel("terrain/"+f.Terrain.T.Name, 10, FieldGUIY*ControlPanelSY, IconS, IconS)
	var aI = 0
	for a, q := range f.Terrain.Resources.Artifacts {
		ArtifactsToControlPanel(p, aI, a, q, ArtifactsGUIY*ControlPanelSY)
		aI++
	}
	cp.SetDynamicPanel(p)
}
