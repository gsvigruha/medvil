package controller

import (
	"medvil/model/navigation"
	"medvil/view/gui"
)

var FieldGUIY = 0.15

func FieldToControlPanel(cp *ControlPanel, f *navigation.Field) {
	p := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: ControlPanelDynamicPanelSY}
	p.AddTextureLabel("terrain/"+f.Terrain.T.Name, 24, FieldGUIY*ControlPanelSY, LargeIconS, LargeIconS)
	if f.Deposit != nil {
		p.AddTextureLabel("terrain/"+f.Deposit.T.Name, 24, FieldGUIY*ControlPanelSY+LargeIconD, 120, 80)
		p.AddTextLabel(ArtifactQStr(f.Deposit.Q), 160, FieldGUIY*ControlPanelSY+LargeIconD+40)
	}
	var aI = 0
	for a, q := range f.Terrain.Resources.Artifacts {
		ArtifactsToControlPanel(p, aI, a, q, ArtifactsGUIY*ControlPanelSY)
		aI++
	}
	cp.SetDynamicPanel(p)
}
