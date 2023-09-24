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
		p.AddLabel(&gui.CustomImageLabel{RenderFn: func(cv *canvas.Canvas) {
			if f.Deposit != nil {
				cv.DrawImage("texture/terrain/"+f.Deposit.T.Name+".png", 24, FieldGUIY*ControlPanelSY+LargeIconD, 120, 80)
			}
		}})
		p.AddTextLabel(ArtifactQStr(f.Deposit.Q), 160, FieldGUIY*ControlPanelSY+LargeIconD+40)
	}
	if f.Plant != nil {
		if f.Plant.T.TreeT != nil {
			p.AddImageLabel("infra/"+strings.Replace(f.Plant.T.Name, " ", "_", -1), 24, FieldGUIY*ControlPanelSY+LargeIconD, LargeIconD, LargeIconD, gui.ImageLabelStyleRegular)
		} else {
			p.AddLabel(&gui.CustomImageLabel{RenderFn: func(cv *canvas.Canvas) {
				if f.Plant != nil {
					cv.DrawImage("texture/terrain/"+f.Plant.T.Name+".png", 24, FieldGUIY*ControlPanelSY+LargeIconD, 120, 108)
				}
			}})
		}
	}
	if f.Animal != nil {
		p.AddLabel(&gui.CustomImageLabel{RenderFn: func(cv *canvas.Canvas) {
			if f.Animal != nil {
				cv.DrawImage("texture/terrain/"+f.Animal.T.Name+"_0.png", 24, FieldGUIY*ControlPanelSY+LargeIconD, 64, 64)
			}
		}})
		p.AddTextLabel(strconv.Itoa(int(f.Animal.AgeYears(cp.C.Map.Calendar)))+" years", 100, FieldGUIY*ControlPanelSY+LargeIconD+32)
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
		ArtifactsToControlPanel(p, aI, a, q, ArtifactsGUIY*ControlPanelSY)
		aI++
	}
	cp.SetDynamicPanel(p)
}
