package controller

import (
	"github.com/tfriedel6/canvas"
	"medvil/model/navigation"
	"medvil/view/gui"
	"strconv"
	"strings"
)

var FieldGUIY = 0.15

func FieldToControlPanel(cp *ControlPanel, f *navigation.Field) {
	p := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: ControlPanelDynamicPanelSY}
	p.AddTextureLabel("terrain/"+f.Terrain.T.Name, 24, FieldGUIY*ControlPanelSY, LargeIconS, LargeIconS)
	if f.Deposit != nil {
		p.AddImageLabel("terrain/"+f.Deposit.T.Name, 24, FieldGUIY*ControlPanelSY+LargeIconD, LargeIconS, LargeIconS, gui.ImageLabelStyleRegular)
		p.AddTextLabel(ArtifactQStr(f.Deposit.Q), 24+LargeIconD, FieldGUIY*ControlPanelSY+LargeIconD+LargeIconS/2)
	}
	if f.Plant != nil {
		if f.Plant.T.TreeT != nil {
			p.AddImageLabel("infra/"+strings.Replace(f.Plant.T.Name, " ", "_", -1), 24, FieldGUIY*ControlPanelSY+LargeIconD, LargeIconS, LargeIconS, gui.ImageLabelStyleRegular)
		} else {
			p.AddLabel(&gui.CustomImageLabel{RenderFn: func(cv *canvas.Canvas) {
				if f.Plant != nil {
					cv.DrawImage("texture/terrain/"+f.Plant.T.Name+".png", 24, FieldGUIY*ControlPanelSY+LargeIconD, LargeIconS, LargeIconS)
				}
			}})
		}
		p.AddTextLabel(strconv.Itoa(int(f.Plant.AgeYears(cp.C.Map.Calendar)))+" years", 24+LargeIconD, FieldGUIY*ControlPanelSY+LargeIconD+LargeIconS/2)
	}
	if f.Animal != nil {
		p.AddLabel(&gui.CustomImageLabel{RenderFn: func(cv *canvas.Canvas) {
			if f.Animal != nil {
				cv.DrawImage("texture/terrain/"+f.Animal.T.Name+"_0.png", 24, FieldGUIY*ControlPanelSY+LargeIconD, LargeIconS, LargeIconS)
			}
		}})
		p.AddTextLabel(strconv.Itoa(int(f.Animal.AgeYears(cp.C.Map.Calendar)))+" years", 24+LargeIconD, FieldGUIY*ControlPanelSY+LargeIconD+LargeIconS/2)
	}
	if f.Road != nil && !f.Road.Construction {
		if f.Road.Broken {
			p.AddTextureLabel("infra/"+f.Road.T.Name+"_broken", 24, FieldGUIY*ControlPanelSY+LargeIconD*2, LargeIconS, LargeIconS)
		} else {
			p.AddTextureLabel("infra/"+f.Road.T.Name, 24, FieldGUIY*ControlPanelSY+LargeIconD*2, LargeIconS, LargeIconS)
		}
	}
	var aI = 0
	for a, q := range f.Terrain.Resources.Artifacts {
		ArtifactsToControlPanel(cp, p, aI, a, q, ArtifactsGUIY*ControlPanelSY)
		aI++
	}
	cp.SetDynamicPanel(p)
}
