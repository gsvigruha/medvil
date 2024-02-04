package controller

import (
	"github.com/tfriedel6/canvas"
	"math"
	"medvil/model/artifacts"
	"medvil/model/navigation"
	"medvil/view/gui"
	"path/filepath"
	"strconv"
	"strings"
)

var FieldGUIY = 0.15

func plantDeathRateStr(rate float64) string {
	return strconv.Itoa(int(math.Pow(1.0-rate, 30*9) * 100))
}

func treeDeathRateStr(rate float64) string {
	return strconv.Itoa(int(math.Pow(1.0-rate, 30*12*10) * 100))
}

func artifactQStr(q uint16) string {
	var qStr = strconv.Itoa(int(q))
	if q == artifacts.InfiniteQuantity {
		qStr = "infinite"
	}
	return qStr
}

func FieldToControlPanel(cp *ControlPanel, f *navigation.Field) {
	p := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: ControlPanelDynamicPanelSY}
	p.AddTextureLabel("terrain/"+f.Terrain.T.Name, 24, FieldGUIY*ControlPanelSY, LargeIconS, LargeIconS)
	if f.Deposit != nil {
		p.AddImageLabel("terrain/"+f.Deposit.T.Name, 24, FieldGUIY*ControlPanelSY+LargeIconD, LargeIconS, LargeIconS, gui.ImageLabelStyleRegular)
		p.AddTextLabel(artifactQStr(f.Deposit.Q)+" "+f.Deposit.T.A.Name, 24+LargeIconD, FieldGUIY*ControlPanelSY+LargeIconD+LargeIconS/2)
	} else {
		if f.Plantable(false) {
			p.AddTextLabel("Tree soil quality: "+treeDeathRateStr(cp.C.Map.TreeDeathRate(f))+"%", 24+LargeIconD, FieldGUIY*ControlPanelSY+LargeIconS*0.8)
		}
		if f.Arable() {
			p.AddTextLabel("Plant soil quality: "+plantDeathRateStr(cp.C.Map.PlantDeathRate(f))+"%", 24+LargeIconD, FieldGUIY*ControlPanelSY+LargeIconS*0.4)
		}
	}
	if f.Plant != nil {
		if f.Plant.T.TreeT != nil {
			p.AddImageLabel("infra/"+strings.Replace(f.Plant.T.Name, " ", "_", -1), 24, FieldGUIY*ControlPanelSY+LargeIconD, LargeIconS, LargeIconS, gui.ImageLabelStyleRegular)
		} else {
			p.AddLabel(&gui.CustomImageLabel{RenderFn: func(cv *canvas.Canvas) {
				if f.Plant != nil {
					cv.DrawImage(filepath.FromSlash("texture/terrain/"+f.Plant.T.Name+".png"), 24, FieldGUIY*ControlPanelSY+LargeIconD, LargeIconS, LargeIconS)
				}
			}})
		}
		p.AddTextLabel(strconv.Itoa(int(f.Plant.AgeYears(cp.C.Map.Calendar)))+" years", 24+LargeIconD, FieldGUIY*ControlPanelSY+LargeIconD+LargeIconD/2)
	}
	if f.Animal != nil {
		p.AddLabel(&gui.CustomImageLabel{RenderFn: func(cv *canvas.Canvas) {
			if f.Animal != nil {
				cv.DrawImage(filepath.FromSlash("texture/terrain/"+f.Animal.T.Name+"_0.png"), 24, FieldGUIY*ControlPanelSY+LargeIconD, LargeIconS, LargeIconS)
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
		p.AddTextLabel(f.Road.T.DisplayName+", Speed: "+strconv.FormatFloat(f.GetSpeed(), 'f', -1, 64), 24+LargeIconD, FieldGUIY*ControlPanelSY+LargeIconD*2+LargeIconS/2)
	}
	var aI = 0
	for a, q := range f.Terrain.Resources.Artifacts {
		ArtifactsToControlPanel(cp, p, aI, a, q, ArtifactsGUIY*ControlPanelSY)
		aI++
	}
	cp.SetDynamicPanel(p)
}
