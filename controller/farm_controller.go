package controller

import (
	"github.com/tfriedel6/canvas"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/social"
	"medvil/renderer"
	"medvil/view/gui"
)

const FarmFieldUseTypeDisallocate uint8 = 255

type FarmController struct {
	householdPanel *gui.Panel
	farmPanel      *gui.Panel
	UseType        uint8
	farm           *social.Farm
	cp             *ControlPanel
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
	HouseholdToControlPanel(cp, hp, farm.Household)
	fc := &FarmController{householdPanel: hp, farmPanel: fp, farm: farm, UseType: economy.FarmFieldUseTypeBarren, cp: cp}

	hcy := HouseholdControllerGUIBottomY * ControlPanelSY
	fp.AddButton(&LandUseButton{
		b:       gui.ButtonGUI{Icon: "clear_land", X: float64(24 + IconW*0), Y: hcy, SX: IconS, SY: IconS},
		luc:     fc,
		useType: economy.FarmFieldUseTypeBarren,
	})
	fp.AddButton(&LandUseButton{
		b:       gui.ButtonGUI{Icon: "artifacts/grain", X: float64(24 + IconW*1), Y: hcy, SX: IconS, SY: IconS},
		luc:     fc,
		useType: economy.FarmFieldUseTypeWheat,
	})
	fp.AddButton(&LandUseButton{
		b:       gui.ButtonGUI{Icon: "artifacts/vegetable", X: float64(24 + IconW*2), Y: hcy, SX: IconS, SY: IconS},
		luc:     fc,
		useType: economy.FarmFieldUseTypeVegetables,
	})
	fp.AddButton(&LandUseButton{
		b:       gui.ButtonGUI{Icon: "artifacts/fruit", X: float64(24 + IconW*3), Y: hcy, SX: IconS, SY: IconS},
		luc:     fc,
		useType: economy.FarmFieldUseTypeOrchard,
	})
	fp.AddButton(&LandUseButton{
		b:       gui.ButtonGUI{Icon: "artifacts/sheep", X: float64(24 + IconW*0), Y: hcy + float64(IconH), SX: IconS, SY: IconS},
		luc:     fc,
		useType: economy.FarmFieldUseTypePasture,
	})
	fp.AddButton(&LandUseButton{
		b:       gui.ButtonGUI{Icon: "artifacts/log", X: float64(24 + IconW*1), Y: hcy + float64(IconH), SX: IconS, SY: IconS},
		luc:     fc,
		useType: economy.FarmFieldUseTypeForestry,
	})
	fp.AddButton(&LandUseButton{
		b:       gui.ButtonGUI{Icon: "artifacts/reed", X: float64(24 + IconW*2), Y: hcy + float64(IconH), SX: IconS, SY: IconS},
		luc:     fc,
		useType: economy.FarmFieldUseTypeReed,
	})
	fp.AddButton(&LandUseButton{
		b:       gui.ButtonGUI{Icon: "artifacts/herb", X: float64(24 + IconW*3), Y: hcy + float64(IconH), SX: IconS, SY: IconS},
		luc:     fc,
		useType: economy.FarmFieldUseTypeHerb,
	})
	fp.AddButton(&LandUseButton{
		b:       gui.ButtonGUI{Icon: "cancel", X: float64(24 + IconW*0), Y: hcy + float64(IconH*2), SX: IconS, SY: IconS},
		luc:     fc,
		useType: FarmFieldUseTypeDisallocate,
	})
	fc.RefreshLandUseButtons()

	cp.SetDynamicPanel(fc)
	cp.C.ClickHandler = fc
}

func (fc *FarmController) RefreshLandUseButtons() {
	landDist := fc.farm.GetLandDistribution()
	for _, b := range fc.farmPanel.Buttons {
		if lub, ok := b.(*LandUseButton); ok {
			lub.cnt = landDist[lub.useType]
		}
	}
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
	HouseholdToControlPanel(fc.cp, fc.householdPanel, fc.farm.Household)
}

func (fc *FarmController) GetActiveFields(c *Controller, rf *renderer.RenderedField) []navigation.FieldWithContext {
	fields := fc.farm.GetFields()
	if fc.farm.FieldUsableFor(fc.cp.C.Map, rf.F, fc.UseType) && !rf.F.Allocated {
		fields = append(fields, social.FarmLand{
			X:       rf.F.X,
			Y:       rf.F.Y,
			UseType: fc.UseType,
			F:       rf.F,
		})
	}
	return fields
}

func (fc *FarmController) HandleClick(c *Controller, rf *renderer.RenderedField) bool {
	for i := range fc.farm.Land {
		l := &fc.farm.Land[i]
		if l.F.X == rf.F.X && l.F.Y == rf.F.Y {
			if fc.UseType == FarmFieldUseTypeDisallocate {
				// Disallocate land
				fc.farm.Land = append(fc.farm.Land[:i], fc.farm.Land[i+1:]...)
				rf.F.Allocated = false
			} else if fc.farm.FieldUsableFor(c.Map, l.F, fc.UseType) {
				l.UseType = fc.UseType
			}
			fc.RefreshLandUseButtons()
			return true
		}
	}
	if fc.UseType != FarmFieldUseTypeDisallocate && !rf.F.Allocated && fc.farm.FieldUsableFor(c.Map, rf.F, fc.UseType) && fc.farm.FieldWithinDistance(rf.F) {
		fc.farm.Land = append(fc.farm.Land, social.FarmLand{
			X:       rf.F.X,
			Y:       rf.F.Y,
			UseType: fc.UseType,
			F:       rf.F,
		})
		rf.F.Allocated = true
		fc.RefreshLandUseButtons()
		return true
	}
	return false
}
