package controller

import (
	"github.com/tfriedel6/canvas"
	"medvil/model/economy"
	"medvil/model/social"
	"medvil/renderer"
	"medvil/view/gui"
)

const HouseholdControllerSY = 500

type FarmController struct {
	householdPanel *gui.Panel
	farmPanel      *gui.Panel
	UseType        uint8
	farm           *social.Farm
}

func (fc *FarmController) GetUseType() uint8 {
	return fc.UseType
}

func (fc *FarmController) SetUseType(ut uint8) {
	fc.UseType = ut
}

func FarmToControlPanel(cp *ControlPanel, farm *social.Farm) {
	hp := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: HouseholdControllerSY}
	fp := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop + HouseholdControllerSY, SX: ControlPanelSX, SY: ControlPanelDynamicPanelSY - HouseholdControllerSY}
	HouseholdToControlPanel(hp, &farm.Household)
	fc := &FarmController{householdPanel: hp, farmPanel: fp, farm: farm, UseType: economy.FarmFieldUseTypeBarren}

	fp.AddButton(LandUseButton{
		b:       gui.ButtonGUI{Texture: "terrain/grass", X: float64(10), Y: float64(HouseholdControllerGUIBottomY), SX: 32, SY: 32},
		luc:     fc,
		useType: economy.FarmFieldUseTypeBarren,
	})
	fp.AddButton(LandUseButton{
		b:       gui.ButtonGUI{Icon: "artifacts/grain", X: float64(50), Y: float64(HouseholdControllerGUIBottomY), SX: 32, SY: 32},
		luc:     fc,
		useType: economy.FarmFieldUseTypeWheat,
	})
	fp.AddButton(LandUseButton{
		b:       gui.ButtonGUI{Icon: "artifacts/vegetable", X: float64(90), Y: float64(HouseholdControllerGUIBottomY), SX: 32, SY: 32},
		luc:     fc,
		useType: economy.FarmFieldUseTypeVegetables,
	})
	fp.AddButton(LandUseButton{
		b:       gui.ButtonGUI{Icon: "artifacts/fruit", X: float64(130), Y: float64(HouseholdControllerGUIBottomY), SX: 32, SY: 32},
		luc:     fc,
		useType: economy.FarmFieldUseTypeOrchard,
	})
	fp.AddButton(LandUseButton{
		b:       gui.ButtonGUI{Icon: "artifacts/log", X: float64(170), Y: float64(HouseholdControllerGUIBottomY), SX: 32, SY: 32},
		luc:     fc,
		useType: economy.FarmFieldUseTypeForestry,
	})

	cp.SetDynamicPanel(fc)
	cp.C.ClickHandler = fc
}

func (fc *FarmController) CaptureClick(x, y float64) {
	fc.householdPanel.CaptureClick(x, y)
	fc.farmPanel.CaptureClick(x, y)
}

func (fc *FarmController) Render(cv *canvas.Canvas) {
	fc.householdPanel.Render(cv)
	fc.farmPanel.Render(cv)
}

func (fc *FarmController) Clear() {}

func (fc *FarmController) Refresh() {
	fc.householdPanel.Clear()
	HouseholdToControlPanel(fc.householdPanel, &fc.farm.Household)
}

func (fc *FarmController) HandleClick(c *Controller, rf *renderer.RenderedField) bool {
	var owns = false
	for i := range fc.farm.Land {
		l := &fc.farm.Land[i]
		if l.F.X == rf.F.X && l.F.Y == rf.F.Y {
			if fc.UseType != economy.FarmFieldUseTypeBarren {
				if l.F.Arable() {
					l.UseType = fc.UseType
				}
			} else {
				// Disallocate land
				fc.farm.Land = append(fc.farm.Land[:i], fc.farm.Land[i+1:]...)
				rf.F.Allocated = false
			}
			owns = true
			break
		}
	}
	if !owns && !rf.F.Allocated && fc.UseType != economy.FarmFieldUseTypeBarren {
		fc.farm.Land = append(fc.farm.Land, social.FarmLand{
			X:       rf.F.X,
			Y:       rf.F.Y,
			UseType: fc.UseType,
			F:       rf.F,
		})
		rf.F.Allocated = true
	}
	return true
}
