package controller

import (
	"github.com/tfriedel6/canvas"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/social"
	"medvil/renderer"
	"medvil/view/gui"
)

type MineController struct {
	householdPanel *gui.Panel
	minePanel      *gui.Panel
	UseType        uint8
	mine           *social.Mine
	cp             *ControlPanel
}

func (mc *MineController) GetUseType() uint8 {
	return mc.UseType
}

func (mc *MineController) SetUseType(ut uint8) {
	mc.UseType = ut
}

func MineToControlPanel(cp *ControlPanel, mine *social.Mine) {
	hp := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: HouseholdControllerSY}
	mp := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop + HouseholdControllerSY, SX: ControlPanelSX, SY: ControlPanelDynamicPanelSY - HouseholdControllerSY}
	HouseholdToControlPanel(cp, hp, mine.Household)
	mc := &MineController{householdPanel: hp, minePanel: mp, mine: mine, UseType: economy.MineFieldUseTypeNone, cp: cp}

	hcy := HouseholdControllerGUIBottomY * ControlPanelSY
	mp.AddButton(LandUseButton{
		b:       gui.ButtonGUI{Texture: "terrain/grass", X: float64(24), Y: hcy, SX: IconS, SY: IconS},
		luc:     mc,
		useType: economy.MineFieldUseTypeNone,
	})
	mp.AddButton(LandUseButton{
		b:       gui.ButtonGUI{Icon: "artifacts/stone", X: float64(24 + IconW*1), Y: hcy, SX: IconS, SY: IconS},
		luc:     mc,
		useType: economy.MineFieldUseTypeStone,
	})
	mp.AddButton(LandUseButton{
		b:       gui.ButtonGUI{Icon: "artifacts/clay", X: float64(24 + IconW*2), Y: hcy, SX: IconS, SY: IconS},
		luc:     mc,
		useType: economy.MineFieldUseTypeClay,
	})
	mp.AddButton(LandUseButton{
		b:       gui.ButtonGUI{Icon: "artifacts/iron_ore", X: float64(24 + IconW*3), Y: hcy, SX: IconS, SY: IconS},
		luc:     mc,
		useType: economy.MineFieldUseTypeIron,
	})
	mp.AddButton(LandUseButton{
		b:       gui.ButtonGUI{Icon: "artifacts/gold_ore", X: float64(24 + IconW*4), Y: hcy, SX: IconS, SY: IconS},
		luc:     mc,
		useType: economy.MineFieldUseTypeGold,
	})

	cp.SetDynamicPanel(mc)
	cp.C.ClickHandler = mc
}

func (mc *MineController) CaptureClick(x, y float64) {
	mc.householdPanel.CaptureClick(x, y)
	mc.minePanel.CaptureClick(x, y)
}

func (mc *MineController) Render(cv *canvas.Canvas) {
	mc.householdPanel.Render(cv)
	mc.minePanel.Render(cv)
}

func (mc *MineController) Clear() {}

func (mc *MineController) Refresh() {
	mc.householdPanel.Clear()
	HouseholdToControlPanel(mc.cp, mc.householdPanel, mc.mine.Household)
}

func (mc *MineController) GetActiveFields(c *Controller, rf *renderer.RenderedField) []navigation.FieldWithContext {
	fields := mc.mine.GetFields()
	if social.CheckMineUseType(mc.UseType, rf.F) {
		fields = append(fields, social.MineLand{
			X:       rf.F.X,
			Y:       rf.F.Y,
			UseType: mc.UseType,
			F:       rf.F,
		})
	}
	return fields
}

func (mc *MineController) HandleClick(c *Controller, rf *renderer.RenderedField) bool {
	for i := range mc.mine.Land {
		l := &mc.mine.Land[i]
		if l.F.X == rf.F.X && l.F.Y == rf.F.Y {
			if mc.UseType != economy.MineFieldUseTypeNone {
				if social.CheckMineUseType(mc.UseType, l.F) {
					l.UseType = mc.UseType
				}
			} else {
				// Disallocate land
				mc.mine.Land = append(mc.mine.Land[:i], mc.mine.Land[i+1:]...)
				rf.F.Allocated = false
			}
			return true
		}
	}
	if social.CheckMineUseType(mc.UseType, rf.F) && !rf.F.Allocated && mc.UseType != economy.MineFieldUseTypeNone {
		if social.CheckMineUseType(mc.UseType, rf.F) && mc.mine.FieldWithinDistance(rf.F) {
			mc.mine.Land = append(mc.mine.Land, social.MineLand{
				X:       rf.F.X,
				Y:       rf.F.Y,
				UseType: mc.UseType,
				F:       rf.F,
			})
			rf.F.Allocated = true
			return true
		}
	}
	return false
}
