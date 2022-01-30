package controller

import (
	"github.com/tfriedel6/canvas"
	"image/color"
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

type LandUseButton struct {
	b       gui.ButtonGUI
	fc      *FarmController
	useType uint8
}

func (b LandUseButton) Click() {
	b.fc.UseType = b.useType
}

func (b LandUseButton) Render(cv *canvas.Canvas) {
	b.b.Render(cv)
	if b.fc.UseType != b.useType {
		cv.SetFillStyle(color.RGBA{R: 64, G: 0, B: 0, A: 128})
		cv.FillRect(b.b.X, b.b.Y, 32, 32)
	}
}

func (b LandUseButton) Contains(x float64, y float64) bool {
	return b.b.Contains(x, y)
}

func FarmToControlPanel(cp *ControlPanel, farm *social.Farm) {
	hp := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: HouseholdControllerSY}
	fp := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop + HouseholdControllerSY, SX: ControlPanelSX, SY: ControlPanelDynamicPanelSY - HouseholdControllerSY}
	HouseholdToControlPanel(hp, &farm.Household)
	fc := &FarmController{householdPanel: hp, farmPanel: fp, farm: farm, UseType: economy.FarmFieldUseTypeBarren}

	fp.AddButton(LandUseButton{
		b:       gui.ButtonGUI{Texture: "terrain/grass", X: float64(10), Y: float64(HouseholdControllerGUIBottomY), SX: 32, SY: 32},
		fc:      fc,
		useType: economy.FarmFieldUseTypeBarren,
	})
	fp.AddButton(LandUseButton{
		b:       gui.ButtonGUI{Texture: "terrain/grain", X: float64(50), Y: float64(HouseholdControllerGUIBottomY), SX: 32, SY: 32},
		fc:      fc,
		useType: economy.FarmFieldUseTypeWheat,
	})
	fp.AddButton(LandUseButton{
		b:       gui.ButtonGUI{Texture: "terrain/vegetables", X: float64(90), Y: float64(HouseholdControllerGUIBottomY), SX: 32, SY: 32},
		fc:      fc,
		useType: economy.FarmFieldUseTypeVegetables,
	})
	fp.AddButton(LandUseButton{
		b:       gui.ButtonGUI{Icon: "artifacts/fruit", X: float64(140), Y: float64(HouseholdControllerGUIBottomY), SX: 32, SY: 32},
		fc:      fc,
		useType: economy.FarmFieldUseTypeOrchard,
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
